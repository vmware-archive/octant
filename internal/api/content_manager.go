/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vmware-tanzu/octant/internal/event"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	RequestSetContentPath = "setContentPath"
	RequestSetNamespace   = "setNamespace"
)

// ContentManagerOption is an option for configuring ContentManager.
type ContentManagerOption func(manager *ContentManager)

// ContentGenerateFunc is a function that generates content. It returns `rerun=true`
// if the action should be be immediately rerun.
type ContentGenerateFunc func(ctx context.Context, state octant.State) (component.ContentResponse, bool, error)

// WithContentGenerator configures the content generate function.
func WithContentGenerator(fn ContentGenerateFunc) ContentManagerOption {
	return func(manager *ContentManager) {
		manager.contentGenerateFunc = fn
	}
}

// WithContentGeneratorPoller configures the poller.
func WithContentGeneratorPoller(poller Poller) ContentManagerOption {
	return func(manager *ContentManager) {
		manager.poller = poller
	}
}

// ContentManager manages content for websockets.
type ContentManager struct {
	moduleManager       module.ManagerInterface
	logger              log.Logger
	contentGenerateFunc ContentGenerateFunc
	poller              Poller
	updateContentCh     chan struct{}
}

// NewContentManager creates an instance of ContentManager.
func NewContentManager(moduleManager module.ManagerInterface, logger log.Logger, options ...ContentManagerOption) *ContentManager {
	cm := &ContentManager{
		moduleManager:   moduleManager,
		logger:          logger,
		poller:          NewInterruptiblePoller("content"),
		updateContentCh: make(chan struct{}, 1),
	}
	cm.contentGenerateFunc = cm.generateContent

	for _, option := range options {
		option(cm)
	}

	return cm
}

var _ StateManager = (*ContentManager)(nil)

// Start starts the manager.
func (cm *ContentManager) Start(ctx context.Context, state octant.State, s OctantClient) {
	defer func() {
		close(cm.updateContentCh)
	}()

	updateCancel := state.OnContentPathUpdate(func(contentPath string) {
		cm.updateContentCh <- struct{}{}
	})
	defer updateCancel()

	cm.poller.Run(ctx, cm.updateContentCh, cm.runUpdate(state, s), event.DefaultScheduleDelay)
}

func (cm *ContentManager) runUpdate(state octant.State, s OctantClient) PollerFunc {
	return func(ctx context.Context) (bool, error) {
		contentPath := state.GetContentPath()
		if contentPath == "" {
			return false, nil
		}

		contentResponse, _, err := cm.contentGenerateFunc(ctx, state)
		if err != nil {
			cm.logger.
				WithErr(err).
				With("content-path", contentPath).
				Errorf("generate content")
			return false, err
		}

		if ctx.Err() == nil {
			s.Send(CreateContentEvent(contentResponse, state.GetNamespace(), contentPath, state.GetQueryParams()))
		}

		return false, nil
	}
}

func (cm *ContentManager) generateContent(ctx context.Context, state octant.State) (component.ContentResponse, bool, error) {
	contentPath := state.GetContentPath()
	logger := cm.logger.With("contentPath", contentPath)

	now := time.Now()
	defer func() {
		logger.With("elapsed", time.Since(now)).Debugf("generating content")
	}()

	m, ok := cm.moduleManager.ModuleForContentPath(contentPath)
	if !ok {
		return component.EmptyContentResponse, false, fmt.Errorf("unable to find module for content path %q", contentPath)
	}
	modulePath := strings.TrimPrefix(contentPath, m.Name())
	options := module.ContentOptions{
		LabelSet: FiltersToLabelSet(state.GetFilters()),
	}
	contentResponse, err := m.Content(ctx, modulePath, options)
	if err != nil {
		if nfe, ok := err.(notFound); ok && nfe.NotFound() {
			logger.Debugf("path not found, redirecting to parent")
			state.SetContentPath(notFoundRedirectPath(contentPath))
			return component.EmptyContentResponse, true, nil
		} else {
			return component.EmptyContentResponse, false, fmt.Errorf("generate content: %w", err)
		}
	}

	return contentResponse, false, nil
}

// Handlers returns a slice of client request handlers.
func (cm *ContentManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: RequestSetContentPath,
			Handler:     cm.SetContentPath,
		},
		{
			RequestType: RequestSetNamespace,
			Handler:     cm.SetNamespace,
		},
	}
}

// SetQueryParams sets the current query params.
func (cm *ContentManager) SetQueryParams(state octant.State, payload action.Payload) error {
	if params, ok := payload["params"].(map[string]interface{}); ok {
		// handle filters
		if filters, ok := params["filters"]; ok {
			list, err := FiltersFromQueryParams(filters)
			if err != nil {
				return fmt.Errorf("extract filters from query params: %w", err)
			}
			state.SetFilters(list)
		}
	}

	return nil
}

// SetNamespace sets the current namespace.
func (cm *ContentManager) SetNamespace(state octant.State, payload action.Payload) error {
	namespace, err := payload.String("namespace")
	if err != nil {
		return fmt.Errorf("extract namespace from payload: %w", err)
	}
	state.SetNamespace(namespace)
	return nil
}

// SetContentPath sets the current content path.
func (cm *ContentManager) SetContentPath(state octant.State, payload action.Payload) error {
	contentPath, err := payload.String("contentPath")
	if err != nil {
		return fmt.Errorf("extract contentPath from payload: %w", err)
	}
	if err := cm.SetQueryParams(state, payload); err != nil {
		return fmt.Errorf("extract query params from payload: %w", err)
	}

	state.SetContentPath(contentPath)
	return nil
}

type notFound interface {
	NotFound() bool
	Path() string
}

// CreateContentEvent creates a content event.
func CreateContentEvent(contentResponse component.ContentResponse, namespace, contentPath string, queryParams map[string][]string) octant.Event {
	return octant.Event{
		Type: octant.EventTypeContent,
		Data: map[string]interface{}{
			"content":     contentResponse,
			"namespace":   namespace,
			"contentPath": contentPath,
			"queryParams": queryParams,
		},
	}
}

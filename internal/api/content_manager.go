/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/vmware-tanzu/octant/internal/util/json"

	ocontext "github.com/vmware-tanzu/octant/internal/context"
	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/config"
	oevent "github.com/vmware-tanzu/octant/pkg/event"

	oerrors "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/event"
	internalLog "github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	RequestSetContentPath = "action.octant.dev/setContentPath"
)

// ContentManagerOption is an option for configuring ContentManager.
type ContentManagerOption func(manager *ContentManager)

// ContentGenerateFunc is a function that generates content. It returns `rerun=true`
// if the action should be be immediately rerun.
type ContentGenerateFunc func(ctx context.Context, state octant.State) (Content, bool, error)

// Content is a content to be sent to clients.
type Content struct {
	// Response is a content response.
	Response component.ContentResponse
	// Path is the path of the content.
	Path string
}

// Checksum return's the content's checksum.
func (c *Content) Checksum() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return fmt.Sprintf("%s", sha256.Sum256(b))
}

var (
	emptyContent = Content{
		Response: component.EmptyContentResponse,
	}
)

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
	ctx                 context.Context
	dashConfig          config.Dash
	moduleManager       module.ManagerInterface
	logger              log.Logger
	contentGenerateFunc ContentGenerateFunc
	poller              Poller
	updateContentCh     chan struct{}
}

// NewContentManager creates an instance of ContentManager.
func NewContentManager(moduleManager module.ManagerInterface, dashConfig config.Dash, logger log.Logger, options ...ContentManagerOption) *ContentManager {
	cm := &ContentManager{
		moduleManager:   moduleManager,
		dashConfig:      dashConfig,
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
func (cm *ContentManager) Start(ctx context.Context, state octant.State, s api.OctantClient) {
	cm.ctx = ctx
	logger := internalLog.From(ctx)
	logger.Debugf("starting content manager")

	defer func() {
		logger.Debugf("stopping content manager")
		close(cm.updateContentCh)
	}()

	ctx, cancel := context.WithCancel(ctx)

	updateCancel := state.OnContentPathUpdate(func(contentPath string) {
		cm.updateContentCh <- struct{}{}
	})

	go func() {
		<-s.StopCh()
		updateCancel()
		cancel()
	}()

	cm.poller.Run(ctx, cm.updateContentCh, cm.runUpdate(state, s), event.DefaultScheduleDelay)
}

func (cm *ContentManager) runUpdate(state octant.State, s api.OctantClient) PollerFunc {
	previousChecksum := ""

	return func(ctx context.Context) bool {
		contentPath := state.GetContentPath()
		if contentPath == "" {
			return false
		}

		content, _, err := cm.contentGenerateFunc(ctx, state)
		if err != nil {
			var ae *oerrors.AccessError
			if errors.As(err, &ae) {
				if ae.Name() == oerrors.OctantAccessError {
					return false
				}
			}
			if ctx.Err() == context.Canceled {
				return false
			}
			cm.logger.
				WithErr(err).
				With("content-path", contentPath).
				Errorf("generate content")
			return false
		}

		checksum := content.Checksum()
		if checksum == previousChecksum {
			return false
		}
		previousChecksum = checksum

		if ctx.Err() == nil {
			if content.Path == state.GetContentPath() {
				s.Send(CreateContentEvent(content.Response, state.GetNamespace(), contentPath, state.GetQueryParams()))
			}

		}

		return false
	}
}

func (cm *ContentManager) generateContent(ctx context.Context, state octant.State) (Content, bool, error) {
	contentPath := state.GetContentPath()
	logger := cm.logger.With("contentPath", contentPath)

	now := time.Now()
	defer func() {
		logger.With("elapsed", time.Since(now)).Debugf("generating content")
	}()

	m, ok := cm.moduleManager.ModuleForContentPath(contentPath)
	if !ok {
		return emptyContent, false, fmt.Errorf("unable to find module for content path %q", contentPath)
	}
	modulePath := strings.TrimPrefix(contentPath, m.Name())
	options := module.ContentOptions{
		LabelSet: FiltersToLabelSet(state.GetFilters()),
	}

	pluginState := ocontext.ClientState{
		ClientID:    state.GetClientID(),
		Filters:     state.GetFilters(),
		Namespace:   state.GetNamespace(),
		ContextName: cm.dashConfig.CurrentContext(),
	}
	ctx = ocontext.WithClientState(ctx, pluginState)

	contentResponse, err := m.Content(ctx, modulePath, options)
	if err != nil {
		if nfe, ok := err.(notFound); ok && nfe.NotFound() {
			logger.Debugf("path not found")
			contentResponse = notFoundPage(contentPath)
		} else {
			return emptyContent, false, fmt.Errorf("generate content: %w", err)
		}
	}

	title := GenerateBreadcrumb(cm, contentPath, state, m, options)
	if len(title) > 0 {
		contentResponse.Title = title
	}

	content := Content{
		Response: contentResponse,
		Path:     contentPath,
	}
	return content, false, nil
}

// Handlers returns a slice of client request handlers.
func (cm *ContentManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: RequestSetContentPath,
			Handler:     cm.SetContentPath,
		},
		{
			RequestType: action.RequestSetNamespace,
			Handler:     cm.SetNamespace,
		},
		{
			RequestType: CheckLoading,
			Handler:     cm.Loaded,
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
	return state.Dispatch(cm.ctx, action.RequestSetNamespace, payload)
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

// Loaded is no-op once content is serving
func (cm *ContentManager) Loaded(state octant.State, payload action.Payload) error {
	return nil
}

type notFound interface {
	NotFound() bool
	Path() string
}

// CreateContentEvent creates a content event.
func CreateContentEvent(contentResponse component.ContentResponse, namespace, contentPath string, queryParams map[string][]string) oevent.Event {
	return oevent.Event{
		Type: oevent.EventTypeContent,
		Data: map[string]interface{}{
			"content":     contentResponse,
			"namespace":   namespace,
			"contentPath": contentPath,
			"queryParams": queryParams,
		},
	}
}

func notFoundPage(contentPath string) component.ContentResponse {
	title := component.TitleFromString("Not Found")
	cr := component.NewContentResponse(title)

	// TODO change periodically with something fun
	wiki := moduloIndex(contentPath, [][]string{
		{"https://en.wikipedia.org/wiki/Octant_(plane_geometry)", "Octant (plane geometry) via Wikipedia"},
		{"https://en.wikipedia.org/wiki/Octant_(solid_geometry)", "Octant (solid geometry) via Wikipedia"},
		{"https://en.wikipedia.org/wiki/Octant_(circle)", "Octant (circle) via Wikipedia"},
		{"https://en.wikipedia.org/wiki/Octant_(instrument)", "Octant (instrument) via Wikipedia"},
		{"https://en.wikipedia.org/wiki/Octans", "Octant (constellation) via Wikipedia"},
		{"https://en.wikipedia.org/wiki/Octant_(band)", "Octant (band) via Wikipedia"},
	})

	text := component.NewMarkdownText(fmt.Sprintf(`
The requested page was not found. The resource may have been deleted.

You can:
  * Wait for the resource to be created
  * Use the navigation links to go to a new page
  * Go [up one level](#/%s)
  * Use the [back button](javascript:window.history.back()) to return to the previous page
  * Learn more about [%s](%s)
`, notFoundRedirectPath(contentPath), wiki[1], wiki[0]))
	text.EnableTrustedContent()

	body := &component.List{}
	body.Title = title
	body.Add(text)
	cr.Add(body)

	return *cr
}

func moduloIndex(key string, options [][]string) []string {
	h := fnv.New32a()
	h.Write([]byte(key))
	i := int(h.Sum32()) % len(options)
	return options[i]
}

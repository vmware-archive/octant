/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/event"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/octant"
)

// HelperStateManagerOption is an option for configuration HelperManager
type HelperStateManagerOption func(manager *HelperStateManager)

// HelperGenerateFunc is a function which generates a helper event
type HelperGenerateFunc func(ctx context.Context, state octant.State) (octant.Event, error)

// WithHelperGenerator sets the helper generator
func WithHelperGenerator(fn HelperGenerateFunc) HelperStateManagerOption {
	return func(manager *HelperStateManager) {
		manager.helperGenerateFunc = fn
	}
}

// WithHelperGeneratorPoll generates the poller
func WithHelperGeneratorPoll(poller Poller) HelperStateManagerOption {
	return func(manager *HelperStateManager) {
		manager.poller = poller
	}
}

// HelperStateManager manages buildInfo
type HelperStateManager struct {
	dashConfig         config.Dash
	helperGenerateFunc HelperGenerateFunc
	poller             Poller
}

var _StateManager = (*HelperStateManager)(nil)

// NewHelperStateManager creates an instance of HelperStateManager
func NewHelperStateManager(dashConfig config.Dash, options ...HelperStateManagerOption) *HelperStateManager {
	hm := &HelperStateManager{
		dashConfig: dashConfig,
		poller:     NewInterruptiblePoller("buildInfo"),
	}

	hm.helperGenerateFunc = hm.generateContexts

	for _, option := range options {
		option(hm)
	}

	return hm
}

// Handlers returns a slice of handlers
func (h *HelperStateManager) Handlers() []octant.ClientRequestHandler {
	return nil
}

// Start starts the manager
func (h *HelperStateManager) Start(ctx context.Context, state octant.State, client OctantClient) {
	h.poller.Run(ctx, nil, h.runUpdate(state, client), event.DefaultScheduleDelay)
}

func (h *HelperStateManager) runUpdate(state octant.State, client OctantClient) PollerFunc {
	var previous []byte

	return func(ctx context.Context) bool {
		logger := log.From(ctx)

		ev, err := h.helperGenerateFunc(ctx, state)
		if err != nil {
			logger.WithErr(err).Errorf("generate helper buildInfo")
			return false
		}

		if ctx.Err() == nil {
			cur, err := json.Marshal(ev)
			if err != nil {
				logger.WithErr(err).Errorf("unable to marshal buildInfo")
				return false
			}

			if bytes.Compare(previous, cur) != 0 {
				previous = cur
				client.Send(ev)
			}
		}

		return false
	}
}

func (h *HelperStateManager) generateContexts(ctx context.Context, state octant.State) (octant.Event, error) {
	generator, err := h.initGenerator(state)
	if err != nil {
		return octant.Event{}, err
	}
	return generator.Event(ctx)
}

func (h *HelperStateManager) initGenerator(state octant.State) (*event.HelperGenerator, error) {
	return event.NewHelperGenerator(h.dashConfig), nil
}

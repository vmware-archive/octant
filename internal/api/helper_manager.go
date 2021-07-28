/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"

	"github.com/vmware-tanzu/octant/pkg/api"
	oevent "github.com/vmware-tanzu/octant/pkg/event"

	"github.com/vmware-tanzu/octant/internal/event"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/config"
)

// HelperStateManagerOption is an option for configuration HelperManager
type HelperStateManagerOption func(manager *HelperStateManager)

// HelperGenerateFunc is a function which generates a helper event
type HelperGenerateFunc func(ctx context.Context) ([]oevent.Event, error)

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

var _ StateManager = (*HelperStateManager)(nil)

// NewHelperStateManager creates an instance of HelperStateManager
func NewHelperStateManager(dashConfig config.Dash, options ...HelperStateManagerOption) *HelperStateManager {
	hm := &HelperStateManager{
		dashConfig: dashConfig,
		poller:     NewInterruptiblePoller("helperManager"),
	}

	hm.helperGenerateFunc = hm.generateEvents

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
func (h *HelperStateManager) Start(ctx context.Context, state octant.State, client api.OctantClient) {
	h.poller.Run(ctx, nil, h.runUpdate(state, client), event.DefaultScheduleDelay)
}

func (h *HelperStateManager) runUpdate(state octant.State, client api.OctantClient) PollerFunc {
	var buildInfoGenerated, kubeConfigPathGenerated bool

	return func(ctx context.Context) bool {
		logger := log.From(ctx)

		events, err := h.helperGenerateFunc(ctx)
		if err != nil {
			logger.WithErr(err).Errorf("helper manager runUpdate")
			return false
		}

		buildInfoEvent, _ := oevent.FindEvent(events, oevent.EventTypeBuildInfo)
		kubeConfigPathEvent, _ := oevent.FindEvent(events, oevent.EventTypeKubeConfigPath)

		if ctx.Err() == nil {

			if !buildInfoGenerated {
				buildInfoGenerated = !buildInfoGenerated
				client.Send(buildInfoEvent)
			}

			if !kubeConfigPathGenerated {
				kubeConfigPathGenerated = !kubeConfigPathGenerated
				client.Send(kubeConfigPathEvent)
			}
		}

		return false
	}
}

func (h *HelperStateManager) generateEvents(ctx context.Context) ([]oevent.Event, error) {
	generator := event.NewHelperGenerator(h.dashConfig)
	return generator.Events(ctx)
}

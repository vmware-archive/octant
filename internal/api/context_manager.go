/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"bytes"
	"context"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/vmware-tanzu/octant/pkg/api"
	oevent "github.com/vmware-tanzu/octant/pkg/event"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/event"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/config"
)

// ContextManagerOption is an option for configuring ContextManager.
type ContextManagerOption func(manager *ContextManager)

// ContextGenerateFunc is a function which generates a context event.
type ContextGenerateFunc func(ctx context.Context, state octant.State) ([]oevent.Event, error)

// WithContextGenerator sets the context generator.
func WithContextGenerator(fn ContextGenerateFunc) ContextManagerOption {
	return func(manager *ContextManager) {
		manager.contextGenerateFunc = fn
	}
}

// WithContextGeneratorPoll generates the poller.
func WithContextGeneratorPoll(poller Poller) ContextManagerOption {
	return func(manager *ContextManager) {
		manager.poller = poller
	}
}

// ContextManager manages context.
type ContextManager struct {
	ctx                 context.Context
	dashConfig          config.Dash
	contextGenerateFunc ContextGenerateFunc
	poller              Poller
}

var _ StateManager = (*ContextManager)(nil)

// NewContextManager creates an instances of ContextManager.
func NewContextManager(dashConfig config.Dash, options ...ContextManagerOption) *ContextManager {
	cm := &ContextManager{
		dashConfig: dashConfig,
		poller:     NewInterruptiblePoller("context"),
	}

	cm.contextGenerateFunc = cm.generateContexts

	for _, option := range options {
		option(cm)
	}

	return cm
}

// Handlers returns a slice of handlers.
func (c *ContextManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: action.RequestSetContext,
			Handler:     c.SetContext,
		},
	}
}

// SetContext sets the current context.
func (c *ContextManager) SetContext(state octant.State, payload action.Payload) error {
	requestedContext, err := payload.String("requestedContext")
	if err != nil {
		return errors.Wrap(err, "extract requested context from payload")
	}
	state.SetContext(requestedContext)
	state.Dispatch(c.ctx, action.RequestSetContext, action.Payload{"contextName": requestedContext})
	return nil
}

// Start starts the manager.
func (c *ContextManager) Start(ctx context.Context, state octant.State, s api.OctantClient) {
	c.poller.Run(ctx, nil, c.runUpdate(state, s), event.DefaultScheduleDelay)
	c.ctx = ctx
}

func (c *ContextManager) runUpdate(state octant.State, s api.OctantClient) PollerFunc {
	var previous []byte

	logger := c.dashConfig.Logger()
	return func(ctx context.Context) bool {
		evs, err := c.contextGenerateFunc(ctx, state)
		if err != nil {
			logger.WithErr(err).Errorf("generate contexts")
		}

		if ctx.Err() == nil {
			for _, ev := range evs {
				cur, err := json.Marshal(ev)
				if err != nil {
					logger.WithErr(err).Errorf("unable to marshal context")
					return false
				}

				if bytes.Compare(previous, cur) != 0 {
					previous = cur
					s.Send(ev)
				}
			}
		}

		return false
	}
}

func (c *ContextManager) generateContexts(ctx context.Context, state octant.State) ([]oevent.Event, error) {
	generator, err := c.initGenerator(state)
	if err != nil {
		return []oevent.Event{}, err
	}
	return generator.Events(ctx)
}

func (c *ContextManager) initGenerator(state octant.State) (*event.ContextsGenerator, error) {
	return event.NewContextsGenerator(c.dashConfig), nil
}

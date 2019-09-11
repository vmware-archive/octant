/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"sync"
	"time"

	"github.com/vmware/octant/internal/event"
	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/octant"
)

// NavigationManager manages the navigation tree.
type NavigationManager struct {
	moduleManager module.ManagerInterface
}

var _ StateManager = (*NavigationManager)(nil)

// NewNavigationManager creates an instance of NavigationManager.
func NewNavigationManager(moduleManager module.ManagerInterface) *NavigationManager {
	return &NavigationManager{
		moduleManager: moduleManager,
	}
}

// Handlers returns nil.
func (n NavigationManager) Handlers() []octant.ClientRequestHandler {
	return nil
}

// Start starts the manager. It periodically generates navigation updates.
func (n *NavigationManager) Start(ctx context.Context, state octant.State, s OctantClient) {
	mu := sync.Mutex{}
	updateNamespaceCh := make(chan struct{}, 1)
	updateCancel := state.OnNamespaceUpdate(func(_ string) {
		mu.Lock()
		defer mu.Unlock()
		updateNamespaceCh <- struct{}{}
	})
	defer updateCancel()

	var generateCtx context.Context
	var cancel context.CancelFunc

	timer := time.NewTimer(0)
	done := false
	for !done {
		select {
		case <-ctx.Done():
			done = true
			break
		case <-updateNamespaceCh:
			if cancel != nil {
				cancel()
				cancel = nil
			}
			timer.Reset(0)
		case <-timer.C:
			go func() {
				generateCtx, cancel = context.WithCancel(ctx)
				defer func() {
					if cancel != nil {
						cancel()
						cancel = nil
					}
				}()

				generator := n.initGenerator(state)
				ev, err := generator.Event(generateCtx)
				if err != nil {
					// do something with this error?
					return
				}

				if ctx.Err() == nil {
					s.Send(ev)
					timer.Reset(generator.ScheduleDelay())
				}
			}()
		}
	}

	timer.Stop()
}

func (n *NavigationManager) initGenerator(state octant.State) *event.NavigationGenerator {
	return &event.NavigationGenerator{
		Modules:     n.moduleManager.Modules(),
		Namespace:   state.GetNamespace(),
		DefaultPath: state.GetContentPath(),
	}
}

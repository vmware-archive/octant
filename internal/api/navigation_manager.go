/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	oevent "github.com/vmware-tanzu/octant/pkg/event"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/vmware-tanzu/octant/internal/event"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/navigation"
)

// NavigationManagerConfig is configuration of NavigationManager.
type NavigationManagerConfig interface {
	ModuleManager() module.ManagerInterface
}

// NavigationManagerConfig is an option for configuration NavigationManager.
type NavigationManagerOption func(n *NavigationManager)

// NavigationGeneratorFunc is a function that generates a navigation tree.
type NavigationGeneratorFunc func(ctx context.Context, state octant.State, config NavigationManagerConfig) ([]navigation.Navigation, error)

// WithNavigationGenerator configures the navigation generator function.
func WithNavigationGenerator(fn NavigationGeneratorFunc) NavigationManagerOption {
	return func(n *NavigationManager) {
		n.navigationGeneratorFunc = fn
	}
}

// WithNavigationGeneratorPoller configures the poller.
func WithNavigationGeneratorPoller(poller Poller) NavigationManagerOption {
	return func(n *NavigationManager) {
		n.poller = poller
	}
}

// NavigationManager manages the navigation tree.
type NavigationManager struct {
	config                  NavigationManagerConfig
	navigationGeneratorFunc NavigationGeneratorFunc
	poller                  Poller
}

var _ StateManager = (*NavigationManager)(nil)

// NewNavigationManager creates an instance of NavigationManager.
func NewNavigationManager(config NavigationManagerConfig, options ...NavigationManagerOption) *NavigationManager {
	n := &NavigationManager{
		config:                  config,
		poller:                  NewInterruptiblePoller("navigation"),
		navigationGeneratorFunc: NavigationGenerator,
	}

	for _, option := range options {
		option(n)
	}

	return n
}

// Handlers returns nil.
func (n NavigationManager) Handlers() []octant.ClientRequestHandler {
	return nil
}

// Start starts the manager. It periodically generates navigation updates.
func (n *NavigationManager) Start(ctx context.Context, state octant.State, s OctantClient) {
	ch := make(chan struct{}, 1)
	defer func() {
		close(ch)
	}()

	n.poller.Run(ctx, ch, n.runUpdate(state, s), event.DefaultScheduleDelay)
}

func (n *NavigationManager) runUpdate(state octant.State, client OctantClient) PollerFunc {
	var previous []byte

	return func(ctx context.Context) bool {
		logger := log.From(ctx)

		entries, err := n.navigationGeneratorFunc(ctx, state, n.config)
		if err != nil {
			logger.WithErr(err).Errorf("load namespaces")
			return false
		}

		if ctx.Err() == nil {
			cur, err := json.Marshal(entries)
			if err != nil {
				logger.WithErr(err).Errorf("unable to marshal navigation entries")
				return false
			}

			if bytes.Compare(previous, cur) != 0 {
				previous = cur
				client.Send(CreateNavigationEvent(entries, state.GetContentPath()))
			}

		}

		return false
	}
}

// NavigationGenerator generates a navigation tree given a set of modules and a namespace.
func NavigationGenerator(ctx context.Context, state octant.State, config NavigationManagerConfig) ([]navigation.Navigation, error) {
	if state == nil {
		return nil, errors.New("state is nil")
	}

	if config == nil {
		return nil, errors.New("navigation config is nil")
	}

	modules := config.ModuleManager().Modules()
	namespace := state.GetNamespace()

	var sections []navigation.Navigation

	lookup := make(map[string][]navigation.Navigation)
	var mu sync.Mutex

	var g errgroup.Group

	for i := range modules {
		m := modules[i]
		g.Go(func() error {
			contentPath := m.ContentPath()
			navList, err := m.Navigation(ctx, namespace, contentPath)
			if err != nil {
				return fmt.Errorf("unable to generate navigation for module %s: %v", m.Name(), err)
			}

			mu.Lock()
			defer mu.Unlock()
			lookup[m.Name()] = navList
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for _, m := range modules {
		sections = append(sections, lookup[m.Name()]...)
	}

	return sections, nil
}

// CreateNavigationEvent creates a navigation event.
func CreateNavigationEvent(sections []navigation.Navigation, defaultPath string) oevent.Event {
	return oevent.Event{
		Type: oevent.EventTypeNavigation,
		Data: map[string]interface{}{
			"sections":    sections,
			"defaultPath": defaultPath,
		},
	}
}

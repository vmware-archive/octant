/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"path"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/pkg/navigation"
)

type navigationResponse struct {
	Sections    []navigation.Navigation `json:"sections"`
	DefaultPath string                  `json:"defaultPath"`
}

// NavigationGenerator generates navigation events.
type NavigationGenerator struct {
	// Modules is a list of modules to generate navigation from.
	Modules []module.Module

	// Namespace is the current namespace
	Namespace string

	// RunEvery is how often the event generator should be run.
	RunEvery time.Duration

	// DefaultPath is Octant's default path.
	DefaultPath string
}

var _ octant.Generator = (*NavigationGenerator)(nil)

// Event generates a navigation event.
func (g *NavigationGenerator) Event(ctx context.Context) (octant.Event, error) {
	ans := newAPINavSections(g.Modules)

	ns, err := ans.Sections(ctx, g.Namespace)
	if err != nil {
		return octant.Event{}, err
	}

	nr := navigationResponse{
		Sections:    ns,
		DefaultPath: g.DefaultPath,
	}

	return octant.Event{
		Type: octant.EventTypeNavigation,
		Data: nr,
	}, nil
}

// ScheduleDelay returns how long to delay before running this generator again.
func (g *NavigationGenerator) ScheduleDelay() time.Duration {
	return DefaultScheduleDelay
}

// Name return the name of this generator.
func (NavigationGenerator) Name() string {
	return "navigation"
}

type apiNavSections struct {
	modules []module.Module
}

func newAPINavSections(modules []module.Module) *apiNavSections {
	return &apiNavSections{
		modules: modules,
	}
}

func (ans *apiNavSections) Sections(ctx context.Context, namespace string) ([]navigation.Navigation, error) {
	var sections []navigation.Navigation

	lookup := make(map[string][]navigation.Navigation)
	var mu sync.Mutex

	var g errgroup.Group

	for i := range ans.modules {
		m := ans.modules[i]
		g.Go(func() error {
			contentPath := path.Join("/content", m.ContentPath())
			navList, err := m.Navigation(ctx, namespace, contentPath)
			if err != nil {
				return err
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

	for _, m := range ans.modules {
		sections = append(sections, lookup[m.Name()]...)
	}

	return sections, nil
}

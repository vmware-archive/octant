package event

import (
	"context"
	"encoding/json"
	"path"
	"time"

	"github.com/heptio/developer-dash/internal/octant"
	"github.com/heptio/developer-dash/internal/module"
)

type navigationResponse struct {
	Sections []octant.Navigation `json:"sections,omitempty"`
}

// NavigationGenerator generates navigation events.
type NavigationGenerator struct {
	// Modules is a list of modules to generate navigation from.
	Modules []module.Module

	// Namespace is the current namespace
	Namespace string

	// RunEvery is how often the event generator should be run.
	RunEvery time.Duration
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
		Sections: ns,
	}

	data, err := json.Marshal(nr)
	if err != nil {
		return octant.Event{}, err
	}

	return octant.Event{
		Type: octant.EventTypeNavigation,
		Data: data,
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

func (ans *apiNavSections) Sections(ctx context.Context, namespace string) ([]octant.Navigation, error) {
	var sections []octant.Navigation

	for _, m := range ans.modules {
		contentPath := path.Join("/content", m.ContentPath())
		navList, err := m.Navigation(ctx, namespace, contentPath)
		if err != nil {
			return nil, err
		}

		sections = append(sections, navList...)
	}

	return sections, nil
}

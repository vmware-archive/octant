package describer

import (
	"context"

	"github.com/heptio/developer-dash/pkg/view/component"
)

type StubDescriber struct {
	path       string
	components []component.Component
}

func NewStubDescriber(p string, components ...component.Component) *StubDescriber {
	return &StubDescriber{
		path:       p,
		components: components,
	}
}
func (d *StubDescriber) Describe(context.Context, string, string, Options) (component.ContentResponse, error) {
	return component.ContentResponse{
		Components: d.components,
	}, nil
}

func (d *StubDescriber) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(d.path, d),
	}
}

func NewEmptyDescriber(p string) *StubDescriber {
	return &StubDescriber{
		path: p,
	}
}

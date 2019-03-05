package overview

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/heptio/developer-dash/internal/cluster"
)

func Test_realGenerator_Generate(t *testing.T) {
	textOther := component.NewText("other")
	textFoo := component.NewText("foo")
	textSub := component.NewText("sub")

	describers := []Describer{
		newStubDescriber("/other", textOther),
		newStubDescriber("/foo", textFoo),
		newStubDescriber("/sub/(?P<name>.*?)", textSub),
	}

	var pathFilters []pathFilter
	for _, d := range describers {
		pathFilters = append(pathFilters, d.PathFilters()...)
	}

	cases := []struct {
		name     string
		path     string
		expected component.ContentResponse
		isErr    bool
	}{
		{
			name:     "dynamic content",
			path:     "/foo",
			expected: component.ContentResponse{ViewComponents: []component.ViewComponent{textFoo}},
		},
		{
			name:  "invalid path",
			path:  "/missing",
			isErr: true,
		},
		{
			name: "sub path",
			path: "/sub/foo",
			expected: component.ContentResponse{
				ViewComponents: []component.ViewComponent{textSub},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			clusterClient := clusterfake.NewMockClientInterface(controller)
			c := cachefake.NewMockCache(controller)

			di := clusterfake.NewMockDiscoveryInterface(controller)

			ctx := context.Background()
			pm := newPathMatcher()
			for _, pf := range pathFilters {
				pm.Register(ctx, pf)
			}

			g, err := newGenerator(c, di, pm, clusterClient, nil)
			require.NoError(t, err)

			cResponse, err := g.Generate(ctx, tc.path, "/prefix", "default", GeneratorOptions{})
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, cResponse)
		})
	}
}

type stubDescriber struct {
	path       string
	components []component.ViewComponent
}

func newStubDescriber(p string, components ...component.ViewComponent) *stubDescriber {
	return &stubDescriber{
		path:       p,
		components: components,
	}
}

func newEmptyDescriber(p string) *stubDescriber {
	return &stubDescriber{
		path: p,
	}
}

func (d *stubDescriber) Describe(context.Context, string, string, cluster.ClientInterface, DescriberOptions) (component.ContentResponse, error) {
	return component.ContentResponse{
		ViewComponents: d.components,
	}, nil
}

func (d *stubDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

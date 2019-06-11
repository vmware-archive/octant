package overview

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	clusterFake "github.com/heptio/developer-dash/internal/cluster/fake"
	configFake "github.com/heptio/developer-dash/internal/config/fake"
	"github.com/heptio/developer-dash/internal/describer"
	objectStoreFake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func Test_realGenerator_Generate(t *testing.T) {
	textOther := component.NewText("other")
	textFoo := component.NewText("foo")
	textSub := component.NewText("sub")

	describers := []describer.Describer{
		describer.NewStubDescriber("/other", textOther),
		describer.NewStubDescriber("/foo", textFoo),
		describer.NewStubDescriber("/sub/(?P<name>.*?)", textSub),
	}

	var PathFilters []describer.PathFilter
	for _, d := range describers {
		PathFilters = append(PathFilters, d.PathFilters()...)
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
			expected: component.ContentResponse{Components: []component.Component{textFoo}},
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
				Components: []component.Component{textSub},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashConfig := configFake.NewMockDash(controller)

			clusterClient := clusterFake.NewMockClientInterface(controller)
			dashConfig.EXPECT().ClusterClient().Return(clusterClient).AnyTimes()

			discoveryInterface := clusterFake.NewMockDiscoveryInterface(controller)
			clusterClient.EXPECT().DiscoveryClient().Return(discoveryInterface, nil).AnyTimes()

			objectStore := objectStoreFake.NewMockObjectStore(controller)
			dashConfig.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

			ctx := context.Background()
			pathMatcher := describer.NewPathMatcher("module")
			for _, pf := range PathFilters {
				pathMatcher.Register(ctx, pf)
			}

			g, err := newGenerator(pathMatcher, dashConfig)
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

type emptyComponent struct{}

var _ component.Component = (*emptyComponent)(nil)

func (c *emptyComponent) GetMetadata() component.Metadata {
	return component.Metadata{
		Type: "empty",
	}
}

func (c *emptyComponent) SetAccessor(string) {
	// no-op
}

func (c *emptyComponent) IsEmpty() bool {
	return true
}

func (c *emptyComponent) String() string {
	return ""
}

func (c *emptyComponent) LessThan(interface{}) bool {
	return false
}

func (c emptyComponent) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})

	return json.Marshal(m)
}

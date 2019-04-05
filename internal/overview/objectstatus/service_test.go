package objectstatus

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/pkg/cacheutil"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_service(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T, *cachefake.MockCache) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, c *cachefake.MockCache) runtime.Object {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Endpoints",
					Name:       "stateful",
				}

				endpoints := testutil.LoadObjectFromFile(t, "endpoints_ok.yaml")

				c.EXPECT().Get(gomock.Any(), gomock.Eq(key)).
					Return(testutil.ToUnstructured(t, endpoints), nil)

				objectFile := "service_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusOK,
				Details:    component.TitleFromString("Service is OK"),
			},
		},
		{
			name: "no endpoint subsets",
			init: func(t *testing.T, c *cachefake.MockCache) runtime.Object {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Endpoints",
					Name:       "stateful",
				}

				endpoints := testutil.LoadObjectFromFile(t, "endpoints_no_subsets.yaml")

				c.EXPECT().Get(gomock.Any(), gomock.Eq(key)).
					Return(testutil.ToUnstructured(t, endpoints), nil)

				objectFile := "service_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details:    component.TitleFromString("Service has no endpoints"),
			},
		},
		{
			name: "object is nil",
			init: func(t *testing.T, c *cachefake.MockCache) runtime.Object {
				return nil
			},
			isErr: true,
		},
		{
			name: "object is not a daemon set",
			init: func(t *testing.T, c *cachefake.MockCache) runtime.Object {
				return &unstructured.Unstructured{}
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			c := cachefake.NewMockCache(controller)

			object := tc.init(t, c)

			ctx := context.Background()
			status, err := service(ctx, object, c)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}

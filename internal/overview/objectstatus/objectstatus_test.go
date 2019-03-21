package objectstatus

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/heptio/developer-dash/internal/cache"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_status(t *testing.T) {
	deployObjectStatus := ObjectStatus{
		nodeStatus: component.NodeStatusOK,
		Details:    component.Title(component.NewText("apps/v1 Deployment is OK")),
	}

	lookup := statusLookup{
		{apiVersion: "v1", kind: "Object"}: func(context.Context, runtime.Object, cache.Cache) (ObjectStatus, error) {
			return deployObjectStatus, nil
		},
	}

	cases := []struct {
		name     string
		object   runtime.Object
		lookup   statusLookup
		expected ObjectStatus
		isErr    bool
	}{
		{
			name:     "in general",
			object:   testutil.CreateDeployment("deployment"),
			lookup:   lookup,
			expected: deployObjectStatus,
		},
		{
			name:   "nil object",
			object: nil,
			lookup: lookup,
			isErr:  true,
		},
		{
			name:   "nil lookup",
			object: testutil.CreateDeployment("deployment"),
			lookup: nil,
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			c := cachefake.NewMockCache(controller)

			ctx := context.Background()
			got, err := status(ctx, tc.object, c, tc.lookup)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}

func Test_ObjectStatus_AddDetail(t *testing.T) {
	os := ObjectStatus{}
	os.AddDetail("detail")

	expected := component.TitleFromString("detail")
	assert.Equal(t, expected, os.Details)
}

func Test_ObjectStatus_AddDetailf(t *testing.T) {
	os := ObjectStatus{}
	os.AddDetailf("detail %d", 1)

	expected := component.TitleFromString("detail 1")
	assert.Equal(t, expected, os.Details)
}

func Test_ObjectStatus_SetError(t *testing.T) {
	os := ObjectStatus{}
	os.SetError()
	assert.Equal(t, component.NodeStatusError, os.Status())
}

func Test_ObjectStatus_SetWarning(t *testing.T) {
	os := ObjectStatus{}
	os.SetWarning()
	assert.Equal(t, component.NodeStatusWarning, os.Status())

	os.SetError()
	os.SetWarning()
	assert.Equal(t, component.NodeStatusError, os.Status())
}

func Test_ObjectStatus_Default(t *testing.T) {
	os := ObjectStatus{}

	expected := component.NodeStatusOK
	assert.Equal(t, expected, os.Status())
}

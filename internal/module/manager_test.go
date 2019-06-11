package module_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"

	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/module/fake"
)

func TestManager(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	clusterClient := clusterfake.NewMockClientInterface(controller)

	manager, err := module.NewManager(clusterClient, "default", log.NopLogger())
	require.NoError(t, err)

	modules := manager.Modules()
	require.NoError(t, err)
	require.Len(t, modules, 0)

	m := fake.NewMockModule(controller)
	m.EXPECT().Start().Return(nil)
	m.EXPECT().Stop()
	m.EXPECT().SetNamespace("other").Return(nil)

	manager.Register(m)
	require.NoError(t, manager.Load())

	modules = manager.Modules()
	require.NoError(t, err)
	require.Len(t, modules, 1)

	manager.SetNamespace("other")
	manager.Unload()
}

func TestManager_ObjectPath(t *testing.T) {
	cases := []struct {
		name       string
		apiVersion string
		kind       string
		isErr      bool
		expected   string
	}{
		{
			name:       "exists",
			apiVersion: "group/version",
			kind:       "kind",
			expected:   "/path",
		},
		{
			name:       "does not exist",
			apiVersion: "v1",
			kind:       "kind",
			isErr:      true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			clusterClient := clusterfake.NewMockClientInterface(controller)

			manager, err := module.NewManager(clusterClient, "default", log.NopLogger())
			require.NoError(t, err)

			modules := manager.Modules()
			require.NoError(t, err)
			require.Len(t, modules, 0)

			m := fake.NewMockModule(controller)
			m.EXPECT().Start().Return(nil)
			supportedGVK := []schema.GroupVersionKind{
				{Group: "group", Version: "version", Kind: "kind"},
			}
			m.EXPECT().SupportedGroupVersionKind().Return(supportedGVK)
			m.EXPECT().
				GroupVersionKindPath("namespace", "group/version",  "kind", "name").
				Return("/path", nil).
				AnyTimes()

			manager.Register(m)
			require.NoError(t, manager.Load())

			got, err := manager.ObjectPath("namespace", tc.apiVersion, tc.kind, "name")

			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

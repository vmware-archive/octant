package module_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/stretchr/testify/require"

	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/pkg/view/component"
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

	manager.Register(&stubModule{})
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
			expected:   "/foo/bar",
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

			manager.Register(&stubModule{})
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

type stubModule struct{}

var _ module.Module = (*stubModule)(nil)

func (m *stubModule) Name() string {
	return "stub-module"
}

func (m *stubModule) ContentPath() string {
	panic("not implemented")
}

func (m *stubModule) Handler(root string) http.Handler {
	panic("not implemented")
}

func (m *stubModule) Navigation(ctx context.Context, namespace, root string) ([]clustereye.Navigation, error) {
	panic("not implemented")
}

func (m *stubModule) SetNamespace(namespace string) error {
	return nil
}

func (m *stubModule) Start() error {
	return nil
}

func (m *stubModule) Stop() {
}

func (m *stubModule) Content(ctx context.Context, contentPath string, prefix string, namespace string, opts module.ContentOptions) (component.ContentResponse, error) {
	panic("not implemented")
}

func (m *stubModule) Handlers(ctx context.Context) map[string]http.Handler {
	return make(map[string]http.Handler)
}

func (m *stubModule) SupportedGroupVersionKind() []schema.GroupVersionKind {
	return []schema.GroupVersionKind{
		{Group: "group", Version: "version", Kind: "kind"},
	}
}

func (m *stubModule) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	return "/foo/bar", nil
}

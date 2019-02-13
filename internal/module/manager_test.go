package module

import (
	"context"
	"net/http"
	"testing"

	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestManager(t *testing.T) {
	scheme := runtime.NewScheme()
	objects := []runtime.Object{
		newUnstructured("apps/v1", "Deployment", "default", "deploy3"),
	}

	clusterClient, err := fake.NewClient(scheme, nil, objects)
	require.NoError(t, err)

	manager, err := NewManager(clusterClient, "default", log.NopLogger())
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

func newUnstructured(apiVersion, kind, namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"namespace": namespace,
				"name":      name,
			},
		},
	}
}

type stubModule struct{}

func (m *stubModule) Name() string {
	return "stub-module"
}

func (m *stubModule) ContentPath() string {
	panic("not implemented")
}

func (m *stubModule) Handler(root string) http.Handler {
	panic("not implemented")
}

func (m *stubModule) Navigation(namespace, root string) (*hcli.Navigation, error) {
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

func (m *stubModule) Content(ctx context.Context, contentPath string, prefix string, namespace string) (component.ContentResponse, error) {
	panic("not implemented")
}

func (m *stubModule) Handlers() map[string]http.Handler {
	return make(map[string]http.Handler)
}

package module_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/require"
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

func (m *stubModule) Navigation(ctx context.Context, namespace, root string) (*clustereye.Navigation, error) {
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

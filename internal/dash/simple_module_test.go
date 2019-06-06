package dash

import (
	"context"
	"net/http"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/pkg/view/component"
)

// dashModule is a fake that implements overview.Interface.
// TODO: replace and remove them with module.Module uses mockgen
type dashModule struct{}

var _ module.Module = (*dashModule)(nil)

// newDashModule creates an instance of dashModule.
func newDashModule() *dashModule {
	return &dashModule{}
}

// Name is the module name.
func (m *dashModule) Name() string {
	return "overview"
}

// Content generates content
func (m *dashModule) Content(ctx context.Context, contentPath, prefix, namespace string, opts module.ContentOptions) (component.ContentResponse, error) {
	return component.ContentResponse{}, nil
}

// ContentPath returns the content path for mounting this module.
func (m *dashModule) ContentPath() string {
	return "/overview"
}

// Navigation is a no-op.
func (m *dashModule) Navigation(ctx context.Context, namespace, root string) ([]clustereye.Navigation, error) {
	return []clustereye.Navigation{}, nil
}

// SetNamespace sets the namespace for this module. It is a no-op.
func (m *dashModule) SetNamespace(namespace string) error {
	return nil
}

// Start starts the module. It is a no-op.
func (m *dashModule) Start() error {
	return nil
}

// Stop stops the module. It is a no-op.
func (m *dashModule) Stop() {
}

// Handlers returns an empty set of handlers.
func (m *dashModule) Handlers(ctx context.Context) map[string]http.Handler {
	return make(map[string]http.Handler)
}

func (m *dashModule) SupportedGroupVersionKind() []schema.GroupVersionKind {
	panic("implement me")
}

func (m *dashModule) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	panic("implement me")
}

func (m *dashModule) AddCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	panic("implement me")
}

func (m *dashModule) RemoveCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	panic("implement me")
}

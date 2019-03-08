package dash

import (
	"context"
	"net/http"

	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/view/component"
)

// dashModule is a fake that implements overview.Interface.
type dashModule struct{}

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
func (m *dashModule) Navigation(ctx context.Context, namespace, root string) (*hcli.Navigation, error) {
	return nil, nil
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

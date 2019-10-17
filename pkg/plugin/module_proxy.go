package plugin

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// ModuleProxy is a proxy that satisfies Octant module requirements. It allows plugins to behave as if they
// are internal modules.
type ModuleProxy struct {
	Metadata   *Metadata
	PluginName string
	Service    ModuleService
}

var _ module.Module = (*ModuleProxy)(nil)

// NewModuleProxy creates a ModuleProxy instance.
func NewModuleProxy(pluginName string, metadata *Metadata, service ModuleService) (*ModuleProxy, error) {
	if metadata == nil {
		return nil, errors.New("metadata is nil")
	}

	return &ModuleProxy{
		PluginName: pluginName,
		Metadata:   metadata,
		Service:    service,
	}, nil
}

// Name returns the module's name. It is the same as the plugin's metadata name.
func (m *ModuleProxy) Name() string {
	return m.Metadata.Name
}

func (m *ModuleProxy) ClientRequestHandlers() []octant.ClientRequestHandler {
	return nil
}

// Content returns content from the plugin. Plugins are expected to handle paths appropriately.
func (m *ModuleProxy) Content(ctx context.Context, contentPath string, opts module.ContentOptions) (component.ContentResponse, error) {
	return m.Service.Content(ctx, contentPath)
}

func (m *ModuleProxy) ContentPath() string {
	return m.Name()
}

// Navigation returns navigation from the plugin.
func (m *ModuleProxy) Navigation(ctx context.Context, namespace, root string) ([]navigation.Navigation, error) {
	topLevel, err := m.Service.Navigation(ctx)
	if err != nil {
		return nil, err
	}

	return []navigation.Navigation{topLevel}, nil
}

// SetNamespace is a no-op
func (ModuleProxy) SetNamespace(namespace string) error {
	return nil
}

// Start is a no-op
func (ModuleProxy) Start() error {
	return nil
}

// Stop is a no-op
func (ModuleProxy) Stop() {
}

// SetContext is a no-op
func (ModuleProxy) SetContext(ctx context.Context, contextName string) error {
	return nil
}

// Generators is a no-op
func (ModuleProxy) Generators() []octant.Generator {
	return []octant.Generator{}
}

// SupportedGroupVersionKind is currently a no-op. In the future this will allow plugins
// to handle paths for GVKs.
func (ModuleProxy) SupportedGroupVersionKind() []schema.GroupVersionKind {
	return []schema.GroupVersionKind{}
}

// GroupVersionKindPath is currently a no-op. In the future this will allow plugins
// to handle paths for GVKs.
func (ModuleProxy) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	return "", nil
}

// AddCRD is a no-op
func (ModuleProxy) AddCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

// RemoveCRD is a no-op
func (ModuleProxy) RemoveCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

// ResetCRDs is a no-op
func (ModuleProxy) ResetCRDs(ctx context.Context) error {
	return nil
}

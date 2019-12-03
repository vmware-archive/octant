package pane

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/generator"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Options are options for configuring Module.
type Options struct {
	DashConfig config.Dash
}

// Module is an applications module.
type Module struct {
	Options
	pathMatcher *describer.PathMatcher
}

var _ module.Module = (*Module)(nil)

// New creates an instance of Module.
func New(ctx context.Context, options Options) *Module {
	pm := describer.NewPathMatcher("overview-ext")

	appDescriber := NewPaneDescriber(options.DashConfig)
	for _, pf := range appDescriber.PathFilters() {
		pm.Register(ctx, pf)
	}

	return &Module{
		Options:     options,
		pathMatcher: pm,
	}
}

// Name is the name of the module.
func (m Module) Name() string {
	return "overview-ext"
}

// ClientRequestHandlers are client handlers for the module.
func (m Module) ClientRequestHandlers() []octant.ClientRequestHandler {
	return nil
}

// Content generates content for a content path.
func (m *Module) Content(ctx context.Context, contentPath string, opts module.ContentOptions) (component.ContentResponse, error) {
	g, err := generator.NewGenerator(m.pathMatcher, m.DashConfig)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	return g.Generate(ctx, contentPath, generator.Options{})
}

// ContentPath returns the root content path for the module.
func (m *Module) ContentPath() string {
	return m.Name()
}

// Navigation generates navigation entries for the module.
func (m *Module) Navigation(ctx context.Context, namespace, root string) ([]navigation.Navigation, error) {
	return []navigation.Navigation{}, nil
}

// SetNamespace sets the module's namespace.
func (m Module) SetNamespace(namespace string) error {
	return nil
}

// Start does nothing.
func (m Module) Start() error {
	return nil
}

// Stop does nothing.
func (m Module) Stop() {
}

// SetContext does nothing.
func (m Module) SetContext(ctx context.Context, contextName string) error {
	return nil
}

// Generators does nothing.
func (m Module) Generators() []octant.Generator {
	return nil
}

// SupportedGroupVersionKind does nothing.
func (m Module) SupportedGroupVersionKind() []schema.GroupVersionKind {
	return nil
}

// GroupVersionKindPath does nothing.
func (m Module) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	return "", errors.Errorf("not supported")
}

// AddCRD does nothing.
func (m Module) AddCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

// RemoveCRD does nothing.
func (m Module) RemoveCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

// ResetCRDs does nothing.
func (m Module) ResetCRDs(ctx context.Context) error {
	return nil
}

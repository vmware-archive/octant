package module

import (
	"context"
	"net/http"

	"github.com/heptio/developer-dash/internal/sugarloaf"
	"github.com/heptio/developer-dash/pkg/view/component"
	"k8s.io/apimachinery/pkg/labels"
)

// ContentOptions are additional options for content generation
type ContentOptions struct {
	LabelSet *labels.Set
}

// Module is an sugarloaf plugin.
type Module interface {
	// Name is the name of the module.
	Name() string
	// Handlers are additional handlers for the module
	Handlers(ctx context.Context) map[string]http.Handler
	// Content generates content for a path.
	Content(ctx context.Context, contentPath, prefix, namespace string, opts ContentOptions) (component.ContentResponse, error)
	// ContentPath will be used to construct content paths.
	ContentPath() string
	// Navigation returns navigation entries for this module.
	Navigation(ctx context.Context, namespace, root string) (*sugarloaf.Navigation, error)
	// SetNamespace is called when the current namespace changes.
	SetNamespace(namespace string) error
	// Start starts the module.
	Start() error
	// Stop stops the module.
	Stop()
}

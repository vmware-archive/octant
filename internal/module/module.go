package module

import (
	"context"

	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/view/component"
)

// Module is an hcli plugin.
type Module interface {
	// Name is the name of the module.
	Name() string
	// Content generates content for a path.
	Content(ctx context.Context, contentPath, prefix, namespace string) (component.ContentResponse, error)
	// ContentPath will be used to construct content paths.
	ContentPath() string
	// Navigation returns navigation entries for this module.
	Navigation(namespace, root string) (*hcli.Navigation, error)
	// SetNamespace is called when the current namespace changes.
	SetNamespace(namespace string) error
	// Start starts the module.
	Start() error
	// Stop stops the module.
	Stop()
}

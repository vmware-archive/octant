package module

import (
	"net/http"

	"github.com/heptio/developer-dash/internal/hcli"
)

// Module is an hcli plugin.
type Module interface {
	// Name is the name of the module.
	Name() string
	// ContentPath will be used to construct content paths.
	ContentPath() string
	// Handler is the HTTP handler for this module.
	Handler(root string) http.Handler
	// Navigation returns navigation entries for this module.
	Navigation(root string) (*hcli.Navigation, error)
	// SetNamespace is called when the current namespace changes.
	SetNamespace(namespace string) error
	// Start starts the module.
	Start() error
	// Stop stops the module.
	Stop()
}

package module

import (
	"net/http"

	"github.com/heptio/developer-dash/internal/hcli"
)

// Module is an hcli plugin.
type Module interface {
	Name() string
	ContentPath() string
	Handler(root string) http.Handler
	Navigation(root string) (*hcli.Navigation, error)
	SetNamespace(namespace string) error
	Start() error
	Stop()
}

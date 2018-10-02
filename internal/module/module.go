package module

import (
	"net/http"

	"github.com/heptio/developer-dash/internal/hcli"
)

// Module is an hcli plugin.
type Module interface {
	ContentPath() string
	Handler(root string) http.Handler
	Navigation(root string) (*hcli.Navigation, error)
	Start() error
	Stop()
}

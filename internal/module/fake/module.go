package fake

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/log"
)

// Module is a fake module.
type Module struct {
	name   string
	logger log.Logger
}

// NewModule creates an instance of Module.
func NewModule(name string, logger log.Logger) *Module {
	return &Module{
		name:   name,
		logger: logger,
	}
}

// Name is the name of the module.
func (m *Module) Name() string {
	return m.name
}

// ContentPath is the path to the module's content.
func (m *Module) ContentPath() string {
	return fmt.Sprintf("/%s", m.name)
}

// Handler is a HTTP handler for the module.
func (m *Module) Handler(prefix string) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	s := router.PathPrefix(prefix).Subrouter()
	s.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "root")
	}))

	s.Handle("/nested", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, m.name)
	}))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Errorf("fake module path not found: %s", r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	})

	return router

}

// Navigation returns navigation entries for the module.
func (m *Module) Navigation(prefix string) (*hcli.Navigation, error) {
	nav := &hcli.Navigation{
		Path:  prefix,
		Title: m.name,
	}

	return nav, nil
}

// SetNamespace sets the current namespace.
func (m *Module) SetNamespace(namespace string) error {
	return nil
}

// Start doesn't do anything.
func (m *Module) Start() error {
	return nil
}

// Stop doesn't do anything.
func (m *Module) Stop() {
}

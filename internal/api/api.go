package api

import (
	"encoding/json"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/module"
)

// Service is an API service.
type Service interface {
	RegisterModule(module.Module) error
	Handler() *mux.Router
}

type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type notFoundResponse struct {
	Error errorResponse `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	r := &notFoundResponse{
		Error: errorResponse{
			Code:    code,
			Message: message,
		},
	}

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Printf("encoding response error: %v", err)
	}
}

// API is the API for the dashboard client
type API struct {
	nsClient cluster.NamespaceInterface
	sections []*hcli.Navigation
	prefix   string

	modules map[string]http.Handler
}

// New creates an instance of API.
func New(prefix string, nsClient cluster.NamespaceInterface) *API {
	return &API{
		prefix:   prefix,
		nsClient: nsClient,
		modules:  make(map[string]http.Handler),
	}
}

// Handler returns a HTTP handler for the service.
func (a *API) Handler() *mux.Router {
	router := mux.NewRouter()
	s := router.PathPrefix(a.prefix).Subrouter()

	namespacesService := newNamespaces(a.nsClient)
	s.Handle("/namespaces", namespacesService)

	navigationService := newNavigation(a.sections)
	s.Handle("/navigation", navigationService)

	for p, h := range a.modules {
		s.PathPrefix(p).Handler(h)
	}

	s.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("api handler not found: %s", r.URL.String())
		respondWithError(w, http.StatusNotFound, "not found")
	})

	return router
}

// RegisterModule registers a module with the API service.
func (a *API) RegisterModule(m module.Module) error {
	contentPath := path.Join("/content", m.ContentPath())
	log.Printf("Registering content path %s", contentPath)
	a.modules[contentPath] = m.Handler(path.Join(a.prefix, contentPath))

	nav, err := m.Navigation(contentPath)
	if err != nil {
		return err
	}

	a.sections = append(a.sections, nav)

	return nil
}

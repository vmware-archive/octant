package api

import (
	"encoding/json"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/go-telemetry/pkg/telemetry"
	"github.com/pkg/errors"
)

// Service is an API service.
type Service interface {
	RegisterModule(module.Module) error
	Handler() *mux.Router
}

type errorMessage struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type errorResponse struct {
	Error errorMessage `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, message string) error {
	r := &errorResponse{
		Error: errorMessage{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		return errors.Errorf("encoding response: %v", err)
	}
	return nil
}

// API is the API for the dashboard client
type API struct {
	nsClient        cluster.NamespaceInterface
	moduleManager   module.ManagerInterface
	sections        []*hcli.Navigation
	prefix          string
	logger          log.Logger
	telemetryClient telemetry.Interface

	modules map[string]http.Handler
}

func (a *API) telemetryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		msDuration := int64(time.Since(startTime) / time.Millisecond)
		if a.telemetryClient != nil {
			go a.telemetryClient.With(telemetry.Labels{"endpoint": r.URL.Path, "client.useragent": r.Header.Get("User-Agent")}).SendEvent("dash.api", telemetry.Measurements{"count": 1, "duration": msDuration})
		}
	})
}

// New creates an instance of API.
func New(prefix string, nsClient cluster.NamespaceInterface, moduleManager module.ManagerInterface, logger log.Logger, telemetryClient telemetry.Interface) *API {
	return &API{
		prefix:          prefix,
		nsClient:        nsClient,
		moduleManager:   moduleManager,
		modules:         make(map[string]http.Handler),
		logger:          logger,
		telemetryClient: telemetryClient,
	}
}

// Handler returns a HTTP handler for the service.
func (a *API) Handler() *mux.Router {
	router := mux.NewRouter()
	router.Use(a.telemetryMiddleware)
	s := router.PathPrefix(a.prefix).Subrouter()

	namespacesService := newNamespaces(a.nsClient, a.logger)
	s.Handle("/namespaces", namespacesService).Methods(http.MethodGet)

	navigationService := newNavigation(a.sections, a.logger)
	s.Handle("/navigation", navigationService).Methods(http.MethodGet)

	namespaceUpdateService := newNamespace(a.moduleManager, a.logger)
	s.HandleFunc("/namespace", namespaceUpdateService.update).Methods(http.MethodPost)
	s.HandleFunc("/namespace", namespaceUpdateService.read).Methods(http.MethodGet)

	for p, h := range a.modules {
		s.PathPrefix(p).Handler(h)
	}

	s.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Errorf("api handler not found: %s", r.URL.String())
		if err := respondWithError(w, http.StatusNotFound, "not found"); err != nil {
			a.logger.Errorf("responding: %v", err)
		}
	})

	return router
}

// RegisterModule registers a module with the API service.
func (a *API) RegisterModule(m module.Module) error {
	contentPath := path.Join("/content", m.ContentPath())
	a.logger.Debugf("registering content path %s", contentPath)
	a.modules[contentPath] = m.Handler(path.Join(a.prefix, contentPath))

	nav, err := m.Navigation(contentPath)
	if err != nil {
		return err
	}

	a.sections = append(a.sections, nav)

	return nil
}

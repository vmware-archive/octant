package api

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/mime"
	"github.com/heptio/developer-dash/internal/module"
)

func serveAsJSON(w http.ResponseWriter, v interface{}, logger log.Logger) {
	w.Header().Set("Content-Type", mime.JSONContentType)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Errorf("encoding JSON response: %v", err)
	}
}

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

func respondWithError(w http.ResponseWriter, code int, message string, logger log.Logger) {
	r := &errorResponse{
		Error: errorMessage{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", mime.JSONContentType)

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		logger.Errorf("encoding JSON response: %v", err)
	}
}

// API is the API for the dashboard client
type API struct {
	nsClient      cluster.NamespaceInterface
	infoClient    cluster.InfoInterface
	moduleManager module.ManagerInterface
	prefix        string
	logger        log.Logger

	modulePaths map[string]module.Module
	modules     []module.Module
}

func (a *API) telemetryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// New creates an instance of API.
func New(prefix string, nsClient cluster.NamespaceInterface, infoClient cluster.InfoInterface, moduleManager module.ManagerInterface, logger log.Logger) *API {
	return &API{
		prefix:        prefix,
		nsClient:      nsClient,
		infoClient:    infoClient,
		moduleManager: moduleManager,
		modulePaths:   make(map[string]module.Module),
		logger:        logger,
	}
}

// Handler returns a HTTP handler for the service.
func (a *API) Handler() *mux.Router {
	router := mux.NewRouter()
	router.Use(a.telemetryMiddleware)
	s := router.PathPrefix(a.prefix).Subrouter()

	namespacesService := newNamespaces(a.nsClient, a.logger)
	s.Handle("/namespaces", namespacesService).Methods(http.MethodGet)

	ans := newAPINavSections(a.modules)

	navigationService := newNavigation(ans, a.logger)
	s.Handle("/navigation", navigationService).Methods(http.MethodGet)

	namespaceUpdateService := newNamespace(a.moduleManager, a.logger)
	s.HandleFunc("/namespace", namespaceUpdateService.update).Methods(http.MethodPost)
	s.HandleFunc("/namespace", namespaceUpdateService.read).Methods(http.MethodGet)

	infoService := newClusterInfo(a.infoClient, a.logger)
	s.Handle("/cluster-info", infoService)

	contentService := &contentHandler{
		modulePaths: a.modulePaths,
		logger:      a.logger,
		prefix:      a.prefix,
	}
	s.Handle(`/content/{rest:[a-zA-Z0-9=\-\/]+}`, contentService).Methods(http.MethodGet)

	s.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Errorf("api handler not found: %s", r.URL.String())
		respondWithError(w, http.StatusNotFound, "not found", a.logger)
	})

	return router
}

// RegisterModule registers a module with the API service.
func (a *API) RegisterModule(m module.Module) error {
	contentPath := path.Join("/content", m.ContentPath())
	a.logger.Debugf("registering content path %s", contentPath)
	a.modulePaths[contentPath] = m

	return nil
}

type apiNavSections struct {
	modules []module.Module
}

func newAPINavSections(modules []module.Module) *apiNavSections {
	return &apiNavSections{
		modules: modules,
	}
}

func (ans *apiNavSections) Sections() ([]*hcli.Navigation, error) {
	var sections []*hcli.Navigation

	for _, m := range ans.modules {
		contentPath := path.Join("/content", m.ContentPath())
		nav, err := m.Navigation(contentPath)
		if err != nil {
			return nil, err
		}

		sections = append(sections, nav)
	}

	return sections, nil
}

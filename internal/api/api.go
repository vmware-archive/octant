package api

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/mime"
	"github.com/heptio/developer-dash/internal/module"
)

//go:generate mockgen -destination=./fake/mock_cluster_client.go -package=fake github.com/heptio/developer-dash/internal/api ClusterClient
//go:generate mockgen -destination=./fake/mock_service.go -package=fake github.com/heptio/developer-dash/internal/api Service

var (
	// acceptedHosts are the hosts this api will answer for.
	acceptedHosts = []string{
		"localhost",
		"127.0.0.1",
	}
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
	Handler(ctx context.Context) (*mux.Router, error)
}

type errorMessage struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type errorResponse struct {
	Error errorMessage `json:"error,omitempty"`
}

// RespondWithError responds with an error message.
func RespondWithError(w http.ResponseWriter, code int, message string, logger log.Logger) {
	r := &errorResponse{
		Error: errorMessage{
			Code:    code,
			Message: message,
		},
	}

	logger.With(
		"code", code,
		"message", message,
	).Infof("unable to serve")

	w.Header().Set("Content-Type", mime.JSONContentType)

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		logger.Errorf("encoding JSON response: %v", err)
	}
}

type ClusterClient interface {
	NamespaceClient() (cluster.NamespaceInterface, error)
	InfoClient() (cluster.InfoInterface, error)
}

// API is the API for the dashboard client
type API struct {
	ctx           context.Context
	clusterClient ClusterClient
	moduleManager module.ManagerInterface
	prefix        string
	logger        log.Logger

	modulePaths map[string]module.Module
	modules     []module.Module
}

var _ Service = (*API)(nil)

// New creates an instance of API.
func New(ctx context.Context, prefix string, clusterClient ClusterClient, moduleManager module.ManagerInterface, logger log.Logger) *API {
	return &API{
		ctx:           ctx,
		prefix:        prefix,
		clusterClient: clusterClient,
		moduleManager: moduleManager,
		modulePaths:   make(map[string]module.Module),
		logger:        logger,
	}
}

// Handler returns a HTTP handler for the service.
func (a *API) Handler(ctx context.Context) (*mux.Router, error) {
	router := mux.NewRouter()
	router.Use(rebindHandler(acceptedHosts))

	s := router.PathPrefix(a.prefix).Subrouter()

	nsClient, err := a.clusterClient.NamespaceClient()
	if err != nil {
		return nil, errors.Wrap(err, "retrieve namespace client")
	}

	infoClient, err := a.clusterClient.InfoClient()
	if err != nil {
		return nil, errors.Wrap(err, "retrieve cluster info client")
	}

	namespacesService := newNamespaces(nsClient, a.logger)
	s.Handle("/namespaces", namespacesService).Methods(http.MethodGet)

	ans := newAPINavSections(a.modules)

	navigationService := newNavigation(ans, a.logger)
	// Support no namespace (default) or specifying namespace in path
	s.Handle("/navigation", navigationService).Methods(http.MethodGet)
	s.Handle("/navigation/namespace/{namespace}", navigationService).Methods(http.MethodGet)

	namespaceUpdateService := newNamespace(a.moduleManager, a.logger)
	s.HandleFunc("/namespace", namespaceUpdateService.update).Methods(http.MethodPost)
	s.HandleFunc("/namespace", namespaceUpdateService.read).Methods(http.MethodGet)

	infoService := newClusterInfo(infoClient, a.logger)
	s.Handle("/cluster-info", infoService)

	// Register content routes
	contentService := &contentHandler{
		nsClient:    nsClient,
		modulePaths: a.modulePaths,
		modules:     a.modules,
		logger:      a.logger,
		prefix:      a.prefix,
	}

	if err := contentService.RegisterRoutes(ctx, s); err != nil {
		a.logger.WithErr(err).Errorf("register routers")
	}

	s.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Errorf("api handler not found: %s", r.URL.String())
		RespondWithError(w, http.StatusNotFound, "not found", a.logger)
	})

	return router, nil
}

// RegisterModule registers a module with the API service.
func (a *API) RegisterModule(m module.Module) error {
	contentPath := path.Join("/content", m.ContentPath())
	a.logger.With("contentPath", contentPath).Debugf("registering content path")
	a.modulePaths[contentPath] = m
	a.modules = append(a.modules, m)

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

func (ans *apiNavSections) Sections(ctx context.Context, namespace string) ([]clustereye.Navigation, error) {
	var sections []clustereye.Navigation

	for _, m := range ans.modules {
		contentPath := path.Join("/content", m.ContentPath())
		navList, err := m.Navigation(ctx, namespace, contentPath)
		if err != nil {
			return nil, err
		}

		sections = append(sections, navList...)
	}

	return sections, nil
}

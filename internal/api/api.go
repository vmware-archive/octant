/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/mime"
	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/pkg/navigation"
)

//go:generate mockgen -destination=./fake/mock_cluster_client.go -package=fake github.com/vmware/octant/internal/api ClusterClient
//go:generate mockgen -destination=./fake/mock_service.go -package=fake github.com/vmware/octant/internal/api Service

const (
	// ListenerAddrKey is the environment variable for the Octant listener address.
	ListenerAddrKey = "OCTANT_LISTENER_ADDR"
	AcceptedHostsKey = "OCTANT_ACCEPTED_HOSTS"
	// PathPrefix is a string for the api path prefix.
	PathPrefix          = "/api/v1"
	defaultListenerAddr = "127.0.0.1:7777"
)

func acceptedHosts() []string {
	hosts := []string{
		"localhost",
		"127.0.0.1",
	}
	if customHosts := os.Getenv(AcceptedHostsKey); customHosts != "" {
		allowedHosts := strings.Split(customHosts, ",")
		hosts = append(hosts, allowedHosts...)
	}
	
	listenerAddr := ListenerAddr()
	host, _, err := net.SplitHostPort(listenerAddr)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse OCTANT_LISTENER_ADDR: %s", listenerAddr))
	}

	hosts = append(hosts, host)
	return hosts
}

// ListenerAddr returns the default listener address if OCTANT_LISTENER_ADDR is not set.
func ListenerAddr() string {
	listenerAddr := defaultListenerAddr
	if customListenerAddr := os.Getenv(ListenerAddrKey); customListenerAddr != "" {
		listenerAddr = customListenerAddr
	}
	return listenerAddr
}

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
	ForceUpdate() error
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
	ctx              context.Context
	clusterClient    ClusterClient
	moduleManager    module.ManagerInterface
	actionDispatcher ActionDispatcher
	prefix           string
	logger           log.Logger

	modulePaths   map[string]module.Module
	modules       []module.Module
	forceUpdateCh chan bool
}

var _ Service = (*API)(nil)

// New creates an instance of API.
func New(ctx context.Context, prefix string, clusterClient ClusterClient, moduleManager module.ManagerInterface, actionDispatcher ActionDispatcher, logger log.Logger) *API {
	return &API{
		ctx:              ctx,
		prefix:           prefix,
		clusterClient:    clusterClient,
		moduleManager:    moduleManager,
		actionDispatcher: actionDispatcher,
		modulePaths:      make(map[string]module.Module),
		logger:           logger,
		forceUpdateCh:    make(chan bool, 1),
	}
}

func (a *API) ForceUpdate() error {
	a.forceUpdateCh <- true
	return nil
}

// Handler returns a HTTP handler for the service.
func (a *API) Handler(ctx context.Context) (*mux.Router, error) {
	router := mux.NewRouter()
	router.Use(rebindHandler(acceptedHosts()))

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

	navigationService := newNavigationHandler(ans, a.logger)
	// Support no namespace (default) or specifying namespace in path
	s.Handle("/navigationHandler", navigationService).Methods(http.MethodGet)
	s.Handle("/navigationHandler/namespace/{namespace}", navigationService).Methods(http.MethodGet)

	namespaceUpdateService := newNamespace(a.moduleManager, a.logger)
	s.HandleFunc("/namespace", namespaceUpdateService.update).Methods(http.MethodPost)
	s.HandleFunc("/namespace", namespaceUpdateService.read).Methods(http.MethodGet)

	infoService := newClusterInfo(infoClient, a.logger)
	s.Handle("/cluster-info", infoService)

	actionService := newAction(a.logger, a.actionDispatcher)
	s.Handle("/action", actionService)

	// Register content routes
	contentService := &contentHandler{
		nsClient:      nsClient,
		modulePaths:   a.modulePaths,
		modules:       a.modules,
		logger:        a.logger,
		prefix:        a.prefix,
		forceUpdateCh: a.forceUpdateCh,
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

func (ans *apiNavSections) Sections(ctx context.Context, namespace string) ([]navigation.Navigation, error) {
	var sections []navigation.Navigation

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

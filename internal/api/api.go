/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/gorilla/mux"

	"github.com/spf13/viper"

	"github.com/vmware-tanzu/octant/internal/mime"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/log"
)

//go:generate mockgen -destination=./fake/mock_service.go -package=fake github.com/vmware-tanzu/octant/internal/api Service

const (
	// ListenerAddrKey is the environment variable for the Octant listener address.
	ListenerAddrKey  = "listener-addr"
	AcceptedHostsKey = "accepted-hosts"
	// PathPrefix is a string for the api path prefix.
	PathPrefix          = "/api/v1"
	defaultListenerAddr = "127.0.0.1:7777"
)

func AcceptedHosts() []string {
	hosts := []string{
		"localhost",
		"127.0.0.1",
	}
	if customHosts := viper.GetString(AcceptedHostsKey); customHosts != "" {
		allowedHosts := strings.Split(customHosts, ",")
		hosts = append(hosts, allowedHosts...)
	}

	listenerAddr := getListenerAddr()
	host, _, err := net.SplitHostPort(listenerAddr)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse OCTANT_LISTENER_ADDR: %s", listenerAddr))
	}

	hosts = append(hosts, host)
	return hosts
}

// Listener returns the default listener if OCTANT_LISTENER_ADDR is not set.
func Listener() (net.Listener, error) {
	listenerAddr := getListenerAddr()
	conn, err := net.DialTimeout("tcp", listenerAddr, time.Millisecond*500)
	if err != nil {
		return net.Listen("tcp", listenerAddr)
	}
	_ = conn.Close()
	return nil, fmt.Errorf("tcp %s: dial: already in use", listenerAddr)
}

func getListenerAddr() string {
	listenerAddr := defaultListenerAddr
	if customListenerAddr := viper.GetString(ListenerAddrKey); customListenerAddr != "" {
		listenerAddr = customListenerAddr
	}
	return listenerAddr
}

// Service is an API service.
type Service interface {
	Handler(ctx context.Context) (http.Handler, error)
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

// API is the API for the dashboard client
type API struct {
	ctx              context.Context
	moduleManager    module.ManagerInterface
	actionDispatcher api.ActionDispatcher
	prefix           string
	dashConfig       config.Dash
	logger           log.Logger
	scManager        *api.StreamingConnectionManager

	modulePaths   map[string]module.Module
	modules       []module.Module
	forceUpdateCh chan bool
}

var _ Service = (*API)(nil)

// New creates an instance of API.
func New(ctx context.Context, prefix string, actionDispatcher api.ActionDispatcher, streamingConnectionManager *api.StreamingConnectionManager, dashConfig config.Dash) *API {
	logger := dashConfig.Logger().With("component", "api")
	return &API{
		ctx:              ctx,
		prefix:           prefix,
		actionDispatcher: actionDispatcher,
		modulePaths:      make(map[string]module.Module),
		dashConfig:       dashConfig,
		logger:           logger,
		forceUpdateCh:    make(chan bool, 1),
		scManager:        streamingConnectionManager,
	}
}

func (a *API) ForceUpdate() error {
	a.forceUpdateCh <- true
	return nil
}

// Handler returns a HTTP handler for the service.
func (a *API) Handler(ctx context.Context) (http.Handler, error) {
	if a.dashConfig == nil {
		return nil, fmt.Errorf("missing dashConfig")
	}
	router := mux.NewRouter()
	router.Use(rebindHandler(ctx, AcceptedHosts()))

	s := router.PathPrefix(a.prefix).Subrouter()

	s.Handle("/stream", streamService(a.scManager, a.dashConfig))

	s.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Errorf("api handler not found: %s", r.URL.String())
		RespondWithError(w, http.StatusNotFound, "not found", a.logger)
	})

	return router, nil
}

// LoadingAPI is an API for startup modules to run
type LoadingAPI struct {
	ctx              context.Context
	moduleManager    module.ManagerInterface
	actionDispatcher api.ActionDispatcher
	prefix           string
	logger           log.Logger
	scManager        *api.StreamingConnectionManager

	modulePaths   map[string]module.Module
	modules       []module.Module
	forceUpdateCh chan bool
}

var _ Service = (*LoadingAPI)(nil)

// NewLoadingAPI creates an instance of LoadingAPI
func NewLoadingAPI(ctx context.Context, prefix string, actionDispatcher api.ActionDispatcher, websocketClientManager *api.StreamingConnectionManager, logger log.Logger) *LoadingAPI {
	logger = logger.With("component", "loading api")
	return &LoadingAPI{
		ctx:              ctx,
		prefix:           prefix,
		actionDispatcher: actionDispatcher,
		modulePaths:      make(map[string]module.Module),
		logger:           logger,
		forceUpdateCh:    make(chan bool, 1),
		scManager:        websocketClientManager,
	}
}

func (l *LoadingAPI) ForceUpdate() error {
	l.forceUpdateCh <- true
	return nil
}

// Handler contains a list of handlers
func (l *LoadingAPI) Handler(ctx context.Context) (http.Handler, error) {
	router := mux.NewRouter()
	router.Use(rebindHandler(ctx, AcceptedHosts()))

	s := router.PathPrefix(l.prefix).Subrouter()

	s.Handle("/stream", loadingStreamService(l.scManager))

	return router, nil
}

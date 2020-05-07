/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/mime"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/pkg/log"
)

//go:generate mockgen -destination=./fake/mock_service.go -package=fake github.com/vmware-tanzu/octant/internal/api Service

const (
	// ListenerAddrKey is the environment variable for the Octant listener address.
	ListenerAddrKey = "listener-addr"

	// AcceptedHostsKey the list of accepted hosts to be allowed to connect over web sockets
	AcceptedHostsKey = "accepted-hosts"

	// AcceptLocalIPKey if no accepted hosts are specified this flag will add all the local IP addresses
	//starting with `192.168.1.` to the accepted hosts
	AcceptLocalIPKey = "accept-local-ip"
	// PathPrefix is a string for the api path prefix.
	PathPrefix          = "/api/v1"
	defaultListenerAddr = "127.0.0.1:7777"
)

func acceptedHosts() []string {
	hosts := []string{
		"localhost",
		"127.0.0.1",
	}
	customHosts := viper.GetString(AcceptedHostsKey)
	if customHosts != "" {
		allowedHosts := strings.Split(customHosts, ",")
		hosts = append(hosts, allowedHosts...)
	} else {
		if viper.GetBool(AcceptLocalIPKey) {
			for i := 0; i < 256; i++ {
				hosts = append(hosts, fmt.Sprintf("192.168.1.%d", i))
			}
		}
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
	if customListenerAddr := viper.GetString(ListenerAddrKey); customListenerAddr != "" {
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
	actionDispatcher ActionDispatcher
	prefix           string
	dashConfig       config.Dash
	logger           log.Logger

	modulePaths   map[string]module.Module
	modules       []module.Module
	forceUpdateCh chan bool
}

var _ Service = (*API)(nil)

// New creates an instance of API.
func New(ctx context.Context, prefix string, actionDispatcher ActionDispatcher, dashConfig config.Dash) *API {
	logger := dashConfig.Logger().With("component", "api")
	return &API{
		ctx:              ctx,
		prefix:           prefix,
		actionDispatcher: actionDispatcher,
		modulePaths:      make(map[string]module.Module),
		dashConfig:       dashConfig,
		logger:           logger,
		forceUpdateCh:    make(chan bool, 1),
	}
}

func (a *API) ForceUpdate() error {
	a.forceUpdateCh <- true
	return nil
}

// Handler returns a HTTP handler for the service.
func (a *API) Handler(ctx context.Context) (http.Handler, error) {
	router := mux.NewRouter()
	router.Use(rebindHandler(ctx, acceptedHosts()))

	s := router.PathPrefix(a.prefix).Subrouter()

	manager := NewWebsocketClientManager(ctx, a.actionDispatcher)
	go manager.Run(ctx)

	s.Handle("/stream", websocketService(manager, a.dashConfig))

	s.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Errorf("api handler not found: %s", r.URL.String())
		RespondWithError(w, http.StatusNotFound, "not found", a.logger)
	})

	return router, nil
}

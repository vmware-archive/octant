/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package websockets

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"

	internalAPI "github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/util/path_util"

	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/google/uuid"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/config"
)

var (
	reContentPathNamespace = regexp.MustCompile(`^/namespace/(?P<namespace>[^/]+)/?`)
)

func defaultStateManagers(clientID string, dashConfig config.Dash) []api.StateManager {
	logger := dashConfig.Logger().With("client-id", clientID)

	return []api.StateManager{
		internalAPI.NewContentManager(dashConfig.ModuleManager(), dashConfig, logger),
		internalAPI.NewHelperStateManager(dashConfig),
		internalAPI.NewFilterManager(),
		internalAPI.NewNavigationManager(dashConfig),
		internalAPI.NewNamespacesManager(dashConfig),
		internalAPI.NewContextManager(dashConfig),
		internalAPI.NewActionRequestManager(dashConfig),
		internalAPI.NewTerminalStateManager(dashConfig),
		internalAPI.NewPodLogsStateManager(dashConfig),
	}
}

type atomicString struct {
	mu sync.RWMutex
	s  string
}

func newStringValue(initial string) *atomicString {
	return &atomicString{
		mu: sync.RWMutex{},
		s:  initial,
	}
}

func (s *atomicString) get() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.s
}

func (s *atomicString) set(v string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s = v
}

// WebsocketStateOption is an option for configuring WebsocketState.
type WebsocketStateOption func(w *WebsocketState)

// WebsocketStateManagers configures WebsocketState's state managers.
func WebsocketStateManagers(managers []api.StateManager) WebsocketStateOption {
	return func(w *WebsocketState) {
		w.managers = managers
	}
}

// WebsocketState manages state for a websocket client.
type WebsocketState struct {
	dashConfig         config.Dash
	wsClient           api.OctantClient
	contentPath        *atomicString
	namespace          *atomicString
	filters            []octant.Filter
	contentPathUpdates map[string]octant.ContentPathUpdateFunc
	namespaceUpdates   map[string]octant.NamespaceUpdateFunc

	mu               sync.RWMutex
	managers         []api.StateManager
	actionDispatcher api.ActionDispatcher

	startCtx           context.Context
	managersCancelFunc context.CancelFunc
}

var _ octant.State = (*WebsocketState)(nil)

// NewWebsocketState creates an instance of WebsocketState.
func NewWebsocketState(dashConfig config.Dash, actionDispatcher api.ActionDispatcher, wsClient api.OctantClient, options ...WebsocketStateOption) *WebsocketState {
	defaultNamespace := dashConfig.DefaultNamespace()

	w := &WebsocketState{
		dashConfig:         dashConfig,
		wsClient:           wsClient,
		contentPathUpdates: make(map[string]octant.ContentPathUpdateFunc),
		namespaceUpdates:   make(map[string]octant.NamespaceUpdateFunc),
		namespace:          newStringValue(defaultNamespace),
		contentPath:        newStringValue(""),
		filters:            make([]octant.Filter, 0),
		actionDispatcher:   actionDispatcher,
	}

	for _, option := range options {
		option(w)
	}

	if len(w.managers) < 1 {
		w.managers = defaultStateManagers(wsClient.ID(), dashConfig)
	}

	return w
}

func NewTemporaryWebsocketState(actionDispatcher api.ActionDispatcher, wsClient api.OctantClient, options ...WebsocketStateOption) *WebsocketState {
	w := &WebsocketState{
		wsClient:           wsClient,
		contentPathUpdates: make(map[string]octant.ContentPathUpdateFunc),
		namespaceUpdates:   make(map[string]octant.NamespaceUpdateFunc),
		actionDispatcher:   actionDispatcher,
	}

	for _, option := range options {
		option(w)
	}

	if len(w.managers) < 1 {
		w.managers = []api.StateManager{
			internalAPI.NewLoadingManager(),
		}
	}

	return w
}

// Start starts WebsocketState by starting all associated StateManagers.
func (c *WebsocketState) Start(ctx context.Context) {
	for i := range c.managers {
		go c.managers[i].Start(ctx, c, c.wsClient)
	}
}

// Handlers returns all the handlers for WebsocketState.
func (c *WebsocketState) Handlers() []octant.ClientRequestHandler {
	var handlers []octant.ClientRequestHandler

	for _, manager := range c.managers {
		handlers = append(handlers, manager.Handlers()...)
	}

	return handlers
}

// Dispatch dispatches a message.
func (c *WebsocketState) Dispatch(ctx context.Context, actionName string, payload action.Payload) error {
	return c.actionDispatcher.Dispatch(ctx, c, actionName, payload)
}

// SetContentPath sets the content path.
func (c *WebsocketState) SetContentPath(contentPath string) {
	if contentPath == "" {
		contentPath = path_util.NamespacedPath("overview", c.namespace.get())
		contentPath = path.Join("overview", "namespace", c.namespace.get())
	} else if c.contentPath.get() == contentPath {
		return
	}

	c.dashConfig.Logger().With(
		"contentPath", contentPath).
		Debugf("setting content path")

	c.contentPath.set(contentPath)

	m, ok := c.dashConfig.ModuleManager().ModuleForContentPath(contentPath)
	if !ok {
		c.dashConfig.Logger().
			With("contentPath", contentPath).
			Warnf("unable to find module for content path")
	} else {
		modulePath := strings.TrimPrefix(contentPath, m.Name())
		match := reContentPathNamespace.FindStringSubmatch(modulePath)
		result := make(map[string]string)
		if len(match) > 0 {
			for i, name := range reContentPathNamespace.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}

			if result["namespace"] != "" {
				c.SetNamespace(result["namespace"])
			}
		}
	}

	for _, fn := range c.contentPathUpdates {
		fn(contentPath)
	}

}

// GetContentPath returns the content path.
func (c *WebsocketState) GetContentPath() string {
	return c.contentPath.get()
}

// OnContentPathUpdate registers a function that will be called when the content path changes.
func (c *WebsocketState) OnContentPathUpdate(fn octant.ContentPathUpdateFunc) octant.UpdateCancelFunc {
	c.mu.Lock()
	defer c.mu.Unlock()

	id, _ := uuid.NewUUID()
	c.contentPathUpdates[id.String()] = fn

	cancelFunc := func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.contentPathUpdates, id.String())
	}

	return cancelFunc
}

// SetNamespace sets the namespace.
func (c *WebsocketState) SetNamespace(namespace string) {
	cur := c.namespace.get()
	if namespace == cur {
		return
	}

	c.dashConfig.Logger().
		With("namespace", namespace).
		Debugf("setting namespace")
	c.namespace.set(namespace)

	newPath := updateContentPathNamespace(c.contentPath.get(), namespace)
	if newPath != c.contentPath.get() {
		c.SetContentPath(newPath)
	}

	for _, fn := range c.namespaceUpdates {
		fn(namespace)
	}
}

// GetNamespace gets the namespace.
func (c *WebsocketState) GetNamespace() string {
	return c.namespace.get()
}

// OnNamespaceUpdate registers a function that will be run when the namespace changes.
func (c *WebsocketState) OnNamespaceUpdate(fn octant.NamespaceUpdateFunc) octant.UpdateCancelFunc {
	c.mu.Lock()
	defer c.mu.Unlock()

	id, _ := uuid.NewUUID()
	c.namespaceUpdates[id.String()] = fn

	cancelFunc := func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.namespaceUpdates, id.String())
	}

	return cancelFunc
}

// AddFilter adds a content filter.
func (c *WebsocketState) AddFilter(filter octant.Filter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.filters {
		if c.filters[i].IsEqual(filter) {
			return
		}
	}

	c.filters = append(c.filters, filter)
}

// RemoveFilter removes a content filter.
func (c *WebsocketState) RemoveFilter(filter octant.Filter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var newFilters []octant.Filter

	for i := range c.filters {
		if c.filters[i].IsEqual(filter) {
			continue
		}
		newFilters = append(newFilters, c.filters[i])
	}

	c.filters = newFilters
}

// GetFilters returns all filters.
func (c *WebsocketState) GetFilters() []octant.Filter {
	filters := make([]octant.Filter, len(c.filters))
	copy(filters, c.filters)

	sort.Slice(filters, func(i, j int) bool {
		return filters[i].Key < filters[j].Key
	})

	return filters
}

func (c *WebsocketState) SetFilters(filters []octant.Filter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.filters = filters
}

// SetContext sets the Kubernetes context.
func (c *WebsocketState) SetContext(requestedContext string) {
	c.dashConfig.SetContextChosenInUI(true)

	if err := c.dashConfig.UseContext(context.TODO(), requestedContext); err != nil {
		c.dashConfig.Logger().WithErr(err).Errorf("update context")
	}

	c.SetNamespace(c.dashConfig.DefaultNamespace())

	for _, fn := range c.contentPathUpdates {
		fn(c.GetContentPath())
	}

	c.wsClient.Send(CreateAlertUpdate(action.CreateAlert(
		action.AlertTypeInfo,
		fmt.Sprintf("Changing context to %s", requestedContext),
		action.DefaultAlertExpiration,
	)))
}

func (c *WebsocketState) GetQueryParams() map[string][]string {
	filters := c.filters

	c.wsClient.Send(CreateFiltersUpdate(filters))

	queryParams := map[string][]string{}

	var filterList []string
	for _, filter := range filters {
		filterList = append(filterList, filter.ToQueryParam())
	}
	if len(filterList) > 0 {
		queryParams["filters"] = filterList
	}

	return queryParams
}

// SendAlert sends an alert to the websocket client.
func (c *WebsocketState) SendAlert(alert action.Alert) {
	c.wsClient.Send(CreateAlertUpdate(alert))
}

func (c *WebsocketState) GetClientID() string {
	if c.wsClient == nil {
		return ""
	}
	return c.wsClient.ID()
}

func updateContentPathNamespace(in, namespace string) string {
	parts := strings.Split(in, "/")
	if in == "" {
		return ""
	}

	if len(parts) > 1 && parts[1] == "namespace" {
		parts[2] = namespace
		return path.Join(parts...)
	}
	return in
}

// CreateFiltersUpdate creates a filters update event.
func CreateFiltersUpdate(filters []octant.Filter) event.Event {
	if filters == nil {
		filters = make([]octant.Filter, 0)
	}
	return event.CreateEvent(event.EventTypeFilters, action.Payload{
		"filters": filters,
	})
}

// CreateAlertUpdate creates an alert update event.
func CreateAlertUpdate(alert action.Alert) event.Event {
	return event.CreateEvent(event.EventTypeAlert, action.Payload{
		"type":       alert.Type,
		"message":    alert.Message,
		"expiration": alert.Expiration,
	})
}

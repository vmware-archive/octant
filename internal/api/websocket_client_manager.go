/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/vmware/octant/internal/config"
)

//go:generate mockgen -destination=./fake/mock_client_manager.go -package=fake github.com/vmware/octant/internal/api ClientManager

// ClientManager is an interface for managing clients.
type ClientManager interface {
	Run(ctx context.Context)
	ClientFromRequest(dashConfig config.Dash, w http.ResponseWriter, r *http.Request) (*WebsocketClient, error)
}

type clientMeta struct {
	cancelFunc context.CancelFunc
	client     *WebsocketClient
}

// WebsocketClientManager is a client manager for websockets.
type WebsocketClientManager struct {
	// clients is the currently registered clients.
	clients map[*WebsocketClient]context.CancelFunc

	// Register registers requests from clients.
	register chan *clientMeta

	// unregister unregisters request from clients.
	unregister       chan *WebsocketClient
	ctx              context.Context
	actionDispatcher ActionDispatcher
}

var _ ClientManager = (*WebsocketClientManager)(nil)

// NewWebsocketClientManager creates an instance of WebsocketClientManager.
func NewWebsocketClientManager(ctx context.Context, dispatcher ActionDispatcher) *WebsocketClientManager {
	return &WebsocketClientManager{
		ctx:              ctx,
		clients:          make(map[*WebsocketClient]context.CancelFunc),
		register:         make(chan *clientMeta),
		unregister:       make(chan *WebsocketClient),
		actionDispatcher: dispatcher,
	}
}

// Run runs the manager. It manages multiple websocket clients.
func (m *WebsocketClientManager) Run(ctx context.Context) {
	done := false
	for !done {
		select {
		case <-ctx.Done():
			done = true
		case meta := <-m.register:
			m.clients[meta.client] = meta.cancelFunc
		case client := <-m.unregister:
			if cancelFunc, ok := m.clients[client]; ok {
				cancelFunc()
				delete(m.clients, client)
			}
		}
	}
}

// ClientFromRequest creates a websocket client from a http request.
func (m *WebsocketClientManager) ClientFromRequest(dashConfig config.Dash, w http.ResponseWriter, r *http.Request) (*WebsocketClient, error) {
	clientID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(m.ctx)
	client := NewWebsocketClient(ctx, conn, dashConfig, m.actionDispatcher, clientID)
	m.register <- &clientMeta{
		cancelFunc: func() {
			cancel()
			m.unregister <- client
		},
		client: client,
	}

	return client, nil
}

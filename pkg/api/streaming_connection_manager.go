/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package api

import (
	"context"
	"net/http"

	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/event"
)

var _ ClientManager = (*StreamingConnectionManager)(nil)

// StreamingConnectionManager is a client manager for streams.
type StreamingConnectionManager struct {
	clientFactory StreamingClientFactory

	// clients is the currently registered clients.
	clients map[StreamingClient]context.CancelFunc

	// Register registers requests from clients.
	register chan *clientMeta

	// unregister unregisters request from clients.
	unregister chan StreamingClient

	// list populates a client list
	requestList chan bool
	recvList    chan []StreamingClient

	ctx              context.Context
	actionDispatcher ActionDispatcher
}

type clientMeta struct {
	cancelFunc context.CancelFunc
	client     StreamingClient
}

// NewStreamingConnectionManager creates an instance of WebsocketClientManager.
func NewStreamingConnectionManager(ctx context.Context, dispatcher ActionDispatcher, clientFactory StreamingClientFactory) *StreamingConnectionManager {
	return &StreamingConnectionManager{
		ctx:              ctx,
		clients:          make(map[StreamingClient]context.CancelFunc),
		register:         make(chan *clientMeta),
		unregister:       make(chan StreamingClient),
		requestList:      make(chan bool),
		recvList:         make(chan []StreamingClient),
		actionDispatcher: dispatcher,
		clientFactory:    clientFactory,
	}
}

func (m *StreamingConnectionManager) Clients() []StreamingClient {
	m.requestList <- true
	clients := <-m.recvList
	return clients
}

// Run runs the manager. It manages multiple websocket clients.
func (m *StreamingConnectionManager) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case meta := <-m.register:
			m.clients[meta.client] = meta.cancelFunc
		case client := <-m.unregister:
			if cancelFunc, ok := m.clients[client]; ok {
				cancelFunc()
				delete(m.clients, client)
			}
		case <-m.requestList:
			clients := []StreamingClient{}
			for client := range m.clients {
				clients = append(clients, client)
			}
			m.recvList <- clients
		}
	}
}

// ClientFromRequest creates a websocket client from a http request.
func (m *StreamingConnectionManager) ClientFromRequest(dashConfig config.Dash, w http.ResponseWriter, r *http.Request) (StreamingClient, error) {

	client, cancel, err := m.clientFactory.NewConnection(w, r, m, dashConfig)
	if err != nil {
		return nil, err
	}
	m.register <- &clientMeta{
		cancelFunc: func() {
			cancel()
			m.unregister <- client
		},
		client: client,
	}

	return client, nil
}

func (m *StreamingConnectionManager) TemporaryClientFromLoadingRequest(w http.ResponseWriter, r *http.Request) (StreamingClient, error) {
	client, cancel, err := m.clientFactory.NewTemporaryConnection(w, r, m)
	if err != nil {
		return nil, err
	}
	m.register <- &clientMeta{
		cancelFunc: func() {
			cancel()
			m.unregister <- client
		},
		client: client,
	}

	return client, nil
}

func (m *StreamingConnectionManager) Get(id string) event.WSEventSender {
	for _, client := range m.Clients() {
		if id == client.ID() {
			return client
		}
	}
	return nil
}

func (m *StreamingConnectionManager) Context() context.Context {
	return m.ctx
}

func (m *StreamingConnectionManager) ActionDispatcher() ActionDispatcher {
	return m.actionDispatcher
}

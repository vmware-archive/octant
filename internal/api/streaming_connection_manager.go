/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"net/http"

	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/google/uuid"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/octant"
)

//go:generate mockgen -destination=./fake/mock_client_manager.go -package=fake github.com/vmware-tanzu/octant/internal/api ClientManager
//go:generate mockgen -destination=./fake/mock_client_factory.go -package=fake github.com/vmware-tanzu/octant/internal/api StreamingClientFactory
//go:generate mockgen -destination=./fake/mock_streaming_client.go -package=fake github.com/vmware-tanzu/octant/internal/api StreamingClient

// ClientManager is an interface for managing clients.
type ClientManager interface {
	Run(ctx context.Context)
	Clients() []StreamingClient
	ClientFromRequest(dashConfig config.Dash, w http.ResponseWriter, r *http.Request) (StreamingClient, error)
	TemporaryClientFromLoadingRequest(w http.ResponseWriter, r *http.Request) (StreamingClient, error)
	Get(id string) event.WSEventSender
}

type clientMeta struct {
	cancelFunc context.CancelFunc
	client     StreamingClient
}

type StreamingClientFactory interface {
	NewConnection(uuid.UUID, http.ResponseWriter, *http.Request, *StreamingConnectionManager, config.Dash) (StreamingClient, context.CancelFunc, error)
	NewTemporaryConnection(uuid.UUID, http.ResponseWriter, *http.Request, *StreamingConnectionManager) (StreamingClient, context.CancelFunc, error)
}

// StreamingClient is the interface responsible for sending and receiving
// streaming data to a users session, usually in a browser.
type StreamingClient interface {
	OctantClient

	Receive() (StreamRequest, error)

	Handlers() map[string][]octant.ClientRequestHandler
	State() octant.State
}

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

var _ ClientManager = (*StreamingConnectionManager)(nil)

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
	clientID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	client, cancel, err := m.clientFactory.NewConnection(clientID, w, r, m, dashConfig)
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
	clientID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	client, cancel, err := m.clientFactory.NewTemporaryConnection(clientID, w, r, m)
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

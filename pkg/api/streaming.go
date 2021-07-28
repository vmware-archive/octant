/*
   Copyright (c) 2019 the Octant contributors. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package api

import (
	"context"
	"net/http"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/event"
)

//go:generate mockgen -destination=./fake/mock_client_manager.go -package=fake github.com/vmware-tanzu/octant/pkg/api ClientManager
//go:generate mockgen -destination=./fake/mock_client_factory.go -package=fake github.com/vmware-tanzu/octant/pkg/api StreamingClientFactory
//go:generate mockgen -destination=./fake/mock_streaming_client.go -package=fake github.com/vmware-tanzu/octant/pkg/api StreamingClient
//go:generate mockgen -destination=./fake/mock_octant_client.go -package=fake github.com/vmware-tanzu/octant/pkg/api OctantClient

// ClientManager is an interface for managing clients.
type ClientManager interface {
	Run(ctx context.Context)
	Clients() []StreamingClient
	ClientFromRequest(dashConfig config.Dash, w http.ResponseWriter, r *http.Request) (StreamingClient, error)
	TemporaryClientFromLoadingRequest(w http.ResponseWriter, r *http.Request) (StreamingClient, error)
	Get(id string) event.WSEventSender
	Context() context.Context
	ActionDispatcher() ActionDispatcher
}

type StreamRequest struct {
	Type    string         `json:"type"`
	Payload action.Payload `json:"payload"`
}

// OctantClient is the interface responsible for sending streaming data to a
// users session, usually in a browser.
type OctantClient interface {
	Send(event.Event)
	ID() string
	StopCh() <-chan struct{}
}

// StreamingClient is the interface responsible for sending and receiving
// streaming data to a users session, usually in a browser.
type StreamingClient interface {
	OctantClient

	Receive() (StreamRequest, error)

	Handlers() map[string][]octant.ClientRequestHandler
	State() octant.State
}

type StreamingClientFactory interface {
	NewConnection(http.ResponseWriter, *http.Request, ClientManager, config.Dash) (StreamingClient, context.CancelFunc, error)
	NewTemporaryConnection(http.ResponseWriter, *http.Request, ClientManager) (StreamingClient, context.CancelFunc, error)
}

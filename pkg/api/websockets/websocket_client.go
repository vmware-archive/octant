/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package websockets

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/vmware-tanzu/octant/internal/util/json"
	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/log"

	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"

	"github.com/vmware-tanzu/octant/internal/config"
	internalLog "github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
)

const (
	// writeWait is the time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// pongWait is the time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// pingPeriod is the how often the client will send pings to peer with this period.
	// Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize is the maximum message size allowed from peer.
	maxMessageSize = 2 * 1024 * 1024 // 2MiB
)

// WebsocketClient manages websocket clients.
type WebsocketClient struct {
	conn       *websocket.Conn
	send       chan event.Event
	dashConfig config.Dash
	logger     log.Logger
	ctx        context.Context
	cancel     context.CancelFunc
	manager    api.ClientManager

	isOpen   atomic.Value
	state    octant.State
	handlers map[string][]octant.ClientRequestHandler
	id       uuid.UUID
	stopCh   chan struct{}
}

// NewWebsocketClient creates an instance of WebsocketClient.
func NewWebsocketClient(ctx context.Context, conn *websocket.Conn, manager api.ClientManager, dashConfig config.Dash, actionDispatcher api.ActionDispatcher, id uuid.UUID) *WebsocketClient {
	logger := dashConfig.Logger().With("component", "websocket-client", "client-id", id.String())
	ctx = internalLog.WithLoggerContext(ctx, logger)

	ctx, cancel := context.WithCancel(ctx)

	client := &WebsocketClient{
		ctx:        ctx,
		cancel:     cancel,
		conn:       conn,
		id:         id,
		send:       make(chan event.Event),
		manager:    manager,
		dashConfig: dashConfig,
		logger:     logger,
		handlers:   make(map[string][]octant.ClientRequestHandler),
		stopCh:     make(chan struct{}, 1),
	}

	state := NewWebsocketState(dashConfig, actionDispatcher, client)
	go state.Start(ctx)

	client.state = state
	for _, handler := range state.Handlers() {
		client.RegisterHandler(handler)
	}

	for _, handler := range dashConfig.ModuleManager().ClientRequestHandlers() {
		client.RegisterHandler(handler)
	}

	logger.Debugf("created websocket client")

	return client
}

// NewTemporaryWebsocketClient creates an instance of WebsocketClient
func NewTemporaryWebsocketClient(ctx context.Context, conn *websocket.Conn, manager api.ClientManager, actionDispatcher api.ActionDispatcher, id uuid.UUID) *WebsocketClient {
	ctx, cancel := context.WithCancel(ctx)
	logger := internalLog.From(ctx)

	client := &WebsocketClient{
		ctx:      ctx,
		cancel:   cancel,
		conn:     conn,
		id:       id,
		send:     make(chan event.Event),
		manager:  manager,
		logger:   logger,
		handlers: make(map[string][]octant.ClientRequestHandler),
		stopCh:   make(chan struct{}, 1),
	}

	state := NewTemporaryWebsocketState(actionDispatcher, client)
	go state.Start(ctx)

	client.state = state
	for _, handler := range state.Handlers() {
		client.RegisterHandler(handler)
	}

	return client
}

// ID returns the ID of the websocket client.
func (c *WebsocketClient) ID() string {
	return c.id.String()
}

func (c *WebsocketClient) Handlers() map[string][]octant.ClientRequestHandler {
	return c.handlers
}

func (c *WebsocketClient) readPump() {
	defer func() {
		c.isOpen.Store(false)
		c.logger.Debugf("closing read pump")
	}()

	go func() {
		<-c.ctx.Done()
		if err := c.conn.Close(); err != nil {
			c.logger.WithErr(err).Errorf("Close websocket connection")
		}

	}()

	c.isOpen.Store(true)

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.logger.WithErr(err).Errorf("Set websocket read deadline")
		return
	}
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		request, err := c.Receive()
		if err != nil {
			if errors.IsFatalStreamError(err) {
				c.cancel()
				break
			}

			continue
		}

		if err := handleStreamingMessage(c, request); err != nil {
			c.logger.WithErr(err).Errorf("Handle websocket message")
		}
	}

	close(c.stopCh)
}

func handleStreamingMessage(client api.StreamingClient, request api.StreamRequest) error {
	handlers, ok := client.Handlers()[request.Type]
	if !ok {
		return handleUnknownRequest(client, request)
	}

	var g errgroup.Group

	for _, handler := range handlers {
		g.Go(func() error {
			return handler.Handler(client.State(), request.Payload)
		})
	}

	if err := g.Wait(); err != nil {
		client.Send(event.CreateEvent("handlerError", action.Payload{
			"requestType": request.Type,
			"error":       err.Error(),
		}))

	}

	return nil
}

func handleUnknownRequest(client api.OctantClient, request api.StreamRequest) error {
	message := "unknown request"
	if request.Type != "" {
		message = fmt.Sprintf("unknown request %s", request.Type)

	}
	m := map[string]interface{}{
		"message": message,
		"payload": request.Payload,
	}
	client.Send(event.CreateEvent(event.EventTypeUnknown, m))
	return nil
}

func (c *WebsocketClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.logger.Debugf("closing write pump")
	}()

	done := false
	for !done {
		select {
		case <-c.ctx.Done():
			done = true
			break
		case response, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.WithErr(err).Errorf("Update websocket write deadline")
				return
			}
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			data, err := json.Marshal(response)
			if err != nil {
				c.logger.WithErr(err).Errorf("Marshal websocket response")
				return
			}
			if _, err := w.Write(data); err != nil {
				c.logger.WithErr(err).Errorf("Write websocket response")
				return
			}

			if err := w.Close(); err != nil {
				c.logger.WithErr(err).Errorf("Close websocket writer")
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.WithErr(err).Errorf("Set websocket write deadline")
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.WithErr(err).Errorf("Send websocket ping")
				return
			}
		}
	}
}

func (c *WebsocketClient) Send(ev event.Event) {
	v := c.isOpen.Load()
	if v != nil && v.(bool) {
		c.send <- ev
	}
}

func (c *WebsocketClient) Receive() (api.StreamRequest, error) {
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(
			err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived,
		) {
			c.logger.WithErr(err).Errorf("Unhandled websocket error")
		}
		return api.StreamRequest{}, errors.FatalStreamError(err)
	}

	var request api.StreamRequest
	if err := json.Unmarshal(message, &request); err != nil {
		c.logger.WithErr(err).Errorf("Unmarshaling websocket message")
		return api.StreamRequest{}, err
	}

	return request, nil
}

func (c *WebsocketClient) State() octant.State {
	return c.state
}

// StopCh returns the client's stop channel. It will be closed when the WebsocketClient is closed.
func (c *WebsocketClient) StopCh() <-chan struct{} {
	return c.stopCh
}

func (c *WebsocketClient) RegisterHandler(handler octant.ClientRequestHandler) {
	c.handlers[handler.RequestType] = append(c.handlers[handler.RequestType], handler)
}

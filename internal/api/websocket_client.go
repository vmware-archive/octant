/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/log"
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
	maxMessageSize = 512
)

// WebsocketClient manages websocket clients.
type WebsocketClient struct {
	conn       *websocket.Conn
	send       chan octant.Event
	dashConfig config.Dash
	logger     log.Logger
	ctx        context.Context
	cancel     context.CancelFunc
	manager    *WebsocketClientManager

	isOpen   bool
	state    octant.State
	handlers map[string][]octant.ClientRequestHandler
	id       uuid.UUID
}

var _ OctantClient = (*WebsocketClient)(nil)

// NewWebsocketClient creates an instance of WebsocketClient.
func NewWebsocketClient(ctx context.Context, conn *websocket.Conn, manager *WebsocketClientManager, dashConfig config.Dash, actionDispatcher ActionDispatcher, id uuid.UUID) *WebsocketClient {
	logger := dashConfig.Logger().With("component", "websocket-client", "client-id", id.String())

	ctx, cancel := context.WithCancel(ctx)

	client := &WebsocketClient{
		ctx:        ctx,
		cancel:     cancel,
		conn:       conn,
		id:         id,
		send:       make(chan octant.Event),
		manager:    manager,
		dashConfig: dashConfig,
		logger:     logger,
		handlers:   make(map[string][]octant.ClientRequestHandler),
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

	return client
}

// ID returns the ID of the websocket client.
func (c *WebsocketClient) ID() string {
	return c.id.String()
}

func (c *WebsocketClient) readPump() {
	defer func() {
		c.isOpen = false
	}()

	go func() {
		<-c.ctx.Done()
		if err := c.conn.Close(); err != nil {
			c.logger.WithErr(err).Errorf("Close websocket connection")
		}

	}()

	c.isOpen = true

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.logger.WithErr(err).Errorf("Set websocket read deadline")
		return
	}
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				c.logger.WithErr(err).Errorf("Unhandled websocket error")
			}
			c.cancel()
			break
		}

		if err := c.handle(message); err != nil {
			c.logger.WithErr(err).Errorf("Handle websocket message")
		}
	}
}

func (c *WebsocketClient) handle(message []byte) error {
	var request websocketRequest
	if err := json.Unmarshal(message, &request); err != nil {
		return err
	}

	handlers, ok := c.handlers[request.Type]
	if !ok {
		return c.handleUnknownRequest(request)
	}

	var g errgroup.Group

	for _, handler := range handlers {
		g.Go(func() error {
			return handler.Handler(c.state, request.Payload)
		})
	}

	if err := g.Wait(); err != nil {
		c.Send(CreateEvent("handlerError", action.Payload{
			"requestType": request.Type,
			"error":       err.Error(),
		}))

	}

	return nil
}

func (c *WebsocketClient) handleUnknownRequest(request websocketRequest) error {
	message := "unknown request"
	if request.Type != "" {
		message = fmt.Sprintf("unknown request %s", request.Type)

	}
	m := map[string]interface{}{
		"message": message,
		"payload": request.Payload,
	}
	c.Send(CreateEvent(octant.EventTypeUnknown, m))
	return nil
}

func (c *WebsocketClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
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

func (c *WebsocketClient) Send(ev octant.Event) {
	if c.isOpen {
		c.send <- ev
	}
}

func (c *WebsocketClient) RegisterHandler(handler octant.ClientRequestHandler) {
	c.handlers[handler.RequestType] = append(c.handlers[handler.RequestType], handler)
}

type websocketRequest struct {
	Type    string         `json:"type"`
	Payload action.Payload `json:"payload"`
}

func CreateEvent(eventType octant.EventType, fields action.Payload) octant.Event {
	return octant.Event{
		Type: eventType,
		Data: fields,
	}
}

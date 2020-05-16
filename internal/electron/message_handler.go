/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package electron

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/davecgh/go-spew/spew"

	"github.com/vmware-tanzu/octant/internal/log"
)

// MessageHandler handles a message.
type MessageHandler interface {
	// Key is the key for this handler.
	Key() string
	// Handle processes the message.
	Handle(ctx context.Context, in json.RawMessage) (interface{}, error)
}

// MessageIn is an incoming message.
type MessageIn struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// MessageOut is an outgoing message.
type MessageOut struct {
	Name    string      `json:"name"`
	Payload interface{} `json:"payload,omitempty"`
}

// CreateMessageOut creates an outgoing message.
func CreateMessageOut(name string, payload interface{}) *MessageOut {
	return &MessageOut{
		Name:    name + ".callback",
		Payload: payload,
	}
}

// MessageListener listens for messages.
type MessageListener interface {
	// Handle handles a message using one of the registered message handlers.
	Handle(ctx context.Context, w *astilectron.Window, message MessageIn) (interface{}, error)
	// Register registers a message handler.
	Register(handlers ...MessageHandler)
	// Unregister unregisters a message handler by key.
	Unregister(key string)
}

// DefaultMessageListener listens for messages.
type DefaultMessageListener struct {
	handlers map[string]MessageHandler

	mu sync.RWMutex
}

var _ MessageListener = &DefaultMessageListener{}

// NewMessageListener creates an instance of DefaultMessageListener.
func NewMessageListener() *DefaultMessageListener {
	m := DefaultMessageListener{
		handlers: make(map[string]MessageHandler),
	}

	return &m
}

// Handle handles a message.
func (m *DefaultMessageListener) Handle(ctx context.Context, _ *astilectron.Window, message MessageIn) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	h, ok := m.handlers[message.Name]
	if !ok {
		return nil, nil
	}

	p, err := h.Handle(ctx, message.Payload)
	if err != nil {
		return nil, fmt.Errorf("handle message of type %s: %w", message.Name, err)
	}

	return p, nil
}

// Register registers a handler.
func (m *DefaultMessageListener) Register(handlers ...MessageHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, handler := range handlers {
		m.handlers[handler.Key()] = handler
	}

}

// Unregister unregisters a listener key.
func (m *DefaultMessageListener) Unregister(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.handlers, key)
}

func handleMessage(ctx context.Context, w *astilectron.Window, listener MessageListener, logger astikit.SeverityLogger) astilectron.ListenerMessage {
	return func(m *astilectron.EventMessage) interface{} {
		var in MessageIn
		if err := m.Unmarshal(&in); err != nil {
			logger.Errorf("unmarshaling message %+v failed: %v", *m, err)
			return nil
		}

		p, err := listener.Handle(ctx, w, in)
		if err != nil {
			logger.Errorf("handling message %+v failed: %v", *m, err)
		}

		if p != nil {
			out := CreateMessageOut(in.Name, p)
			if err != nil {
				out.Name = "error"
			}

			return out
		}

		return nil
	}
}

// CallbackMessage represents a bootstrap message callback
type CallbackMessage func(m *MessageIn)

// SendMessage sends a message
func SendMessage(ctx context.Context, w *astilectron.Window, name string, payload interface{}, cs ...CallbackMessage) {
	var callbacks []astilectron.CallbackMessage
	for _, c := range cs {
		callbacks = append(callbacks, func(e *astilectron.EventMessage) {
			var m *MessageIn
			if e != nil {
				m = &MessageIn{}
				if err := e.Unmarshal(m); err != nil {
					return
				}
			}
			c(m)
		})
	}

	if err := w.SendMessage(MessageOut{Name: name, Payload: payload}, callbacks...); err != nil {
		logger := log.From(ctx)
		logger.WithErr(err).
			With("name", name,
				"payload", spew.Sdump(payload)).
			Errorf("send message")
	}
}

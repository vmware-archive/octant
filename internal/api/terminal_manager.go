/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/event"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
)

const RequestTerminalScrollback = "sendTerminalScrollback"

type terminalStateManager struct {
	config         config.Dash
	poller         Poller
	sendScrollback map[string]bool

	mu sync.Mutex
}

type terminalOutput struct {
	Scrollback string `json:"scrollback,omitempty"`
	Line       string `json:"line,omitempty"`
	New        bool   `json:"new,omitempty"`
}

var _ StateManager = (*terminalStateManager)(nil)

func NewTerminalStateManager(config config.Dash) *terminalStateManager {
	return &terminalStateManager{
		config:         config,
		poller:         NewInterruptiblePoller("terminal"),
		sendScrollback: map[string]bool{},
	}
}

// Handlers returns a slice of handlers.
func (c *terminalStateManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: RequestTerminalScrollback,
			Handler:     c.SendTerminalScrollback,
		},
	}
}

// SetContext sets the current context.
func (c *terminalStateManager) SendTerminalScrollback(state octant.State, payload action.Payload) error {
	terminalID, err := payload.String("terminalID")
	if err != nil {
		return errors.Wrap(err, "extract terminal ID from payload")
	}
	c.setSendScrollback(terminalID, true)
	return nil
}

func (s *terminalStateManager) setSendScrollback(id string, v bool) {
	s.mu.Lock()
	s.sendScrollback[id] = v
	s.mu.Unlock()
}

func (s *terminalStateManager) Start(ctx context.Context, state octant.State, client OctantClient) {
	ch := make(chan struct{}, 1)
	defer func() {
		close(ch)
	}()
	s.poller.Run(ctx, ch, s.runUpdate(state, client), event.TerminalStreamDelay)
}

func (s *terminalStateManager) runUpdate(state octant.State, client OctantClient) PollerFunc {
	return func(ctx context.Context) bool {
		tm := s.config.TerminalManager()
		for _, t := range tm.List(ctx) {
			line, err := t.Read(ctx, s.config.Logger())
			if err != nil {
				//TODO: report error directly to Terminal
				s.config.Logger().Errorf("%s", err)
			}

			sendScrollback, ok := s.sendScrollback[t.ID()]
			if line == nil && (!ok || !sendScrollback) {
				return false
			}

			key := t.Key()
			eventType := octant.EventType(fmt.Sprintf("terminals/namespace/%s/pod/%s/container/%s/%s", key.Namespace, key.Name, t.Container(), t.ID()))
			data := terminalOutput{Line: string(line)}

			if ok && sendScrollback {
				data.Scrollback = string(t.Scrollback())
				s.setSendScrollback(t.ID(), false)
			}
			terminalEvent := octant.Event{
				Type: eventType,
				Data: data,
				Err:  nil,
			}
			client.Send(terminalEvent)
		}
		return false
	}
}

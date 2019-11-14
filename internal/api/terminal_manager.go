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

const (
	RequestTerminalScrollback = "sendTerminalScrollback"
	RequestTerminalCommand    = "sendTerminalCommand"
	RequestTerminalResize     = "sendTerminalResize"
)

type terminalStateManager struct {
	config         config.Dash
	poller         Poller
	sendScrollback map[string]bool

	commands map[string][]string
	resize   map[string][]uint16

	mu sync.Mutex
}

type terminalOutput struct {
	Scrollback []byte `json:"scrollback,omitempty"`
	Line       []byte `json:"line,omitempty"`
}

var _ StateManager = (*terminalStateManager)(nil)

func NewTerminalStateManager(cfg config.Dash) *terminalStateManager {
	return &terminalStateManager{
		config:         cfg,
		poller:         NewInterruptiblePoller("terminal"),
		sendScrollback: map[string]bool{},
		resize:         map[string][]uint16{},
	}
}

// Handlers returns a slice of handlers.
func (c *terminalStateManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: RequestTerminalScrollback,
			Handler:     c.SendTerminalScrollback,
		},
		{
			RequestType: RequestTerminalCommand,
			Handler:     c.SendTerminalCommand,
		},
		{
			RequestType: RequestTerminalResize,
			Handler:     c.SendTerminalResize,
		},
	}
}

func (c *terminalStateManager) SendTerminalResize(state octant.State, payload action.Payload) error {
	terminalID, err := payload.String("terminalID")
	if err != nil {
		return errors.Wrap(err, "extract terminal ID from payload")
	}

	rows, err := payload.Uint16("rows")
	if err != nil {
		return errors.Wrap(err, "extract rows from payload")
	}

	cols, err := payload.Uint16("cols")
	if err != nil {
		return errors.Wrap(err, "extract cols from payload")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.resize[terminalID] = []uint16{rows, cols}
	return nil
}

func (c *terminalStateManager) SendTerminalCommand(state octant.State, payload action.Payload) error {
	terminalID, err := payload.String("terminalID")
	if err != nil {
		return errors.Wrap(err, "extract terminal ID from payload")
	}

	command, err := payload.String("command")
	if err != nil {
		return errors.Wrap(err, "extract command from payload")
	}

	c.appendCommand(terminalID, command)
	return nil
}

func (s *terminalStateManager) appendCommand(id, command string) {
	s.mu.Lock()
	_, ok := s.commands[id]
	if !ok {
		s.commands = map[string][]string{id: []string{command}}
	} else {
		s.commands[id] = append(s.commands[id], command)
	}
	s.mu.Unlock()
}

func (s *terminalStateManager) SendTerminalScrollback(state octant.State, payload action.Payload) error {
	terminalID, err := payload.String("terminalID")
	if err != nil {
		return errors.Wrap(err, "extract terminal ID from payload")
	}
	s.setSendScrollback(terminalID, true)
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
			line, err := t.Read(ctx)
			if err != nil {
				//TODO: report error directly to Terminal
				s.config.Logger().Errorf("%s", err)
			}

			size, ok := s.resize[t.ID()]
			if ok {
				t.Resize(ctx, size[0], size[1])
				s.mu.Lock()
				delete(s.resize, t.ID())
				s.mu.Unlock()
			}

			sendScrollback, ok := s.sendScrollback[t.ID()]
			if line == nil && (!ok || !sendScrollback) {
				commands, ok := s.commands[t.ID()]
				if ok && len(commands) != 0 {
					var command string
					s.mu.Lock()
					command, s.commands[t.ID()] = commands[len(commands)-1], commands[:len(commands)-1]
					s.mu.Unlock()
					t.Exec(ctx, command)
				}
				return false
			}

			key := t.Key()
			eventType := octant.EventType(fmt.Sprintf("terminals/namespace/%s/pod/%s/container/%s/%s", key.Namespace, key.Name, t.Container(), t.ID()))
			data := terminalOutput{Line: line}

			if ok && sendScrollback {
				data.Scrollback = t.Scrollback()
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

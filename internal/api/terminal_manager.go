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
	config config.Dash
	poller Poller

	sendScrollback sync.Map
	commands       sync.Map
	resize         sync.Map
}

type terminalOutput struct {
	Scrollback []byte `json:"scrollback,omitempty"`
	Line       []byte `json:"line,omitempty"`
}

var _ StateManager = (*terminalStateManager)(nil)

// NewTerminalStateManager returns a terminal state manager.
func NewTerminalStateManager(cfg config.Dash) StateManager {
	return &terminalStateManager{
		config: cfg,
		poller: NewInterruptiblePoller("terminal"),
	}
}

// Handlers returns a slice of handlers.
func (s *terminalStateManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: RequestTerminalScrollback,
			Handler:     s.SendTerminalScrollback,
		},
		{
			RequestType: RequestTerminalCommand,
			Handler:     s.SendTerminalCommand,
		},
		{
			RequestType: RequestTerminalResize,
			Handler:     s.SendTerminalResize,
		},
	}
}

func (s *terminalStateManager) SendTerminalResize(state octant.State, payload action.Payload) error {
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

	s.resize.Store(terminalID, []uint16{cols, rows})
	return nil
}

func (s *terminalStateManager) SendTerminalCommand(state octant.State, payload action.Payload) error {
	terminalID, err := payload.String("terminalID")
	if err != nil {
		return errors.Wrap(err, "extract terminal ID from payload")
	}

	command, err := payload.String("command")
	if err != nil {
		return errors.Wrap(err, "extract command from payload")
	}

	s.appendCommand(terminalID, command)
	return nil
}

func (s *terminalStateManager) appendCommand(id, command string) {
	var cmds []string
	commands, ok := s.commands.Load(id)
	if !ok {
		cmds = []string{command}
	} else {
		cmds = commands.([]string)
		cmds = append(cmds, command)
	}
	s.commands.Store(id, cmds)
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
	s.sendScrollback.Store(id, v)
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

			size, ok := s.resize.Load(t.ID())
			if ok {
				val := size.([]uint16)
				s.resize.Delete(t.ID())
				t.Resize(ctx, val[0], val[1])
			}

			sendScrollback, ok := s.sendScrollback.Load(t.ID())
			if line == nil && (!ok || !sendScrollback.(bool)) {
				commands, ok := s.commands.Load(t.ID())
				if ok {
					cmds := commands.([]string)
					if len(cmds) != 0 {
						var command string
						command = cmds[len(cmds)-1]
						s.commands.Store(t.ID(), cmds[:len(cmds)-1])
						t.Exec(ctx, command)
					}
				}
				return false
			}

			key := t.Key()
			eventType := octant.EventType(fmt.Sprintf("terminals/namespace/%s/pod/%s/container/%s/%s", key.Namespace, key.Name, t.Container(), t.ID()))
			data := terminalOutput{Line: line}

			if ok && sendScrollback.(bool) {
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

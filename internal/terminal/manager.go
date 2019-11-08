/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import "context"

//go:generate mockgen -source=manager.go -destination=./fake/mock_interface.go -package=fake github.com/vmware-tanzu/octant/internal/terminal TerminalManager

type terminalManager struct {
	terminals map[string]Terminal
}

var _ TerminalManager = (*terminalManager)(nil)

// NewTerminalManager creates a concrete TerminalMananger
func NewTerminalManager(ctx context.Context) (TerminalManager, error) {
	tm := &terminalManager{
		terminals: map[string]Terminal{},
	}
	return tm, nil
}

func (tm *terminalManager) Create(ctx context.Context) Terminal {
	t := NewTerminal(ctx)
	tm.terminals[t.ID(ctx)] = t
	return t
}

func (tm *terminalManager) Get(ctx context.Context, id string) (Terminal, bool) {
	v, ok := tm.terminals[id]
	return v, ok
}

func (tm *terminalManager) List(ctx context.Context) []Terminal {
	terminals := make([]Terminal, len(tm.terminals))
	for _, terminal := range tm.terminals {
		terminals = append(terminals, terminal)
	}
	return terminals
}

func (tm *terminalManager) StopAll(ctx context.Context) error {
	for _, terminal := range tm.terminals {
		terminal.Stop(ctx)
	}
	return nil
}

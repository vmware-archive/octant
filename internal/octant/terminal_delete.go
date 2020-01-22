/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/terminal"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/store"
)

type terminalDelete struct {
	logger          log.Logger
	objectStore     store.Store
	terminalManager terminal.Manager
}

var _ action.Dispatcher = (*terminalDelete)(nil)

// NewTerminalDelete creates a new terminal delete action dispatcher.
func NewTerminalDelete(logger log.Logger, objectStore store.Store, terminalManager terminal.Manager) action.Dispatcher {
	td := &terminalDelete{
		objectStore:     objectStore,
		terminalManager: terminalManager,
	}
	td.logger = logger.With("actionName", td.ActionName())
	return td
}

func (t *terminalDelete) ActionName() string {
	return "overview/deleteTerminal"
}

func (t *terminalDelete) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	terminalID, err := payload.String("terminalID")
	if err != nil {
		return err
	}
	t.terminalManager.Delete(terminalID)
	return nil
}

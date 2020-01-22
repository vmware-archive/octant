/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/terminal"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// TerminalCommandExec command executor.
type TerminalCommandExec struct {
	logger          log.Logger
	objectStore     store.Store
	terminalManager terminal.Manager
}

var _ action.Dispatcher = (*TerminalCommandExec)(nil)

// NewTerminalCommandExec creates an instance of TerminalCommandExec.
func NewTerminalCommandExec(logger log.Logger, objectStore store.Store, terminalManager terminal.Manager) *TerminalCommandExec {
	tce := &TerminalCommandExec{
		objectStore:     objectStore,
		terminalManager: terminalManager,
	}
	tce.logger = logger.With("actionName", tce.ActionName())
	return tce
}

// ActionName returns the name of this action.
func (t *TerminalCommandExec) ActionName() string {
	return "overview/commandExec"
}

// Handle executing a command.
func (t *TerminalCommandExec) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	t.logger.With("payload", payload).Debugf("received action payload")
	request, err := terminalExecFromPayload(payload)
	if err != nil {
		return errors.Wrap(err, "handle terminal exec")
	}
	t.logger.Debugf("%s", request)

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}
	t.logger.Debugf("%s", key)

	_, err = t.terminalManager.Create(ctx, t.logger, key, request.container, request.command)
	if err != nil {
		t.logger.Errorf("%s", err)
		return errors.Wrap(err, "terminal manager create")
	}

	return nil
}

type terminalExecRequest struct {
	container string
	command   string
	tty       bool
}

func terminalExecFromPayload(payload action.Payload) (*terminalExecRequest, error) {
	var err error
	t := &terminalExecRequest{tty: false}

	t.container, err = payload.String("containerName")
	if err != nil {
		return nil, err
	}

	t.command, err = payload.String("containerCommand")
	if err != nil {
		return nil, err
	}

	return t, nil
}

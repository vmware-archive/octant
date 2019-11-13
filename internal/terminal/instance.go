/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/store"
)

//go:generate mockgen -source=terminal.go -destination=./fake/mock_interface.go -package=fake github.com/vmware-tanzu/octant/internal/terminal Terminal

type instance struct {
	id        uuid.UUID
	key       store.Key
	createdAt time.Time
	stdout    io.ReadWriter
	stderr    io.ReadWriter
	stdin     io.Reader

	container  string
	command    string
	scrollback bytes.Buffer

	stdoutPiper *io.PipeReader
	stdoutPipew *io.PipeWriter
}

var _ Instance = (*instance)(nil)

// NewTerminal creates a concrete Terminal
func NewTerminalInstance(ctx context.Context, key store.Key, container, command string) Instance {
	t := &instance{
		id:        uuid.New(),
		key:       key,
		createdAt: time.Now(),
		stdout:    &bytes.Buffer{},
		stderr:    &bytes.Buffer{},
		stdin:     nil,
		container: container,
		command:   command,
	}

	return t
}

func (t *instance) Stream(ctx context.Context, logger log.Logger) {
	logger.Debugf("starting exec stream for %s", t.id.String())
	t.stdoutPiper, t.stdoutPipew = io.Pipe()

	go func() {
		defer t.stdoutPipew.Close()
		if _, err := io.Copy(t.stdoutPipew, t.stdout); err != nil {
			logger.Errorf("%v", err)
		}
	}()
}

func (t *instance) Read(ctx context.Context, logger log.Logger) ([]byte, error) {
	if t.stdoutPiper == nil {
		return nil, errors.New("stdout is nil, call Stream before Read")
	}
	buf := make([]byte, 256)
	n, err := t.stdoutPiper.Read(buf)
	if err != nil {
		if err == io.EOF {
			line := buf[:n]
			if string(line) == "" {
				return nil, nil
			}
			if _, err := t.scrollback.Write(line); err != nil {
				return nil, err
			}
			return line, nil
		}
		return nil, err
	}
	b := buf[:n]
	if _, err := t.scrollback.Write(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (t *instance) Stop(ctx context.Context) {}
func (t *instance) Key() store.Key           { return t.key }
func (t *instance) Scrollback() []byte       { return t.scrollback.Bytes() }
func (t *instance) ID() string               { return t.id.String() }
func (t *instance) Container() string        { return t.container }
func (t *instance) Command() string          { return t.command }
func (t *instance) CreatedAt() time.Time     { return t.createdAt }
func (t *instance) Stdin() io.Reader         { return t.stdin }
func (t *instance) Stdout() io.Writer        { return t.stdout }
func (t *instance) Stderr() io.Writer        { return t.stderr }

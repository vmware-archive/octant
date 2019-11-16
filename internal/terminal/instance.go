/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"bytes"
	"context"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/store"
	"k8s.io/client-go/tools/remotecommand"
)

//go:generate mockgen -source=terminal.go -destination=./fake/mock_interface.go -package=fake github.com/vmware-tanzu/octant/internal/terminal Terminal

type pty struct {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue

	logger   log.Logger
	commands []string

	out        io.ReadWriter
	rows, cols uint16
	mu         sync.Mutex
}

func (p *pty) Write(b []byte) (int, error) {
	return p.out.Write(b)
}

func (p *pty) Read(b []byte) (int, error) {
	c, ok := p.pop()
	if !ok {
		return 0, nil
	}
	return copy(b, c), nil
}

func (p *pty) Next() *remotecommand.TerminalSize {
	return &remotecommand.TerminalSize{Width: p.cols, Height: p.rows}
}

func (p *pty) resize(rows, cols uint16) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.rows = rows
	p.cols = cols
}

func (p *pty) push(command string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.commands = append(p.commands, command)
}

func (p *pty) pop() ([]byte, bool) {
	var command string
	if len(p.commands) > 0 {
		p.mu.Lock()
		defer p.mu.Unlock()
		command, p.commands = p.commands[len(p.commands)-1], p.commands[:len(p.commands)-1]
		return []byte(command), true
	}
	return []byte{}, false
}

func (p *pty) stdout() io.Reader {
	return p.out
}

type instance struct {
	id        uuid.UUID
	key       store.Key
	createdAt time.Time

	container  string
	command    string
	scrollback bytes.Buffer

	pty *pty
	tty bool

	logger log.Logger
}

var _ Instance = (*instance)(nil)

// NewTerminal creates a concrete Terminal
func NewTerminalInstance(ctx context.Context, logger log.Logger, key store.Key, container, command string, tty bool) Instance {
	t := &instance{
		id:        uuid.New(),
		key:       key,
		createdAt: time.Now(),
		container: container,
		command:   command,
		pty:       &pty{logger: logger, out: &bytes.Buffer{}},
		tty:       tty,
		logger:    logger,
	}

	return t
}

func (t *instance) Resize(ctx context.Context, cols, rows uint16) {
	t.pty.resize(rows, cols)
}

func (t *instance) Read(ctx context.Context) ([]byte, error) {
	if t.pty == nil {
		return nil, nil
	}
	buf := make([]byte, 4096)
	n, err := t.pty.stdout().Read(buf)
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

func (t *instance) Exec(ctx context.Context, command string) error {
	if t.pty == nil {
		return errors.New("can not execute command, no stdin")
	}
	t.pty.push(command)
	return nil
}

func (t *instance) Stop(ctx context.Context) {}
func (t *instance) Key() store.Key           { return t.key }
func (t *instance) Scrollback() []byte       { return t.scrollback.Bytes() }
func (t *instance) ID() string               { return t.id.String() }
func (t *instance) Container() string        { return t.container }
func (t *instance) Command() string          { return t.command }
func (t *instance) CreatedAt() time.Time     { return t.createdAt }
func (t *instance) Stdin() io.Reader         { return t.pty }
func (t *instance) Stdout() io.Writer        { return t.pty }
func (t *instance) Stderr() io.Writer {
	if t.tty {
		return nil
	}
	return t.pty
}

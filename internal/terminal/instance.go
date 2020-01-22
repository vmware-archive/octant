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

//go:generate mockgen -source=instance.go -destination=./fake/mock_instance.go -package=fake github.com/vmware-tanzu/octant/internal/terminal Instance

// Instance defines the interface to a single exec instance.
type Instance interface {
	ID() string
	Key() store.Key
	Container() string
	Command() string
	Scrollback() []byte

	Read(size int) ([]byte, error)
	Write(key []byte) error
	Resize(cols, rows uint16)

	Stop()
	Active() bool
	SetExitMessage(string)
	ExitMessage() string
	CreatedAt() time.Time

	PTY() PTY
}

type PTY interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

type pty struct {
	ctx      context.Context
	cancelFn context.CancelFunc

	logger       log.Logger
	keystroke    chan []byte
	resize       chan remotecommand.TerminalSize
	activityFunc func()

	out  io.ReadWriter
	size *remotecommand.TerminalSize

	mu sync.RWMutex
}

var _ PTY = (*pty)(nil)

// Write writes bytes to the "stdout/err" buffer.
func (p *pty) Write(b []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	defer p.activityFunc()

	return p.out.Write(b)
}

// Read reads bytes from stdin (p.keystroke) and copies them to the stdin buffer.
func (p *pty) Read(b []byte) (int, error) {
	select {
	case <-p.ctx.Done():
		p.cancelFn()

		if p.ctx.Err() != nil {
			if p.ctx.Err() == context.Canceled {
				return 0, io.ErrClosedPipe
			}

			return 0, io.ErrUnexpectedEOF
		}

		return 0, io.ErrClosedPipe
	default:
		break
	}

	key, ok := <-p.keystroke
	if !ok {
		return 0, nil
	}

	defer p.activityFunc()

	return copy(b, key), nil
}

// Next creates a new TerminalSize based on resize events.
func (p *pty) Next() *remotecommand.TerminalSize {
	select {
	case <-p.ctx.Done():
		return nil
	default:
		break
	}

	size, ok := <-p.resize
	if !ok {
		return nil
	}

	return &size
}

// ReadStdout reads from the buffer that Write writes stdout/err to.
func (p *pty) ReadStdout(buf []byte) (int, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.out.Read(buf)
}

type instance struct {
	ctx context.Context

	id        uuid.UUID
	key       store.Key
	createdAt time.Time

	container   string
	command     string
	exitMessage string

	scrollback bytes.Buffer

	pty *pty

	logger log.Logger
}

var _ Instance = (*instance)(nil)

// NewTerminalInstance creates a concrete Terminal
func NewTerminalInstance(ctx context.Context, logger log.Logger, key store.Key, container, command string, activityChan chan Instance) Instance {
	ctx, cancelFn := context.WithCancel(ctx)

	termPty := &pty{
		ctx:       ctx,
		cancelFn:  cancelFn,
		logger:    logger,
		out:       &bytes.Buffer{},
		keystroke: make(chan []byte, 25),
		resize:    make(chan remotecommand.TerminalSize, 2),
		size:      &remotecommand.TerminalSize{},
	}

	t := &instance{
		ctx:       ctx,
		id:        uuid.New(),
		key:       key,
		createdAt: time.Now(),
		container: container,
		command:   command,
		pty:       termPty,
		logger:    logger,
	}

	termPty.activityFunc = func() {
		activityChan <- t
	}

	return t
}

func (t *instance) Resize(cols, rows uint16) {
	t.pty.resize <- remotecommand.TerminalSize{
		Width:  cols,
		Height: rows,
	}
}

// Read attempts to read from the stdout bytes.Buffer. As a side-effect
// of calling Read any data that is read is also appended to the internal
// scollback buffer that can be retrieved by calling Scrollback.
func (t *instance) Read(size int) ([]byte, error) {
	if t.pty == nil {
		return nil, nil
	}

	buf := make([]byte, size)

	n, err := t.pty.ReadStdout(buf)
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

// Write sends the passed in key to the stdin of the instance.
// If the instance is not a TTY, Write will return an error.
func (t *instance) Write(key []byte) error {
	if t.pty == nil {
		return errors.New("can not execute command, no stdin")
	}
	t.pty.keystroke <- key

	return nil
}

// SetExitMessage sets tne exit message for the terminal instance.
func (t *instance) SetExitMessage(m string) { t.exitMessage = m }

// ExitMessage returns the exit message for the terminal instance.
func (t *instance) ExitMessage() string { return t.exitMessage }

// Active returns if the terminal is currently active.
func (t *instance) Active() bool {
	select {
	case <-t.ctx.Done():
		return false
	default:
		return true
	}
}

// Stop stops the terminal from attempting to read/write to stdout/in streams.
// Calling stop will also cause the PTY to return an io.ErrClosedPipe from the PTY
// Read command.
func (t *instance) Stop() { t.pty.cancelFn() }

// Key returns the store.Key for the Pod that this terminal is associated with.
func (t *instance) Key() store.Key { return t.key }

// Scrollback returns the scrollback buffer for the terminal instance. Scrollback buffer
// is populated by calling Read.
func (t *instance) Scrollback() []byte { return t.scrollback.Bytes() }

// ID returns the ID for the termianl. This is a UUID returned as a string.
func (t *instance) ID() string { return t.id.String() }

// Container returns the container name that the terminal is associated with.
func (t *instance) Container() string { return t.container }

// Command returns the command that was used to stat this terminal.
func (t *instance) Command() string { return t.command }

// CreatedAt returns the date/time this terminal was created.
func (t *instance) CreatedAt() time.Time { return t.createdAt }

func (t *instance) PTY() PTY {
	return t.pty
}

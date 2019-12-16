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
	"k8s.io/client-go/tools/remotecommand"
)

//go:generate mockgen -source=instance.go -destination=./fake/mock_instance.go -package=fake github.com/vmware-tanzu/octant/internal/terminal Instance

// Instance defines the interface to a single exec instance.
type Instance interface {
	ID() string
	Key() store.Key
	Container() string
	Command() string
	TTY() bool
	Scrollback() []byte

	Read(size int) ([]byte, error)
	Write(key []byte) error
	Resize(cols, rows uint16)

	Stop()
	Active() bool
	SetExitMessage(string)
	ExitMessage() string
	CreatedAt() time.Time

	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	SizeQueue() remotecommand.TerminalSizeQueue
}

type pty struct {
	ctx      context.Context
	cancelFn context.CancelFunc

	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue

	logger    log.Logger
	keystroke chan []byte
	resize    chan []uint16

	out  io.ReadWriter
	size *remotecommand.TerminalSize
}

func (p *pty) Write(b []byte) (int, error) {
	return p.out.Write(b)
}

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
	case key := <-p.keystroke:
		return copy(b, key), nil
	default:
		return 0, nil
	}
}

func (p *pty) Next() *remotecommand.TerminalSize {
	select {
	case <-p.ctx.Done():
		return nil
	case size := <-p.resize:
		p.size.Width, p.size.Height = size[0], size[1]
		return p.size
	default:
		return p.size
	}
}

func (p *pty) stdout() io.Reader {
	return p.out
}

type instance struct {
	ctx context.Context

	id        uuid.UUID
	key       store.Key
	createdAt time.Time

	container   string
	command     string
	exitMessage string
	scrollback  bytes.Buffer

	pty *pty
	tty bool

	logger log.Logger
}

var _ Instance = (*instance)(nil)

// NewTerminalInstance creates a concrete Terminal
func NewTerminalInstance(ctx context.Context, logger log.Logger, key store.Key, container, command string, tty bool) Instance {
	ctx, cancelFn := context.WithCancel(ctx)

	termPty := &pty{
		ctx:       ctx,
		cancelFn:  cancelFn,
		logger:    logger,
		out:       &bytes.Buffer{},
		keystroke: make(chan []byte, 25),
		resize:    make(chan []uint16),
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
		tty:       tty,
		logger:    logger,
	}

	return t
}

func (t *instance) Resize(cols, rows uint16) {
	// TODO: fix this, currently when uncommented interactive terminals do not work.
	// t.pty.resize <- []uint16{cols, rows}
}

// Read attempts to read from the stdout bytes.Buffer. As a side-effect
// of calling Read any data that is read is also appended to the internal
// scollback buffer that can be retrieved by calling Scrollback.
func (t *instance) Read(size int) ([]byte, error) {
	if t.pty == nil {
		return nil, nil
	}

	buf := make([]byte, size)

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

// Active returns if the terminal is currently active. Active terminals
// are non-TTY commands that are still streaming output OR tty terminals
// that have not been exited.
func (t *instance) Active() bool { return t.ctx.Err() == nil }

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

// TTY returns a boolean if this terminal was started as a TTY.
func (t *instance) TTY() bool { return t.tty }

// CreatedAt returns the date/time this terminal was created.
func (t *instance) CreatedAt() time.Time { return t.createdAt }

func (t *instance) SizeQueue() remotecommand.TerminalSizeQueue { return t.pty }
func (t *instance) Stdin() io.Reader                           { return t.pty }
func (t *instance) Stdout() io.Writer                          { return t.pty }
func (t *instance) Stderr() io.Writer {
	if t.tty {
		return nil
	}

	return t.pty
}

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

	Read() ([]byte, error)
	Exec(key []byte) error
	Resize(cols, rows uint16)

	Stop()
	SetExitMessage(string)
	ExitMessage() string
	CreatedAt() time.Time

	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

type pty struct {
	ctx context.Context
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue

	logger    log.Logger
	commands  []string
	keystroke chan []byte
	resize    chan []uint16

	out        io.ReadWriter
	rows, cols uint16
	mu         sync.Mutex
}

func (p *pty) Write(b []byte) (int, error) {
	return p.out.Write(b)
}

func (p *pty) Read(b []byte) (int, error) {
	select {
	case <-p.ctx.Done():
		return 0, io.ErrClosedPipe
	case key := <-p.keystroke:
		return copy(b, key), nil
	default:
		if p.ctx.Err() != nil {
			if p.ctx.Err() == context.Canceled {
				return 0, io.ErrClosedPipe
			}
			return 0, io.ErrUnexpectedEOF
		}
		return 0, nil
	}
}

func (p *pty) Next() *remotecommand.TerminalSize {
	select {
	case size := <-p.resize:
		p.cols, p.rows = size[0], size[1]
		return &remotecommand.TerminalSize{Width: p.cols, Height: p.rows}
	default:
		return &remotecommand.TerminalSize{Width: p.cols, Height: p.rows}
	}
}

func (p *pty) stdout() io.Reader {
	return p.out
}

type instance struct {
	ctx      context.Context
	cancelFn context.CancelFunc

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
		logger:    logger,
		out:       &bytes.Buffer{},
		keystroke: make(chan []byte, 25),
		resize:    make(chan []uint16, 1),
	}

	t := &instance{
		ctx:       ctx,
		cancelFn:  cancelFn,
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
	//TODO figure out why sending terminal size breaks stream.
	//t.pty.resize <- []uint16{cols, rows}
}

func (t *instance) Read() ([]byte, error) {
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

func (t *instance) Exec(key []byte) error {
	if t.pty == nil {
		return errors.New("can not execute command, no stdin")
	}
	t.pty.keystroke <- key
	return nil
}

func (t *instance) SetExitMessage(m string) { t.exitMessage = m }
func (t *instance) ExitMessage() string     { return t.exitMessage }
func (t *instance) Stop()                   { t.cancelFn() }
func (t *instance) Key() store.Key          { return t.key }
func (t *instance) Scrollback() []byte      { return t.scrollback.Bytes() }
func (t *instance) ID() string              { return t.id.String() }
func (t *instance) Container() string       { return t.container }
func (t *instance) Command() string         { return t.command }
func (t *instance) CreatedAt() time.Time    { return t.createdAt }
func (t *instance) Stdin() io.Reader        { return t.pty }
func (t *instance) Stdout() io.Writer       { return t.pty }
func (t *instance) Stderr() io.Writer {
	if t.tty {
		return nil
	}
	return t.pty
}

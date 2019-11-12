/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/vmware-tanzu/octant/internal/log"
)

//go:generate mockgen -source=terminal.go -destination=./fake/mock_interface.go -package=fake github.com/vmware-tanzu/octant/internal/terminal Terminal

type instance struct {
	id        uuid.UUID `json:"id"`
	createdAt time.Time `json:"createdAt"`
	stdout    io.ReadWriter
	stderr    io.ReadWriter
	stdin     io.Reader

	scrollback []string
}

var _ Instance = (*instance)(nil)

// NewTerminal creates a concrete Terminal
func NewTerminalInstance(ctx context.Context) Instance {
	t := &instance{
		id:        uuid.New(),
		createdAt: time.Now(),
		stdout:    &bytes.Buffer{},
		stderr:    &bytes.Buffer{},
		stdin:     nil,
	}
	t.stream(ctx)

	return t
}

func (t *instance) ID(ctx context.Context) string {
	return t.id.String()
}

func (t *instance) Exec(ctx context.Context) error {
	return nil
}

func (t *instance) Scrollback(ctx context.Context) []string {
	return t.scrollback
}

func (t *instance) stream(ctx context.Context) {
	logger := log.From(ctx)

	scanner := bufio.NewScanner(t.stdout)
	go func() {
		for ctx.Err() == nil && scanner.Scan() {
			t.scrollback = append(t.scrollback, scanner.Text())
		}
		if scanner.Err() != nil {
			logger.With("TerminalInstance").Errorf("%s", scanner.Err())
			t.Stop(ctx)
		}
	}()
}

func (t *instance) Stop(ctx context.Context) error {
	return nil
}

func (t *instance) CreatedAt(ctx context.Context) time.Time {
	return t.createdAt
}

func (t *instance) Stdin() io.Reader  { return t.stdin }
func (t *instance) Stdout() io.Writer { return t.stdout }
func (t *instance) Stderr() io.Writer { return t.stderr }

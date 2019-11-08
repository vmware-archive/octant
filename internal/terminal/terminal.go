/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=terminal.go -destination=./fake/mock_interface.go -package=fake github.com/vmware-tanzu/octant/internal/terminal Terminal

type terminal struct {
	id uuid.UUID
	createdAt time.Time
}

var _ Terminal = (*terminal)(nil)

// NewTerminal creates a concrete Terminal
func NewTerminal(ctx context.Context) Terminal {
	t := &terminal{
		id: uuid.New(),
		createdAt: time.Now(),
	}
	return t
}

func (t *terminal) ID(ctx context.Context) string {
	return t.id.String()
}

func (t *terminal) Exec(ctx context.Context) error {
	return nil
}

func (t *terminal) Stream(ctx context.Context) (io.ReadCloser, error) {
	return nil, nil
}

func (t *terminal) Stop(ctx context.Context) error {
	return nil
}

func (t *terminal) CreatedAt(ctx context.Context) time.Time {
	return t.createdAt
}
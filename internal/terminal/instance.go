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

type instance struct {
	id        uuid.UUID `json:"id"`
	createdAt time.Time `json:"createdAt"`
}

var _ Instance = (*instance)(nil)

// NewTerminal creates a concrete Terminal
func NewTerminalInstance(ctx context.Context) Instance {
	t := &instance{
		id:        uuid.New(),
		createdAt: time.Now(),
	}
	return t
}

func (t *instance) ID(ctx context.Context) string {
	return t.id.String()
}

func (t *instance) Exec(ctx context.Context) error {
	return nil
}

func (t *instance) Stream(ctx context.Context) (io.ReadCloser, error) {
	return nil, nil
}

func (t *instance) Stop(ctx context.Context) error {
	return nil
}

func (t *instance) CreatedAt(ctx context.Context) time.Time {
	return t.createdAt
}

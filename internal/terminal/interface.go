/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"
	"io"
	"time"
)

// Terminal defines the interface to a single terminal.
type Terminal interface {
	ID(ctx context.Context) string
	Exec(ctx context.Context) error
	Stream(ctx context.Context) (io.ReadCloser, error)
	Stop(ctx context.Context) error
	CreatedAt(ctx context.Context) time.Time
}

// TerminalManager defines the interface for querying terminals.
type TerminalManager interface {
	List(ctx context.Context) []Terminal
	Get(ctx context.Context, ID string) (Terminal, bool)
	Create(ctx context.Context) Terminal
	StopAll(ctx context.Context) error
}

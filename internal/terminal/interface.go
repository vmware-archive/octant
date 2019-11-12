/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"
	"io"
	"time"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// Instance defines the interface to a single exec instance.
type Instance interface {
	ID() string
	Key() store.Key
	Container() string
	Command() string
	Scrollback() []string

	Stream(ctx context.Context, logger log.Logger)
	Stop(ctx context.Context)
	CreatedAt() time.Time

	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

// Manager defines the interface for querying terminal instance.
type Manager interface {
	List(ctx context.Context) []Instance
	Get(ctx context.Context, ID string) (Instance, bool)
	Create(ctx context.Context, logger log.Logger, key store.Key, container string, command string) (Instance, error)
	StopAll(ctx context.Context) error
}

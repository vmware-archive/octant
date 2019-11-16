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
	Scrollback() []byte

	Read() ([]byte, error)
	Exec(key []byte) error
	Resize(cols, rows uint16)

	Stop()
	CreatedAt() time.Time

	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

// Manager defines the interface for querying terminal instance.
type Manager interface {
	List() []Instance
	Get(ID string) (Instance, bool)
	Delete(id string)
	Create(ctx context.Context, logger log.Logger, key store.Key, container string, command string, tty bool) (Instance, error)
	StopAll() error
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"
	"io"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Instance defines the interface to a single exec instance.
type Instance interface {
	ID(ctx context.Context) string
	Exec(ctx context.Context) error
	Stream(ctx context.Context) (io.ReadCloser, error)
	Stop(ctx context.Context) error
	CreatedAt(ctx context.Context) time.Time
}

// Manager defines the interface for querying terminal instance.
type Manager interface {
	List(ctx context.Context) []Instance
	Get(ctx context.Context, ID string) (Instance, bool)
	Create(ctx context.Context, gvk schema.GroupVersionKind, name string, namespace string, container string, command string) (Instance, error)
	StopAll(ctx context.Context) error
}

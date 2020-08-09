/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"

	"github.com/dop251/goja"
)

//go:generate mockgen -destination=./fake/mock_dashboard_client.go -package=fake . DashboardClientFunction

// DashboardClientFactory is an interface for a factory that creates dashboard clients.
type DashboardClientFactory interface {
	// Create creates a dashboard clients in a goja value.
	Create(ctx context.Context, vm *goja.Runtime) goja.Value
}

// DashboardClientFunction is a function in the the dashboard client.
type DashboardClientFunction interface {
	// Name returns the name of the function.
	Name() string
	// Call generates a function that executes the function.
	Call(ctx context.Context, vm *goja.Runtime) func(c goja.FunctionCall) goja.Value
}

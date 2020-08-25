/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package javascript

import (
	"context"
	"fmt"

	"github.com/dop251/goja"

	"github.com/vmware-tanzu/octant/internal/octant"
)

// ModularDashboardClientFactory is a modular octant.DashboardClientFactory. It configures
// itself based on functions passed when it is initialized.
type ModularDashboardClientFactory struct {
	functions []octant.DashboardClientFunction
}

var _ octant.DashboardClientFactory = &ModularDashboardClientFactory{}

// NewModularDashboardClientFactory creates an instance of ModularDashboardClientFactory.
func NewModularDashboardClientFactory(functions []octant.DashboardClientFunction) *ModularDashboardClientFactory {
	m := &ModularDashboardClientFactory{
		functions: functions,
	}
	return m
}

// Create creates a dashboard client javascript value.
func (m *ModularDashboardClientFactory) Create(ctx context.Context, vm *goja.Runtime) goja.Value {
	obj := vm.NewObject()

	for _, fn := range m.functions {
		if err := obj.Set(fn.Name(), fn.Call(ctx, vm)); err != nil {
			return vm.NewGoError(err)
		}
	}

	return obj
}

// OctantClient is a client for interacting with Octant.
type OctantClient interface {
	octant.LinkGenerator
	octant.Storage
}

// DefaultFunctions are the default functions for the ModularDashboardClientFactory.
func DefaultFunctions(octantClient OctantClient) []octant.DashboardClientFunction {
	return []octant.DashboardClientFunction{
		NewDashboardGet(octantClient),
		NewDashboardList(octantClient),
		NewDashboardUpdate(octantClient),
		NewDashboardDelete(octantClient),
		NewDashboardRefPath(octantClient),
		NewDashboardCreateOrUpdateFromYAML(octantClient),
	}
}

// panicMessage creates a message for a panic given an error and an optional reason.
// If the reason is blank, it will be omitted.
func panicMessage(vm *goja.Runtime, err error, reason string) goja.Value {
	if reason == "" {
		return vm.ToValue(err.Error())
	}

	return vm.ToValue(fmt.Sprintf("%s: %s", reason, err.Error()))
}

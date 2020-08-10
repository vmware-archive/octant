/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package javascript

import (
	"context"
	"errors"

	"github.com/dop251/goja"

	"github.com/vmware-tanzu/octant/internal/octant"
)

// DashboardUpdate is a function that updates YAML. The text can send one
// or more objects.
type DashboardUpdate struct {
	storage octant.Storage
}

var _ octant.DashboardClientFunction = &DashboardUpdate{}

// NewDashboardUpdate creates an instance of DashboardUpdate.
func NewDashboardUpdate(storage octant.Storage) *DashboardUpdate {
	d := &DashboardUpdate{
		storage: storage,
	}
	return d
}

// Name returns the name of this function. It will always return "Update".
func (d *DashboardUpdate) Name() string {
	return "Update"
}

// Call create a function call that sends YAML to the Kubernetes cluster. If the the cluster
// rejects the YAML, it will throw a javascript exception.
func (d *DashboardUpdate) Call(ctx context.Context, vm *goja.Runtime) func(c goja.FunctionCall) goja.Value {
	return func(c goja.FunctionCall) goja.Value {
		namespace := c.Argument(0).String()
		update := c.Argument(1).String()

		results, err := d.storage.ObjectStore().CreateOrUpdateFromYAML(ctx, namespace, update)
		if err != nil {
			panic(panicMessage(vm, errors.New("no YAML was supplied"), ""))
		}

		return vm.ToValue(results)
	}
}

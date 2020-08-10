/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package javascript

import (
	"context"

	"github.com/dop251/goja"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// DashboardGet is a function that gets an object by key.
type DashboardGet struct {
	storage octant.Storage
}

var _ octant.DashboardClientFunction = &DashboardGet{}

// NewDashboardGet creates an instance of DashboardGet.
func NewDashboardGet(storage octant.Storage) *DashboardGet {
	d := &DashboardGet{
		storage: storage,
	}
	return d
}

// Name returns the name of this function. It will always return "Get".
func (d *DashboardGet) Name() string {
	return "Get"
}

// Call creates a function call that gets an object by key. If the key is invalid, or if the
// get is unsuccessful, it will throw a javascript exception.
func (d *DashboardGet) Call(ctx context.Context, vm *goja.Runtime) func(c goja.FunctionCall) goja.Value {
	return func(c goja.FunctionCall) goja.Value {
		var key store.Key
		obj := c.Argument(0).ToObject(vm)

		// This will never error since &key is a pointer to a type.
		_ = vm.ExportTo(obj, &key)

		u, err := d.storage.ObjectStore().Get(ctx, key)
		if err != nil {
			panic(panicMessage(vm, err, ""))
		}

		return vm.ToValue(u.Object)
	}
}

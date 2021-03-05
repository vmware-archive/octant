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

// DashboardDelete is a function for deleting an object by key.
type DashboardDelete struct {
	storage octant.Storage
}

var _ octant.DashboardClientFunction = &DashboardDelete{}

// NewDashboardDelete creates an instance of DashboardDelete.
func NewDashboardDelete(storage octant.Storage) *DashboardDelete {
	d := &DashboardDelete{
		storage: storage,
	}
	return d
}

// Name returns the name of this function. It will always return "Delete".
func (d *DashboardDelete) Name() string {
	return "Delete"
}

// Call creates a function call that deletes an object by key. If the key is invalid, or if the
// delete is unsuccessful, it will throw a javascript exception.
func (d *DashboardDelete) Call(ctx context.Context, vm *goja.Runtime) func(c goja.FunctionCall) goja.Value {
	return func(c goja.FunctionCall) goja.Value {
		newCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		var key store.Key
		obj := c.Argument(0).ToObject(vm)

		// This will never error since &key is a pointer to a type.
		_ = vm.ExportTo(obj, &key)

		metadataArg := c.Argument(1)
		if !goja.IsUndefined(metadataArg) {
			var metadata map[string]string
			metadataObj := metadataArg.ToObject(vm)

			// This will not error as js plugins restrict this type
			// and we handle both cases
			_ = vm.ExportTo(metadataObj, &metadata)
			for k, val := range metadata {
				newCtx = context.WithValue(newCtx, DashboardMetadataKey(k), val)
			}
		}

		if err := d.storage.ObjectStore().Delete(newCtx, key); err != nil {
			panic(panicMessage(vm, err, ""))
		}
		return goja.Undefined()
	}
}

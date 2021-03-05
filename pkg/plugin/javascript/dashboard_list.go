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
	"github.com/vmware-tanzu/octant/pkg/store"
)

// DashboardList is a function that lists objects by key.
type DashboardList struct {
	storage octant.Storage
}

var _ octant.DashboardClientFunction = &DashboardList{}

// NewDashboardList creates an instance of DashboardList.
func NewDashboardList(storage octant.Storage) *DashboardList {
	d := &DashboardList{
		storage: storage,
	}
	return d
}

// Name returns the name of this function. It will always return "List".
func (d *DashboardList) Name() string {
	return "List"
}

// Call creates a function call that lists objects by key. If the key is invalid, or if the
// list is unsuccessful, it will throw a javascript exception.
func (d *DashboardList) Call(ctx context.Context, vm *goja.Runtime) func(c goja.FunctionCall) goja.Value {
	return func(c goja.FunctionCall) goja.Value {
		newCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		m := map[string]interface{}{}

		var key store.Key
		obj := c.Argument(0).ToObject(vm)

		// This will never error since &m is a pointer to a type.
		_ = vm.ExportTo(obj, &m)

		key, err := store.KeyFromPayload(m)
		if err != nil {
			panicMessage(vm, fmt.Errorf("key is invalid: %w", err), "")
		}

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

		u, _, err := d.storage.ObjectStore().List(newCtx, key)
		if err != nil {
			panic(panicMessage(vm, err, ""))
		}

		items := make([]interface{}, len(u.Items))
		for i := 0; i < len(u.Items); i++ {
			items[i] = u.Items[i].Object
		}

		return vm.ToValue(items)
	}
}

/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package javascript

import (
	"context"
	"fmt"

	"github.com/dop251/goja"

	"github.com/vmware-tanzu/octant/pkg/store"
)

type dashboardClient struct {
	objectStore store.Store
	vm          *goja.Runtime
	ctx         context.Context
}

func CreateDashClientObject(ctx context.Context, objStore store.Store, vm *goja.Runtime) goja.Value {
	d := dashboardClient{
		ctx:         ctx,
		objectStore: objStore,
		vm:          vm,
	}

	obj := d.vm.NewObject()
	if err := obj.Set("Get", d.Get); err != nil {
		return d.vm.NewGoError(err)
	}
	if err := obj.Set("List", d.List); err != nil {
		return d.vm.NewGoError(err)
	}
	if err := obj.Set("Create", d.Create); err != nil {
		return d.vm.NewGoError(err)
	}
	if err := obj.Set("Update", d.Create); err != nil {
		return d.vm.NewGoError(err)
	}
	if err := obj.Set("Delete", d.Delete); err != nil {
		return d.vm.NewGoError(err)
	}
	return obj
}

func (d *dashboardClient) Delete(c goja.FunctionCall) goja.Value {
	var key store.Key
	obj := c.Argument(0).ToObject(d.vm)
	if err := d.vm.ExportTo(obj, &key); err != nil {
		return d.vm.NewTypeError(fmt.Errorf("dashboardClient.Delete: %w", err))
	}

	if err := d.objectStore.Delete(d.ctx, key); err != nil {
		return d.vm.NewGoError(err)
	}
	return goja.Undefined()
}

func (d *dashboardClient) Get(c goja.FunctionCall) goja.Value {
	var key store.Key
	obj := c.Argument(0).ToObject(d.vm)
	if err := d.vm.ExportTo(obj, &key); err != nil {
		return d.vm.NewGoError(fmt.Errorf("dashboardClient.Get: %w", err))
	}

	u, err := d.objectStore.Get(d.ctx, key)
	if err != nil {
		return d.vm.NewGoError(err)
	}

	return d.vm.ToValue(u.Object)
}

func (d *dashboardClient) List(c goja.FunctionCall) goja.Value {
	var key store.Key
	obj := c.Argument(0).ToObject(d.vm)
	if err := d.vm.ExportTo(obj, &key); err != nil {
		return d.vm.NewGoError(fmt.Errorf("dashboardClient.List: %w", err))
	}

	u, _, err := d.objectStore.List(d.ctx, key)
	if err != nil {
		return d.vm.NewGoError(err)
	}

	items := make([]interface{}, len(u.Items))
	for i := 0; i < len(u.Items); i++ {
		items[i] = u.Items[i].Object
	}

	return d.vm.ToValue(items)
}

func (d *dashboardClient) Create(c goja.FunctionCall) goja.Value {
	namespace := c.Argument(0).String()
	update := c.Argument(1).String()

	if namespace == "" {
		return d.vm.NewTypeError(fmt.Errorf("create/update: invalid namespace"))
	}

	if update == "" {
		return d.vm.NewTypeError(fmt.Errorf("create/update: empty yaml"))
	}

	results, err := d.objectStore.CreateOrUpdateFromYAML(d.ctx, namespace, update)
	if err != nil {
		return d.vm.NewTypeError(fmt.Errorf("create/update: %w", err))
	}

	return d.vm.ToValue(results)
}

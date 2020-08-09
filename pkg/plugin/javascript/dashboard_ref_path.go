/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package javascript

import (
	"context"

	"github.com/dop251/goja"

	"github.com/vmware-tanzu/octant/internal/octant"
)

// DashboardRefPath is a function that returns the path for a ref.
type DashboardRefPath struct {
	linkGenerator octant.LinkGenerator
}

var _ octant.DashboardClientFunction = &DashboardRefPath{}

// NewDashboardRefPath creates an instance of DashboardRefPath.
func NewDashboardRefPath(linkGenerator octant.LinkGenerator) *DashboardRefPath {
	d := &DashboardRefPath{
		linkGenerator: linkGenerator,
	}
	return d
}

// Name returns the name of this function. It will always return "RefPath".
func (d *DashboardRefPath) Name() string {
	return "RefPath"
}

type ref struct {
	Namespace  string `json:"namespace"`
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
}

// Call create a function call generates a path for a ref. If the operation is unsuccessful,
// it will throw a javascript exception.
func (d *DashboardRefPath) Call(_ context.Context, vm *goja.Runtime) func(c goja.FunctionCall) goja.Value {
	return func(c goja.FunctionCall) goja.Value {
		obj := c.Argument(0).ToObject(vm)

		var r ref

		// This will never error since &ref is a pointer to a type.
		_ = vm.ExportTo(obj, &r)

		p, err := d.linkGenerator.ObjectPath(r.Namespace, r.APIVersion, r.Kind, r.Name)
		if err != nil {
			panic(panicMessage(vm, err, "dashboardClient.RefPath"))
		}

		return vm.ToValue(p)
	}
}

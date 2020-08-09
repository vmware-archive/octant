/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package javascript

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/dop251/goja"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/octant/fake"
)

type functionRunner struct {
	wantErr bool
}

func (fr *functionRunner) run(ctx context.Context, t *testing.T, fn octant.DashboardClientFunction, call string) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	obj := vm.NewObject()
	callFn := fn.Call(ctx, vm)
	require.NoError(t, obj.Set(fn.Name(), callFn), "set call function")

	vm.Set("dashClient", obj)
	_, err := vm.RunString(call)
	spew.Dump(err)
	if fr.wantErr {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)
}

// CreateFn creates a function.
func createFn(ctx context.Context, ctrl *gomock.Controller, vm *goja.Runtime, name string) octant.DashboardClientFunction {
	fn := fake.NewMockDashboardClientFunction(ctrl)
	fn.EXPECT().Name().Return(name)
	fn.EXPECT().Call(ctx, vm)
	return fn
}

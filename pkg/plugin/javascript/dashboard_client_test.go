/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package javascript

import (
	"context"
	"sort"
	"testing"

	"github.com/dop251/goja"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/octant"
)

func TestModularDashboardClientFactory_Create(t *testing.T) {
	type ctorArgs struct {
		functions func(ctx context.Context, ctrl *gomock.Controller, vm *goja.Runtime) []octant.DashboardClientFunction
	}
	tests := []struct {
		name     string
		ctorArgs ctorArgs
		want     []string
	}{
		{
			name: "functions are valid",
			ctorArgs: ctorArgs{
				functions: func(ctx context.Context, ctrl *gomock.Controller, vm *goja.Runtime) []octant.DashboardClientFunction {
					return []octant.DashboardClientFunction{
						createFn(ctx, ctrl, vm, "Fn2"),
						createFn(ctx, ctrl, vm, "Fn1"),
					}
				},
			},
			want: []string{"Fn1", "Fn2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			vm := goja.New()

			f := NewModularDashboardClientFactory(tt.ctorArgs.functions(ctx, ctrl, vm))

			got := f.Create(ctx, vm).Export()

			m, ok := got.(map[string]interface{})
			if !ok {
				t.Error("....")
			}

			var gotFns []string
			for k := range m {
				gotFns = append(gotFns, k)
			}
			sort.Strings(gotFns)
			require.Equal(t, tt.want, gotFns)
		})
	}

}

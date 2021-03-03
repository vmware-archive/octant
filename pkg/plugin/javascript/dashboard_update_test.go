/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package javascript

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/octant/fake"
	fake2 "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestDashboardUpdate_Name(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := fake.NewMockStorage(ctrl)

	d := NewDashboardUpdate(storage)

	want := "Update"
	got := d.Name()

	require.Equal(t, want, got)
}

func TestDashboardUpdate_Call(t *testing.T) {
	type ctorArgs struct {
		storage func(ctx context.Context, ctrl *gomock.Controller) octant.Storage
	}
	tests := []struct {
		name     string
		ctorArgs ctorArgs
		call     string
		wantErr  bool
	}{
		{
			name: "in general",
			ctorArgs: ctorArgs{
				storage: func(ctx context.Context, ctrl *gomock.Controller) octant.Storage {
					objectStore := fake2.NewMockStore(ctrl)
					ctx = context.WithValue(ctx, DashboardMetadataKey("foo"), "baz")
					ctx = context.WithValue(ctx, DashboardMetadataKey("foo"), "bar")
					ctx = context.WithValue(ctx, DashboardMetadataKey("qux"), "quuux")
					ctx = context.WithValue(ctx, DashboardMetadataKey("qux"), "quux")

					objectStore.EXPECT().
						CreateOrUpdateFromYAML(ContextType, "test", "create-yaml").
						Return([]string{"test"}, nil).
						Do(func(c context.Context, _, _ string) {
							require.Equal(t, "bar", c.Value(DashboardMetadataKey("foo")))
							require.Equal(t, "quux", c.Value(DashboardMetadataKey("qux")))
						})

					storage := fake.NewMockStorage(ctrl)
					storage.EXPECT().ObjectStore().Return(objectStore).AnyTimes()
					return storage
				},
			},
			call: `dashClient.Update('test', 'create-yaml',{"foo": "bar", "qux": "quux"})`,
		},
		{
			name: "create fails",
			ctorArgs: ctorArgs{
				storage: func(ctx context.Context, ctrl *gomock.Controller) octant.Storage {
					objectStore := fake2.NewMockStore(ctrl)
					objectStore.EXPECT().
						CreateOrUpdateFromYAML(ContextType, "test", "create-yaml").
						Return(nil, errors.New("error"))

					storage := fake.NewMockStorage(ctrl)
					storage.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

					return storage
				},
			},
			call:    `dashClient.Update('test', 'create-yaml')`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			d := NewDashboardUpdate(tt.ctorArgs.storage(ctx, ctrl))

			runner := functionRunner{wantErr: tt.wantErr}
			runner.run(ctx, t, d, tt.call)

		})
	}
}

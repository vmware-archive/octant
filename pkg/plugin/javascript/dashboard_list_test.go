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
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	fake2 "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestDashboardList_Name(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := fake.NewMockStorage(ctrl)

	d := NewDashboardList(storage)

	want := "List"
	got := d.Name()

	require.Equal(t, want, got)
}

func TestDashboardList_Call(t *testing.T) {
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
						List(ContextType, store.Key{
							Namespace:  "test",
							APIVersion: "v1",
							Kind:       "Pod",
						}).
						Return(testutil.ToUnstructuredList(t, testutil.CreatePod("pod")), false, nil).
						Do(func(c context.Context, _ store.Key) {
							require.Equal(t, "bar", c.Value(DashboardMetadataKey("foo")))
							require.Equal(t, "quux", c.Value(DashboardMetadataKey("qux")))
						})

					storage := fake.NewMockStorage(ctrl)
					storage.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

					return storage
				},
			},
			call: `dashClient.List({namespace:'test', apiVersion: 'v1', kind:'Pod'},{"foo": "bar", "qux": "quux"})`,
		},
		{
			name: "list fails",
			ctorArgs: ctorArgs{
				storage: func(ctx context.Context, ctrl *gomock.Controller) octant.Storage {
					objectStore := fake2.NewMockStore(ctrl)

					objectStore.EXPECT().
						List(ContextType, store.Key{
							Namespace:  "test",
							APIVersion: "v1",
							Kind:       "ReplicaSet",
						}).
						Return(nil, false, errors.New("error"))
					storage := fake.NewMockStorage(ctrl)
					storage.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

					return storage
				},
			},
			call:    `dashClient.List({namespace:'test', apiVersion: 'v1', kind:'ReplicaSet'})`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			d := NewDashboardList(tt.ctorArgs.storage(ctx, ctrl))

			runner := functionRunner{wantErr: tt.wantErr}
			runner.run(ctx, t, d, tt.call)

		})
	}
}

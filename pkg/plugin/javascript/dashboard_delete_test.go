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
	"github.com/vmware-tanzu/octant/pkg/store"
	fake2 "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestDashboardDelete_Name(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := fake.NewMockStorage(ctrl)

	d := NewDashboardDelete(storage)

	want := "Delete"
	got := d.Name()

	require.Equal(t, want, got)
}

func TestDashboardDelete_Call(t *testing.T) {
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

					objectStore.EXPECT().
						Delete(ContextType, store.Key{
							Namespace:  "test",
							APIVersion: "v1",
							Kind:       "ReplicaSet",
							Name:       "my-replica-set"}).
						Return(nil)

					storage := fake.NewMockStorage(ctrl)
					storage.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

					return storage
				},
			},
			call: `dashClient.Delete({namespace:'test', apiVersion: 'v1', kind:'ReplicaSet', name: 'my-replica-set'})`,
		},
		{
			name: "with arbitrary metadata",
			ctorArgs: ctorArgs{
				storage: func(ctx context.Context, ctrl *gomock.Controller) octant.Storage {
					objectStore := fake2.NewMockStore(ctrl)
					ctx = context.WithValue(ctx, DashboardMetadataKey("foo"), "baz")
					ctx = context.WithValue(ctx, DashboardMetadataKey("foo"), "bar")
					ctx = context.WithValue(ctx, DashboardMetadataKey("qux"), "quuux")
					ctx = context.WithValue(ctx, DashboardMetadataKey("qux"), "quux")

					objectStore.EXPECT().
						Delete(ContextType, store.Key{
							Namespace:  "test",
							APIVersion: "v1",
							Kind:       "ReplicaSet",
							Name:       "my-replica-set"}).
						Return(nil).
						Do(func(c context.Context, _ store.Key) {
							require.Equal(t, "bar", c.Value(DashboardMetadataKey("foo")))
							require.Equal(t, "quux", c.Value(DashboardMetadataKey("qux")))
						})

					storage := fake.NewMockStorage(ctrl)
					storage.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

					return storage
				},
			},
			call: `dashClient.Delete({namespace:'test', apiVersion: 'v1', kind:'ReplicaSet', name: 'my-replica-set'},{"foo": "bar", "qux": "quux"})`,
		},
		{
			name: "delete fails",
			ctorArgs: ctorArgs{
				storage: func(ctx context.Context, ctrl *gomock.Controller) octant.Storage {
					objectStore := fake2.NewMockStore(ctrl)
					objectStore.EXPECT().
						Delete(ContextType, store.Key{
							Namespace:  "test",
							APIVersion: "v1",
							Kind:       "ReplicaSet",
							Name:       "my-replica-set"}).
						Return(errors.New("error"))

					storage := fake.NewMockStorage(ctrl)
					storage.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

					return storage
				},
			},
			call:    `dashClient.Delete({namespace:'test', apiVersion: 'v1', kind:'ReplicaSet', name: 'my-replica-set'})`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			d := NewDashboardDelete(tt.ctorArgs.storage(ctx, ctrl))

			runner := functionRunner{wantErr: tt.wantErr}
			runner.run(ctx, t, d, tt.call)

		})
	}
}

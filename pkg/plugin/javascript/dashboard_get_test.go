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
	"github.com/vmware-tanzu/octant/pkg/plugin/api"
	"github.com/vmware-tanzu/octant/pkg/store"
	fake2 "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestDashboardGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := fake.NewMockStorage(ctrl)

	d := NewDashboardGet(storage)

	want := "Get"
	got := d.Name()

	require.Equal(t, want, got)
}

func TestDashboardGet_Call(t *testing.T) {
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
						Get(ContextType, store.Key{
							Namespace:  "test",
							APIVersion: "v1",
							Kind:       "Pod",
							Name:       "pod"}).
						Return(testutil.ToUnstructured(t, testutil.CreatePod("pod")), nil)

					storage := fake.NewMockStorage(ctrl)
					storage.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

					return storage
				},
			},
			call: `dashClient.Get({namespace:'test', apiVersion: 'v1', kind:'Pod', name: 'pod'})`,
		},
		{
			name: "with arbitrary metadata",
			ctorArgs: ctorArgs{
				storage: func(ctx context.Context, ctrl *gomock.Controller) octant.Storage {
					objectStore := fake2.NewMockStore(ctrl)
					ctx = context.WithValue(ctx, api.DashboardMetadataKey("foo"), "baz")
					ctx = context.WithValue(ctx, api.DashboardMetadataKey("foo"), "bar")
					ctx = context.WithValue(ctx, api.DashboardMetadataKey("qux"), "quuux")
					ctx = context.WithValue(ctx, api.DashboardMetadataKey("qux"), "quux")

					objectStore.EXPECT().
						Get(ContextType, store.Key{
							Namespace:  "test",
							APIVersion: "v1",
							Kind:       "Pod",
							Name:       "pod"}).
						Return(testutil.ToUnstructured(t, testutil.CreatePod("pod")), nil).
						Do(func(c context.Context, _ store.Key) {
							require.Equal(t, "bar", c.Value(api.DashboardMetadataKey("foo")))
							require.Equal(t, "quux", c.Value(api.DashboardMetadataKey("qux")))
						})

					storage := fake.NewMockStorage(ctrl)
					storage.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

					return storage
				},
			},
			call: `dashClient.Get({namespace:'test', apiVersion: 'v1', kind:'Pod', name: 'pod'},{"foo": "bar", "qux": "quux"})`,
		},
		{
			name: "delete fails",
			ctorArgs: ctorArgs{
				storage: func(ctx context.Context, ctrl *gomock.Controller) octant.Storage {
					objectStore := fake2.NewMockStore(ctrl)
					objectStore.EXPECT().
						Get(ContextType, store.Key{
							Namespace:  "test",
							APIVersion: "v1",
							Kind:       "Pod",
							Name:       "pod"}).
						Return(nil, errors.New("error"))

					storage := fake.NewMockStorage(ctrl)
					storage.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

					return storage
				},
			},
			call:    `dashClient.Get({namespace:'test', apiVersion: 'v1', kind:'Pod', name: 'pod'})`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			d := NewDashboardGet(tt.ctorArgs.storage(ctx, ctrl))

			runner := functionRunner{wantErr: tt.wantErr}
			runner.run(ctx, t, d, tt.call)

		})
	}
}

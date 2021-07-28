/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"

	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	internalErr "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/store"
	objectStoreFake "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestNewDefaultCRDWatcher_requires_object_store(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	client := clusterFake.NewMockClientInterface(controller)

	ctx := context.Background()
	_, err := NewDefaultCRDWatcher(ctx, client, nil, nil)
	require.Error(t, err)
}

func TestDefaultCRDWatcher_Watch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	client := clusterFake.NewMockClientInterface(controller)

	objectStore := objectStoreFake.NewMockStore(controller)
	objectStore.EXPECT().
		Watch(ctx, crdKey, gomock.Any()).
		DoAndReturn(func(_ context.Context, key store.Key, c *cache.ResourceEventHandlerFuncs) error {
			assert.NotNil(t, c.AddFunc)
			assert.NotNil(t, c.DeleteFunc)
			return nil
		})
	errorStore, err := internalErr.NewErrorStore()
	require.NoError(t, err)

	watcher, err := NewDefaultCRDWatcher(ctx, client, objectStore, errorStore)
	require.NoError(t, err)

	watchConfig := &config.CRDWatchConfig{
		Add:          func(ctx context.Context, object *unstructured.Unstructured) {},
		Delete:       func(ctx context.Context, object *unstructured.Unstructured) {},
		IsNamespaced: false,
	}

	require.NoError(t, watcher.AddConfig(watchConfig))
	require.NoError(t, watcher.Watch(ctx))
}

func TestDefaultCRDWatcher_Watch_failure(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	client := clusterFake.NewMockClientInterface(controller)

	objectStore := objectStoreFake.NewMockStore(controller)
	objectStore.EXPECT().
		Watch(ctx, crdKey, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ store.Key, c *cache.ResourceEventHandlerFuncs) error {
			return errors.New("failure")
		})
	errorStore, err := internalErr.NewErrorStore()
	require.NoError(t, err)

	watcher, err := NewDefaultCRDWatcher(ctx, client, objectStore, errorStore)
	require.NoError(t, err)

	watchConfig := &config.CRDWatchConfig{
		Add:          func(ctx context.Context, object *unstructured.Unstructured) {},
		Delete:       func(ctx context.Context, object *unstructured.Unstructured) {},
		IsNamespaced: false,
	}

	require.NoError(t, watcher.AddConfig(watchConfig))
	err = watcher.Watch(ctx)
	require.Error(t, err)
}

func Test_performWatch(t *testing.T) {
	object := &unstructured.Unstructured{
		Object: map[string]interface{}{},
	}

	tests := []struct {
		name       string
		canPerform func(t *testing.T) func(*unstructured.Unstructured) bool
		handler    func(t *testing.T) config.ObjectHandler
		object     interface{}
	}{
		{
			name: "in general",
			canPerform: func(t *testing.T) func(*unstructured.Unstructured) bool {
				return func(u *unstructured.Unstructured) bool {
					assert.Equal(t, object, u)
					return true
				}
			},
			handler: func(t *testing.T) config.ObjectHandler {
				return func(_ context.Context, u *unstructured.Unstructured) {
					assert.Equal(t, object, u)
				}
			},
			object: object,
		},
		{
			name: "object was not unstructured",
			canPerform: func(t *testing.T) func(*unstructured.Unstructured) bool {
				return func(u *unstructured.Unstructured) bool {
					return true
				}
			},
			handler: func(t *testing.T) config.ObjectHandler {
				return func(_ context.Context, u *unstructured.Unstructured) {
				}
			},
			object: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			f := performWatch(ctx, test.canPerform(t), test.handler(t))

			f(test.object)
		})
	}
}

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

	"github.com/heptio/developer-dash/internal/config"
	objectStoreFake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
)

func TestNewDefaultCRDWatcher_requires_object_store(t *testing.T) {
	_, err := NewDefaultCRDWatcher(nil)
	require.Error(t, err)
}

func TestDefaultCRDWatcher_Watch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	objectStore := objectStoreFake.NewMockObjectStore(controller)
	objectStore.EXPECT().
		Watch(ctx, crdKey, gomock.Any()).
		DoAndReturn(func(_ context.Context, key objectstoreutil.Key, c *cache.ResourceEventHandlerFuncs) error {
			assert.NotNil(t, c.AddFunc)
			assert.NotNil(t, c.DeleteFunc)
			return nil
		})

	watcher, err := NewDefaultCRDWatcher(objectStore)
	require.NoError(t, err)

	watchConfig := &config.CRDWatchConfig{
		Add:          func(ctx context.Context, object *unstructured.Unstructured) {},
		Delete:       func(ctx context.Context, object *unstructured.Unstructured) {},
		IsNamespaced: false,
	}

	err = watcher.Watch(ctx, watchConfig)
	require.NoError(t, err)
}

func TestDefaultCRDWatcher_Watch_failure(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	objectStore := objectStoreFake.NewMockObjectStore(controller)
	objectStore.EXPECT().
		Watch(ctx, crdKey, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ objectstoreutil.Key, c *cache.ResourceEventHandlerFuncs) error {
			return errors.New("failure")
		})

	watcher, err := NewDefaultCRDWatcher(objectStore)
	require.NoError(t, err)

	watchConfig := &config.CRDWatchConfig{
		Add:          func(ctx context.Context, object *unstructured.Unstructured) {},
		Delete:       func(ctx context.Context, object *unstructured.Unstructured) {},
		IsNamespaced: false,
	}

	err = watcher.Watch(ctx, watchConfig)
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
		object interface{}
	}{
		{
			name: "in general",
			canPerform: func(t *testing.T) func(*unstructured.Unstructured) bool {
				return func(u *unstructured.Unstructured) bool {
					assert.Equal(t, object, u)
					return true
				}
			},
			handler:    func(t *testing.T) config.ObjectHandler {
				return func(_ context.Context, u *unstructured.Unstructured)	 {
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
			handler:    func(t *testing.T) config.ObjectHandler {
				return func(_ context.Context, u *unstructured.Unstructured)	 {
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

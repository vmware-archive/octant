package describer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kcache "k8s.io/client-go/tools/cache"

	"github.com/heptio/developer-dash/internal/config"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
)

// DefaultCRDWatcher is the default CRD watcher.
type DefaultCRDWatcher struct {
	objectStore objectstore.ObjectStore
}

var _ config.CRDWatcher = (*DefaultCRDWatcher)(nil)

// NewDefaultCRDWatcher creates an instance of DefaultCRDWatcher.
func NewDefaultCRDWatcher(objectStore objectstore.ObjectStore) (*DefaultCRDWatcher, error) {
	if objectStore == nil {
		return nil, errors.New("object store is nil")
	}

	return &DefaultCRDWatcher{
		objectStore: objectStore,
	}, nil
}

var (
	crdKey = objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}
)

// Watch watches for CRDs given a configuration.
func (cw *DefaultCRDWatcher) Watch(ctx context.Context, watchConfig *config.CRDWatchConfig) error {
	if watchConfig == nil {
		return errors.New("watch config is nil")
	}

	handler := &kcache.ResourceEventHandlerFuncs{}

	if watchConfig.Add != nil {
		handler.AddFunc = performWatch(ctx, watchConfig.CanPerform, watchConfig.Add)
	}

	if watchConfig.Delete != nil {
		handler.DeleteFunc = performWatch(ctx, watchConfig.CanPerform, watchConfig.Delete)
	}

	if err := cw.objectStore.Watch(ctx, crdKey, handler); err != nil {
		return errors.WithMessage(err, "crd watcher has failed")
	}

	return nil
}

func performWatch(ctx context.Context, canPerform func(*unstructured.Unstructured) bool, handler config.ObjectHandler) func(object interface{}) {
	return func(object interface{}) {
		u, ok := object.(*unstructured.Unstructured)
		if !ok {
			logger := log.From(ctx)
			logger.
				With("object-type", fmt.Sprintf("%T", object)).
				Warnf("crd watcher received a non dynamic object")
			return
		}

		if canPerform(u) {
			handler(ctx, u)
		}
	}
}

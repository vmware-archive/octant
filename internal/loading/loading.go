package loading

import (
	"context"

	"github.com/vmware/octant/internal/describer"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/pkg/store"
)

func IsObjectLoading(ctx context.Context, namespace string, resource *describer.Resource, objectStore store.Store) bool {
	logger := log.From(ctx)

	if resource == nil {
		logger.Debugf("can't determine if a nil object is loading")
		return false
	}

	key := resource.ObjectStoreKey
	key.Namespace = namespace

	return objectStore.IsLoading(ctx, key)
}

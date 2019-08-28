package configuration

import (
	"context"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/store"
)

type ObjectDeleter struct {
	logger log.Logger
	store  store.Store
}

func NewObjectDeleter(logger log.Logger, clusterClient store.Store) *ObjectDeleter {
	return &ObjectDeleter{
		logger: logger.With("action", octant.ActionDeleteObject),
		store:  clusterClient,
	}
}

func (d *ObjectDeleter) ActionName() string {
	return octant.ActionDeleteObject
}

func (d *ObjectDeleter) Handle(ctx context.Context, payload action.Payload) error {
	d.logger.With("payload", payload).Debugf("deleting object")

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}

	return d.store.Delete(ctx, key)
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kcache "k8s.io/client-go/tools/cache"

	"github.com/vmware/octant/internal/config"
	internalErr "github.com/vmware/octant/internal/errors"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/pkg/store"
)

// DefaultCRDWatcher is the default CRD watcher.
type DefaultCRDWatcher struct {
	objectStore store.Store
	errorStore  internalErr.ErrorStore

	accessError sync.Once
	mu          sync.Mutex
}

var _ config.CRDWatcher = (*DefaultCRDWatcher)(nil)

// NewDefaultCRDWatcher creates an instance of DefaultCRDWatcher.
func NewDefaultCRDWatcher(ctx context.Context, objectStore store.Store, errorStore internalErr.ErrorStore) (*DefaultCRDWatcher, error) {
	if objectStore == nil {
		return nil, errors.New("object store is nil")
	}

	cw := &DefaultCRDWatcher{
		objectStore: objectStore,
		errorStore:  errorStore,
	}

	objectStore.RegisterOnUpdate(func(newObjectStore store.Store) {
		cw.mu.Lock()
		defer cw.mu.Unlock()

		cw.objectStore = newObjectStore

		logger := log.From(ctx)
		logger.Debugf("default crd watcher updated object store")
	})

	return cw, nil
}

var (
	crdKey = store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}
)

// Watch watches for CRDs given a configuration.
func (cw *DefaultCRDWatcher) Watch(ctx context.Context, watchConfig *config.CRDWatchConfig) error {
	if watchConfig == nil {
		return errors.New("watch config is nil")
	}

	cw.mu.Lock()
	defer cw.mu.Unlock()

	handler := &kcache.ResourceEventHandlerFuncs{}

	if watchConfig.Add != nil {
		handler.AddFunc = performWatch(ctx, watchConfig.CanPerform, watchConfig.Add)
	}

	if watchConfig.Delete != nil {
		handler.DeleteFunc = performWatch(ctx, watchConfig.CanPerform, watchConfig.Delete)
	}

	err := cw.objectStore.Watch(ctx, crdKey, handler)
	if err != nil {
		aErr, ok := err.(*internalErr.AccessError)
		if ok {
			found := cw.errorStore.Add(aErr)
			// Log if we have not seen this access error before.
			if !found {
				logger := log.From(ctx)
				logger.WithErr(aErr).Errorf("access denied")
			}
			return nil
		}
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

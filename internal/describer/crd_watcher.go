/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kcache "k8s.io/client-go/tools/cache"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/config"
	internalErr "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// DefaultCRDWatcher is the default CRD watcher.
type DefaultCRDWatcher struct {
	objectStore   store.Store
	clusterClient cluster.ClientInterface
	errorStore    internalErr.ErrorStore

	watchConfigs map[string]*config.CRDWatchConfig

	accessError sync.Once
	mu          sync.Mutex
}

var _ config.CRDWatcher = (*DefaultCRDWatcher)(nil)

// NewDefaultCRDWatcher creates an instance of DefaultCRDWatcher.
func NewDefaultCRDWatcher(ctx context.Context, clusterClient cluster.ClientInterface, objectStore store.Store, errorStore internalErr.ErrorStore) (*DefaultCRDWatcher, error) {
	if objectStore == nil {
		return nil, errors.New("object store is nil")
	}

	cw := &DefaultCRDWatcher{
		objectStore:   objectStore,
		clusterClient: clusterClient,
		errorStore:    errorStore,
		watchConfigs:  make(map[string]*config.CRDWatchConfig),
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
func (cw *DefaultCRDWatcher) Watch(ctx context.Context) error {
	logger := log.From(ctx)

	handler := &kcache.ResourceEventHandlerFuncs{
		AddFunc: func(object interface{}) {
			cw.mu.Lock()
			defer cw.mu.Unlock()

			logger := logger.With("crdwatcher", "add")

			u, ok := object.(*unstructured.Unstructured)
			if !ok {
				logger.
					With("object-type", fmt.Sprintf("%T", object)).
					Warnf("crd watcher received a non dynamic object")
				return
			}

			cw.clusterClient.ResetMapper()

			for _, watchConfig := range cw.watchConfigs {
				if watchConfig.CanPerform(u) {
					watchConfig.Add(ctx, u)
				}
			}
		},
		DeleteFunc: func(object interface{}) {
			cw.mu.Lock()
			defer cw.mu.Unlock()

			logger := logger.With("crdwatcher", "delete")

			u, ok := object.(*unstructured.Unstructured)
			if !ok {
				logger.
					With("object-type", fmt.Sprintf("%T", object)).
					Warnf("crd watcher received a non dynamic object")
				return
			}

			cw.clusterClient.ResetMapper()

			for _, watchConfig := range cw.watchConfigs {
				if watchConfig.CanPerform(u) {
					watchConfig.Delete(ctx, u)
				}
			}

			list, err := kubernetes.CRDResources(u)
			if err != nil {
				logger.WithErr(err).Errorf("unable to get group/version/kinds for CRD")

			}

			if err := cw.objectStore.Unwatch(ctx, list...); err != nil {
				logger.WithErr(err).Errorf("unable to unwatch CRD")
				return
			}

		},
	}

	err := cw.objectStore.Watch(ctx, crdKey, handler)
	if err != nil {
		aErr, ok := err.(*internalErr.AccessError)
		if ok {
			found := cw.errorStore.Add(aErr)
			// Log if we have not seen this access error before.
			if !found {
				logger.WithErr(aErr).Errorf("access denied")
			}
			return nil
		}
		return errors.WithMessage(err, "crd watcher has failed")
	}

	return nil
}

// AddConfig adds watch config to the watcher.
func (cw *DefaultCRDWatcher) AddConfig(watchConfig *config.CRDWatchConfig) error {
	if watchConfig == nil {
		return fmt.Errorf("watch config is nil")
	}

	cw.mu.Lock()
	defer cw.mu.Unlock()

	id := uuid.New().String()
	cw.watchConfigs[id] = watchConfig

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

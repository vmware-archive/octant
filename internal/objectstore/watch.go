/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstore

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kcache "k8s.io/client-go/tools/cache"

	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/third_party/k8s.io/client-go/dynamic/dynamicinformer"
)

// WatchOpt is an option for configuration Watch.
type WatchOpt func(*Watch)

// Watch is a cache which watches the cluster for updates to known objects. It wraps a dynamic cache
// by default. Since the cache knows about all cluster updates, a majority of operations for listing
// and getting objects can happen in local memory instead of requiring a network request.
type Watch struct {
	initFactoryFunc func(context.Context, cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error)
	initBackendFunc func(watch *Watch) (store.Store, error)
	client          cluster.ClientInterface
	stopCh          <-chan struct{}
	cancelFunc      context.CancelFunc
	factories       *factoriesCache
	watchedGVKs     *watchedGVKsCache
	cachedObjects   *cachedObjectsCache
	handlers        map[string]map[schema.GroupVersionKind]watchEventHandler

	backendObjectStore store.Store

	onClientUpdate chan store.Store
	updateFns      []store.UpdateFn
}

var _ store.Store = (*Watch)(nil)

func initWatchBackend(w *Watch) (store.Store, error) {
	backendObjectStore, err := NewDynamicCache(w.client, w.stopCh, func(d *DynamicCache) {
		d.initFactoryFunc = func(ctx context.Context, client cluster.ClientInterface, namespace string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			factory, ok := w.factories.get(namespace)

			if !ok {
				if err := w.HasAccess(ctx, store.Key{Namespace: metav1.NamespaceAll}, "watch"); err != nil {
					factory, err = w.initFactoryFunc(ctx, w.client, namespace)
					if err != nil {
						return nil, err
					}
				} else {
					factory, ok = w.factories.get("")
					if !ok {
						return nil, errors.New("no default DynamicInformerFactory found")
					}
				}
			}

			w.factories.set(namespace, factory)

			return factory, nil
		}
	})
	if err != nil {
		return nil, errors.Wrap(err, "initial dynamic cache")
	}

	return backendObjectStore, nil
}

// NewWatch create an instance of new watch. By default, it will create a dynamic cache as its
// backend.
func NewWatch(ctx context.Context, client cluster.ClientInterface, options ...WatchOpt) (*Watch, error) {
	c := &Watch{
		initFactoryFunc: initDynamicSharedInformerFactory,
		initBackendFunc: initWatchBackend,
		client:          client,
		factories:       initFactoriesCache(),
		watchedGVKs:     initWatchedGVKsCache(),
		cachedObjects:   initCachedObjectsCache(),
		handlers:        make(map[string]map[schema.GroupVersionKind]watchEventHandler),
		onClientUpdate:  make(chan store.Store, 10),
	}

	for _, option := range options {
		option(c)
	}

	if err := c.bootstrap(ctx, false); err != nil {
		return nil, err
	}

	return c, nil
}

func (w *Watch) bootstrap(ctx context.Context, forceBackendInit bool) error {
	logger := log.From(ctx)
	logger.With("backend-init", forceBackendInit).Debugf("bootstrapping")

	if forceBackendInit {
		w.factories = initFactoriesCache()
		w.watchedGVKs = initWatchedGVKsCache()
		w.cachedObjects = initCachedObjectsCache()
		w.handlers = make(map[string]map[schema.GroupVersionKind]watchEventHandler)
	}

	ctx, cancel := context.WithCancel(ctx)
	w.cancelFunc = cancel
	w.stopCh = ctx.Done()

	if _, ok := w.factories.get(""); !ok {
		factory, err := w.initFactoryFunc(ctx, w.client, "")
		if err != nil {
			return errors.Wrap(err, "initialize dynamic shared informer factory")
		}
		w.factories.set("", factory)
	}

	if w.backendObjectStore == nil || forceBackendInit {
		backendObjectStore, err := w.initBackendFunc(w)
		if err != nil {
			return errors.Wrap(err, "initial dynamic cache")
		}

		w.backendObjectStore = backendObjectStore
	}

	nsKey := store.Key{APIVersion: "v1", Kind: "Namespace"}
	nsHandler := &nsUpdateHandler{
		watch:  w,
		logger: log.From(ctx),
	}
	if err := w.Watch(ctx, nsKey, nsHandler); err != nil {
		return errors.Wrap(err, "create namespace watcher")
	}

	return nil
}

// HasAccess access to objects using a key
func (w *Watch) HasAccess(ctx context.Context, key store.Key, verb string) error {
	return w.backendObjectStore.HasAccess(ctx, key, verb)
}

// List lists objects using a key.
func (w *Watch) List(ctx context.Context, key store.Key) ([]*unstructured.Unstructured, error) {
	ctx, span := trace.StartSpan(ctx, "watchCacheList")
	defer span.End()

	if w.backendObjectStore == nil {
		return nil, errors.New("backend object store is nil")
	}

	// TODO: find out why this doesn't work with watch.
	logger := log.From(ctx)
	if err := w.backendObjectStore.HasAccess(ctx, key, "list"); err != nil {
		logger.Errorf("check access failed: %v", err)
		return []*unstructured.Unstructured{}, nil
	}

	gvk := key.GroupVersionKind()
	if w.isKeyCached(key) {
		var filteredObjects []*unstructured.Unstructured

		var selector = labels.Everything()
		if key.Selector != nil {
			selector = key.Selector.AsSelector()
		}

		cachedObjects := w.cachedObjects.list(key.Namespace, gvk)

		for _, object := range cachedObjects {
			if key.Namespace == object.GetNamespace() {
				objectLabels := labels.Set(object.GetLabels())
				if selector.Matches(objectLabels) {
					filteredObjects = append(filteredObjects, object)
				}
			}
		}

		return filteredObjects, nil
	}

	updateCh := make(chan watchEvent)
	deleteCh := make(chan watchEvent)

	go w.handleUpdates(key, updateCh, deleteCh)

	objects, err := w.backendObjectStore.List(ctx, key)
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		w.cachedObjects.update(key.Namespace, gvk, object)
	}

	if err := w.createEventHandler(ctx, key, updateCh, deleteCh); err != nil {
		return nil, errors.Wrap(err, "create event handler")
	}

	w.flagGVKWatched(key, gvk)

	return objects, nil
}

// Get gets an object using a key.
func (w *Watch) Get(ctx context.Context, key store.Key) (*unstructured.Unstructured, error) {
	ctx, span := trace.StartSpan(ctx, "watchCacheGet")
	defer span.End()

	if w.backendObjectStore == nil {
		return nil, errors.New("backend cached is nil")
	}

	logger := log.From(ctx)
	if err := w.backendObjectStore.HasAccess(ctx, key, "get"); err != nil {
		logger.Errorf("check access failed: %v", err)
		u := unstructured.Unstructured{}
		return &u, nil
	}

	gvk := key.GroupVersionKind()

	if w.isKeyCached(key) {
		cachedObjects := w.cachedObjects.list(key.Namespace, gvk)

		for _, object := range cachedObjects {
			if key.Namespace == object.GetNamespace() &&
				key.Name == object.GetName() {
				return object, nil
			}
		}

		// TODO: handle not found case
		return nil, nil
	}

	updateCh := make(chan watchEvent)
	deleteCh := make(chan watchEvent)

	go w.handleUpdates(key, updateCh, deleteCh)

	object, err := w.backendObjectStore.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	w.cachedObjects.update(key.Namespace, gvk, object)

	if err := w.createEventHandler(ctx, key, updateCh, deleteCh); err != nil {
		return nil, errors.Wrap(err, "create event handler")
	}

	w.flagGVKWatched(key, gvk)

	return object, nil
}

// Watch watches the cluster given a key and a handler.
func (w *Watch) Watch(ctx context.Context, key store.Key, handler kcache.ResourceEventHandler) error {
	if w.backendObjectStore == nil {
		return errors.New("backend object store is nil")
	}
	return w.backendObjectStore.Watch(ctx, key, handler)
}

func (w *Watch) isKeyCached(key store.Key) bool {
	return w.watchedGVKs.isWatched(key.Namespace, key.GroupVersionKind())
}

func (w *Watch) handleUpdates(key store.Key, updateCh, deleteCh chan watchEvent) {
	defer close(updateCh)
	defer close(deleteCh)

	done := false
	for !done {
		select {
		case <-w.stopCh:
			done = true
		case event := <-updateCh:
			w.cachedObjects.update(key.Namespace, event.gvk, event.object)
		case event := <-deleteCh:
			w.cachedObjects.delete(key.Namespace, event.gvk, event.object)
		}
	}
}

func (w *Watch) createEventHandler(ctx context.Context, key store.Key, updateCh, deleteCh chan watchEvent) error {
	handler := &watchEventHandler{
		gvk: key.GroupVersionKind(),
		updateFunc: func(event watchEvent) {
			if event.object == nil {
				return
			}

			updateCh <- event
		},
		deleteFunc: func(event watchEvent) {
			if event.object == nil {
				return
			}

			deleteCh <- event
		},
	}

	if w.client == nil {
		return errors.New("cluster client is nil")
	}
	gvk := key.GroupVersionKind()
	gvr, err := w.client.Resource(gvk.GroupKind())
	if err != nil {
		return errors.Wrap(err, "client resource")
	}

	factory, ok := w.factories.get(key.Namespace)

	if !ok {
		if err := w.HasAccess(ctx, store.Key{Namespace: metav1.NamespaceAll}, "watch"); err != nil {
			factory, err = w.initFactoryFunc(ctx, w.client, key.Namespace)
			if err != nil {
				return err
			}
		} else {
			factory, ok = w.factories.get("")
			if !ok {
				return errors.New("no default DynamicInformerFactory found")
			}
		}
	}

	w.factories.set(key.Namespace, factory)

	informer, err := currentInformer(gvr, factory, w.stopCh)
	if err != nil {
		return errors.Wrapf(err, "find informer for key %s", key)
	}

	informer.Informer().AddEventHandler(handler)

	return nil
}

func (w *Watch) flagGVKWatched(key store.Key, gvk schema.GroupVersionKind) {
	w.watchedGVKs.setWatched(key.Namespace, gvk)
}

type watchEvent struct {
	object *unstructured.Unstructured
	gvk    schema.GroupVersionKind
}

type watchEventHandler struct {
	gvk        schema.GroupVersionKind
	updateFunc func(event watchEvent)
	deleteFunc func(event watchEvent)
}

var _ kcache.ResourceEventHandler = (*watchEventHandler)(nil)

func (h *watchEventHandler) OnAdd(obj interface{}) {
	if object, ok := obj.(*unstructured.Unstructured); ok {
		event := watchEvent{object: object, gvk: h.gvk}
		h.updateFunc(event)
	}
}

func (h *watchEventHandler) OnUpdate(oldObj, newObj interface{}) {
	if object, ok := newObj.(*unstructured.Unstructured); ok {
		event := watchEvent{object: object, gvk: h.gvk}
		h.updateFunc(event)
	}

}

func (h *watchEventHandler) OnDelete(obj interface{}) {
	if object, ok := obj.(*unstructured.Unstructured); ok {
		event := watchEvent{object: object, gvk: h.gvk}
		h.deleteFunc(event)
	}
}

var nsGVK = schema.GroupVersionKind{Version: "v1", Kind: "Namespace"}

type nsUpdateHandler struct {
	watch  *Watch
	logger log.Logger
}

var _ kcache.ResourceEventHandler = (*nsUpdateHandler)(nil)

func (h *nsUpdateHandler) OnAdd(obj interface{}) {
	if h.watch.initFactoryFunc == nil {
		return
	}

	if object, ok := obj.(*unstructured.Unstructured); ok && object.GroupVersionKind().String() == nsGVK.String() {
		factory, err := h.watch.initFactoryFunc(context.Background(), h.watch.client, object.GetName())
		if err != nil {
			h.logger.WithErr(err).Errorf("create namespace factory")
			return
		}

		h.logger.With("namespace", object.GetName()).Debugf("adding factory for namespace")
		h.watch.factories.set(object.GetName(), factory)
	}
}

func (h *nsUpdateHandler) OnUpdate(oldObj, newObj interface{}) {
}

func (h *nsUpdateHandler) OnDelete(obj interface{}) {
	if h.watch.initFactoryFunc == nil {
		return
	}

	if object, ok := obj.(*unstructured.Unstructured); ok && object.GroupVersionKind().String() == nsGVK.String() {
		h.watch.factories.delete(object.GetName())
		h.logger.With("namespace", object.GetName()).Debugf("removed factory for namespace")
	}
}

// UpdateClusterClient updates the cluster client.
func (w *Watch) UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error {
	logger := log.From(ctx)
	logger.Debugf("watch is updating its cluster client")

	w.cancelFunc()

	w.client = client
	if err := w.bootstrap(ctx, true); err != nil {
		return err
	}

	for _, fn := range w.updateFns {
		fn(w)
	}

	w.onClientUpdate <- w

	return nil
}

func (w *Watch) RegisterOnUpdate(fn store.UpdateFn) {
	w.updateFns = append(w.updateFns, fn)
}

// Update defers the update to the backend store.
func (w *Watch) Update(ctx context.Context, key store.Key, updater func(*unstructured.Unstructured) error) error {
	return w.backendObjectStore.Update(ctx, key, updater)
}

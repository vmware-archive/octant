/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstore

import (
	"context"
	"sync"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	kcache "k8s.io/client-go/tools/cache"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/pkg/store"
	"github.com/heptio/developer-dash/third_party/k8s.io/client-go/dynamic/dynamicinformer"
)

// WatchOpt is an option for configuration Watch.
type WatchOpt func(*Watch)

// Watch is a cache which watches the cluster for updates to known objects. It wraps a dynamic cache
// by default. Since the cache knows about all cluster updates, a majority of operations for listing
// and getting objects can happen in local memory instead of requiring a network request.
type Watch struct {
	initFactoryFunc func(cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error)
	initBackendFunc func(watch *Watch) (store.Store, error)
	client          cluster.ClientInterface
	stopCh          <-chan struct{}
	cancelFunc      context.CancelFunc
	factories       map[string]dynamicinformer.DynamicSharedInformerFactory
	watchedGVKs     map[string]map[schema.GroupVersionKind]bool
	cachedObjects   map[string]map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured
	handlers        map[string]map[schema.GroupVersionKind]watchEventHandler

	backendObjectStore store.Store
	gvkLock            sync.Mutex
	objectLock         sync.RWMutex

	onClientUpdate chan store.Store
	updateFns      []store.UpdateFn
}

var _ store.Store = (*Watch)(nil)

func initWatchBackend(w *Watch) (store.Store, error) {
	backendObjectStore, err := NewDynamicCache(w.client, w.stopCh, func(d *DynamicCache) {
		d.initFactoryFunc = func(client cluster.ClientInterface, namespace string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return w.factories[namespace], nil
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
		factories:       make(map[string]dynamicinformer.DynamicSharedInformerFactory),
		watchedGVKs:     make(map[string]map[schema.GroupVersionKind]bool),
		cachedObjects:   make(map[string]map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured),
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
		w.factories = make(map[string]dynamicinformer.DynamicSharedInformerFactory)
		w.watchedGVKs = make(map[string]map[schema.GroupVersionKind]bool)
		w.cachedObjects = make(map[string]map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured)
		w.handlers = make(map[string]map[schema.GroupVersionKind]watchEventHandler)
	}

	ctx, cancel := context.WithCancel(ctx)
	w.cancelFunc = cancel
	w.stopCh = ctx.Done()

	namespaceClient, err := w.client.NamespaceClient()
	if err != nil {
		return errors.Wrap(err, "client namespace")
	}

	namespaces, err := namespaceClient.Names()
	if err != nil {
		namespaces = []string{namespaceClient.InitialNamespace()}
	}
	namespaces = append(namespaces, "")

	for _, namespace := range namespaces {
		factory, err := w.initFactoryFunc(w.client, namespace)
		if err != nil {
			return errors.Wrap(err, "initialize dynamic shared informer factory")
		}

		if _, ok := w.factories[namespace]; !ok {
			w.factories[namespace] = factory
		}
		if _, ok := w.watchedGVKs[namespace]; !ok {
			w.watchedGVKs[namespace] = make(map[schema.GroupVersionKind]bool)
		}
		if _, ok := w.cachedObjects[namespace]; !ok {
			w.cachedObjects[namespace] = make(map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured)
		}
		if _, ok := w.handlers[namespace]; !ok {
			w.handlers[namespace] = make(map[schema.GroupVersionKind]watchEventHandler)
		}
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
func (w *Watch) HasAccess(key store.Key, verb string) error {
	return w.backendObjectStore.HasAccess(key, verb)
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
	if err := w.backendObjectStore.HasAccess(key, "list"); err != nil {
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

		w.objectLock.RLock()
		defer w.objectLock.RUnlock()
		cachedObjects := w.cachedObjects[key.Namespace][gvk]
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

	w.objectLock.Lock()
	if _, ok := w.cachedObjects[key.Namespace]; !ok {
		w.cachedObjects[key.Namespace] = make(map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured)
	}
	w.cachedObjects[key.Namespace][gvk] = make(map[types.UID]*unstructured.Unstructured)
	for _, object := range objects {
		w.cachedObjects[key.Namespace][gvk][object.GetUID()] = object
	}
	w.objectLock.Unlock()

	if err := w.createEventHandler(key, updateCh, deleteCh); err != nil {
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
	if err := w.backendObjectStore.HasAccess(key, "get"); err != nil {
		logger.Errorf("check access failed: %v", err)
		u := unstructured.Unstructured{}
		return &u, nil
	}

	gvk := key.GroupVersionKind()

	if w.isKeyCached(key) {
		w.objectLock.RLock()
		defer w.objectLock.RUnlock()
		cachedObjects := w.cachedObjects[key.Namespace][gvk]
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

	w.objectLock.Lock()
	w.cachedObjects[key.Namespace][gvk] = make(map[types.UID]*unstructured.Unstructured)
	w.cachedObjects[key.Namespace][gvk][object.GetUID()] = object
	w.objectLock.Unlock()

	if err := w.createEventHandler(key, updateCh, deleteCh); err != nil {
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
	w.gvkLock.Lock()
	defer w.gvkLock.Unlock()

	gvk := key.GroupVersionKind()

	_, ok := w.watchedGVKs[key.Namespace][gvk]
	return ok
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
			w.objectLock.Lock()
			w.cachedObjects[key.Namespace][event.gvk][event.object.GetUID()] = event.object
			w.objectLock.Unlock()
		case event := <-deleteCh:
			w.objectLock.Lock()
			delete(w.cachedObjects[key.Namespace][event.gvk], event.object.GetUID())
			w.objectLock.Unlock()
		}
	}
}

func (w *Watch) createEventHandler(key store.Key, updateCh, deleteCh chan watchEvent) error {
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

	factory, ok := w.factories[key.Namespace]
	if !ok {
		return errors.Errorf("no informer factory found for %s", key.Namespace)
	}

	informer, err := currentInformer(gvr, factory, w.stopCh)
	if err != nil {
		return errors.Wrapf(err, "find informer for key %s", key)
	}

	informer.Informer().AddEventHandler(handler)

	return nil
}

func (w *Watch) flagGVKWatched(key store.Key, gvk schema.GroupVersionKind) {
	w.gvkLock.Lock()
	defer w.gvkLock.Unlock()
	if _, ok := w.watchedGVKs[key.Namespace]; !ok {
		w.watchedGVKs[key.Namespace] = make(map[schema.GroupVersionKind]bool)
	}
	w.watchedGVKs[key.Namespace][gvk] = true
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
		factory, err := h.watch.initFactoryFunc(h.watch.client, object.GetName())
		if err != nil {
			h.logger.WithErr(err).Errorf("create namespace factory")
			return
		}

		h.logger.With("namespace", object.GetName()).Debugf("adding factory for namespace")
		h.watch.factories[object.GetName()] = factory
	}
}

func (h *nsUpdateHandler) OnUpdate(oldObj, newObj interface{}) {
}

func (h *nsUpdateHandler) OnDelete(obj interface{}) {
	if h.watch.initFactoryFunc == nil {
		return
	}

	if object, ok := obj.(*unstructured.Unstructured); ok && object.GroupVersionKind().String() == nsGVK.String() {
		delete(h.watch.factories, object.GetName())
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

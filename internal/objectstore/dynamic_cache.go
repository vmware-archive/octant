/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstore

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opencensus.io/trace"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	kcache "k8s.io/client-go/tools/cache"
	kretry "k8s.io/client-go/util/retry"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/store"
)

const (
	// defaultMutableResync is the resync period for informers.
	defaultInformerResync = time.Second * 180

	// initialInformerSyncTimeout
	initialInformerSyncTimeout = time.Second * 10
)

func initInformerFactory(ctx context.Context, client cluster.ClientInterface, namespace string) (InformerFactory, error) {
	return newInformerFactory(ctx.Done(), client, defaultInformerResync, namespace), nil
}

// DynamicCacheOpt is an option for configuration DynamicCache.
type DynamicCacheOpt func(*DynamicCache)

// Access sets the Resource Access cache for a DynamicCache.
func Access(resourceAccess ResourceAccess) DynamicCacheOpt {
	return func(dc *DynamicCache) {
		dc.setResourceAccess(resourceAccess)
	}
}

// DynamicCache is a cache based on the dynamic shared informer factory.
type DynamicCache struct {
	initFactoryFunc func(context.Context, cluster.ClientInterface, string) (InformerFactory, error)
	factories       *factoriesCache
	informerSynced  *informerSynced
	backoffMap      sync.Map
	client          cluster.ClientInterface
	seenGVKs        *seenGVKsCache
	access          ResourceAccess
	updateFns       []store.UpdateFn
	updateMu        sync.Mutex

	syncTimeoutFunc func(context.Context, store.Key, chan bool)
	waitForSyncFunc func(context.Context, store.Key, *DynamicCache, informers.GenericInformer, chan bool)
}

func syncTimeout(ctx context.Context, key store.Key, done chan bool) {
	logger := log.From(ctx).With("key", key)
	timer := time.NewTimer(initialInformerSyncTimeout)
	select {
	case <-timer.C:
		logger.Debugf("cache has taken more than %s seconds to sync", initialInformerSyncTimeout)
	case <-done:
		timer.Stop()
	}
}

func waitForSync(ctx context.Context, key store.Key, dc *DynamicCache, informer informers.GenericInformer, done chan bool) {
	now := time.Now()
	logger := log.From(ctx).With("key", key)
	msg := "informer cache has synced"
	kcache.WaitForCacheSync(ctx.Done(), informer.Informer().HasSynced)
	<-time.After(100 * time.Millisecond)
	logger.With("elapsed", time.Since(now)).
		Debugf(msg)
	dc.informerSynced.setSynced(key, true)
	done <- true
}

var _ store.Store = (*DynamicCache)(nil)

// NewDynamicCache creates an instance of DynamicCache.
func NewDynamicCache(ctx context.Context, client cluster.ClientInterface, options ...DynamicCacheOpt) (*DynamicCache, error) {
	c := &DynamicCache{
		initFactoryFunc: initInformerFactory,
		syncTimeoutFunc: syncTimeout,
		waitForSyncFunc: waitForSync,
		client:          client,
		seenGVKs:        initSeenGVKsCache(),
		informerSynced:  initInformerSynced(),
	}

	for _, option := range options {
		option(c)
	}

	logger := log.From(ctx).With("component", "DynamicCache")

	c.factories = initFactoriesCache()
	go initStatusCheck(ctx.Done(), logger, c.factories)

	factory, err := c.initFactoryFunc(context.Background(), client, "")
	if err != nil {
		return nil, fmt.Errorf("initialize dynamic shared informer factory: %w", err)
	}

	c.factories.set("", factory)

	return c, nil
}

type lister interface {
	List(selector kLabels.Selector) ([]kruntime.Object, error)
}

func (dc *DynamicCache) setResourceAccess(resourceAccess ResourceAccess) {
	dc.access = resourceAccess
}

func (dc *DynamicCache) currentInformer(ctx context.Context, key store.Key) (informers.GenericInformer, bool, error) {
	if dc.client == nil {
		return nil, false, errors.New("cluster client is nil")
	}

	gvk := key.GroupVersionKind()
	factory, ok := dc.factories.get(key.Namespace)
	if !ok {
		if err := dc.access.HasAccess(ctx, store.Key{Namespace: metav1.NamespaceAll}, "watch"); err != nil {
			factory, err = dc.initFactoryFunc(ctx, dc.client, key.Namespace)
			if err != nil {
				return nil, false, fmt.Errorf("check access watch all namespaces: %w", err)
			}
		} else {
			factory, ok = dc.factories.get("")
			if !ok {
				return nil, false, errors.New("no default DynamicInformerFactory found")
			}
		}

		dc.factories.set(key.Namespace, factory)
	}

	informer, err := factory.ForResource(gvk)
	if err != nil {
		return nil, false, fmt.Errorf("find informer for %s: %w", gvk, err)
	}

	dc.checkKeySynced(ctx, informer, key)
	dc.seenGVKs.setSeen(key.Namespace, gvk, true)

	return informer, dc.informerSynced.hasSynced(key), nil
}

func (dc *DynamicCache) checkKeySynced(ctx context.Context, informer informers.GenericInformer, key store.Key) {
	dc.updateMu.Lock()
	defer dc.updateMu.Unlock()

	if dc.seenGVKs.hasSeen(key.Namespace, key.GroupVersionKind()) ||
		(dc.informerSynced.hasSeen(key) && dc.informerSynced.hasSynced(key)) {
		return
	}

	done := make(chan bool, 1)
	go dc.waitForSyncFunc(ctx, key, dc, informer, done)
	go dc.syncTimeoutFunc(ctx, key, done)
}

// backoff returns the Duration of the next backoff. A duration of 0 is returned if
// backoff is already waiting.
func (dc *DynamicCache) backoff(ctx context.Context, key store.Key) time.Duration {
	newEntry := newBackoffEntry(key, defaultBackoff)
	actual, _ := dc.backoffMap.LoadOrStore(key.GroupVersionKind(), newEntry)
	entry := actual.(backoffer)
	if entry.isWaiting() {
		return 0
	}

	entry.setWaiting(true)
	dc.backoffMap.Store(key.GroupVersionKind(), entry)
	t := entry.wait()

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			logger := log.From(ctx)
			defer func() {
				entry.setWaiting(false)
				dc.backoffMap.Store(key.GroupVersionKind(), entry)
			}()
			unwatchErr := dc.Unwatch(ctx, key.GroupVersionKind())
			if unwatchErr != nil {
				logger.Errorf("unwatch: %w", unwatchErr)
			}
			logger.Infof("backing off for %s (%s)", key.GroupVersionKind(), t)
			<-time.After(t)
		}
	}()

	return t
}

func (dc *DynamicCache) isBackingOff(ctx context.Context, key store.Key) bool {
	actual, ok := dc.backoffMap.Load(key.GroupVersionKind())
	if !ok {
		return false
	}
	entry := actual.(backoffer)
	return entry.isWaiting()
}

// List lists objects.
func (dc *DynamicCache) List(ctx context.Context, key store.Key) (*unstructured.UnstructuredList, bool, error) {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:list")
	defer span.End()

	if dc.isBackingOff(ctx, key) {
		return &unstructured.UnstructuredList{}, false, nil
	}

	if err := dc.access.HasAccess(ctx, key, "list"); err != nil {
		if meta.IsNoMatchError(err) {
			return &unstructured.UnstructuredList{}, false, nil
		}
		if !dc.isBackingOff(ctx, key) {
			dc.backoff(ctx, key)
		}
		return nil, false, fmt.Errorf("check access to list %s: %w", key, err)
	}

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("namespace", key.Namespace),
		trace.StringAttribute("apiVersion", key.APIVersion),
		trace.StringAttribute("kind", key.Kind),
	}, "list key")

	return dc.listFromInformer(ctx, key)
}

func (dc *DynamicCache) listFromInformer(ctx context.Context, key store.Key) (*unstructured.UnstructuredList, bool, error) {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:list:informer")
	defer span.End()

	informer, hasSynced, err := dc.currentInformer(ctx, key)
	if err != nil {
		return nil, false, fmt.Errorf("retrieving informer for %+v: %w", key, err)
	}

	if !hasSynced {
		list, err := dc.listFromDynamicClient(ctx, key)
		return list, false, err
	}

	var l lister
	if key.Namespace == "" {
		l = informer.Lister()
	} else {
		l = informer.Lister().ByNamespace(key.Namespace)
	}

	var selector = kLabels.Everything()
	if key.Selector != nil {
		selector = key.Selector.AsSelector()
	}

	objects, err := l.List(selector)
	if err != nil {
		return nil, false, fmt.Errorf("listing %v: %w", key, err)
	}

	list := &unstructured.UnstructuredList{}
	for i := range objects {
		list.Items = append(list.Items, *objects[i].(*unstructured.Unstructured))
	}

	return list, !dc.informerSynced.hasSynced(key), nil
}

func (dc *DynamicCache) listFromDynamicClient(ctx context.Context, key store.Key) (*unstructured.UnstructuredList, error) {
	_, span := trace.StartSpan(ctx, "dynamicCache:list:informer")
	defer span.End()

	var selector = kLabels.Everything()
	if key.Selector != nil {
		selector = key.Selector.AsSelector()
	}

	dynamicClient, err := dc.client.DynamicClient()
	if err != nil {
		return nil, err
	}

	gvr, _, err := dc.client.Resource(key.GroupVersionKind().GroupKind())
	if err != nil {
		return nil, err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	if key.Namespace == "" {
		return dynamicClient.Resource(gvr).List(ctx, listOptions)
	}

	return dynamicClient.Resource(gvr).Namespace(key.Namespace).List(ctx, listOptions)
}

type getter interface {
	Get(string) (kruntime.Object, error)
}

// Get retrieves a single object.
func (dc *DynamicCache) Get(ctx context.Context, key store.Key) (*unstructured.Unstructured, error) {
	ctx, span := trace.StartSpan(ctx, "dynamicCacheGet")
	defer span.End()

	if dc.isBackingOff(ctx, key) {
		return &unstructured.Unstructured{}, nil
	}

	if err := dc.access.HasAccess(ctx, key, "get"); err != nil {
		if meta.IsNoMatchError(err) {
			return &unstructured.Unstructured{}, nil
		}
		if !dc.isBackingOff(ctx, key) {
			dc.backoff(ctx, key)
		}
		return nil, fmt.Errorf("check access for get to %s: %w", key, err)
	}

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("namespace", key.Namespace),
		trace.StringAttribute("apiVersion", key.APIVersion),
		trace.StringAttribute("kind", key.Kind),
		trace.StringAttribute("name", key.Name),
	}, "get key")

	object, err := dc.getFromInformer(ctx, key)
	if err != nil {
		if kerrors.IsNotFound(err) {
			object, err := dc.getFromDynamicClient(ctx, key)
			if err != nil {
				return nil, err
			}
			return object, nil
		}

		return nil, err
	}

	return object, nil

}

func (dc *DynamicCache) getFromInformer(ctx context.Context, key store.Key) (*unstructured.Unstructured, error) {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:get:informer")
	defer span.End()

	informer, hasSynced, err := dc.currentInformer(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("retrieving informer for %v: %w", key, err)
	}

	if !hasSynced {
		return dc.getFromDynamicClient(ctx, key)
	}

	var g getter
	if key.Namespace == "" {
		g = informer.Lister()
	} else {
		g = informer.Lister().ByNamespace(key.Namespace)
	}

	object, err := g.Get(key.Name)
	if err != nil {
		return nil, err
	}
	return object.(*unstructured.Unstructured), nil
}

func (dc *DynamicCache) getFromDynamicClient(ctx context.Context, key store.Key) (*unstructured.Unstructured, error) {
	_, span := trace.StartSpan(ctx, "dynamicCache:get:dynamicClient")
	defer span.End()

	dynamicClient, err := dc.client.DynamicClient()
	if err != nil {
		return nil, err
	}

	gvr, _, err := dc.client.Resource(key.GroupVersionKind().GroupKind())
	if err != nil {
		return nil, err
	}

	if key.Namespace == "" {
		return dynamicClient.Resource(gvr).Get(ctx, key.Name, metav1.GetOptions{})
	}
	return dynamicClient.Resource(gvr).Namespace(key.Namespace).Get(ctx, key.Name, metav1.GetOptions{})
}

// Watch watches the cluster for an event and performs actions with the
// supplied handler.
func (dc *DynamicCache) Watch(ctx context.Context, key store.Key, handler kcache.ResourceEventHandler) error {
	if dc.isBackingOff(ctx, key) {
		return nil
	}

	if err := dc.access.HasAccess(ctx, key, "watch"); err != nil {
		if meta.IsNoMatchError(err) {
			return nil
		}
		if !dc.isBackingOff(ctx, key) {
			dc.backoff(ctx, key)
		}
		return fmt.Errorf("check access to watch %s: %w", key, err)
	}

	informer, _, err := dc.currentInformer(ctx, key)
	if err != nil {
		return fmt.Errorf("retrieving informer for %s: %w", key, err)
	}

	informer.Informer().AddEventHandler(handler)
	return nil
}

// Unwatch un-watches a key by stopping it's informer.
func (dc *DynamicCache) Unwatch(ctx context.Context, groupVersionKinds ...schema.GroupVersionKind) error {
	for _, namespace := range dc.factories.keys() {
		factory, ok := dc.factories.get(namespace)
		if ok {
			for _, groupVersionKind := range groupVersionKinds {
				factory.Delete(groupVersionKind)
			}
		}
	}

	return nil
}

// Delete deletes an object from the cluster using a key.
func (dc *DynamicCache) Delete(ctx context.Context, key store.Key) error {
	_, span := trace.StartSpan(ctx, "dynamicCache:delete")
	defer span.End()

	if dc.isBackingOff(ctx, key) {
		return nil
	}

	if err := dc.access.HasAccess(ctx, key, "delete"); err != nil {
		if meta.IsNoMatchError(err) {
			return nil
		}
		if !dc.isBackingOff(ctx, key) {
			dc.backoff(ctx, key)
		}
		return fmt.Errorf("check access to delete %s: %w", key, err)
	}

	dynamicClient, err := dc.client.DynamicClient()
	if err != nil {
		return err
	}

	gvr, _, err := dc.client.Resource(key.GroupVersionKind().GroupKind())
	if err != nil {
		return err
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	if key.Namespace == "" {
		return dynamicClient.Resource(gvr).Delete(ctx, key.Name, deleteOptions)
	}

	return dynamicClient.Resource(gvr).Namespace(key.Namespace).Delete(ctx, key.Name, deleteOptions)
}

// UpdateClusterClient updates the cluster client.
func (dc *DynamicCache) UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error {
	logger := log.From(ctx)
	logger.Debugf("updating its cluster client")

	dc.updateMu.Lock()
	dc.client = client
	dc.factories.reset()
	dc.seenGVKs.reset()
	dc.informerSynced.reset()
	dc.access = NewResourceAccess(client)
	dc.updateMu.Unlock()

	for _, fn := range dc.updateFns {
		fn(dc)
	}

	return nil
}

// RegisterOnUpdate registers a function that will be called when the store updates it's client.
// TODO: investigate if this needed since object store isn't replaced, it's client is. (GH#496)
func (dc *DynamicCache) RegisterOnUpdate(fn store.UpdateFn) {
	dc.updateFns = append(dc.updateFns, fn)
}

func (dc *DynamicCache) Update(ctx context.Context, key store.Key, updater func(*unstructured.Unstructured) error) error {
	if updater == nil {
		return errors.New("can't update object")
	}

	if dc.isBackingOff(ctx, key) {
		return nil
	}

	if err := dc.access.HasAccess(ctx, key, "update"); err != nil {
		if meta.IsNoMatchError(err) {
			return nil
		}
		if !dc.isBackingOff(ctx, key) {
			dc.backoff(ctx, key)
		}
		return fmt.Errorf("check access to update %s: %w", key, err)
	}

	err := kretry.RetryOnConflict(kretry.DefaultRetry, func() error {
		object, err := dc.Get(ctx, key)
		if err != nil {
			return err
		}

		if object == nil {
			return errors.New("object not found")
		}

		gvk := object.GroupVersionKind()

		gvr, _, err := dc.client.Resource(gvk.GroupKind())
		if err != nil {
			return err
		}

		dynamicClient, err := dc.client.DynamicClient()
		if err != nil {
			return err
		}

		if err := updater(object); err != nil {
			return fmt.Errorf("unable to update object: %w", err)
		}

		client := dynamicClient.Resource(gvr).Namespace(object.GetNamespace())

		_, err = client.Update(ctx, object, metav1.UpdateOptions{})
		return err
	})

	return err
}

func (dc *DynamicCache) IsLoading(ctx context.Context, key store.Key) bool {
	return !dc.informerSynced.hasSynced(key)
}

// Create creates an object in the cluster.
// Note: test coverage of DynamicCache is slim.
func (dc *DynamicCache) Create(ctx context.Context, object *unstructured.Unstructured) error {
	_, span := trace.StartSpan(ctx, "dynamicCache:create")
	defer span.End()

	key, err := store.KeyFromObject(object)
	if err != nil {
		return fmt.Errorf("key from object: %w", err)
	}

	if dc.isBackingOff(ctx, key) {
		return nil
	}

	if err := dc.access.HasAccess(ctx, key, "create"); err != nil {
		if meta.IsNoMatchError(err) {
			return nil
		}
		if !dc.isBackingOff(ctx, key) {
			dc.backoff(ctx, key)
		}
		return fmt.Errorf("check access to create %s: %w", key, err)
	}

	dynamicClient, err := dc.client.DynamicClient()
	if err != nil {
		return err
	}

	gvr, _, err := dc.client.Resource(key.GroupVersionKind().GroupKind())
	if err != nil {
		return err
	}

	createOptions := metav1.CreateOptions{}

	if key.Namespace == "" {
		_, err := dynamicClient.Resource(gvr).Create(ctx, object, createOptions)
		return err
	}

	_, err = dynamicClient.Resource(gvr).Namespace(key.Namespace).Create(ctx, object, createOptions)
	return err
}

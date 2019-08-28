/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstore

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	authorizationv1 "k8s.io/api/authorization/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	kcache "k8s.io/client-go/tools/cache"
	kretry "k8s.io/client-go/util/retry"

	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/util/retry"
	"github.com/vmware/octant/pkg/store"
)

const (
	// defaultMutableResync is the resync period for informers.
	defaultInformerResync = time.Second * 180

	// initialInformerSyncTimeout
	initialInformerSyncTimeout = time.Second * 10
)

func initDynamicSharedInformerFactory(ctx context.Context, client cluster.ClientInterface, namespace string) (dynamicinformer.DynamicSharedInformerFactory, error) {
	dynamicClient, err := client.DynamicClient()
	if err != nil {
		return nil, err
	}
	return dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, defaultInformerResync, namespace, nil), nil
}

type accessKey struct {
	Namespace string
	Group     string
	Resource  string
	Verb      string
}
type accessMap map[accessKey]bool

// DynamicCacheOpt is an option for configuration DynamicCache.
type DynamicCacheOpt func(*DynamicCache)

// DynamicCache is a cache based on the dynamic shared informer factory.
type DynamicCache struct {
	initFactoryFunc func(context.Context, cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error)
	factories       *factoriesCache
	informerSynced  *informerSynced
	client          cluster.ClientInterface
	stopCh          <-chan struct{}
	seenGVKs        *seenGVKsCache
	access          *accessCache
	updateFns       []store.UpdateFn
	updateMu        sync.Mutex

	syncTimeoutFunc func(context.Context, store.Key, chan bool)
	waitForSyncFunc func(context.Context, store.Key, *DynamicCache, informers.GenericInformer, <-chan struct{}, chan bool)

	useDynamicClient bool
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

func waitForSync(ctx context.Context, key store.Key, dc *DynamicCache, informer informers.GenericInformer, stopCh <-chan struct{}, done chan bool) {
	now := time.Now()
	logger := log.From(ctx).With("key", key)
	kcache.WaitForCacheSync(stopCh, informer.Informer().HasSynced)
	dc.informerSynced.setSynced(key, true)
	logger.With("elapsed", time.Since(now)).
		Debugf("informer cache has synced")
	done <- true
}

var _ store.Store = (*DynamicCache)(nil)

// NewDynamicCache creates an instance of DynamicCache.
func NewDynamicCache(ctx context.Context, client cluster.ClientInterface, options ...DynamicCacheOpt) (*DynamicCache, error) {
	c := &DynamicCache{
		initFactoryFunc:  initDynamicSharedInformerFactory,
		syncTimeoutFunc:  syncTimeout,
		waitForSyncFunc:  waitForSync,
		client:           client,
		stopCh:           ctx.Done(),
		seenGVKs:         initSeenGVKsCache(),
		informerSynced:   initInformerSynced(),
		useDynamicClient: os.Getenv("OCTANT_USE_DYNAMIC_CLIENT") == "1",
	}

	for _, option := range options {
		option(c)
	}

	if c.access == nil {
		c.access = initAccessCache()
	}

	logger := log.From(ctx).With("component", "DynamicCache")

	c.factories = initFactoriesCache()
	go initStatusCheck(ctx.Done(), logger, c.factories)

	if c.useDynamicClient {
		logger.Debugf("using dynamic client instead of informer")
	}

	factory, err := c.initFactoryFunc(context.Background(), client, "")
	if err != nil {
		return nil, errors.Wrap(err, "initialize dynamic shared informer factory")
	}

	c.factories.set("", factory)

	return c, nil
}

type lister interface {
	List(selector kLabels.Selector) ([]kruntime.Object, error)
}

func (dc *DynamicCache) fetchAccess(key accessKey, verb string) (bool, error) {
	k8sClient, err := dc.client.KubernetesClient()
	if err != nil {
		return false, errors.Wrap(err, "client kubernetes")
	}

	authClient := k8sClient.AuthorizationV1()
	sar := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace: key.Namespace,
				Group:     key.Group,
				Resource:  key.Resource,
				Verb:      verb,
			},
		},
	}

	review, err := authClient.SelfSubjectAccessReviews().Create(sar)
	if err != nil {
		return false, errors.Wrap(err, "client auth")
	}
	return review.Status.Allowed, nil
}

// HasAccess returns an error if the current user does not have access to perform the verb action
// for the given key.
func (dc *DynamicCache) HasAccess(ctx context.Context, key store.Key, verb string) error {
	_, span := trace.StartSpan(ctx, "dynamicCacheHasAccess")
	defer span.End()

	gvk := key.GroupVersionKind()

	if gvk.GroupKind().Empty() {
		return errors.Errorf("unable to check access for key %s", key.String())
	}

	gvr, err := dc.client.Resource(gvk.GroupKind())
	if err != nil {
		return errors.Wrap(err, "client resource")
	}

	aKey := accessKey{
		Namespace: key.Namespace,
		Group:     gvr.Group,
		Resource:  gvr.Resource,
		Verb:      verb,
	}

	access, ok := dc.access.get(aKey)

	if !ok {
		span.Annotate([]trace.Attribute{}, "fetch access start")
		val, err := dc.fetchAccess(aKey, verb)
		if err != nil {
			return errors.Wrapf(err, "fetch access: %+v", aKey)
		}

		dc.access.set(aKey, val)
		access = val
		span.Annotate([]trace.Attribute{}, "fetch access finish")
	}

	if !access {
		return errors.Errorf("denied %+v", aKey)
	}

	return nil
}

func (dc *DynamicCache) currentInformer(ctx context.Context, key store.Key) (informers.GenericInformer, error) {
	if dc.client == nil {
		return nil, errors.New("cluster client is nil")
	}

	gvk := key.GroupVersionKind()
	gvr, err := dc.client.Resource(gvk.GroupKind())
	if err != nil {
		return nil, errors.Wrap(err, "client resource")
	}

	factory, ok := dc.factories.get(key.Namespace)
	if !ok {
		if err := dc.HasAccess(ctx, store.Key{Namespace: metav1.NamespaceAll}, "watch"); err != nil {
			factory, err = dc.initFactoryFunc(ctx, dc.client, key.Namespace)
			if err != nil {
				return nil, err
			}
		} else {
			factory, ok = dc.factories.get("")
			if !ok {
				return nil, errors.New("no default DynamicInformerFactory found")
			}
		}

		dc.factories.set(key.Namespace, factory)
	}

	informer := factory.ForResource(gvr)
	factory.Start(dc.stopCh)

	dc.updateMu.Lock()
	dc.checkKeySynced(ctx, dc.stopCh, informer, key)
	dc.updateMu.Unlock()

	if dc.seenGVKs.hasSeen(key.Namespace, gvk) {
		return informer, nil
	}

	dc.seenGVKs.setSeen(key.Namespace, gvk, true)

	return informer, nil
}

func (dc *DynamicCache) checkKeySynced(ctx context.Context, stopCh <-chan struct{}, informer informers.GenericInformer, key store.Key) {
	if dc.seenGVKs.hasSeen(key.Namespace, key.GroupVersionKind()) ||
		(dc.informerSynced.hasSeen(key) && dc.informerSynced.hasSynced(key)) {
		return
	}

	done := make(chan bool, 1)
	go dc.waitForSyncFunc(ctx, key, dc, informer, stopCh, done)
	go dc.syncTimeoutFunc(ctx, key, done)
}

// List lists objects.
func (dc *DynamicCache) List(ctx context.Context, key store.Key) (*unstructured.UnstructuredList, bool, error) {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:list")
	defer span.End()

	if err := dc.HasAccess(ctx, key, "list"); err != nil {
		if meta.IsNoMatchError(err) {
			return &unstructured.UnstructuredList{}, false, nil
		}
		return nil, false, errors.Wrapf(err, "list access forbidden to %+v", key)
	}

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("namespace", key.Namespace),
		trace.StringAttribute("apiVersion", key.APIVersion),
		trace.StringAttribute("kind", key.Kind),
	}, "list key")

	if dc.useDynamicClient {
		list, err := dc.listFromDynamicClient(ctx, key)
		return list, true, err
	}

	return dc.listFromInformer(ctx, key)
}

func (dc *DynamicCache) listFromInformer(ctx context.Context, key store.Key) (*unstructured.UnstructuredList, bool, error) {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:list:informer")
	defer span.End()

	informer, err := dc.currentInformer(ctx, key)
	if err != nil {
		return nil, false, errors.Wrapf(err, "retrieving informer for %+v", key)
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
		return nil, false, errors.Wrapf(err, "listing %v", key)
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

	gvr, err := dc.client.Resource(key.GroupVersionKind().GroupKind())
	if err != nil {
		return nil, err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	if key.Namespace == "" {
		return dynamicClient.Resource(gvr).List(listOptions)
	}

	return dynamicClient.Resource(gvr).Namespace(key.Namespace).List(listOptions)
}

type getter interface {
	Get(string) (kruntime.Object, error)
}

// Get retrieves a single object.
func (dc *DynamicCache) Get(ctx context.Context, key store.Key) (*unstructured.Unstructured, bool, error) {
	ctx, span := trace.StartSpan(ctx, "dynamicCacheGet")
	defer span.End()

	if err := dc.HasAccess(ctx, key, "get"); err != nil {
		return nil, false, errors.Wrapf(err, "get access forbidden to %+v", key)
	}

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("namespace", key.Namespace),
		trace.StringAttribute("apiVersion", key.APIVersion),
		trace.StringAttribute("kind", key.Kind),
		trace.StringAttribute("name", key.Name),
	}, "get key")

	var object *unstructured.Unstructured
	var err error

	if dc.useDynamicClient {
		object, err = dc.getFromDynamicClient(ctx, key)
	} else {
		object, err = dc.getFromInformer(ctx, key)
	}

	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return object, true, nil

}

func (dc *DynamicCache) getFromInformer(ctx context.Context, key store.Key) (*unstructured.Unstructured, error) {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:get:informer")
	defer span.End()

	informer, err := dc.currentInformer(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving informer for %v", key)
	}

	var g getter
	if key.Namespace == "" {
		g = informer.Lister()
	} else {
		g = informer.Lister().ByNamespace(key.Namespace)
	}

	var retryCount int64

	var object kruntime.Object
	retryErr := retry.Retry(3, time.Second, func() error {
		object, err = g.Get(key.Name)
		if err != nil {
			if !kerrors.IsNotFound(err) {
				retryCount++
				return retry.Stop(errors.Wrap(err, "lister Get"))
			}
			return err
		}

		return nil
	})

	if retryCount > 0 {
		span.Annotate([]trace.Attribute{
			trace.Int64Attribute("retryCount", retryCount),
		}, "get retried")
	}

	if retryErr != nil {
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

	gvr, err := dc.client.Resource(key.GroupVersionKind().GroupKind())
	if err != nil {
		return nil, err
	}

	if key.Namespace == "" {
		return dynamicClient.Resource(gvr).Get(key.Name, metav1.GetOptions{})
	}
	return dynamicClient.Resource(gvr).Namespace(key.Namespace).Get(key.Name, metav1.GetOptions{})
}

// Watch watches the cluster for an event and performs actions with the
// supplied handler.
func (dc *DynamicCache) Watch(ctx context.Context, key store.Key, handler kcache.ResourceEventHandler) error {
	logger := log.From(ctx)
	if err := dc.HasAccess(ctx, key, "watch"); err != nil {
		logger.Errorf("check access failed: %v, access forbidden to %+v", key)
		return nil
	}

	informer, err := dc.currentInformer(ctx, key)
	if err != nil {
		return errors.Wrapf(err, "retrieving informer for %s", key)
	}

	informer.Informer().AddEventHandler(handler)
	return nil
}

// UpdateClusterClient updates the cluster client.
func (dc *DynamicCache) UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error {
	logger := log.From(ctx)
	logger.Debugf("updating its cluster client")

	dc.updateMu.Lock()
	dc.client = client
	dc.factories.reset()
	dc.updateMu.Unlock()

	for _, fn := range dc.updateFns {
		fn(dc)
	}

	return nil
}

func (dc *DynamicCache) RegisterOnUpdate(fn store.UpdateFn) {
	dc.updateFns = append(dc.updateFns, fn)
}

func (dc *DynamicCache) Update(ctx context.Context, key store.Key, updater func(*unstructured.Unstructured) error) error {
	if updater == nil {
		return errors.New("can't update object")
	}

	err := kretry.RetryOnConflict(kretry.DefaultRetry, func() error {
		object, found, err := dc.Get(ctx, key)
		if err != nil {
			return err
		}

		if !found {
			return errors.Errorf("object not found")
		}

		gvk := object.GroupVersionKind()

		gvr, err := dc.client.Resource(gvk.GroupKind())
		if err != nil {
			return err
		}

		dynamicClient, err := dc.client.DynamicClient()
		if err != nil {
			return err
		}

		if err := updater(object); err != nil {
			return errors.Wrap(err, "unable to update object")
		}

		client := dynamicClient.Resource(gvr).Namespace(object.GetNamespace())

		_, err = client.Update(object, metav1.UpdateOptions{})
		return err
	})

	return err
}

func (dc *DynamicCache) IsLoading(ctx context.Context, key store.Key) bool {
	return !dc.informerSynced.hasSynced(key)
}

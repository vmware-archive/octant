package overview

import (
	"context"
	"sync"
	"time"

	"github.com/heptio/developer-dash/internal/log"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// InformerCacheOpt is an option for configuring memory cache.
type InformerCacheOpt func(*InformerCache)

// InformerCacheNotificationOpt sets a channel that will receive a notification
// every time cache performs an add/delete.
// The done channel can be used to cancel notifications that are blocked.
func InformerCacheNotificationOpt(ch chan<- CacheNotification, done <-chan struct{}) InformerCacheOpt {
	return func(c *InformerCache) {
		c.notifyCh = ch
		c.stopCh = done
	}
}

// InformerCacheLoggerOpt sets a logger for the cache
func InformerCacheLoggerOpt(logger log.Logger) InformerCacheOpt {
	return func(c *InformerCache) {
		c.logger = logger
	}
}

type informerKey struct {
	namespace string
	gvk       schema.GroupVersionKind
}

// InformerCache caches
type InformerCache struct {
	client     dynamic.Interface
	restMapper meta.RESTMapper
	informers  map[informerKey]informers.GenericInformer
	logger     log.Logger

	mu             sync.RWMutex
	internalNotify chan CacheNotification
	notifyCh       chan<- CacheNotification
	stopCh         <-chan struct{}
}

var _ Cache = (*InformerCache)(nil)

// NewInformerCache creates a new InformerCache.
func NewInformerCache(stopCh <-chan struct{}, client dynamic.Interface, restMapper meta.RESTMapper, opts ...InformerCacheOpt) *InformerCache {
	c := &InformerCache{
		client:     client,
		restMapper: restMapper,
		stopCh:     stopCh,
		informers:  make(map[informerKey]informers.GenericInformer),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.logger == nil {
		c.logger = log.NopLogger()
	}
	if c.notifyCh != nil {
		c.internalNotify = make(chan CacheNotification)
		go c.runNotifyHandler()
	}
	return c
}

// runNotifyHandler forwards notifications from informer goroutines to
// their ultimate destination, the notifyCh channel.
// Handles shutdown of that channel when stopCh is signalled.
// Run this as a goroutine.
func (c *InformerCache) runNotifyHandler() {
	if c.stopCh == nil || c.notifyCh == nil {
		return
	}

	c.logger.Debugf("notify handler started")
	for {
		// Split selects to since we cannot define priority between incoming notifications
		// and stopCh being signaled. This way we will always check stopCh between incoming notifications.
		select {
		case <-c.stopCh:
			c.logger.Debugf("notify handler stopped")
			close(c.notifyCh)
			c.notifyCh = nil
			return
		default:
		}
		select {
		case event := <-c.internalNotify:
			select {
			case c.notifyCh <- event:
			case <-c.stopCh:
			}
		case <-c.stopCh:
		}
	}
}

func (c *InformerCache) informerForKey(ck CacheKey) (informers.GenericInformer, error) {
	var namespace = ck.Namespace
	if namespace == "" {
		namespace = "default"
	}

	gvk := schema.FromAPIVersionAndKind(ck.APIVersion, ck.Kind)
	key := informerKey{
		namespace: namespace,
		gvk:       gvk,
	}

	{
		// Fastpath
		c.mu.RLock()
		informer, ok := c.informers[key]
		c.mu.RUnlock()
		if ok {
			return informer, nil
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if informer, ok := c.informers[key]; ok {
		return informer, nil
	}

	// Create a new informer here
	restMapping, err := c.restMapper.RESTMapping(gvk.GroupKind())
	if err != nil {
		return nil, errors.Wrapf(err, "mapping %v", gvk.String())
	}
	gvr := restMapping.Resource
	gi := dynamicinformer.NewFilteredDynamicInformer(c.client, gvr, namespace, 180*time.Second, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, nil)

	// Install handlers, start fetching resources
	informer := gi.Informer()
	c.installHandler(informer)
	go informer.Run(c.stopCh) // (as in dynamicSharedInformerFactory.Start())

	// Until upstream issue in the wait package is resolved, we *must* ensure that the
	// stopCh passed to WaitForCacheSync is closed to avoid leaking goroutines spawned within.
	// We create a new channel for this purpose, as we do not want to cancel our factories and watches.
	// See https://github.com/kubernetes/kubernetes/pull/71326
	ctx, cancel := channelContext(c.stopCh)
	defer cancel()

	// Block until cache is synced or context is closed via c.stopCh.
	// Note the context must be closed even after uninterrupted return
	// to ensure cleanup of resources.
	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		return nil, errors.New("shutdown requested")
	}
	c.informers[key] = gi

	return gi, nil
}

// keyForObject returns a CacheKey representing a runtime.Object
func keyForObject(obj interface{}) (CacheKey, error) {
	metaAcc, err := meta.Accessor(obj)
	if err != nil {
		return CacheKey{}, errors.Errorf("fetching metadata accessor: %v", err)
	}
	typeAcc, err := meta.TypeAccessor(obj)
	if err != nil {
		return CacheKey{}, errors.Errorf("fetching type accessor: %v", err)
	}
	return CacheKey{
		Namespace:  metaAcc.GetNamespace(),
		APIVersion: typeAcc.GetAPIVersion(),
		Name:       metaAcc.GetName(),
		Kind:       typeAcc.GetKind(),
	}, nil
}

func (c *InformerCache) sendNotification(obj interface{}, action CacheAction) error {
	if c.internalNotify == nil {
		return nil
	}

	cacheKey, err := keyForObject(obj)
	if err != nil {
		return errors.Wrapf(err, "creating cache key")
	}

	notification := CacheNotification{
		CacheKey: cacheKey,
		Action:   action,
	}

	// Send notification on via runNotifyHandler goroutine
	select {
	case c.internalNotify <- notification:
	case <-c.stopCh:
		c.logger.Debugf("notification channel closed")
	}
	return nil
}

// installHandler installs an event handler on the supplied informer
// The handler will forward cache notifications.
func (c *InformerCache) installHandler(informer cache.SharedInformer) {
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if err := c.sendNotification(obj, CacheStore); err != nil {
				c.logger.Errorf("sending notification: %v", err)
				return
			}
		},
		DeleteFunc: func(obj interface{}) {
			if err := c.sendNotification(obj, CacheDelete); err != nil {
				c.logger.Errorf("sending notification: %v", err)
				return
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if err := c.sendNotification(newObj, CacheUpdate); err != nil {
				c.logger.Errorf("sending notification: %v", err)
				return
			}
		},
	})
}

// channelContext returns a cancellation context that acts as the child of
// a parent channel - the context will close when the returned CancelFunc is
// called or when the parent channel is closed, whichever happens first.
// Note the caller is responsible for *always* calling CancelFunc, otherwise resources
// can be leaked.
func channelContext(parentCh <-chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-parentCh:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

// Retrieve retrieves an object or list of objects from the cluster via cache.
// Blocks if cache needs to be synced.
func (c *InformerCache) Retrieve(key CacheKey) ([]*unstructured.Unstructured, error) {
	if c.restMapper == nil {
		return nil, errors.New("missing RESTMapper")
	}
	if key.Kind == "" {
		return nil, errors.New("kind is required")
	}

	// NOTE: This blocks when setting up new informers
	// TODO: Pass context for timeout
	gi, err := c.informerForKey(key)
	if err != nil {
		return nil, err
	}

	// Handle list operation
	if key.Name == "" {
		// c.logger.With("key", key, "gvk", gvk, "resource", restMapping.Resource).Debugf("listing all objects")
		objs, err := gi.Lister().List(labels.Everything())
		if err != nil {
			return nil, errors.Wrapf(err, "listing")
		}

		ret := make([]*unstructured.Unstructured, len(objs))
		for i, obj := range objs {
			u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
			if err != nil {
				return nil, errors.Wrapf(err, "converting %T to unstructured", obj)
			}
			ret[i] = &unstructured.Unstructured{Object: u}
		}
		return ret, nil
	}

	// Handle get operation
	// c.logger.With("key", key, "gvk", gvk, "resource", restMapping.Resource).Debugf("getting single object: %v", key.Name)
	lister := gi.Lister().ByNamespace(key.Namespace)
	obj, err := lister.Get(key.Name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "fetching %v", key)
	}
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, errors.Wrapf(err, "converting %T to unstructured", obj)
	}
	return []*unstructured.Unstructured{
		&unstructured.Unstructured{Object: u},
	}, nil
}

// Store is not implemented
func (c *InformerCache) Store(obj *unstructured.Unstructured) error {
	return errors.New("not implemented: Store")
}

// Delete is not implemented
func (c *InformerCache) Delete(obj *unstructured.Unstructured) error {
	return errors.New("not implemented: Delete")
}

// Returns events related to the specified object.
// TODO consider reworking this to use EventExpansion.Search(), which
//      utilizes FieldSelectors (involvedObject.uid)
func (c *InformerCache) getEvents(u *unstructured.Unstructured) ([]*unstructured.Unstructured, error) {
	var events []*unstructured.Unstructured

	var eventKey = CacheKey{
		Namespace:  u.GetNamespace(),
		APIVersion: "v1",
		Kind:       "Event",
	}

	allEvents, err := c.Retrieve(eventKey)
	if err != nil {
		return nil, err
	}

	for _, obj := range allEvents {
		event := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, event)
		if err != nil {
			return nil, err
		}

		involvedObject := event.InvolvedObject
		if involvedObject.Namespace == u.GetNamespace() &&
			involvedObject.APIVersion == u.GetAPIVersion() &&
			involvedObject.Kind == u.GetKind() &&
			involvedObject.Name == u.GetName() {
			events = append(events, obj)
		}
	}
	return events, nil
}

// Events returns events for an object.
func (c *InformerCache) Events(u *unstructured.Unstructured) ([]*unstructured.Unstructured, error) {
	return c.getEvents(u)
}

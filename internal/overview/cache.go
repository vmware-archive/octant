package overview

import (
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Cache stores Kubernetes objects.
type Cache interface {
	Store(obj *unstructured.Unstructured) error
	Retrieve(key CacheKey) ([]*unstructured.Unstructured, error)
	Delete(obj *unstructured.Unstructured) error

	Events(obj *unstructured.Unstructured) ([]*unstructured.Unstructured, error)
}

// CacheKey is a key for the cache.
type CacheKey struct {
	Namespace  string
	APIVersion string
	Kind       string
	Name       string
}

// MemoryCacheOpt is an option for configuring memory cache.
type MemoryCacheOpt func(*MemoryCache)

// CacheAction is a cache action.
type CacheAction string

const (
	// CacheStore is a store action.
	CacheStore CacheAction = "store"
	// CacheDelete is a delete action.
	CacheDelete CacheAction = "delete"
)

// CacheNotification is a notifcation for a cache.
type CacheNotification struct {
	CacheKey CacheKey
	Action   CacheAction
}

// CacheNotificationOpt sets a channel that will receive a notification
// every time cache performs an add/delete.
// The done channel can be used to cancel notifications that are blocked.
func CacheNotificationOpt(ch chan<- CacheNotification, done <-chan struct{}) MemoryCacheOpt {
	return func(c *MemoryCache) {
		c.notifyCh = ch
		c.notifyDone = done
	}
}

// MemoryCache stores a cache of Kubernetes objects in memory.
type MemoryCache struct {
	store map[CacheKey]*unstructured.Unstructured

	mu         sync.Mutex
	notifyCh   chan<- CacheNotification
	notifyDone <-chan struct{}
}

var _ Cache = (*MemoryCache)(nil)

// NewMemoryCache creates an instance of MemoryCache.
func NewMemoryCache(opts ...MemoryCacheOpt) *MemoryCache {
	mc := &MemoryCache{
		store: make(map[CacheKey]*unstructured.Unstructured),
	}

	for _, opt := range opts {
		opt(mc)
	}

	return mc
}

// Reset resets the cache.
func (mc *MemoryCache) Reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	for k := range mc.store {
		delete(mc.store, k)
	}
}

// Store stores an object to the object.
func (mc *MemoryCache) Store(obj *unstructured.Unstructured) error {
	key := CacheKey{
		Namespace:  obj.GetNamespace(),
		APIVersion: obj.GetAPIVersion(),
		Kind:       obj.GetKind(),
		Name:       obj.GetName(),
	}

	mc.mu.Lock()
	mc.store[key] = obj
	mc.mu.Unlock()

	mc.notify(CacheStore, key)

	return nil
}

// Retrieve retrieves an object from the cache.
func (mc *MemoryCache) Retrieve(key CacheKey) ([]*unstructured.Unstructured, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	var objs []*unstructured.Unstructured

	for k, v := range mc.store {
		if k.Namespace != key.Namespace {
			continue
		}

		if key.APIVersion == "" {
			objs = append(objs, v)
			continue
		}

		if k.APIVersion == key.APIVersion {
			if key.Kind == "" {
				objs = append(objs, v)
				continue
			}

			if k.Kind == key.Kind {
				if key.Name == "" {
					objs = append(objs, v)
					continue
				}

				if k.Name == key.Name {
					objs = append(objs, v)
				}
			}
		}
	}

	return objs, nil
}

// Delete deletes an object from the cache.
func (mc *MemoryCache) Delete(obj *unstructured.Unstructured) error {
	namespace := obj.GetNamespace()
	apiVersion := obj.GetAPIVersion()
	kind := obj.GetKind()
	name := obj.GetName()

	key := CacheKey{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}

	mc.mu.Lock()
	delete(mc.store, key)
	mc.mu.Unlock()

	mc.notify(CacheDelete, key)

	return nil
}

// Events returns events for an object.
func (mc *MemoryCache) Events(u *unstructured.Unstructured) ([]*unstructured.Unstructured, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	var events []*unstructured.Unstructured

	for _, obj := range mc.store {

		if obj.GetAPIVersion() != "v1" && obj.GetKind() != "Event" {
			continue
		}

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

func (mc *MemoryCache) notify(action CacheAction, key CacheKey) {
	if mc.notifyCh == nil {
		return
	}

	select {
	case mc.notifyCh <- CacheNotification{Action: action, CacheKey: key}:
	case <-mc.notifyDone:
	}
}

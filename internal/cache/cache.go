package cache

import (
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Cache stores Kubernetes objects.
type Cache interface {
	Store(obj *unstructured.Unstructured) error
	Retrieve(key Key) ([]*unstructured.Unstructured, error)
	Delete(obj *unstructured.Unstructured) error

	Events(obj *unstructured.Unstructured) ([]*unstructured.Unstructured, error)
}

// Key is a key for the cache.
type Key struct {
	Namespace  string
	APIVersion string
	Kind       string
	Name       string
}

// MemoryCacheOpt is an option for configuring memory cache.
type MemoryCacheOpt func(*MemoryCache)

// Action is a cache action.
type Action string

const (
	// StoreAction is a store action.
	StoreAction Action = "store"
	// DeleteAction is a delete action.
	DeleteAction Action = "delete"
	// UpdateAction is an update action.
	UpdateAction Action = "update"
)

// Notification is a notification for a cache.
type Notification struct {
	CacheKey Key
	Action   Action
}

// NotificationOpt sets a channel that will receive a notification
// every time cache performs an add/delete.
// The done channel can be used to cancel notifications that are blocked.
func NotificationOpt(ch chan<- Notification, done <-chan struct{}) MemoryCacheOpt {
	return func(c *MemoryCache) {
		c.notifyCh = ch
		c.notifyDone = done
	}
}

// MemoryCache stores a cache of Kubernetes objects in memory.
type MemoryCache struct {
	store map[Key]*unstructured.Unstructured

	mu         sync.Mutex
	notifyCh   chan<- Notification
	notifyDone <-chan struct{}
}

var _ Cache = (*MemoryCache)(nil)

// NewMemoryCache creates an instance of MemoryCache.
func NewMemoryCache(opts ...MemoryCacheOpt) *MemoryCache {
	mc := &MemoryCache{
		store: make(map[Key]*unstructured.Unstructured),
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
	key := Key{
		Namespace:  obj.GetNamespace(),
		APIVersion: obj.GetAPIVersion(),
		Kind:       obj.GetKind(),
		Name:       obj.GetName(),
	}

	mc.mu.Lock()
	mc.store[key] = obj
	mc.mu.Unlock()

	mc.notify(StoreAction, key)

	return nil
}

// Retrieve retrieves an object from the cache.
func (mc *MemoryCache) Retrieve(key Key) ([]*unstructured.Unstructured, error) {
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

	key := Key{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}

	mc.mu.Lock()
	delete(mc.store, key)
	mc.mu.Unlock()

	mc.notify(DeleteAction, key)

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

func (mc *MemoryCache) notify(action Action, key Key) {
	if mc.notifyCh == nil {
		return
	}

	select {
	case mc.notifyCh <- Notification{Action: action, CacheKey: key}:
	case <-mc.notifyDone:
	}
}

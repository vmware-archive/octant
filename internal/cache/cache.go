package cache

import (
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

//go:generate mockgen -destination=./fake/mock_cache.go -package=fake github.com/heptio/developer-dash/internal/cache Cache

// Cache stores Kubernetes objects.
type Cache interface {
	Store(obj *unstructured.Unstructured) error
	List(key Key) ([]*unstructured.Unstructured, error)
	Get(key Key) (*unstructured.Unstructured, error)
	Delete(obj *unstructured.Unstructured) error
}

// Key is a key for the cache.
type Key struct {
	Namespace  string
	APIVersion string
	Kind       string
	Name       string
	Selector   labels.Selector
}

// GetAs gets an object from the cache by key.
func GetAs(c Cache, key Key, as interface{}) error {
	u, err := c.Get(key)
	if err != nil {
		return errors.Wrap(err, "get object from cache")
	}

	if u == nil {
		return nil
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, as)
	if err != nil {
		return err
	}

	if err := copyObjectMeta(as, u); err != nil {
		return errors.Wrap(err, "copy object metadata")
	}

	return nil
}

// TODO: see if all the other versions of this function could be replaced
func copyObjectMeta(to interface{}, from *unstructured.Unstructured) error {
	object, ok := to.(metav1.Object)
	if !ok {
		return errors.Errorf("%T is not an object", to)
	}

	t, err := meta.TypeAccessor(object)
	if err != nil {
		return errors.Wrapf(err, "accessing type meta")
	}
	t.SetAPIVersion(from.GetAPIVersion())
	t.SetKind(from.GetObjectKind().GroupVersionKind().Kind)

	object.SetNamespace(from.GetNamespace())
	object.SetName(from.GetName())
	object.SetGenerateName(from.GetGenerateName())
	object.SetUID(from.GetUID())
	object.SetResourceVersion(from.GetResourceVersion())
	object.SetGeneration(from.GetGeneration())
	object.SetSelfLink(from.GetSelfLink())
	object.SetCreationTimestamp(from.GetCreationTimestamp())
	object.SetDeletionTimestamp(from.GetDeletionTimestamp())
	object.SetDeletionGracePeriodSeconds(from.GetDeletionGracePeriodSeconds())
	object.SetLabels(from.GetLabels())
	object.SetAnnotations(from.GetAnnotations())
	object.SetInitializers(from.GetInitializers())
	object.SetOwnerReferences(from.GetOwnerReferences())
	object.SetClusterName(from.GetClusterName())
	object.SetFinalizers(from.GetFinalizers())

	return nil
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

// List retrieves a slice of objects from the cache.
func (mc *MemoryCache) List(key Key) ([]*unstructured.Unstructured, error) {
	if key.Name != "" {
		return nil, errors.Errorf("can't specify a name when listing objects")
	}

	if key.Namespace == "" ||
		key.APIVersion == "" ||
		key.Kind == "" {
		return nil, errors.New("requires namespace, apiVersion, and kind")
	}

	mc.mu.Lock()
	defer mc.mu.Unlock()

	var objects []*unstructured.Unstructured

	for _, v := range mc.store {
		if key.Namespace == v.GetNamespace() &&
			key.APIVersion == v.GetAPIVersion() &&
			key.Kind == v.GetKind() {
			objects = append(objects, v)
		}
	}

	return objects, nil
}

// List retrieves an object from the cache.
func (mc *MemoryCache) Get(key Key) (*unstructured.Unstructured, error) {
	if key.Namespace == "" ||
		key.APIVersion == "" ||
		key.Kind == "" ||
		key.Name == "" {
		return nil, errors.New("requires namespace, apiVersion, kind, and name")
	}
	mc.mu.Lock()
	defer mc.mu.Unlock()

	for _, v := range mc.store {
		if key.Namespace == v.GetNamespace() &&
			key.APIVersion == v.GetAPIVersion() &&
			key.Kind == v.GetKind() &&
			key.Name == v.GetName() {
			return v, nil
		}
	}

	return nil, errors.Errorf("object not found")
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

func (mc *MemoryCache) notify(action Action, key Key) {
	if mc.notifyCh == nil {
		return
	}

	select {
	case mc.notifyCh <- Notification{Action: action, CacheKey: key}:
	case <-mc.notifyDone:
	}
}

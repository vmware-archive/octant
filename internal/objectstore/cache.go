package objectstore

import (
	"context"
	"fmt"
	"sync"

	"go.opencensus.io/trace"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/pkg/store"
)

// ResourceCacheKey creates a key of Namespace and Resource
type ResourceCacheKey struct {
	Namespace string
	Resource  schema.GroupVersionResource
}

// ResourceCache stores a cache of cluster resources.
type ResourceCache struct {
	ctx context.Context

	data  *sync.Map
	reset sync.Mutex
}

// NewResourceCache creates a new uninitialized ResourceCache
func NewResourceCache(ctx context.Context) *ResourceCache {
	ctx, span := trace.StartSpan(ctx, "resCache:NewResourceCache")
	defer span.End()

	rc := &ResourceCache{
		ctx:  ctx,
		data: &sync.Map{},
	}
	return rc
}

// List returns all of the items for a given GroupVersionKind
func (r *ResourceCache) List(ctx context.Context, key ResourceCacheKey) (list *unstructured.UnstructuredList, loading bool, err error) {
	ctx, span := trace.StartSpan(ctx, "resCache:List")
	defer span.End()

	if !r.HasResource(ctx, key) {
		return nil, false, fmt.Errorf("cannot List from cache for uninitialized resource")
	}

	results := unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}

	v, _ := r.data.Load(key)
	itemMap, _ := v.(*sync.Map)

	itemMap.Range(func(_ interface{}, v interface{}) bool {
		item, ok := v.(unstructured.Unstructured)
		if !ok {
			return false
		}
		results.Items = append(results.Items, item)
		return true
	})
	return &results, false, nil
}

// Get gets a single resource from the cache.
func (r *ResourceCache) Get(ctx context.Context, key ResourceCacheKey, getKey store.Key) (*unstructured.Unstructured, error) {
	ctx, span := trace.StartSpan(ctx, "resCache:Get")
	defer span.End()

	if !r.HasResource(ctx, key) {
		return nil, fmt.Errorf("cannot Get from cache for uninitialized resource")
	}

	v, _ := r.data.Load(key)
	itemMap, _ := v.(*sync.Map)

	itemKey := store.Key{
		APIVersion: getKey.APIVersion,
		Kind:       getKey.Kind,
		Name:       getKey.Name,
	}

	v, ok := itemMap.Load(itemKey)
	if !ok {
		return nil, nil
	}

	item, ok := v.(unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("bad item in itemMap")
	}

	return &item, nil
}

// Initialize prepares the cache for a GroupVersionKind and sets the synced resource flag.
func (r *ResourceCache) Initialize(ctx context.Context, key ResourceCacheKey) error {
	ctx, span := trace.StartSpan(ctx, "resCache:Initialize")
	defer span.End()

	if r.HasResource(ctx, key) {
		return fmt.Errorf("resource is already initalized")
	}

	itemMap := &sync.Map{}
	r.data.Store(key, itemMap)
	return nil
}

// HasResource checks if the cache has been intialzied for a GroupVersionKind
func (r *ResourceCache) HasResource(ctx context.Context, key ResourceCacheKey) bool {
	_, span := trace.StartSpan(ctx, "resCache:HasResource")
	defer span.End()

	_, ok := r.data.Load(key)
	return ok
}

// AddMany adds many items to the cache for the GroupVersionResource.
func (r *ResourceCache) AddMany(ctx context.Context, key ResourceCacheKey, items ...unstructured.Unstructured) error {
	ctx, span := trace.StartSpan(ctx, "resCache:AddMany")
	defer span.End()

	if !r.HasResource(ctx, key) {
		return fmt.Errorf("can not add item for unintialized resource, must call Initialize first")
	}

	v, _ := r.data.Load(key)
	itemMap, ok := v.(*sync.Map)
	if !ok {
		return fmt.Errorf("unable to get itemMap from resourceMap")
	}

	for _, item := range items {
		key := store.Key{
			APIVersion: item.GetAPIVersion(),
			Kind:       item.GetKind(),
			Name:       item.GetName(),
		}
		itemMap.Store(key, item)
	}
	return nil
}

// Update updates a single item to the cache for a GroupVersionResource. Update returns the old value.
func (r *ResourceCache) Update(ctx context.Context, key ResourceCacheKey, item unstructured.Unstructured) (unstructured.Unstructured, error) {
	ctx, span := trace.StartSpan(ctx, "resCache:Update")
	defer span.End()

	if !r.HasResource(ctx, key) {
		return unstructured.Unstructured{}, fmt.Errorf("can not add item for unintialized resource, must call Initialize first")
	}

	itemKey := store.Key{
		APIVersion: item.GetAPIVersion(),
		Kind:       item.GetKind(),
		Name:       item.GetName(),
	}

	v, _ := r.data.Load(key)
	itemMap, ok := v.(*sync.Map)
	if !ok {
		return unstructured.Unstructured{}, fmt.Errorf("unable to get itemMap from resourceMap")
	}

	var old unstructured.Unstructured
	if u, err := r.Get(ctx, key, itemKey); err != nil {
		old = *u
	} else {
		old = item
	}

	itemMap.Store(itemKey, item)
	return old, nil
}

// Add adds a single item to the cache for a GroupVersionResource.
func (r *ResourceCache) Add(ctx context.Context, key ResourceCacheKey, item unstructured.Unstructured) error {
	ctx, span := trace.StartSpan(ctx, "resCache:Add")
	defer span.End()

	if !r.HasResource(ctx, key) {
		return fmt.Errorf("can not add item for unintialized resource, must call Initialize first")
	}

	itemKey := store.Key{
		APIVersion: item.GetAPIVersion(),
		Kind:       item.GetKind(),
		Name:       item.GetName(),
	}

	v, _ := r.data.Load(key)
	itemMap, ok := v.(*sync.Map)
	if !ok {
		return fmt.Errorf("unable to get itemMap from resourceMap")
	}
	itemMap.Store(itemKey, item)
	return nil
}

// Delete removes a single resource from the cache for a GroupVersionResource
func (r *ResourceCache) Delete(ctx context.Context, key ResourceCacheKey, item unstructured.Unstructured) error {
	ctx, span := trace.StartSpan(r.ctx, "resCache:Delete")
	defer span.End()

	if !r.HasResource(ctx, key) {
		return fmt.Errorf("can not add item for unintialized resource, must call Initialize first")
	}

	itemKey := store.Key{
		APIVersion: item.GetAPIVersion(),
		Kind:       item.GetKind(),
		Name:       item.GetName(),
	}

	v, _ := r.data.Load(key)
	itemMap, ok := v.(*sync.Map)
	if !ok {
		return nil
	}
	itemMap.Delete(itemKey)
	return nil
}

// Reset clears the resource cache
func (r *ResourceCache) Reset(ctx context.Context) {
	_, span := trace.StartSpan(r.ctx, "resCache:Reset")
	defer span.End()

	r.reset.Lock()
	defer r.reset.Unlock()
	r.data = &sync.Map{}
}

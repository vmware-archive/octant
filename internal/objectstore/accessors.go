package objectstore

import (
	"sync"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

// TODO: investigate sync.Map

type factoriesCache struct {
	factories map[string]dynamicinformer.DynamicSharedInformerFactory

	mu sync.RWMutex
}

func initFactoriesCache() *factoriesCache {
	return &factoriesCache{
		factories: make(map[string]dynamicinformer.DynamicSharedInformerFactory),
	}
}

func (c *factoriesCache) set(key string, value dynamicinformer.DynamicSharedInformerFactory) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.factories[key] = value
}

func (c *factoriesCache) get(key string) (dynamicinformer.DynamicSharedInformerFactory, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.factories[key]
	return v, ok
}

func (c *factoriesCache) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.factories, key)
}

type accessCache struct {
	access accessMap

	mu sync.RWMutex
}

func initAccessCache() *accessCache {
	return &accessCache{
		access: accessMap{},
	}
}

func (c *accessCache) set(key accessKey, value bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.access[key] = value
}

func (c *accessCache) get(key accessKey) (v, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok = c.access[key]
	return v, ok
}

type seenGVKsCache struct {
	seenGVKs map[string]map[schema.GroupVersionKind]bool

	mu sync.RWMutex
}

func initSeenGVKsCache() *seenGVKsCache {
	return &seenGVKsCache{
		seenGVKs: make(map[string]map[schema.GroupVersionKind]bool),
	}
}

func (c *seenGVKsCache) setSeen(key string, groupVersionKind schema.GroupVersionKind, value bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cur, ok := c.seenGVKs[key]
	if !ok {
		cur = make(map[schema.GroupVersionKind]bool)
	}

	cur[groupVersionKind] = value
	c.seenGVKs[key] = cur
}

func (c *seenGVKsCache) hasSeen(key string, groupVersionKind schema.GroupVersionKind) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.seenGVKs[key]
	if !ok {
		return false
	}

	seen, ok := v[groupVersionKind]
	if !ok {
		return false
	}

	return seen
}

type cachedObjectsCache struct {
	cachedObjects map[string]map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured
	mu            sync.RWMutex
}

func initCachedObjectsCache() *cachedObjectsCache {
	return &cachedObjectsCache{
		cachedObjects: make(map[string]map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured),
	}
}

func (c *cachedObjectsCache) list(key string, groupVersionKind schema.GroupVersionKind) []*unstructured.Unstructured {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var list []*unstructured.Unstructured

	gvkList, ok := c.cachedObjects[key]
	if !ok {
		return list
	}

	objectMap, ok := gvkList[groupVersionKind]
	if !ok {
		return list
	}

	for _, object := range objectMap {
		list = append(list, object)
	}

	return list
}

func (c *cachedObjectsCache) update(ns string, groupVersionKind schema.GroupVersionKind, object *unstructured.Unstructured) {
	if object == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	cur, ok := c.cachedObjects[ns]
	if !ok {
		cur = make(map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured)
	}

	curGVK, ok := cur[groupVersionKind]
	if !ok {
		curGVK = make(map[types.UID]*unstructured.Unstructured)
	}

	curGVK[object.GetUID()] = object
	cur[groupVersionKind] = curGVK
	c.cachedObjects[ns] = cur
}

func (c *cachedObjectsCache) delete(ns string, groupVersionKind schema.GroupVersionKind, object *unstructured.Unstructured) {
	if object == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	cur, ok := c.cachedObjects[ns]
	if !ok {
		return
	}

	curGVK, ok := cur[groupVersionKind]
	if !ok {
		return
	}

	delete(curGVK, object.GetUID())
	cur[groupVersionKind] = curGVK
	c.cachedObjects[ns] = cur
}

type watchedGVKsCache struct {
	watchedGVKs map[string]map[schema.GroupVersionKind]bool
	mu          sync.RWMutex
}

func initWatchedGVKsCache() *watchedGVKsCache {
	return &watchedGVKsCache{
		watchedGVKs: make(map[string]map[schema.GroupVersionKind]bool),
	}
}

func (c *watchedGVKsCache) isWatched(key string, groupVersionKind schema.GroupVersionKind) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	gvkMap, ok := c.watchedGVKs[key]
	if !ok {
		return false
	}
	return gvkMap[groupVersionKind]
}

func (c *watchedGVKsCache) setWatched(key string, groupVersionKind schema.GroupVersionKind) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cur, ok := c.watchedGVKs[key]
	if !ok {
		cur = make(map[schema.GroupVersionKind]bool)
	}

	cur[groupVersionKind] = true
	c.watchedGVKs[key] = cur
}

package objectstore

import (
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"

	"github.com/vmware/octant/pkg/store"
)

type informerSynced struct {
	status map[string]bool

	mu sync.RWMutex
}

func initInformerSynced() *informerSynced {
	return &informerSynced{
		status: make(map[string]bool),
	}
}

func (c *informerSynced) setSynced(key store.Key, value bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.status[key.String()] = value
}

func (c *informerSynced) hasSynced(key store.Key) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.status[key.String()]
	if !ok {
		return true
	}

	return v
}

func (c *informerSynced) hasSeen(key store.Key) bool {
	_, ok := c.status[key.String()]
	return ok
}

func (c *informerSynced) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.status {
		delete(c.status, key)
	}
}

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

func (c *factoriesCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name := range c.factories {
		delete(c.factories, name)
	}
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

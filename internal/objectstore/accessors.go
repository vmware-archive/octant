package objectstore

import (
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"

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

type factoriesCache struct {
	factories map[string]InformerFactory

	mu sync.RWMutex
}

func initFactoriesCache() *factoriesCache {
	return &factoriesCache{
		factories: make(map[string]InformerFactory),
	}
}

func (c *factoriesCache) set(key string, value InformerFactory) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.factories[key] = value
}

func (c *factoriesCache) keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var list []string
	for k := range c.factories {
		list = append(list, k)
	}

	return list
}

func (c *factoriesCache) get(key string) (InformerFactory, bool) {
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

func (c *seenGVKsCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.seenGVKs {
		delete(c.seenGVKs, k)
	}
}

type informerContextCache struct {
	cache map[schema.GroupVersionResource]chan struct{}

	mu sync.Mutex
}

func initInformerContextCache() *informerContextCache {
	return &informerContextCache{
		cache: make(map[schema.GroupVersionResource]chan struct{}),
	}
}

func (c *informerContextCache) addChild(key schema.GroupVersionResource) <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan struct{}, 1)
	c.cache[key] = ch
	return ch
}

func (c *informerContextCache) delete(key schema.GroupVersionResource) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if stopCh, ok := c.cache[key]; ok {
		close(stopCh)
		delete(c.cache, key)
	}
}

func (c *informerContextCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, stopCh := range c.cache {
		close(stopCh)
		delete(c.cache, k)
	}
}

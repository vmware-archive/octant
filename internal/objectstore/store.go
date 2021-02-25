package objectstore

import (
	"context"
	"sync"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

// MemoryStore provides the struct for the store interface implementation.
type MemoryStore struct {
	ctx    context.Context
	cancel context.CancelFunc

	client         cluster.ClientInterface
	resourceLookup *sync.Map

	watcher *Watcher
	cache   *ResourceCache
	update  sync.Mutex
}

var _ store.Store = &MemoryStore{}

// NewMemoryStore implements the store interface in memory.
func NewMemoryStore(ctx context.Context, client cluster.ClientInterface) (*MemoryStore, error) {
	c, cf := context.WithCancel(ctx)
	ms := &MemoryStore{
		ctx:            c,
		cancel:         cf,
		client:         client,
		cache:          NewResourceCache(ctx),
		resourceLookup: &sync.Map{},
	}
	ms.watcher = NewWatcher(ctx, ms.client, ms.cache)
	return ms, nil
}

func (m *MemoryStore) cacheKeyFromKey(key store.Key) (ResourceCacheKey, error) {
	res, _, err := m.client.Resource(key.GroupVersionKind().GroupKind())
	if err != nil {
		return ResourceCacheKey{}, err
	}
	cacheKey := ResourceCacheKey{
		Namespace: key.Namespace,
		Resource:  res,
	}
	return cacheKey, nil
}

// List lists all the resources.
func (m *MemoryStore) List(ctx context.Context, key store.Key) (list *unstructured.UnstructuredList, loading bool, err error) {
	cacheKey, err := m.cacheKeyFromKey(key)
	if err != nil {
		return nil, false, err
	}

	if m.cache.HasResource(cacheKey) {
		return m.cache.List(ctx, cacheKey)
	}

	m.cache.Initialize(cacheKey)

	dc, err := m.client.DynamicClient()
	if err != nil {
		return nil, false, err
	}

	listOptions := metav1.ListOptions{}

	var listing *unstructured.UnstructuredList
	if key.Namespace == "" {
		listing, err = dc.Resource(cacheKey.Resource).List(ctx, listOptions)
		if err != nil {
			return nil, false, err
		}
	} else {
		listing, err = dc.Resource(cacheKey.Resource).Namespace(key.Namespace).List(ctx, listOptions)
		if err != nil {
			return nil, false, err
		}
	}

	go func() {
		items := listing.DeepCopy().Items
		m.cache.AddMany(cacheKey, items...)
	}()

	m.watcher.Watch(cacheKey)

	return listing, false, err
}

// Get returns a single resource
func (m *MemoryStore) Get(ctx context.Context, key store.Key) (object *unstructured.Unstructured, err error) {
	cacheKey, err := m.cacheKeyFromKey(key)
	if err != nil {
		return nil, err
	}

	if m.cache.HasResource(cacheKey) {
		return m.cache.Get(ctx, cacheKey, key)
	}

	m.cache.Initialize(cacheKey)

	dc, err := m.client.DynamicClient()
	if err != nil {
		return nil, err
	}

	getOptions := metav1.GetOptions{}
	var item *unstructured.Unstructured
	if key.Namespace == "" {
		item, err = dc.Resource(cacheKey.Resource).Get(ctx, key.Name, getOptions)
		if err != nil {
			return nil, err
		}
	} else {
		item, err = dc.Resource(cacheKey.Resource).Namespace(key.Namespace).Get(ctx, key.Name, getOptions)
		if err != nil {
			return nil, err
		}
	}

	go func() {
		i := *item.DeepCopy()
		m.cache.Add(cacheKey, i)
	}()

	return item, err

}

// Delete deletes a single resource.
func (m *MemoryStore) Delete(ctx context.Context, key store.Key) error {
	cacheKey, err := m.cacheKeyFromKey(key)
	if err != nil {
		return err
	}

	dc, err := m.client.DynamicClient()
	if err != nil {
		return err
	}

	delOptions := metav1.DeleteOptions{}
	if key.Namespace == "" {
		err := dc.Resource(cacheKey.Resource).Delete(ctx, key.Name, delOptions)
		if err != nil {
			return err
		}
	} else {
		err := dc.Resource(cacheKey.Resource).Namespace(key.Namespace).Delete(ctx, key.Name, delOptions)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryStore) Watch(ctx context.Context, key store.Key, handler cache.ResourceEventHandler) error {
	return nil
}
func (m *MemoryStore) Unwatch(ctx context.Context, groupVersionKinds ...schema.GroupVersionKind) error {
	return nil
}
func (m *MemoryStore) UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error {
	m.update.Lock()
	defer m.update.Unlock()

	m.client = client
	m.resourceLookup = &sync.Map{}
	m.cache.Reset()

	m.watcher.StopAll()
	m.watcher = NewWatcher(m.ctx, m.client, m.cache)

	return nil
}
func (m *MemoryStore) RegisterOnUpdate(fn store.UpdateFn) {}
func (m *MemoryStore) Update(ctx context.Context, key store.Key, updater func(*unstructured.Unstructured) error) error {
	return nil
}
func (m *MemoryStore) IsLoading(ctx context.Context, key store.Key) bool {
	return false
}
func (m *MemoryStore) Create(ctx context.Context, object *unstructured.Unstructured) error {
	return nil
}

// CreateOrUpdateFromYAML creates resources in the cluster from YAML input.
// Resources are created in the order they are present in the YAML.
// An error creating a resource halts resource creation.
// A list of created resources is returned. You may have created resources AND a non-nil error.
func (m *MemoryStore) CreateOrUpdateFromYAML(ctx context.Context, namespace, input string) ([]string, error) {
	return []string{}, nil
}

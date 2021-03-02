package objectstore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	kcache "k8s.io/client-go/tools/cache"
	kretry "k8s.io/client-go/util/retry"
	sigyaml "sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// MemoryStore provides the struct for the store interface implementation.
type MemoryStore struct {
	ctx    context.Context
	cancel context.CancelFunc

	client         cluster.ClientInterface
	resourceLookup *sync.Map

	watcher *Watcher
	cache   *ResourceCache

	updateFns []store.UpdateFn
	update    sync.Mutex
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
		return m.cache.List(cacheKey)
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
		return m.cache.Get(cacheKey, key)
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

// Watch watches a resource.
func (m *MemoryStore) Watch(ctx context.Context, key store.Key, handler kcache.ResourceEventHandler) error {
	cacheKey, err := m.cacheKeyFromKey(key)
	if err != nil {
		return err
	}
	m.watcher.AddCallback(cacheKey, handler)
	return nil
}

// Unwatch unwatches a resource.
func (m *MemoryStore) Unwatch(ctx context.Context, groupVersionKinds ...schema.GroupVersionKind) error {
	for _, gvk := range groupVersionKinds {
		key := store.KeyFromGroupVersionKind(gvk)
		cacheKey, err := m.cacheKeyFromKey(key)
		if err != nil {
			return err
		}
		m.watcher.DeleteCallback(cacheKey)
	}
	return nil
}

// UpdateClusterClient resets the store and updates the internal ClusterClient
func (m *MemoryStore) UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error {
	m.update.Lock()
	defer m.update.Unlock()

	m.client = client
	m.resourceLookup = &sync.Map{}
	m.cache.Reset()

	m.watcher.StopAll()
	m.watcher = NewWatcher(m.ctx, m.client, m.cache)

	for _, u := range m.updateFns {
		u(m)
	}

	return nil
}

// RegisterOnUpdate registers callback functions that should be callend when UpdateClusterClient is invoked.
func (m *MemoryStore) RegisterOnUpdate(fn store.UpdateFn) {
	m.update.Lock()
	defer m.update.Unlock()
	m.updateFns = append(m.updateFns, fn)
}

// Update updates a resourec.
func (m *MemoryStore) Update(ctx context.Context, key store.Key, updater func(*unstructured.Unstructured) error) error {
	if updater == nil {
		return fmt.Errorf("can't update object")
	}

	err := kretry.RetryOnConflict(kretry.DefaultRetry, func() error {
		object, err := m.Get(ctx, key)
		if err != nil {
			return err
		}

		if object == nil {
			return fmt.Errorf("object not found")
		}

		cacheKey, err := m.cacheKeyFromKey(key)
		if err != nil {
			return err
		}

		dynamicClient, err := m.client.DynamicClient()
		if err != nil {
			return err
		}

		if err := updater(object); err != nil {
			return fmt.Errorf("unable to update object: %w", err)
		}

		client := dynamicClient.Resource(cacheKey.Resource).Namespace(object.GetNamespace())

		_, err = client.Update(ctx, object, metav1.UpdateOptions{})
		return err
	})

	return err
}

// IsLoading returns if a resource is currently loading.
func (m *MemoryStore) IsLoading(ctx context.Context, key store.Key) bool {
	return false
}

// Create creates a resource.
func (m *MemoryStore) Create(ctx context.Context, object *unstructured.Unstructured) error {
	key, err := store.KeyFromObject(object)
	if err != nil {
		return fmt.Errorf("key from object: %w", err)
	}

	dynamicClient, err := m.client.DynamicClient()
	if err != nil {
		return err
	}

	cacheKey, err := m.cacheKeyFromKey(key)
	if err != nil {
		return err
	}

	createOptions := metav1.CreateOptions{}

	if key.Namespace == "" {
		_, err := dynamicClient.Resource(cacheKey.Resource).Create(ctx, object, createOptions)
		return err
	}

	_, err = dynamicClient.Resource(cacheKey.Resource).Namespace(key.Namespace).Create(ctx, object, createOptions)
	return err
}

// CreateOrUpdateFromHandler encapulates the Apply YAML logic
func CreateOrUpdateFromHandler(
	ctx context.Context, namespace, input string,
	get func(context.Context, store.Key) (*unstructured.Unstructured, error),
	create func(context.Context, *unstructured.Unstructured) error,
	clusterClient cluster.ClientInterface,
) ([]string, error) {
	withDoc := func(cb func(doc map[string]interface{}) error) error {
		d := yaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(input), 4096)
		for {
			doc := map[string]interface{}{}
			if err := d.Decode(&doc); err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("unable to parse yaml: %w", err)
			}
			if len(doc) == 0 {
				// skip empty documents
				continue
			}
			if err := cb(doc); err != nil {
				return err
			}
		}
	}

	logger := log.From(ctx)
	var results []string
	err := withDoc(func(doc map[string]interface{}) error {
		logger.Debugf("apply resource %#v", doc)

		unstructuredObj := &unstructured.Unstructured{Object: doc}
		key, err := store.KeyFromObject(unstructuredObj)
		if err != nil {
			return err
		}

		gvr, namespaced, err := clusterClient.Resource(key.GroupVersionKind().GroupKind())
		if err != nil {
			return fmt.Errorf("unable to discover resource: %w", err)
		}
		if namespaced && key.Namespace == "" {
			unstructuredObj.SetNamespace(namespace)
			key.Namespace = namespace
		}

		if _, err := get(ctx, key); err != nil {
			if !kerrors.IsNotFound(err) {
				// unexpected error
				return fmt.Errorf("unable to get resource: %w", err)
			}

			// create object
			err := create(ctx, &unstructured.Unstructured{Object: doc})
			if err != nil {
				return fmt.Errorf("unable to create resource: %w", err)
			}

			result := fmt.Sprintf("Created %s (%s) %s", key.Kind, key.APIVersion, key.Name)
			if namespaced {
				result = fmt.Sprintf("%s in %s", result, key.Namespace)
			}
			results = append(results, result)

			return nil
		}

		// update object
		unstructuredYaml, err := sigyaml.Marshal(doc)
		if err != nil {
			return fmt.Errorf("unable to marshal resource as yaml: %w", err)
		}
		client, err := clusterClient.DynamicClient()
		if err != nil {
			return fmt.Errorf("unable to get dynamic client: %w", err)
		}

		withForce := true
		if namespaced {
			_, err = client.Resource(gvr).Namespace(key.Namespace).Patch(
				ctx,
				key.Name,
				types.ApplyPatchType,
				unstructuredYaml,
				metav1.PatchOptions{FieldManager: "octant", Force: &withForce},
			)
			if err != nil {
				return fmt.Errorf("unable to patch resource: %w", err)
			}
		} else {
			_, err = client.Resource(gvr).Patch(
				ctx,
				key.Name,
				types.ApplyPatchType,
				unstructuredYaml,
				metav1.PatchOptions{FieldManager: "octant", Force: &withForce},
			)
			if err != nil {
				return fmt.Errorf("unable to patch resource: %w", err)
			}
		}

		result := fmt.Sprintf("Updated %s (%s) %s", key.Kind, key.APIVersion, key.Name)
		if namespaced {
			result = fmt.Sprintf("%s in %s", result, key.Namespace)
		}
		results = append(results, result)

		return nil
	})
	return results, err
}

// CreateOrUpdateFromYAML creates resources in the cluster from YAML input.
// Resources are created in the order they are present in the YAML.
// An error creating a resource halts resource creation.
// A list of created resources is returned. You may have created resources AND a non-nil error.
func (m *MemoryStore) CreateOrUpdateFromYAML(ctx context.Context, namespace, input string) ([]string, error) {
	return CreateOrUpdateFromHandler(ctx, namespace, input, m.Get, m.Create, m.client)
}

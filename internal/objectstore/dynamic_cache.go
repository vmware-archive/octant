package objectstore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"k8s.io/client-go/dynamic"

	"go.opencensus.io/trace"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	sigyaml "sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/store"
)

const resyncPeriod = time.Second * 180

type lister interface {
	// List will return all objects in this namespace
	List(selector labels.Selector) (ret []runtime.Object, err error)
	// Get will attempt to retrieve by namespace and name
	Get(name string) (runtime.Object, error)
}

type DynamicCache struct {
	ctx           context.Context
	cancel        context.CancelFunc
	client        cluster.ClientInterface
	dynamicClient dynamic.Interface
	stopChan      chan struct{}

	informerFactories sync.Map // namespace:DynamicSharedInformerFactory
	knownInformers    sync.Map // SharedIndexInformer:chan struct{}

	removeCh chan schema.GroupVersionResource
}

var _ store.Store = (*DynamicCache)(nil)

type Option func(*DynamicCache)

func NewDynamicCache(ctx context.Context, client cluster.ClientInterface, opts ...Option) (*DynamicCache, error) {
	dynamicClient, err := client.DynamicClient()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)

	dc := &DynamicCache{
		ctx:               ctx,
		cancel:            cancel,
		client:            client,
		dynamicClient:     dynamicClient,
		removeCh:          make(chan schema.GroupVersionResource),
		informerFactories: sync.Map{},
		stopChan:          make(chan struct{}),
	}

	for _, opt := range opts {
		opt(dc)
	}

	go dc.shutdownWorker()

	return dc, nil
}

func (d *DynamicCache) List(ctx context.Context, key store.Key) (list *unstructured.UnstructuredList, loading bool, err error) {
	_, span := trace.StartSpan(ctx, "dynamicCache:List")
	defer span.End()

	resourceLister, err := d.listerForResource(ctx, key)
	if err != nil {
		return nil, false, err
	}

	if key.Selector != nil && key.LabelSelector != nil {
		return nil, false, fmt.Errorf("must provide only one of Key.Selector and Key.LabelSelector")
	}

	var selector = labels.Everything()
	if key.Selector != nil {
		selector = key.Selector.AsSelector()
	} else if key.LabelSelector != nil {
		selector, err = metav1.LabelSelectorAsSelector(key.LabelSelector)
		if err != nil {
			return nil, false, err
		}
	}

	span.AddAttributes(
		trace.StringAttribute("key", fmt.Sprintf("%s", key)),
		trace.StringAttribute("selector", fmt.Sprintf("%s", selector)),
	)

	objs, err := resourceLister.List(selector)
	if err != nil {
		return nil, false, err
	}

	objectCount := len(objs)

	span.AddAttributes(
		trace.Int64Attribute("objectCount", int64(objectCount)),
	)
	ul := &unstructured.UnstructuredList{}
	ul.Items = make([]unstructured.Unstructured, len(objs))

	for i := 0; i < objectCount; i++ {
		u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(objs[i])
		if err != nil {
			return nil, false, err
		}
		ul.Items[i].Object = u
	}

	return ul, false, err
}

func (d *DynamicCache) Get(ctx context.Context, key store.Key) (object *unstructured.Unstructured, err error) {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:Get")
	defer span.End()

	resourceLister, err := d.listerForResource(ctx, key)
	if err != nil {
		return nil, err
	}

	span.AddAttributes(
		trace.StringAttribute("key", fmt.Sprintf("%s", key)),
	)

	obj, err := resourceLister.Get(key.Name)
	if err != nil {
		return nil, err
	}
	return obj.(*unstructured.Unstructured), nil
}

func (d *DynamicCache) Delete(ctx context.Context, key store.Key) error {
	_, span := trace.StartSpan(ctx, "dynamicCache:delete")
	defer span.End()

	if d.dynamicClient == nil {
		return fmt.Errorf("dynamic client is nil")
	}

	gvr, err := d.gvrFromKey(ctx, key)
	if err != nil {
		return err
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	if key.Namespace == "" {
		return d.dynamicClient.Resource(gvr).Delete(ctx, key.Name, deleteOptions)
	}

	return d.dynamicClient.Resource(gvr).Namespace(key.Namespace).Delete(ctx, key.Name, deleteOptions)
}

func (d *DynamicCache) UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error {
	_, span := trace.StartSpan(ctx, "dynamicCache:UpdateClusterClient")
	defer span.End()

	d.stopAllInformers()
	d.client = client

	d.knownInformers, d.informerFactories = sync.Map{}, sync.Map{}

	return nil
}

func (d *DynamicCache) Update(ctx context.Context, key store.Key, updater func(*unstructured.Unstructured) error) error {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:Update")
	defer span.End()

	if updater == nil {
		return fmt.Errorf("can't update object")
	}

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		object, err := d.Get(ctx, key)
		if err != nil {
			return err
		}

		if object == nil {
			return errors.New("object not found")
		}

		gvr, err := d.gvrFromKey(ctx, key)
		if err != nil {
			return err
		}

		if d.dynamicClient == nil {
			return fmt.Errorf("dynamic client is nil")
		}

		if err := updater(object); err != nil {
			return fmt.Errorf("unable to update object: %w", err)
		}
		client := d.dynamicClient.Resource(gvr).Namespace(object.GetNamespace())

		_, err = client.Update(ctx, object, metav1.UpdateOptions{})
		return err
	})

	return err

}

func (d *DynamicCache) IsLoading(ctx context.Context, key store.Key) bool {
	_, span := trace.StartSpan(ctx, "dynamicCache:IsLoading")
	defer span.End()
	return false
}

func (d *DynamicCache) Create(ctx context.Context, object *unstructured.Unstructured) error {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:Create")
	defer span.End()

	key, err := store.KeyFromObject(object)
	if err != nil {
		return fmt.Errorf("key from object: %w", err)
	}

	gvr, err := d.gvrFromKey(ctx, key)
	if err != nil {
		return err
	}

	createOptions := metav1.CreateOptions{}

	if d.dynamicClient == nil {
		return fmt.Errorf("dynamic client is nil")
	}

	if key.Namespace == "" {
		_, err := d.dynamicClient.Resource(gvr).Create(ctx, object, createOptions)
		return err
	}

	_, err = d.dynamicClient.Resource(gvr).Namespace(key.Namespace).Create(ctx, object, createOptions)
	return err
}

func (d *DynamicCache) WaitForCacheSync(ctx context.Context) bool {
	_, span := trace.StartSpan(ctx, "dynamicCache:WaitForCacheSync")
	defer span.End()

	var ret bool

	d.informerFactories.Range(func(k, v interface{}) bool {
		v.(dynamicinformer.DynamicSharedInformerFactory).WaitForCacheSync(d.stopChan)
		return true
	})
	return ret
}

func (d *DynamicCache) informerFactoryForNamespace(namespace string) (dynamicinformer.DynamicSharedInformerFactory, error) {
	if d.dynamicClient == nil {
		return nil, fmt.Errorf("dynamic client is nil")
	}

	var informerFactory dynamicinformer.DynamicSharedInformerFactory
	inf, ok := d.informerFactories.Load(namespace)
	if !ok {
		f := dynamicinformer.NewFilteredDynamicSharedInformerFactory(d.dynamicClient, resyncPeriod, namespace, nil)
		d.informerFactories.Store(namespace, f)
		informerFactory = f.(dynamicinformer.DynamicSharedInformerFactory)
	} else {
		informerFactory = inf.(dynamicinformer.DynamicSharedInformerFactory)
	}

	return informerFactory, nil
}

// Watch creates informers for CRDs
func (d *DynamicCache) Watch(ctx context.Context, key store.Key, handler cache.ResourceEventHandler) error {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:Watch")
	defer span.End()

	logger := log.From(ctx)
	logger.With("dynamicCache", "watch")
	logger.Debugf("creating watch for %s", key)

	gvr, err := d.gvrFromKey(ctx, key)
	if err != nil {
		logger.Warnf("cannot find key: %w", err)
		return nil
	}

	if gvr.Empty() {
		return nil
	}

	stopChan := make(chan struct{})
	informerFactory, err := d.informerFactoryForNamespace(key.Namespace)
	if err != nil {
		return err
	}
	informer := informerFactory.ForResource(gvr).Informer()
	d.knownInformers.LoadOrStore(informer, stopChan)
	informer.SetWatchErrorHandler(d.watchErrorHandler(ctx, gvr))
	if handler != nil {
		informer.AddEventHandlerWithResyncPeriod(handler, resyncPeriod)
	}
	go informer.Run(stopChan)
	span.AddAttributes(trace.StringAttribute("key", fmt.Sprintf("%s", key)))

	return nil
}

// Unwatch needs to remove informers for CRDs: See https://github.com/kubernetes/kubernetes/pull/97214
// TODO: Instead of Watch and Unwatch for a GVK, this should be for adding and removing event handlers
func (d *DynamicCache) Unwatch(_ context.Context, groupVersionKinds ...schema.GroupVersionKind) error {
	return nil
}

func (d *DynamicCache) shutdownWorker() {
	for {
		select {
		case <-d.ctx.Done():
			d.stopAllInformers()
			return
		case <-time.After(time.Millisecond * 500):
			continue
		}
	}
}

func (d *DynamicCache) stopAllInformers() {
	if d.stopChan != nil {
		close(d.stopChan)
		d.stopChan = nil
	}

	d.knownInformers.Range(func(k, v interface{}) bool {
		stopCh := v.(chan struct{})
		close(stopCh)
		return true
	})
	d.knownInformers = sync.Map{}
}

func (d *DynamicCache) listerForResource(ctx context.Context, key store.Key) (lister, error) {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:ListerForResource")
	defer span.End()

	gvr, err := d.gvrFromKey(ctx, key)
	if err != nil {
		return nil, err
	}

	span.AddAttributes(
		trace.StringAttribute("key", fmt.Sprintf("%s", key)),
		trace.StringAttribute("gvr", fmt.Sprintf("%s", gvr)),
	)

	informerFactory, err := d.informerFactoryForNamespace(key.Namespace)
	if err != nil {
		return nil, fmt.Errorf("create informer for namespace %s: %v", key.Namespace, err)
	}

	genericInformer := informerFactory.ForResource(gvr)
	stopChan := make(chan struct{})

	inf := genericInformer.Informer()
	inf.SetWatchErrorHandler(d.watchErrorHandler(ctx, gvr))
	_, loaded := d.knownInformers.LoadOrStore(inf, stopChan)
	if !loaded {
		go inf.Run(stopChan)
	}

	var l lister
	if key.Namespace == "" {
		l = genericInformer.Lister()
	} else {
		l = genericInformer.Lister().ByNamespace(key.Namespace)
	}

	return l, nil
}

func (d *DynamicCache) watchErrorHandler(ctx context.Context, gvr schema.GroupVersionResource) func(*cache.Reflector, error) {
	return func(r *cache.Reflector, err error) {
		_, span := trace.StartSpan(ctx, "dynamicCache:watchErrorHandler")
		defer span.End()

		span.AddAttributes(trace.StringAttribute("gvr", fmt.Sprintf("%s", gvr)))

		logger := log.From(ctx)
		logger.Warnf("unable to start watcher ", err.Error())

		d.removeCh <- gvr
	}
}

func (d *DynamicCache) gvrFromKey(ctx context.Context, key store.Key) (schema.GroupVersionResource, error) {
	_, span := trace.StartSpan(ctx, "dynamicCache:gvrFromKey")
	defer span.End()

	gk := key.GroupVersionKind().GroupKind()
	gvr, _, err := d.client.Resource(key.GroupVersionKind())
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	span.AddAttributes(
		trace.StringAttribute("live", "true"),
		trace.StringAttribute("key", fmt.Sprintf("%s", key)),
		trace.StringAttribute("gvr", fmt.Sprintf("%s %s %s", gvr.Group, gvr.Version, gvr.Resource)),
		trace.StringAttribute("gk", fmt.Sprintf("%s", gk)),
	)
	return gvr, nil
}

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
		gvr, namespaced, err := clusterClient.Resource(key.GroupVersionKind())
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
func (d *DynamicCache) CreateOrUpdateFromYAML(ctx context.Context, namespace, input string) ([]string, error) {
	return CreateOrUpdateFromHandler(ctx, namespace, input, d.Get, d.Create, d.client)
}

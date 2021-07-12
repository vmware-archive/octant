package objectstore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"go.opencensus.io/trace"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
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
	oerrors "github.com/vmware-tanzu/octant/internal/errors"
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
	ctx    context.Context
	cancel context.CancelFunc
	client cluster.ClientInterface

	informerFactory dynamicinformer.DynamicSharedInformerFactory
	knownInformers  sync.Map // gvr:GenericInformer
	unwatched       sync.Map // gvr:bool
	gvrCache        sync.Map // gk:gvr

	removeCh chan schema.GroupVersionResource
	mu       sync.Mutex
}

var _ store.Store = (*DynamicCache)(nil)

type Option func(*DynamicCache)

func WithDynamicSharedInformerFactory(factory dynamicinformer.DynamicSharedInformerFactory) Option {
	return func(d *DynamicCache) {
		d.informerFactory = factory
	}
}

func NewDynamicCache(ctx context.Context, client cluster.ClientInterface, opts ...Option) (*DynamicCache, error) {
	ctx, cancel := context.WithCancel(ctx)
	dynamicClient, err := client.DynamicClient()
	if err != nil {
		cancel()
		return nil, err
	}

	dc := &DynamicCache{
		ctx:            ctx,
		cancel:         cancel,
		client:         client,
		knownInformers: sync.Map{},
		unwatched:      sync.Map{},
		gvrCache:       sync.Map{},
		removeCh:       make(chan schema.GroupVersionResource),
	}

	for _, opt := range opts {
		opt(dc)
	}

	if dc.informerFactory == nil {
		dc.informerFactory = dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, resyncPeriod)
	}

	go dc.worker()

	return dc, nil
}

func (d *DynamicCache) List(ctx context.Context, key store.Key) (list *unstructured.UnstructuredList, loading bool, err error) {
	_, span := trace.StartSpan(ctx, "dynamicCache:List")
	defer span.End()

	resourceLister, err := d.listerForResource(ctx, key)
	if err != nil {
		return nil, false, err
	}

	if resourceLister == nil {
		return nil, false, fmt.Errorf("resourceLister is nil")
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

	dynamicClient, err := d.client.DynamicClient()
	if err != nil {
		return err
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
		return dynamicClient.Resource(gvr).Delete(ctx, key.Name, deleteOptions)
	}

	return dynamicClient.Resource(gvr).Namespace(key.Namespace).Delete(ctx, key.Name, deleteOptions)
}

func (d *DynamicCache) UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error {
	_, span := trace.StartSpan(ctx, "dynamicCache:UpdateClusterClient")
	defer span.End()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.stopAllInformers()

	d.client = client
	dynamicClient, err := client.DynamicClient()
	if err != nil {
		return err
	}

	d.informerFactory = dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, resyncPeriod)
	d.knownInformers = sync.Map{}
	d.unwatched = sync.Map{}

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

		dynamicClient, err := d.client.DynamicClient()
		if err != nil {
			return err
		}

		if err := updater(object); err != nil {
			return fmt.Errorf("unable to update object: %w", err)
		}

		client := dynamicClient.Resource(gvr).Namespace(object.GetNamespace())

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

	dynamicClient, err := d.client.DynamicClient()
	if err != nil {
		return err
	}

	if key.Namespace == "" {
		_, err := dynamicClient.Resource(gvr).Create(ctx, object, createOptions)
		return err
	}

	_, err = dynamicClient.Resource(gvr).Namespace(key.Namespace).Create(ctx, object, createOptions)
	return err
}

func (d *DynamicCache) WaitForCacheSync(ctx context.Context) bool {
	_, span := trace.StartSpan(ctx, "dynamicCache:WaitForCacheSync")
	defer span.End()

	var ret bool
	d.knownInformers.Range(func(k, v interface{}) bool {
		ii := v.(interuptibleInformer)
		ret = cache.WaitForCacheSync(ii.stopCh, ii.informer.Informer().HasSynced)
		return ret
	})
	return ret
}

func (d *DynamicCache) Watch(ctx context.Context, key store.Key, handler cache.ResourceEventHandler) error {
	ctx, span := trace.StartSpan(ctx, "dynamicCache:Watch")
	defer span.End()

	logger := log.From(ctx)
	logger.With("dynamicCache", "watch")
	logger.Debugf("creating watch for %s", key)

	gvr, err := d.gvrFromKey(ctx, key)
	if err != nil {
		return err
	}

	if d.isUnwatched(ctx, gvr) {
		return fmt.Errorf("watcher was unable to start for %s", gvr)
	}

	span.AddAttributes(trace.StringAttribute("key", fmt.Sprintf("%s", key)))

	d.forResource(ctx, gvr, handler)

	return err
}

func (d *DynamicCache) Unwatch(ctx context.Context, groupVersionKinds ...schema.GroupVersionKind) error {
	return nil
}

func (d *DynamicCache) worker() {
	for {
		select {
		case <-d.ctx.Done():
			d.stopAllInformers()
			return
		case gvr := <-d.removeCh:
			d.unwatched.Store(gvr, true)
			v, ok := d.knownInformers.LoadAndDelete(gvr)
			if ok {
				ii := v.(interuptibleInformer)
				ii.Stop()
			}
		case <-time.After(time.Millisecond * 500):
			continue
		}
	}
}

func (d *DynamicCache) stopAllInformers() {
	d.knownInformers.Range(func(k, v interface{}) bool {
		ii := v.(interuptibleInformer)
		ii.Stop()
		return true
	})
}

func (d *DynamicCache) isUnwatched(ctx context.Context, gvr schema.GroupVersionResource) bool {
	_, ok := d.unwatched.Load(gvr)
	return ok
}

func (d *DynamicCache) forResource(ctx context.Context, gvr schema.GroupVersionResource, handler cache.ResourceEventHandler) interuptibleInformer {
	_, span := trace.StartSpan(ctx, "dynamicCache:forResource")
	defer span.End()

	logger := log.From(ctx)
	logger = logger.With("dynamicCache", "forResource")

	v, ok := d.knownInformers.Load(gvr)
	if !ok {
		i := d.informerFactory.ForResource(gvr)
		stopCh := make(chan struct{})
		i.Informer().SetWatchErrorHandler(d.watchErrorHandler(ctx, gvr, stopCh))
		if handler != nil {
			i.Informer().AddEventHandlerWithResyncPeriod(handler, resyncPeriod)
		}

		go func() {
			logger.Debugf("starting informer for %s", gvr)
			i.Informer().Run(stopCh)
			logger.Debugf("stopping informer for %s", gvr)
		}()

		ii := interuptibleInformer{
			stopCh,
			i,
			gvr,
		}
		d.knownInformers.Store(gvr, ii)
		return ii
	}
	ii := v.(interuptibleInformer)
	if handler != nil {
		ii.informer.Informer().AddEventHandlerWithResyncPeriod(handler, resyncPeriod)
	}
	return ii
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

	if d.isUnwatched(ctx, gvr) {
		err = fmt.Errorf("unable to get Lister for %s, watcher was unable to start", gvr)
		return nil, oerrors.NewAccessError(key, "List", err)
	}

	ii := d.forResource(ctx, gvr, nil)

	var l lister
	if key.Namespace == "" {
		l = ii.informer.Lister()
	} else {
		l = ii.informer.Lister().ByNamespace(key.Namespace)
	}

	return l, nil
}

func (d *DynamicCache) watchErrorHandler(ctx context.Context, gvr schema.GroupVersionResource, stopCh chan struct{}) func(*cache.Reflector, error) {
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

	// short circuit for CRD type
	if key.APIVersion == "apiextensions.k8s.io/v1" && key.Kind == "CustomResourceDefinition" {
		return schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}, nil
	}

	gk := key.GroupVersionKind().GroupKind()
	var gvr schema.GroupVersionResource

	v, ok := d.gvrCache.Load(gk)
	if !ok {
		gvr, _, err := d.client.Resource(gk)
		if err != nil {
			return schema.GroupVersionResource{}, err
		}

		if gvr.Resource == "" {
			_, gvr = meta.UnsafeGuessKindToResource(key.GroupVersionKind())
		}

		// do not store an empty GVR
		if gvr.Version == "" && gvr.Resource == "" {
			return schema.GroupVersionResource{}, fmt.Errorf("unable to locate GVR for GVK: %s", key.GroupVersionKind())
		}

		d.gvrCache.Store(gk, gvr)
	} else {
		gvr, _ = v.(schema.GroupVersionResource)
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
func (dc *DynamicCache) CreateOrUpdateFromYAML(ctx context.Context, namespace, input string) ([]string, error) {
	return CreateOrUpdateFromHandler(ctx, namespace, input, dc.Get, dc.Create, dc.client)
}

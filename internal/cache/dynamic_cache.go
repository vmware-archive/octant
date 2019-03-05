package cache

import (
	"fmt"
	"time"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/util/retry"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	kcache "k8s.io/client-go/tools/cache"
)

const (
	// defaultMutableResync is the resync period for informers.
	defaultInformerResync = time.Second * 180
)

func initDynamicSharedInformerFactory(client cluster.ClientInterface) (dynamicinformer.DynamicSharedInformerFactory, error) {
	dynamicClient, err := client.DynamicClient()
	if err != nil {
		return nil, err
	}

	factory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, defaultInformerResync)
	return factory, nil
}

// DynamicCacheOpt is an option for configuration DynamicCache.
type DynamicCacheOpt func(*DynamicCache)

// DynamicCache is a cache based on the dynamic shared informer factory.
type DynamicCache struct {
	initFactoryFunc func(cluster.ClientInterface) (dynamicinformer.DynamicSharedInformerFactory, error)
	factory         dynamicinformer.DynamicSharedInformerFactory
	client          cluster.ClientInterface
	stopCh          <-chan struct{}
}

var _ (Cache) = (*DynamicCache)(nil)

// NewDynamicCache creates an instance of DynamicCache.
func NewDynamicCache(client cluster.ClientInterface, stopCh <-chan struct{}, options ...DynamicCacheOpt) (*DynamicCache, error) {
	c := &DynamicCache{
		initFactoryFunc: initDynamicSharedInformerFactory,
		client:          client,
		stopCh:          stopCh,
	}

	for _, option := range options {
		option(c)
	}

	factory, err := c.initFactoryFunc(client)
	if err != nil {
		return nil, errors.Wrap(err, "initialize dynamic shared informer factory")
	}

	c.factory = factory

	fmt.Println("cache init")

	return c, nil
}

type lister interface {
	List(selector kLabels.Selector) ([]kruntime.Object, error)
}

func (dc *DynamicCache) currentInformer(key Key) (informers.GenericInformer, error) {
	gvk := schema.FromAPIVersionAndKind(key.APIVersion, key.Kind)

	gvr, err := dc.client.Resource(gvk.GroupKind())
	if err != nil {
		return nil, err
	}

	informer := dc.factory.ForResource(gvr)
	dc.factory.Start(dc.stopCh)

	return informer, nil
}

// List lists objects.
func (dc *DynamicCache) List(key Key) ([]*unstructured.Unstructured, error) {
	informer, err := dc.currentInformer(key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving informer for %v", key)
	}

	var l lister
	if key.Namespace == "" {
		l = informer.Lister()
	} else {
		l = informer.Lister().ByNamespace(key.Namespace)
	}

	var selector = kLabels.Everything()
	if key.Selector != nil {
		selector = key.Selector
	}

	objects, err := l.List(selector)
	if err != nil {
		return nil, errors.Wrapf(err, "listing %v", key)
	}

	list := make([]*unstructured.Unstructured, len(objects))
	for i, obj := range objects {
		u, err := kruntime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return nil, errors.Wrapf(err, "converting %T to unstructured", obj)
		}
		list[i] = &unstructured.Unstructured{Object: u}
	}

	return list, nil
}

type getter interface {
	Get(string) (kruntime.Object, error)
}

// Get retrieves a single object.
func (dc *DynamicCache) Get(key Key) (*unstructured.Unstructured, error) {
	informer, err := dc.currentInformer(key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving informer for %v", key)
	}

	var g getter
	if key.Namespace == "" {
		g = informer.Lister()
	} else {
		g = informer.Lister().ByNamespace(key.Namespace)
	}

	var object kruntime.Object
	retryErr := retry.Retry(3, time.Second, func() error {
		object, err = g.Get(key.Name)
		if err != nil {
			if !kerrors.IsNotFound(err) {
				return retry.Stop(errors.Wrap(err, "lister Get"))
			}
			return err
		}

		return nil
	})

	if retryErr != nil {
		return nil, err
	}

	// Verify the selector matches if provided
	if key.Selector != nil {
		accessor := meta.NewAccessor()
		m, err := accessor.Labels(object)
		if err != nil {
			return nil, errors.New("retrieving labels")
		}
		labels := kLabels.Set(m)
		if !key.Selector.Matches(labels) {
			return nil, errors.New("object found but filtered by selector")
		}
	}

	u, err := kruntime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return nil, errors.Wrapf(err, "converting %T to unstructured", object)
	}
	return &unstructured.Unstructured{Object: u}, nil
}

// Watch watches the cluster for an event and performs actions with the
// supplied handler.
func (dc *DynamicCache) Watch(key Key, handler kcache.ResourceEventHandler) error {
	informer, err := dc.currentInformer(key)
	if err != nil {
		return errors.Wrapf(err, "retrieving informer for %s", key)
	}

	informer.Informer().AddEventHandler(handler)
	return nil
}

package objectstore

import (
	"fmt"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	"github.com/vmware-tanzu/octant/internal/cluster"
)

//go:generate mockgen -destination=./fake/mock_informer_factory.go -package=fake github.com/vmware-tanzu/octant/internal/objectstore InformerFactory

// InformerFactory creates informers.
type InformerFactory interface {
	ForResource(gvr schema.GroupVersionKind) (informers.GenericInformer, error)
	Delete(gvr schema.GroupVersionKind)
	WaitForCacheSync(stopCh <-chan struct{}) map[schema.GroupVersionKind]bool
}

type informerFactory struct {
	client        cluster.ClientInterface
	defaultResync time.Duration
	namespace     string

	lock                 sync.Mutex
	informers            map[schema.GroupVersionKind]informers.GenericInformer
	informerErrors       map[schema.GroupVersionKind]error
	tweakListOptions     dynamicinformer.TweakListOptionsFunc
	stopCh               <-chan struct{}
	informerContextCache *informerContextCache
}

var _ InformerFactory = (*informerFactory)(nil)

func newInformerFactory(stopCh <-chan struct{}, client cluster.ClientInterface, defaultResync time.Duration, namespace string) *informerFactory {
	return &informerFactory{
		stopCh:               stopCh,
		client:               client,
		defaultResync:        defaultResync,
		namespace:            namespace,
		informers:            make(map[schema.GroupVersionKind]informers.GenericInformer),
		informerErrors:       make(map[schema.GroupVersionKind]error),
		informerContextCache: initInformerContextCache(),
	}
}

func (f *informerFactory) watchErrorHandler(gvk schema.GroupVersionKind, stopCh chan struct{}) cache.WatchErrorHandler {
	return func(r *cache.Reflector, err error) {
		f.lock.Lock()
		defer f.lock.Unlock()
		f.informerErrors[gvk] = err
		close(stopCh)
	}
}

// ForResource creates an informer and starts it given a group/version/resource.
func (f *informerFactory) ForResource(groupVersionKind schema.GroupVersionKind) (informers.GenericInformer, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	informer, exists := f.informers[groupVersionKind]
	if exists && informer != nil {
		return informer, nil
	}

	stopCh := f.informerContextCache.addChild(groupVersionKind)

	gvr, _, err := f.client.Resource(groupVersionKind.GroupKind())
	if err != nil {
		return nil, fmt.Errorf("unable to find group version resource for group kind %s: %w",
			groupVersionKind.GroupKind(), err)
	}

	dynamicClient, err := f.client.DynamicClient()
	if err != nil {
		return nil, fmt.Errorf("get dynamic client: %w", err)
	}

	genericInformer := dynamicinformer.NewFilteredDynamicInformer(
		dynamicClient,
		gvr,
		f.namespace,
		f.defaultResync,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
		f.tweakListOptions)
	f.informers[groupVersionKind] = genericInformer

	genericInformer.Informer().SetWatchErrorHandler(f.watchErrorHandler(groupVersionKind, stopCh))
	go genericInformer.Informer().Run(stopCh)

	err, exists = f.informerErrors[groupVersionKind]
	if exists && err != nil {
		return genericInformer, err
	}

	return genericInformer, nil
}

// Delete deletes an informer given a a group/version/resource.
func (f *informerFactory) Delete(groupVersionKind schema.GroupVersionKind) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.informerContextCache.delete(groupVersionKind)
	delete(f.informers, groupVersionKind)
	f.informers[groupVersionKind] = nil
}

// WaitForCacheSync waits for all started informers' cache were synced.
func (f *informerFactory) WaitForCacheSync(stopCh <-chan struct{}) map[schema.GroupVersionKind]bool {
	list := func() map[schema.GroupVersionKind]cache.SharedIndexInformer {
		f.lock.Lock()
		defer f.lock.Unlock()

		shared := map[schema.GroupVersionKind]cache.SharedIndexInformer{}
		for informerType, informer := range f.informers {
			shared[informerType] = informer.Informer()
		}
		return shared
	}()

	res := map[schema.GroupVersionKind]bool{}
	for informerType, informer := range list {
		res[informerType] = cache.WaitForCacheSync(stopCh, informer.HasSynced)
	}
	return res
}

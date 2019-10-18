package objectstore

import (
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

//go:generate mockgen -destination=./fake/mock_informer_factory.go -package=fake github.com/vmware-tanzu/octant/internal/objectstore InformerFactory

// InformerFactory creates informers.
type InformerFactory interface {
	ForResource(gvr schema.GroupVersionResource) informers.GenericInformer
	Delete(gvr schema.GroupVersionResource)
	WaitForCacheSync(stopCh <-chan struct{}) map[schema.GroupVersionResource]bool
}

type informerFactory struct {
	client        dynamic.Interface
	defaultResync time.Duration
	namespace     string

	lock                 sync.Mutex
	informers            map[schema.GroupVersionResource]informers.GenericInformer
	tweakListOptions     dynamicinformer.TweakListOptionsFunc
	stopCh               <-chan struct{}
	informerContextCache *informerContextCache
}

var _ InformerFactory = (*informerFactory)(nil)

func newInformerFactory(stopCh <-chan struct{}, client dynamic.Interface, defaultResync time.Duration, namespace string) *informerFactory {
	return &informerFactory{
		stopCh:               stopCh,
		client:               client,
		defaultResync:        defaultResync,
		namespace:            namespace,
		informers:            map[schema.GroupVersionResource]informers.GenericInformer{},
		informerContextCache: initInformerContextCache(),
	}
}

// ForResource creates an informer and starts it given a group/version/resource.
func (f *informerFactory) ForResource(gvr schema.GroupVersionResource) informers.GenericInformer {
	f.lock.Lock()
	defer f.lock.Unlock()

	key := gvr
	informer, exists := f.informers[key]
	if exists && informer != nil {
		return informer
	}

	informer = dynamicinformer.NewFilteredDynamicInformer(f.client, gvr, f.namespace, f.defaultResync, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
	f.informers[key] = informer

	stopCh := f.informerContextCache.addChild(gvr)
	go informer.Informer().Run(stopCh)

	return informer
}

// Delete deletes an informer given a a group/version/resource.
func (f *informerFactory) Delete(gvr schema.GroupVersionResource) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if _, ok := f.informers[gvr]; ok {
		f.informerContextCache.delete(gvr)
		delete(f.informers, gvr)
		f.informers[gvr] = nil
	}
}

// WaitForCacheSync waits for all started informers' cache were synced.
func (f *informerFactory) WaitForCacheSync(stopCh <-chan struct{}) map[schema.GroupVersionResource]bool {
	list := func() map[schema.GroupVersionResource]cache.SharedIndexInformer {
		f.lock.Lock()
		defer f.lock.Unlock()

		shared := map[schema.GroupVersionResource]cache.SharedIndexInformer{}
		for informerType, informer := range f.informers {
			shared[informerType] = informer.Informer()
		}
		return shared
	}()

	res := map[schema.GroupVersionResource]bool{}
	for informType, informer := range list {
		res[informType] = cache.WaitForCacheSync(stopCh, informer.HasSynced)
	}
	return res
}

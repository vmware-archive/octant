package objectstore

import (
	"context"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/vmware-tanzu/octant/internal/cluster"
)

type Watcher struct {
	ctx    context.Context
	client cluster.ClientInterface
	cache  *ResourceCache

	watches *sync.Map
}

func NewWatcher(ctx context.Context, client cluster.ClientInterface, cache *ResourceCache) *Watcher {
	return &Watcher{
		ctx:     ctx,
		client:  client,
		cache:   cache,
		watches: &sync.Map{},
	}
}

func (w *Watcher) StopAll() {
	w.watches.Range(func(k interface{}, v interface{}) bool {
		watch, _ := v.(watch.Interface)
		watch.Stop()
		return true
	})
}

func (w *Watcher) Stop(cacheKey ResourceCacheKey) {
	v, ok := w.watches.Load(cacheKey)
	if ok {
		resWatch, _ := v.(watch.Interface)
		resWatch.Stop()
	}
}

func (w *Watcher) Watch(cacheKey ResourceCacheKey) (bool, error) {
	if _, ok := w.watches.Load(cacheKey); ok {
		return true, nil
	}

	dc, err := w.client.DynamicClient()
	if err != nil {
		return false, err
	}

	listOptions := metav1.ListOptions{}

	var resWatch watch.Interface

	if cacheKey.Namespace == "" {
		resWatch, err = dc.Resource(cacheKey.Resource).Watch(w.ctx, listOptions)
		if err != nil {
			return false, err
		}
	}
	resWatch, err = dc.Resource(cacheKey.Resource).Namespace(cacheKey.Namespace).Watch(w.ctx, listOptions)
	if err != nil {
		return false, err
	}

	w.watches.Store(cacheKey, resWatch)

	go w.handleWatch(cacheKey, resWatch)

	return true, nil
}

func (w *Watcher) handleWatch(cacheKey ResourceCacheKey, resWatch watch.Interface) {
	for event := range resWatch.ResultChan() {
		switch event.Type {
		case watch.Added:
			u, _ := ToUnstructured(event.Object)
			w.cache.Add(cacheKey, u)
		case watch.Modified:
			u, _ := ToUnstructured(event.Object)
			w.cache.Add(cacheKey, u)
		case watch.Deleted:
			u, _ := ToUnstructured(event.Object)
			w.cache.Delete(cacheKey, u)
		default:
			continue
		}
	}
}

func ToUnstructured(object runtime.Object) (unstructured.Unstructured, error) {
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return unstructured.Unstructured{}, err
	}
	return unstructured.Unstructured{Object: unstructuredObj}, nil
}

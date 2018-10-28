package overview

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/davecgh/go-spew/spew"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

// StopFunc tells a watch to stop watching a namespace.
type StopFunc func()

// Watch watches a objects in a namespace.
type Watch interface {
	Start() (StopFunc, error)
}

// ClusterWatch watches a namespace's objects.
type ClusterWatch struct {
	clusterClient cluster.ClientInterface
	cache         Cache
	namespace     string
	logger        log.Logger
}

// NewWatch creates an instance of Watch.
func NewWatch(namespace string, clusterClient cluster.ClientInterface, c Cache, logger log.Logger) *ClusterWatch {
	return &ClusterWatch{
		namespace:     namespace,
		clusterClient: clusterClient,
		cache:         c,
		logger:        logger,
	}
}

// Start starts the watch. It returns a stop function and an error.
func (w *ClusterWatch) Start() (StopFunc, error) {
	resources, err := w.resources()
	if err != nil {
		return nil, err
	}

	var watchers []watch.Interface

	for _, resource := range resources {
		dc, err := w.clusterClient.DynamicClient()
		if err != nil {
			return nil, err
		}

		nri := dc.Resource(resource).Namespace(w.namespace)

		watcher, err := nri.Watch(metav1.ListOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "did not create watcher for %s/%s/%s on %s namespace", resource.Group, resource.Version, resource.Resource, w.namespace)
		}

		watchers = append(watchers, watcher)
	}

	done := make(chan struct{})

	events, shutdownCh := consumeEvents(done, watchers)

	allDone := make(chan interface{})
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		// Forward events to handler.
		// Exits after all watchers are stopped in consumeEvents.
		for event := range events {
			w.eventHandler(event)
		}
		wg.Done()
	}()

	go func() {
		// Block until all watch consumers have finished (in consumeEvents)
		<-shutdownCh
		wg.Done()
	}()

	go func() {
		// Block until fan-in consumer as well as individual watch consumers have completed
		// (above two goroutines)
		wg.Wait()
		close(allDone)
	}()

	stopFn := func() {
		// Signal consumer routines to shutdown. Block until all have finished.
		done <- struct{}{}
		<-allDone
	}

	return stopFn, nil
}

func (w *ClusterWatch) eventHandler(event watch.Event) {
	u, ok := event.Object.(*unstructured.Unstructured)
	if !ok {
		return
	}

	switch t := event.Type; t {
	case watch.Added:
		if err := w.cache.Store(u); err != nil {
			w.logger.Errorf("store object: %v", err)
		}
	case watch.Modified:
		if err := w.cache.Store(u); err != nil {
			w.logger.Errorf("store object: %v", err)
		}
	case watch.Deleted:
		if err := w.cache.Delete(u); err != nil {
			w.logger.Errorf("store object: %v", err)
		}
	case watch.Error:
		w.logger.Errorf("unknown log err: %s", spew.Sdump(event))
	default:
		w.logger.Errorf("unknown event %q", t)
	}
}

func (w *ClusterWatch) resources() ([]schema.GroupVersionResource, error) {
	discoveryClient, err := w.clusterClient.DiscoveryClient()
	if err != nil {
		return nil, err
	}

	// NOTE we may want ServerPreferredResources, but FakeDiscovery does not support it.
	lists, err := discoveryClient.ServerResources()
	if err != nil {
		return nil, err
	}

	var gvrs []schema.GroupVersionResource

	for _, list := range lists {
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			return nil, err
		}

		for _, res := range list.APIResources {
			if !res.Namespaced {
				continue
			}
			if isWatchable(res) {

				gvr := schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: res.Name,
				}

				gvrs = append(gvrs, gvr)

			}
		}
	}

	return gvrs, nil
}

func isWatchable(res metav1.APIResource) bool {
	m := make(map[string]bool)

	for _, v := range res.Verbs {
		m[v] = true
	}

	return m["list"] && m["watch"]
}

// consumeEvents performs fan-in of events from multiple watchers into a single event channel.
// This continues until a message is sent on the provided done channel.
func consumeEvents(done <-chan struct{}, watchers []watch.Interface) (chan watch.Event, chan struct{}) {
	var wg sync.WaitGroup

	wg.Add(len(watchers))

	events := make(chan watch.Event)

	shutdownComplete := make(chan struct{})

	for _, watcher := range watchers {
		// Forward events from each watcher to events channel.
		// Each drainer goroutine ends when its watcher's Stop() method is called,
		// which will have the effect of closing its ResultChan and exiting the range loop.
		go func(watcher watch.Interface) {
			for event := range watcher.ResultChan() {
				events <- event
			}
			wg.Done()
		}(watcher)
	}

	go func() {
		// wait for caller to signal done and
		// start shutting the watcher down
		<-done
		for _, watch := range watchers {
			watch.Stop()
		}
	}()

	go func() {
		// wait for all watchers to exit.
		wg.Wait()
		close(events)
		shutdownComplete <- struct{}{}
	}()

	return events, shutdownComplete
}

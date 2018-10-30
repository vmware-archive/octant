package overview

import (
	"fmt"
	"sync"
	"time"

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
	watchers      []watch.Interface
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

	ec := newEventConsumer(w.logger)

	go func() {
		start := time.Now()
	outer:
		for _, resource := range resources {
			select {
			case <-ec.Done():
				w.logger.Debugf("aborting watch construction loop due to Done event: %p", ec)
				break outer
			default:
			}

			dc, err := w.clusterClient.DynamicClient()
			if err != nil {
				w.logger.With("resource", resource).Errorf("creating dynamicClient: %v", err)
				return
			}

			nri := dc.Resource(resource).Namespace(w.namespace)

			// watchStart := time.Now()
			watcher, err := nri.Watch(metav1.ListOptions{})
			if err != nil {
				w.logger.Errorf("%v", errors.Wrapf(err, "did not create watcher for %s/%s/%s on %s namespace", resource.Group, resource.Version, resource.Resource, w.namespace))
				return
			}
			// w.logger.With("duration", int(time.Since(watchStart)/time.Millisecond), "resource", resource.Resource).Debugf("individual watch")

			if err := ec.Consume([]watch.Interface{watcher}); err != nil {
				// Ownership of the watcher was not taken, we must stop it ourselves.
				watcher.Stop()
			}
		}
		w.logger.With("duration", int(time.Since(start)/time.Millisecond), "watchers", len(resources)).Debugf("top-level watch")
	}()

	forwarderDone := make(chan struct{})
	go func() {
		// Forward events to handler.
		// Exits after all watchers are stopped in consumeEvents.
		w.logger.With("context", "event forwarder").Debugf("started")
		for event := range ec.Events() {
			w.eventHandler(event)
		}
		close(forwarderDone)
		w.logger.With("context", "event forwarder").Debugf("stopped")
	}()

	stopFn := func() {
		// Signal consumer routines to shutdown. Block until all have finished.
		ec.Stop()
		<-forwarderDone
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
	start := time.Now()

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

	w.logger.With("duration", int(time.Since(start)/time.Millisecond), "resources", len(gvrs)).Debugf("resources from dynamic client")

	return gvrs, nil
}

func isWatchable(res metav1.APIResource) bool {
	m := make(map[string]bool)

	for _, v := range res.Verbs {
		m[v] = true
	}

	return m["list"] && m["watch"]
}

// eventConsumer perform fan-in from many watchers into a single event channel
type eventConsumer struct {
	wg       sync.WaitGroup
	mu       sync.Mutex
	logger   log.Logger
	watchers []watch.Interface
	events   chan watch.Event
	once     sync.Once
	done     chan struct{}
}

// newEventConsumer creates a new event consumer. Event consumers perform fan-in from many watchers
// into a single event channel.
func newEventConsumer(logger log.Logger) *eventConsumer {
	return &eventConsumer{
		logger:   logger,
		watchers: make([]watch.Interface, 0),
		events:   make(chan watch.Event),
		done:     make(chan struct{}),
	}
}

// Stop stops all watchers managed by the eventConsumer.
// The events channel (ec.Events()) and done channel (ec.Done()) will be closed to notify consumers of the state change.
func (ec *eventConsumer) Stop() {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	ec.once.Do(func() {
		// Notify clients to stop adding watches (loops calling Consume())
		close(ec.done)
	})

	for _, watch := range ec.watchers {
		watch.Stop()
	}
	ec.logger.Debugf("stopped %d watchers", len(ec.watchers))
	ec.watchers = ec.watchers[:0]
	ec.wg.Wait()
	close(ec.events)
}

// Events returns a channel that can be used to consume events from all watchers
func (ec *eventConsumer) Events() <-chan watch.Event {
	return ec.events
}

// Done returns a channel that when closed indicates the eventConsumer is either
// in the process of shutting down, or has already been shut down.
func (ec *eventConsumer) Done() <-chan struct{} {
	return ec.done
}

var errInterrupted = fmt.Errorf("interrupted")

// Consume will begin consuming events from the provided watchers.
// It is safe to call this method multiple with additional watchers,
// which will be added to the watch consumer list.
// If a non-nil error is returned, ownership of watchers has not been transferred,
// and in that case the caller will be responsible for stopping them themselves.
func (ec *eventConsumer) Consume(watchers []watch.Interface) error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	select {
	case <-ec.done:
		// Stop() has been called. Do not accept ownership of watchers.
		return errInterrupted
	default:
	}

	for _, watcher := range watchers {
		ec.watchers = append(ec.watchers, watcher)
		ec.wg.Add(1)

		// Forward events from each watcher to events channel.
		// Each drainer goroutine ends when its watcher's Stop() method is called,
		// which will have the effect of closing its ResultChan and exiting the range loop.
		go func(watcher watch.Interface) {
			for event := range watcher.ResultChan() {
				ec.events <- event
			}
			ec.wg.Done()
		}(watcher)
	}

	return nil
}

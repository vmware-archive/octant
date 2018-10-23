package overview

import (
	"fmt"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/log"
)

// ClusterOverview is an API for generating a cluster overview.
type ClusterOverview struct {
	client cluster.ClientInterface

	namespace string

	logger log.Logger

	watchFactory func(namespace string, clusterClient cluster.ClientInterface, cache Cache) Watch

	cache  Cache
	stopFn func()

	generator *realGenerator
}

// NewClusterOverview creates an instance of ClusterOverview.
func NewClusterOverview(client cluster.ClientInterface, namespace string, logger log.Logger) *ClusterOverview {
	var opts []MemoryCacheOpt

	if os.Getenv("DASH_VERBOSE_CACHE") != "" {
		ch := make(chan CacheNotification)

		go func() {
			for notif := range ch {
				spew.Dump(notif)
			}
		}()

		opts = append(opts, CacheNotificationOpt(ch))
	}

	cache := NewMemoryCache(opts...)

	var pathFilters []pathFilter
	pathFilters = append(pathFilters, rootDescriber.PathFilters()...)
	pathFilters = append(pathFilters, eventsDescriber.PathFilters()...)

	g := newGenerator(cache, pathFilters, client)

	co := &ClusterOverview{
		namespace: namespace,
		client:    client,
		logger:    logger,
		cache:     cache,
		generator: g,
	}
	co.watchFactory = co.defaultWatchFactory
	return co
}

// Name returns the name for this module.
func (co *ClusterOverview) Name() string {
	return "overview"
}

// ContentPath returns the content path for overview.
func (co *ClusterOverview) ContentPath() string {
	return fmt.Sprintf("/%s", co.Name())
}

// Handler returns a handler for serving overview HTTP content.
func (co *ClusterOverview) Handler(prefix string) http.Handler {
	return newHandler(prefix, co.generator, stream, co.logger)
}

// Namespaces returns a list of namespace names for a cluster.
func (co *ClusterOverview) Namespaces() ([]string, error) {
	nsClient, err := co.client.NamespaceClient()
	if err != nil {
		return nil, err
	}

	return nsClient.Names()
}

// Navigation returns navigation entries for overview.
func (co *ClusterOverview) Navigation(root string) (*hcli.Navigation, error) {
	return navigationEntries(root)
}

// SetNamespace sets the current namespace.
func (co *ClusterOverview) SetNamespace(namespace string) error {
	co.logger.Debugf("setting namespace for overview to %q", namespace)
	if co.stopFn != nil {
		co.stopFn()
	}

	co.namespace = namespace
	return co.Start()
}

// Start starts overview.
func (co *ClusterOverview) Start() error {
	if co.namespace == "" {
		return nil
	}

	co.logger.Debugf("starting cluster overview")

	stopFn, err := co.watch(co.namespace)
	if err != nil {
		return err
	}

	co.stopFn = stopFn

	return nil
}

// Stop stops overview.
func (co *ClusterOverview) Stop() {
	if co.stopFn != nil {
		co.logger.Debugf("stopping cluster overview")

		co.stopFn()
	}
}

func (co *ClusterOverview) watch(namespace string) (StopFunc, error) {
	co.logger.Debugf("watching namespace %s", namespace)

	watch := co.watchFactory(namespace, co.client, co.cache)
	return watch.Start()
}

func (co *ClusterOverview) defaultWatchFactory(namespace string, clusterClient cluster.ClientInterface, cache Cache) Watch {
	return NewWatch(namespace, clusterClient, cache, co.logger)
}

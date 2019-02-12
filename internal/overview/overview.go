package overview

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/overview/container"
	"github.com/heptio/developer-dash/internal/queryer"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/view/component"

	"github.com/davecgh/go-spew/spew"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/pkg/errors"
	"k8s.io/client-go/restmapper"
)

// ClusterOverview is an API for generating a cluster overview.
type ClusterOverview struct {
	client cluster.ClientInterface

	mu sync.Mutex

	logger log.Logger

	cache  cache.Cache
	stopCh chan struct{}

	generator *realGenerator
}

// NewClusterOverview creates an instance of ClusterOverview.
func NewClusterOverview(client cluster.ClientInterface, namespace string, logger log.Logger) (*ClusterOverview, error) {
	stopCh := make(chan struct{})

	var opts []cache.InformerCacheOpt

	if os.Getenv("DASH_VERBOSE_CACHE") != "" {
		ch := make(chan cache.Notification)

		go func() {
			for notif := range ch {
				spew.Dump(notif)
			}
		}()

		opts = append(opts, cache.InformerCacheNotificationOpt(ch, stopCh))
	}

	if client == nil {
		return nil, errors.New("nil cluster client")
	}

	dynamicClient, err := client.DynamicClient()
	if err != nil {
		return nil, errors.Wrapf(err, "creating DynamicClient")
	}
	di, err := client.DiscoveryClient()
	if err != nil {
		return nil, errors.Wrapf(err, "creating DiscoveryClient")
	}

	groupResources, err := restmapper.GetAPIGroupResources(di)
	if err != nil {
		logger.Errorf("discovering APIGroupResources: %v", err)
		return nil, errors.Wrapf(err, "mapping APIGroupResources")
	}
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	opts = append(opts, cache.InformerCacheLoggerOpt(logger))
	informerCache := cache.NewInformerCache(stopCh, dynamicClient, rm, opts...)

	var pathFilters []pathFilter
	pathFilters = append(pathFilters, rootDescriber.PathFilters(namespace)...)
	pathFilters = append(pathFilters, eventsDescriber.PathFilters(namespace)...)

	queryer := queryer.New(informerCache, di)

	g, err := newGenerator(informerCache, queryer, pathFilters, client)
	if err != nil {
		return nil, errors.Wrap(err, "create overview generator")
	}

	co := &ClusterOverview{
		client:    client,
		logger:    logger,
		cache:     informerCache,
		generator: g,
		stopCh:    stopCh,
	}
	return co, nil
}

// Name returns the name for this module.
func (co *ClusterOverview) Name() string {
	return "overview"
}

// ContentPath returns the content path for overview.
func (co *ClusterOverview) ContentPath() string {
	return fmt.Sprintf("/%s", co.Name())
}

// Navigation returns navigation entries for overview.
func (co *ClusterOverview) Navigation(namespace, root string) (*hcli.Navigation, error) {
	nf := NewNavigationFactory(namespace, root)
	return nf.Entries()
}

// SetNamespace sets the current namespace.
func (co *ClusterOverview) SetNamespace(namespace string) error {
	co.logger.With("namespace", namespace, "module", "overview").Debugf("setting namespace (noop)")
	return nil
}

// Start starts overview.
func (co *ClusterOverview) Start() error {
	return nil
}

// Stop stops overview.
func (co *ClusterOverview) Stop() {
	co.mu.Lock()
	defer co.mu.Unlock()
	close(co.stopCh)
	co.stopCh = nil
}

func (co *ClusterOverview) Content(ctx context.Context, contentPath, prefix, namespace string) (component.ContentResponse, error) {
	return co.generator.Generate(ctx, contentPath, prefix, namespace)
}

type logEntry struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Message   string    `json:"message,omitempty"`
}

type logResponse struct {
	Entries []logEntry `json:"entries,omitempty"`
}

func (co *ClusterOverview) Handlers() map[string]http.Handler {
	return map[string]http.Handler{
		"/logs/pod/{pod}/container/{container}": containerLogsHandler(co.client),
	}
}

func containerLogsHandler(clusterClient cluster.ClientInterface) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		containerName := vars["container"]
		podName := vars["pod"]
		namespace := vars["namespace"]

		kubeClient, err := clusterClient.KubernetesClient()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lines := make(chan string)
		done := make(chan bool)

		var entries []logEntry

		go func() {
			for line := range lines {
				parts := strings.SplitN(line, " ", 2)
				logTime, err := time.Parse(time.RFC3339, parts[0])
				if err == nil {
					entries = append(entries, logEntry{
						Timestamp: logTime,
						Message:   parts[1],
					})
				}
			}

			done <- true
		}()

		err = container.Logs(r.Context(), kubeClient, namespace, podName, containerName, lines)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		<-done

		var lr logResponse

		if len(entries) <= 100 {
			lr.Entries = entries
		} else {
			// take last 100 entries from the slice
			lr.Entries = entries[len(entries)-100:]
		}

		if err := json.NewEncoder(w).Encode(&lr); err != nil {
			fmt.Println("oops", err)
		}
	}
}

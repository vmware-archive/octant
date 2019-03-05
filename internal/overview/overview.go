package overview

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/mime"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/overview/container"
	"github.com/heptio/developer-dash/internal/portforward"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/view/component"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/pkg/errors"
)

// ClusterOverview is an API for generating a cluster overview.
type ClusterOverview struct {
	mu sync.Mutex

	client         cluster.ClientInterface
	logger         log.Logger
	cache          cache.Cache
	generator      *realGenerator
	portForwardSvc portforward.PortForwardInterface
}

// NewClusterOverview creates an instance of ClusterOverview.
// TODO: why does cache get passed in here?
func NewClusterOverview(ctx context.Context, client cluster.ClientInterface, c cache.Cache, namespace string, logger log.Logger) (*ClusterOverview, error) {
	if client == nil {
		return nil, errors.New("nil cluster client")
	}

	di, err := client.DiscoveryClient()
	if err != nil {
		return nil, errors.Wrapf(err, "creating DiscoveryClient")
	}

	informerCache, err := cache.NewDynamicCache(client, ctx.Done())
	if err != nil {
		return nil, errors.Wrapf(err, "create cache")
	}

	pm := newPathMatcher()
	for _, pf := range rootDescriber.PathFilters() {
		pm.Register(ctx, pf)
	}

	for _, pf := range eventsDescriber.PathFilters() {
		pm.Register(ctx, pf)
	}

	crdAddFunc := func(pm *pathMatcher, csd *crdSectionDescriber) objectHandler {
		return func(ctx context.Context, object *unstructured.Unstructured) {
			if object == nil {
				return
			}
			addCRD(ctx, object.GetName(), pm, csd)
		}
	}(pm, customResourcesDescriber)
	crdDeleteFunc := func(pm *pathMatcher, csd *crdSectionDescriber) objectHandler {
		return func(ctx context.Context, object *unstructured.Unstructured) {
			if object == nil {
				return
			}
			deleteCRD(ctx, object.GetName(), pm, csd)
		}
	}(pm, customResourcesDescriber)

	go watchCRDs(ctx, informerCache, crdAddFunc, crdDeleteFunc)

	// Port Forwarding
	restClient, err := client.RESTClient()
	if err != nil {
		return nil, errors.Wrap(err, "fetching RESTClient")
	}
	pfOpts := portforward.PortForwardSvcOptions{
		RESTClient: restClient,
		Config:     client.RESTConfig(),
		Cache:      c,
		// TODO -  streams
		PortForwarder: &portforward.DefaultPortForwarder{
			IOStreams: portforward.IOStreams{
				In:     os.Stdin,
				Out:    os.Stdout,
				ErrOut: os.Stderr,
			},
		},
	}
	pfSvc := portforward.NewPortForwardService(ctx, pfOpts, logger)

	g, err := newGenerator(c, di, pm, client, pfSvc)
	if err != nil {
		return nil, errors.Wrap(err, "create overview generator")
	}

	co := &ClusterOverview{
		client:         client,
		logger:         logger,
		cache:          c,
		generator:      g,
		portForwardSvc: pfSvc,
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
	nf := NewNavigationFactory(namespace, root, co.cache)
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
	// NOOP
}

// Content serves content for overview.
func (co *ClusterOverview) Content(ctx context.Context, contentPath, prefix, namespace string, opts module.ContentOptions) (component.ContentResponse, error) {
	ctx = log.WithLoggerContext(ctx, co.logger)
	genOpts := GeneratorOptions{
		Selector:       opts.Selector,
		PortForwardSvc: co.portForwardSvc,
	}
	return co.generator.Generate(ctx, contentPath, prefix, namespace, genOpts)
}

type logEntry struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Message   string    `json:"message,omitempty"`
}

type logResponse struct {
	Entries []logEntry `json:"entries,omitempty"`
}

// Handlers are extra handlers for overview
func (co *ClusterOverview) Handlers() map[string]http.Handler {
	return map[string]http.Handler{
		"/logs/pod/{pod}/container/{container}":                   containerLogsHandler(co.client),
		"/portforward/create/pod/{pod}/port/{port}":               co.portForwardHandler(),
		"/portforward/create/service/{service}/port/{port}":       co.portForwardHandler(),
		"/portforward/create/deployment/{deployment}/port/{port}": co.portForwardHandler(),
		"/portforward/delete/{id}":                                co.portForwardDeleteHandler(),
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

func (co *ClusterOverview) portForwardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if co.portForwardSvc == nil {
			http.Error(w, "portforward service is nil", http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)

		podName := vars["pod"]
		namespace := vars["namespace"]
		// serviceName := vars["service"]
		// deploymentName := vars["deployment"]
		portStr := vars["port"]

		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			http.Error(w, errors.Wrapf(err, "invalid port").Error(), http.StatusInternalServerError)
			return
		}

		var resp portforward.PortForwardCreateResponse
		switch {
		case podName != "":
			gvk := schema.GroupVersionKind{
				Group:   "",
				Version: "v1",
				Kind:    "Pod",
			}
			resp, err = co.portForwardSvc.Create(gvk, podName, namespace, uint16(port))
			if err != nil {
				http.Error(w, errors.Wrapf(err, "creating forwarder for pod").Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", mime.JSONContentType)
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			co.logger.Errorf("encoding JSON response: %v", err)
		}
	}
}

func (co *ClusterOverview) portForwardDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if co.portForwardSvc == nil {
			http.Error(w, "portforward service is nil", http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)

		id := vars["id"]

		co.logger.Debugf("Stopping port forwarder %s", id)
		co.portForwardSvc.StopForwarder(id)

		w.WriteHeader(http.StatusOK)
		http.Redirect(w, r, "/content/overview/portforward", 302)
	}
}

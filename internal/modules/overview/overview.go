package overview

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/heptio/developer-dash/internal/config"
	"github.com/heptio/developer-dash/internal/describer"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/component"
)

type Options struct {
	Namespace  string
	DashConfig config.Dash
}

// Overview is an API for generating a cluster overview.
type Overview struct {
	*clustereye.ObjectPath

	generator  *realGenerator
	dashConfig config.Dash
}

var _ module.Module = (*Overview)(nil)

// New creates an instance of Overview.
func New(ctx context.Context, options Options) (*Overview, error) {
	if options.DashConfig == nil {
		return nil, errors.New("dash configuration is nil")
	}

	if err := options.DashConfig.Validate(); err != nil {
		return nil, errors.Wrap(err, "dash configuration")
	}

	pathMatcher := describer.NewPathMatcher("overview")
	for _, pf := range rootDescriber.PathFilters() {
		pathMatcher.Register(ctx, pf)
	}

	for _, pf := range eventsDescriber.PathFilters() {
		pathMatcher.Register(ctx, pf)
	}

	key := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	objectStore := options.DashConfig.ObjectStore()

	g, err := newGenerator(pathMatcher, options.DashConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create overview generator")
	}

	objectPathConfig := clustereye.ObjectPathConfig{
		ModuleName:     "overview",
		SupportedGVKs:  supportedGVKs,
		PathLookupFunc: gvkPath,
		CRDPathGenFunc: crdPath,
	}
	objectPath, err := clustereye.NewObjectPath(objectPathConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create module object path generator")
	}

	co := &Overview{
		ObjectPath: objectPath,
		generator:  g,
		dashConfig: options.DashConfig,
	}

	crdWatcher := options.DashConfig.CRDWatcher()
	if err := objectStore.HasAccess(key, "watch"); err == nil {
		watchConfig := &config.CRDWatchConfig{
			Add: func(_ *describer.PathMatcher, sectionDescriber *describer.CRDSection) config.ObjectHandler {
				return func(ctx context.Context, object *unstructured.Unstructured) {
					if object == nil {
						return
					}
					describer.AddCRD(ctx, object, pathMatcher, customResourcesDescriber, co)
				}
			}(pathMatcher, customResourcesDescriber),
			Delete: func(_ *describer.PathMatcher, csd *describer.CRDSection) config.ObjectHandler {
				return func(ctx context.Context, object *unstructured.Unstructured) {
					if object == nil {
						return
					}
					describer.DeleteCRD(ctx, object, pathMatcher, customResourcesDescriber, co)
				}
			}(pathMatcher, customResourcesDescriber),
			IsNamespaced: true,
		}

		if err := crdWatcher.Watch(ctx, watchConfig); err != nil {
			return nil, errors.Wrap(err, "create namespaced CRD watcher for overview")
		}
	}

	return co, nil
}

// Name returns the name for this module.
func (co *Overview) Name() string {
	return "overview"
}

// ContentPath returns the content path for overview.
func (co *Overview) ContentPath() string {
	return fmt.Sprintf("/%s", co.Name())
}

// Navigation returns navigation entries for overview.
func (co *Overview) Navigation(ctx context.Context, namespace, root string) ([]clustereye.Navigation, error) {
	navigationEntries := clustereye.NavigationEntries{
		Lookup: navPathLookup,
		EntriesFuncs: map[string]clustereye.EntriesFunc{
			"Workloads":                    workloadEntries,
			"Discovery and Load Balancing": discoAndLBEntries,
			"Config and Storage":           configAndStorageEntries,
			"Custom Resources":             clustereye.CRDEntries,
			"RBAC":                         rbacEntries,
			"Events":                       nil,
		},
		Order: []string{
			"Workloads",
			"Discovery and Load Balancing",
			"Config and Storage",
			"Custom Resources",
			"RBAC",
			"Events",
		},
	}

	objectStore := co.dashConfig.ObjectStore()

	nf := clustereye.NewNavigationFactory(namespace, root, objectStore, navigationEntries)

	entries, err := nf.Generate(ctx, "Overview")
	if err != nil {
		return nil, err
	}

	return []clustereye.Navigation{
		*entries,
	}, nil
}

// Generators allow modules to send events to the frontend.
func (co *Overview) Generators() []clustereye.Generator {
	return []clustereye.Generator{}
}

// SetNamespace sets the current namespace.
func (co *Overview) SetNamespace(namespace string) error {
	co.dashConfig.Logger().With("namespace", namespace, "module", "overview").Debugf("setting namespace (noop)")
	return nil
}

// Start starts overview.
func (co *Overview) Start() error {
	return nil
}

// Stop stops overview.
func (co *Overview) Stop() {
	// NOOP
}

// Content serves content for overview.
func (co *Overview) Content(ctx context.Context, contentPath, prefix, namespace string, opts module.ContentOptions) (component.ContentResponse, error) {
	ctx = log.WithLoggerContext(ctx, co.dashConfig.Logger())
	genOpts := GeneratorOptions{
		LabelSet: opts.LabelSet,
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
func (co *Overview) Handlers(ctx context.Context) map[string]http.Handler {
	return map[string]http.Handler{
		"/logs/pod/{pod}/container/{container}": containerLogsHandler(ctx, co.dashConfig.ClusterClient()),
		"/port-forwards":                        co.portForwardsHandler(),
		"/port-forwards/{id}":                   co.portForwardHandler(),
	}
}

func (co *Overview) portForwardsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := co.dashConfig.PortForwarder()
		logger := co.dashConfig.Logger()

		if svc == nil {
			logger.Errorf("port forward service is nil")
			http.Error(w, "port forward service is nil", http.StatusInternalServerError)
			return
		}

		ctx := log.WithLoggerContext(r.Context(), logger)

		defer func() {
			if cErr := r.Body.Close(); cErr != nil {
				logger.With("err", cErr).Errorf("unable to close port forward request body")
			}
		}()

		switch r.Method {
		case http.MethodPost:
			err := createPortForward(ctx, r.Body, svc, w)
			handlePortForwardError(w, err, logger)
		default:
			api.RespondWithError(
				w,
				http.StatusNotFound,
				fmt.Sprintf("unhandled HTTP method %s", r.Method),
				logger,
			)
		}
	}
}

func (co *Overview) portForwardHandler() http.HandlerFunc {
	logger := co.dashConfig.Logger()

	return func(w http.ResponseWriter, r *http.Request) {
		svc := co.dashConfig.PortForwarder()
		if svc == nil {
			logger.Errorf("port forward service is nil")
			http.Error(w, "port forward service is nil", http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		ctx := log.WithLoggerContext(r.Context(), logger)

		switch r.Method {
		case http.MethodDelete:
			err := deletePortForward(ctx, id, co.dashConfig.PortForwarder(), w)
			handlePortForwardError(w, err, logger)
		default:
			api.RespondWithError(
				w,
				http.StatusNotFound,
				fmt.Sprintf("unhandled HTTP method %s", r.Method),
				logger,
			)
		}
	}
}

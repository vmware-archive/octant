/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware/octant/internal/action"
	"github.com/vmware/octant/internal/api"
	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/describer"
	"github.com/vmware/octant/internal/icon"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

type Options struct {
	Namespace  string
	DashConfig config.Dash
}

// Overview is an API for generating a cluster overview.
type Overview struct {
	*octant.ObjectPath

	generator   *realGenerator
	dashConfig  config.Dash
	contextName string
	logger      log.Logger

	mu sync.Mutex
}

var _ module.Module = (*Overview)(nil)
var _ module.ActionReceiver = (*Overview)(nil)

// New creates an instance of Overview.
func New(ctx context.Context, options Options) (*Overview, error) {
	if options.DashConfig == nil {
		return nil, errors.New("dash configuration is nil")
	}

	if err := options.DashConfig.Validate(); err != nil {
		return nil, errors.Wrap(err, "dash configuration")
	}

	co := &Overview{
		dashConfig: options.DashConfig,
		logger:     options.DashConfig.Logger().With("module", "overview"),
	}

	if err := co.bootstrap(ctx); err != nil {
		return nil, err
	}

	logger := log.From(ctx).With("module", "overview")

	co.dashConfig.ObjectStore().RegisterOnUpdate(func(newObjectStore store.Store) {
		logger.Debugf("object store was updated")
		if err := co.bootstrap(ctx); err != nil {
			logger.WithErr(err).Errorf("updating object store")
		}
	})

	return co, nil
}

func (co *Overview) SetContext(ctx context.Context, contextName string) error {
	co.contextName = contextName
	return nil
}

func (co *Overview) bootstrap(ctx context.Context) error {
	pathMatcher := describer.NewPathMatcher("overview")
	for _, pf := range rootDescriber.PathFilters() {
		pathMatcher.Register(ctx, pf)
	}

	for _, pf := range eventsDescriber.PathFilters() {
		pathMatcher.Register(ctx, pf)
	}

	g, err := newGenerator(pathMatcher, co.dashConfig)
	if err != nil {
		return errors.Wrap(err, "create overview generator")
	}

	objectPathConfig := octant.ObjectPathConfig{
		ModuleName:     "overview",
		SupportedGVKs:  supportedGVKs,
		PathLookupFunc: gvkPath,
		CRDPathGenFunc: crdPath,
	}
	objectPath, err := octant.NewObjectPath(objectPathConfig)
	if err != nil {
		return errors.Wrap(err, "create module object path generator")
	}

	co.ObjectPath = objectPath
	co.generator = g

	key := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	crdWatcher := co.dashConfig.CRDWatcher()
	if err := co.dashConfig.ObjectStore().HasAccess(key, "watch"); err == nil {
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
			return errors.Wrap(err, "create namespaced CRD watcher for overview")
		}
	}

	return nil
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
func (co *Overview) Navigation(ctx context.Context, namespace, root string) ([]octant.Navigation, error) {
	navigationEntries := octant.NavigationEntries{
		Lookup: navPathLookup,
		EntriesFuncs: map[string]octant.EntriesFunc{
			"Workloads":                    workloadEntries,
			"Discovery and Load Balancing": discoAndLBEntries,
			"Config and Storage":           configAndStorageEntries,
			"Custom Resources":             octant.CRDEntries,
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

	nf := octant.NewNavigationFactory(namespace, root, objectStore, navigationEntries)

	entries, err := nf.Generate(ctx, "Overview", icon.Overview, "")
	if err != nil {
		return nil, err
	}

	return []octant.Navigation{
		*entries,
	}, nil
}

// Generators allow modules to send events to the frontend.
func (co *Overview) Generators() []octant.Generator {
	return []octant.Generator{}
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

func (co *Overview) ActionPaths() map[string]action.DispatcherFunc {
	configurationEditor := NewConfigurationEditor(co.logger, co.dashConfig.ObjectStore())

	return map[string]action.DispatcherFunc{
		configurationEditor.ActionName(): configurationEditor.Handle,
	}
}

func roundToInt(val float64) int64 {
	if val < 0 {
		return int64(val - 0.5)
	}
	return int64(val + 0.5)
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/generator"
	internalLog "github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/icon"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type Options struct {
	Namespace  string
	DashConfig config.Dash
}

// Overview is an API for generating a cluster overview.
type Overview struct {
	*octant.ObjectPath

	generator   generator.Interface
	dashConfig  config.Dash
	contextName string
	pathMatcher *describer.PathMatcher
	logger      log.Logger

	watchedCRDs []*unstructured.Unstructured

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

	logger := internalLog.From(ctx).With("module", "overview")

	co.dashConfig.ObjectStore().RegisterOnUpdate(func(newObjectStore store.Store) {
		logger.Debugf("object store was updated")
		if err := co.bootstrap(ctx); err != nil {
			logger.WithErr(err).Errorf("updating object store")
		}
	})

	return co, nil
}

func (co *Overview) SetContext(ctx context.Context, contextName string) error {
	co.mu.Lock()
	defer co.mu.Unlock()

	customResourcesDescriber := describer.NamespacedCRD()
	co.contextName = contextName
	for i := range co.watchedCRDs {
		describer.DeleteCRD(ctx, co.watchedCRDs[i], co.pathMatcher, customResourcesDescriber, co)
	}

	co.watchedCRDs = []*unstructured.Unstructured{}
	crdWatcher := co.dashConfig.CRDWatcher()
	crdWatcher.Watch(ctx)

	return nil
}

func (co *Overview) bootstrap(ctx context.Context) error {
	rootDescriber := describer.NamespacedOverview()

	if err := rootDescriber.Reset(ctx); err != nil {
		return err
	}

	pathMatcher := describer.NewPathMatcher("overview")
	for _, pf := range rootDescriber.PathFilters() {
		pathMatcher.Register(ctx, pf)
	}

	g, err := generator.NewGenerator(pathMatcher, co.dashConfig)
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

	crdWatcher := co.dashConfig.CRDWatcher()

	customResourcesDescriber := describer.NamespacedCRD()

	watchConfig := &config.CRDWatchConfig{
		Add: func(pathMatcher *describer.PathMatcher, sectionDescriber *describer.CRDSection) config.ObjectHandler {
			return func(ctx context.Context, object *unstructured.Unstructured) {
				co.mu.Lock()
				defer co.mu.Unlock()

				if object == nil {
					return
				}
				describer.AddCRD(ctx, object, pathMatcher, customResourcesDescriber, co)
				co.watchedCRDs = append(co.watchedCRDs, object)
			}
		}(pathMatcher, customResourcesDescriber),
		Delete: func(pathMatcher *describer.PathMatcher, csd *describer.CRDSection) config.ObjectHandler {
			return func(ctx context.Context, object *unstructured.Unstructured) {
				co.mu.Lock()
				defer co.mu.Unlock()

				if object == nil {
					return
				}
				describer.DeleteCRD(ctx, object, pathMatcher, customResourcesDescriber, co)
				var list []*unstructured.Unstructured
				for i := range co.watchedCRDs {
					if co.watchedCRDs[i].GetUID() == object.GetUID() {
						continue
					}
					list = append(list, co.watchedCRDs[i])
				}
				co.watchedCRDs = list
			}
		}(pathMatcher, customResourcesDescriber),
		IsNamespaced: true,
	}

	if err := crdWatcher.AddConfig(watchConfig); err != nil {
		return errors.Wrap(err, "create namespaced CRD watcher for overview")
	}

	co.pathMatcher = pathMatcher

	return nil
}

// Name returns the name for this module.
func (co *Overview) Name() string {
	return "overview"
}

func (co *Overview) ClientRequestHandlers() []octant.ClientRequestHandler {
	return nil
}

// ContentPath returns the content path for overview.
func (co *Overview) ContentPath() string {
	return co.Name()
}

// Navigation returns navigation entries for overview.
func (co *Overview) Navigation(ctx context.Context, namespace, root string) ([]navigation.Navigation, error) {
	navigationEntries := octant.NavigationEntries{
		Lookup: navPathLookup,
		EntriesFuncs: map[string]octant.EntriesFunc{
			"Namespace Overview":           nil,
			"Workloads":                    workloadEntries,
			"Discovery and Load Balancing": discoAndLBEntries,
			"Config and Storage":           configAndStorageEntries,
			"Custom Resources":             navigation.CRDEntries,
			"RBAC":                         rbacEntries,
			"Events":                       nil,
		},
		IconMap: map[string]string{
			"Namespace Overview":           icon.Overview,
			"Workloads":                    icon.Workloads,
			"Discovery and Load Balancing": icon.DiscoveryAndLoadBalancing,
			"Config and Storage":           icon.ConfigAndStorage,
			"Custom Resources":             icon.CustomResources,
			"RBAC":                         icon.RBAC,
			"Events":                       icon.Events,
		},
		Order: []string{
			"Namespace Overview",
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

	entries, err := nf.Generate(ctx, "", false)
	if err != nil {
		return nil, err
	}

	return entries, nil
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
func (co *Overview) Content(ctx context.Context, contentPath string, opts module.ContentOptions) (component.ContentResponse, error) {
	ctx = internalLog.WithLoggerContext(ctx, co.dashConfig.Logger())
	genOpts := generator.Options{
		LabelSet: opts.LabelSet,
	}
	return co.generator.Generate(ctx, contentPath, genOpts)
}

// ActionPaths contain the actions this module is responsible for.
func (co *Overview) ActionPaths() map[string]action.DispatcherFunc {
	dispatchers := action.Dispatchers{
		octant.NewDeploymentConfigurationEditor(co.logger, co.dashConfig.ObjectStore()),
		octant.NewContainerEditor(co.dashConfig.ObjectStore()),
		octant.NewServiceConfigurationEditor(co.dashConfig.ObjectStore()),
		octant.NewPortForward(co.logger, co.dashConfig.ObjectStore(), co.dashConfig.PortForwarder()),
		octant.NewPortForwardDelete(co.logger, co.dashConfig.ObjectStore(), co.dashConfig.PortForwarder()),
		octant.NewTerminalCommandExec(co.logger, co.dashConfig.ObjectStore(), co.dashConfig.TerminalManager()),
		octant.NewTerminalDelete(co.logger, co.dashConfig.ObjectStore(), co.dashConfig.TerminalManager()),
		octant.NewCordon(co.dashConfig.ObjectStore(), co.dashConfig.ClusterClient()),
		octant.NewUncordon(co.dashConfig.ObjectStore(), co.dashConfig.ClusterClient()),
		octant.NewCronJobTrigger(co.dashConfig.ObjectStore(), co.dashConfig.ClusterClient()),
		octant.NewObjectUpdaterDispatcher(co.dashConfig.ObjectStore()),
	}

	return dispatchers.ToActionPaths()
}

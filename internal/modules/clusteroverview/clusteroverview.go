/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware/octant/internal/api"
	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/describer"
	"github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/modules/overview/printer"
	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/pkg/icon"
	"github.com/vmware/octant/pkg/navigation"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

// Options are options for ClusterOverview.
type Options struct {
	DashConfig config.Dash
}

// ClusterOverview is a module for the cluster overview.
type ClusterOverview struct {
	*octant.ObjectPath
	Options

	pathMatcher *describer.PathMatcher
}

var _ module.Module = (*ClusterOverview)(nil)

func New(ctx context.Context, options Options) (*ClusterOverview, error) {
	pathMatcher := describer.NewPathMatcher("cluster-overview")
	for _, pf := range rootDescriber.PathFilters() {
		pathMatcher.Register(ctx, pf)
	}

	objectPathConfig := octant.ObjectPathConfig{
		ModuleName:     "cluster-overview",
		SupportedGVKs:  supportedGVKs,
		PathLookupFunc: gvkPath,
		CRDPathGenFunc: crdPath,
	}
	objectPath, err := octant.NewObjectPath(objectPathConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create module object path generator")
	}

	co := &ClusterOverview{
		ObjectPath:  objectPath,
		pathMatcher: pathMatcher,
		Options:     options,
	}

	key := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	objectStore := options.DashConfig.ObjectStore()

	crdWatcher := options.DashConfig.CRDWatcher()
	if err := objectStore.HasAccess(ctx, key, "watch"); err == nil {
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
			IsNamespaced: false,
		}

		if err := crdWatcher.Watch(ctx, watchConfig); err != nil {
			return nil, errors.Wrap(err, "create namespaced CRD watcher for overview")
		}
	}

	return co, nil
}

func (co *ClusterOverview) Name() string {
	return "cluster-overview"
}

func (co *ClusterOverview) Handlers(ctx context.Context) map[string]http.Handler {
	logger := log.From(ctx)

	pfHandler, err := newPortForwardsHandler(logger, co.DashConfig.PortForwarder())
	if err != nil {
		panic(fmt.Sprintf("unable to create port forwards handler: %v", err))
	}

	return map[string]http.Handler{
		"/port-forwards": pfHandler,
	}
}

func (co *ClusterOverview) Content(ctx context.Context, contentPath string, prefix string, namespace string, opts module.ContentOptions) (component.ContentResponse, error) {
	pf, err := co.pathMatcher.Find(contentPath)
	if err != nil {
		if err == describer.ErrPathNotFound {
			return describer.EmptyContentResponse, api.NewNotFoundError(contentPath)
		}
		return describer.EmptyContentResponse, err
	}

	clusterClient := co.DashConfig.ClusterClient()
	objectStore := co.DashConfig.ObjectStore()

	discoveryInterface, err := clusterClient.DiscoveryClient()
	if err != nil {
		return describer.EmptyContentResponse, err
	}

	q := queryer.New(objectStore, discoveryInterface)

	p := printer.NewResource(co.DashConfig)
	if err := printer.AddHandlers(p); err != nil {
		return describer.EmptyContentResponse, errors.Wrap(err, "add print handlers")
	}

	linkGenerator, err := link.NewFromDashConfig(co.DashConfig)
	if err != nil {
		return describer.EmptyContentResponse, err
	}

	loaderFactory := describer.NewObjectLoaderFactory(co.DashConfig)

	options := describer.Options{
		Queryer:  q,
		Fields:   pf.Fields(contentPath),
		Printer:  p,
		LabelSet: opts.LabelSet,
		Dash:     co.DashConfig,
		Link:     linkGenerator,

		LoadObjects: loaderFactory.LoadObjects,
		LoadObject:  loaderFactory.LoadObject,
	}

	cResponse, err := pf.Describer.Describe(ctx, prefix, "", options)
	if err != nil {
		return describer.EmptyContentResponse, err
	}

	return cResponse, nil
}

func (co *ClusterOverview) ContentPath() string {
	return fmt.Sprintf("/%s", co.Name())
}

func (co *ClusterOverview) Navigation(ctx context.Context, namespace string, root string) ([]navigation.Navigation, error) {
	navigationEntries := octant.NavigationEntries{
		Lookup: map[string]string{
			"Custom Resources": "custom-resources",
			"RBAC":             "rbac",
			"Nodes":            "nodes",
			"Port Forwards":    "port-forward",
		},
		EntriesFuncs: map[string]octant.EntriesFunc{
			"Custom Resources": navigation.CRDEntries,
			"RBAC":             rbacEntries,
			"Nodes":            nil,
			"Port Forwards":    nil,
		},
		Order: []string{
			"Custom Resources",
			"RBAC",
			"Nodes",
			"Port Forwards",
		},
	}

	objectStore := co.DashConfig.ObjectStore()

	nf := octant.NewNavigationFactory("", root, objectStore, navigationEntries)

	entries, err := nf.Generate(ctx, "Cluster Overview", icon.ClusterOverview, "", true)
	if err != nil {
		return nil, err
	}

	return []navigation.Navigation{
		*entries,
	}, nil
}

func (co *ClusterOverview) SetNamespace(namespace string) error {
	return nil
}

func (co *ClusterOverview) Start() error {
	return nil
}

func (co *ClusterOverview) Stop() {
}

// Generators allow modules to send events to the frontend.
func (co *ClusterOverview) Generators() []octant.Generator {
	return []octant.Generator{}
}

func rbacEntries(_ context.Context, prefix, _ string, _ store.Store, _ bool) ([]navigation.Navigation, error) {
	neh := navigation.EntriesHelper{}
	neh.Add("Cluster Roles", "cluster-roles", icon.ClusterOverviewClusterRole)
	neh.Add("Cluster Role Bindings", "cluster-role-bindings", icon.ClusterOverviewClusterRoleBinding)

	return neh.Generate(prefix)
}

func (co *ClusterOverview) SetContext(ctx context.Context, contextName string) error {
	return nil
}

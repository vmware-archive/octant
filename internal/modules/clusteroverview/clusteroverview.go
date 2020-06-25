/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/internal/loading"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/printer"
	"github.com/vmware-tanzu/octant/internal/queryer"
	"github.com/vmware-tanzu/octant/pkg/icon"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
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
	watchedCRDs []*unstructured.Unstructured

	mu sync.Mutex
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

	crdWatcher := options.DashConfig.CRDWatcher()
	watchConfig := &config.CRDWatchConfig{
		Add: func(_ *describer.PathMatcher, sectionDescriber *describer.CRDSection) config.ObjectHandler {
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
		Delete: func(_ *describer.PathMatcher, csd *describer.CRDSection) config.ObjectHandler {
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
		IsNamespaced: false,
	}

	if err := crdWatcher.AddConfig(watchConfig); err != nil {
		return nil, errors.Wrap(err, "create cluster scoped CRD watcher for cluster overview")
	}

	return co, nil
}

func (co *ClusterOverview) Name() string {
	return "cluster-overview"
}

func (co *ClusterOverview) ClientRequestHandlers() []octant.ClientRequestHandler {
	return nil
}

func (co *ClusterOverview) Content(ctx context.Context, contentPath string, opts module.ContentOptions) (component.ContentResponse, error) {
	pf, err := co.pathMatcher.Find(contentPath)
	if err != nil {
		if err == describer.ErrPathNotFound {
			return component.EmptyContentResponse, api.NewNotFoundError(contentPath)
		}
		return component.EmptyContentResponse, err
	}

	clusterClient := co.DashConfig.ClusterClient()
	objectStore := co.DashConfig.ObjectStore()

	discoveryInterface, err := clusterClient.DiscoveryClient()
	if err != nil {
		return component.EmptyContentResponse, err
	}

	q := queryer.New(objectStore, discoveryInterface)

	p := printer.NewResource(co.DashConfig)
	if err := printer.AddHandlers(p); err != nil {
		return component.EmptyContentResponse, errors.Wrap(err, "add print handlers")
	}

	linkGenerator, err := link.NewFromDashConfig(co.DashConfig)
	if err != nil {
		return component.EmptyContentResponse, err
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

	cResponse, err := pf.Describer.Describe(ctx, "", options)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	return cResponse, nil
}

func (co *ClusterOverview) ContentPath() string {
	return fmt.Sprintf("%s", co.Name())
}

func (co *ClusterOverview) Navigation(ctx context.Context, _ string, root string) ([]navigation.Navigation, error) {
	navigationEntries := octant.NavigationEntries{
		Lookup: map[string]string{
			"Namespaces":                  "namespaces",
			"Custom Resources":            "custom-resources",
			"Custom Resource Definitions": "custom-resource-definitions",
			"RBAC":                        "rbac",
			"Nodes":                       "nodes",
			"Storage":                     "storage",
			"Port Forwards":               "port-forward",
		},
		EntriesFuncs: map[string]octant.EntriesFunc{
			"Cluster Overview":            nil,
			"Namespaces":                  nil,
			"Custom Resources":            navigation.CRDEntries,
			"Custom Resource Definitions": nil,
			"RBAC":                        rbacEntries,
			"Nodes":                       nil,
			"Storage":                     storageEntries,
			"Port Forwards":               nil,
		},
		IconMap: map[string]string{
			"Cluster Overview":            icon.Overview,
			"Namespaces":                  icon.Namespaces,
			"Custom Resources":            icon.CustomResources,
			"Custom Resource Definitions": icon.CustomResourceDefinition,
			"RBAC":                        icon.RBAC,
			"Nodes":                       icon.Nodes,
			"Storage":                     icon.ConfigAndStorage,
			"Port Forwards":               icon.PortForwards,
		},
		Order: []string{
			"Cluster Overview",
			"Namespaces",
			"Custom Resources",
			"Custom Resource Definitions",
			"RBAC",
			"Nodes",
			"Storage",
			"Port Forwards",
		},
	}

	objectStore := co.DashConfig.ObjectStore()

	nf := octant.NewNavigationFactory("", root, objectStore, navigationEntries)

	entries, err := nf.Generate(ctx, "Cluster", true)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (co *ClusterOverview) SetNamespace(_ string) error {
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

func rbacEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}

	neh.Add("Cluster Roles", "cluster-roles",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ClusterRole), objectStore))
	neh.Add("Cluster Role Bindings", "cluster-role-bindings",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ClusterRoleBinding), objectStore))

	children, err := neh.Generate(prefix, namespace, "")
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func storageEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}

	neh.Add("Persistent Volumes", "persistent-volumes",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.PersistentVolume), objectStore))

	children, err := neh.Generate(prefix, namespace, "")
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func (co *ClusterOverview) SetContext(ctx context.Context, _ string) error {
	co.mu.Lock()
	defer co.mu.Unlock()

	for i := range co.watchedCRDs {
		describer.DeleteCRD(ctx, co.watchedCRDs[i], co.pathMatcher, customResourcesDescriber, co)
	}

	co.watchedCRDs = []*unstructured.Unstructured{}
	return nil
}

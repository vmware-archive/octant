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

	"github.com/vmware/octant/internal/api"
	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/describer"
	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/internal/loading"
	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/internal/printer"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/pkg/action"
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
	objectStore := co.DashConfig.ObjectStore()
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
				describer.DeleteCRD(ctx, object, pathMatcher, customResourcesDescriber, co, objectStore)
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

	if err := crdWatcher.Watch(ctx, watchConfig); err != nil {
		return nil, errors.Wrap(err, "create namespaced CRD watcher for overview")
	}

	return co, nil
}

func (co *ClusterOverview) Name() string {
	return "cluster-overview"
}

func (co *ClusterOverview) ClientRequestHandlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		// TODO: move to overview
		{
			RequestType: "startPortForward",
			Handler: func(state octant.State, payload action.Payload) error {
				req, err := portForwardRequestFromPayload(payload)
				if err != nil {
					return errors.Wrap(err, "convert payload to port forward request")
				}

				_, err = co.DashConfig.PortForwarder().Create(context.TODO(), req.gvk(), req.Name, req.Namespace, req.Port)
				return err
			},
		},
		{
			RequestType: "stopPortForward",
			Handler: func(state octant.State, payload action.Payload) error {
				id, err := payload.String("id")
				if err != nil {
					return errors.Wrap(err, "get port forward id from payload")
				}

				co.DashConfig.PortForwarder().StopForwarder(id)
				return nil
			},
		},
	}
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

func (co *ClusterOverview) Navigation(ctx context.Context, namespace string, root string) ([]navigation.Navigation, error) {
	navigationEntries := octant.NavigationEntries{
		Lookup: map[string]string{
			"Namespaces":       "namespaces",
			"Custom Resources": "custom-resources",
			"RBAC":             "rbac",
			"Nodes":            "nodes",
			"Port Forwards":    "port-forward",
		},
		EntriesFuncs: map[string]octant.EntriesFunc{
			"Namespaces":       nil,
			"Custom Resources": navigation.CRDEntries,
			"RBAC":             rbacEntries,
			"Nodes":            nil,
			"Port Forwards":    nil,
		},
		Order: []string{
			"Namespaces",
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

func rbacEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}
	neh.Add("Cluster Roles", "cluster-roles", icon.ClusterOverviewClusterRole,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ClusterRole), objectStore))
	neh.Add("Cluster Role Bindings", "cluster-role-bindings", icon.ClusterOverviewClusterRoleBinding,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ClusterRoleBinding), objectStore))

	children, err := neh.Generate(prefix)
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func (co *ClusterOverview) SetContext(ctx context.Context, contextName string) error {
	co.mu.Lock()
	defer co.mu.Unlock()

	for i := range co.watchedCRDs {
		describer.DeleteCRD(ctx, co.watchedCRDs[i], co.pathMatcher, customResourcesDescriber, co, co.DashConfig.ObjectStore())
	}

	co.watchedCRDs = []*unstructured.Unstructured{}
	return nil
}

package clusteroverview

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/heptio/developer-dash/internal/config"
	"github.com/heptio/developer-dash/internal/describer"
	"github.com/heptio/developer-dash/internal/link"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/modules/overview/printer"
	"github.com/heptio/developer-dash/internal/modules/overview/resourceviewer"
	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/component"
)

// Options are options for ClusterOverview.
type Options struct {
	DashConfig config.Dash
}

// ClusterOverview is a module for the cluster overview.
type ClusterOverview struct {
	*clustereye.ObjectPath
	Options

	pathMatcher *describer.PathMatcher
}

var _ module.Module = (*ClusterOverview)(nil)

func New(ctx context.Context, options Options) (*ClusterOverview, error) {
	pathMatcher := describer.NewPathMatcher("cluster-overview")
	for _, pf := range rootDescriber.PathFilters() {
		pathMatcher.Register(ctx, pf)
	}

	objectPathConfig := clustereye.ObjectPathConfig{
		ModuleName:     "cluster-overview",
		SupportedGVKs:  supportedGVKs,
		PathLookupFunc: gvkPath,
		CRDPathGenFunc: crdPath,
	}
	objectPath, err := clustereye.NewObjectPath(objectPathConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create module object path generator")
	}

	co := &ClusterOverview{
		ObjectPath:  objectPath,
		pathMatcher: pathMatcher,
		Options:     options,
	}

	key := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	objectStore := options.DashConfig.ObjectStore()

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

	componentCache, err := resourceviewer.NewComponentCache(co.DashConfig)
	if err != nil {
		return describer.EmptyContentResponse, errors.Wrap(err, "create component cache")
	}
	componentCache.SetQueryer(q)

	options := describer.Options{
		Queryer:        q,
		Fields:         pf.Fields(contentPath),
		Printer:        p,
		LabelSet:       opts.LabelSet,
		Dash:           co.DashConfig,
		Link:           linkGenerator,
		ComponentCache: componentCache,

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

func (co *ClusterOverview) Navigation(ctx context.Context, namespace string, root string) ([]clustereye.Navigation, error) {
	navigationEntries := clustereye.NavigationEntries{
		Lookup: map[string]string{
			"Custom Resources": "custom-resources",
			"RBAC":             "rbac",
		},
		EntriesFuncs: map[string]clustereye.EntriesFunc{
			"Custom Resources": clustereye.CRDEntries,
			"RBAC":             rbacEntries,
		},
		Order: []string{
			"Custom Resources",
			"RBAC",
		},
	}

	objectStore := co.DashConfig.ObjectStore()

	nf := clustereye.NewNavigationFactory("", root, objectStore, navigationEntries)

	entries, err := nf.Generate(ctx, "Cluster Overview")
	if err != nil {
		return nil, err
	}

	return []clustereye.Navigation{
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

func rbacEntries(_ context.Context, prefix, _ string, _ objectstore.ObjectStore) ([]clustereye.Navigation, error) {
	return []clustereye.Navigation{
		*clustereye.NewNavigation("Cluster Roles", path.Join(prefix, "cluster-roles")),
		*clustereye.NewNavigation("Cluster Role Bindings", path.Join(prefix, "cluster-role-bindings")),
	}, nil
}

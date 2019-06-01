package clusteroverview

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/heptio/developer-dash/internal/config"
	"github.com/heptio/developer-dash/internal/describer"
	"github.com/heptio/developer-dash/internal/link"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/modules/overview/printer"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/pkg/view/component"
)

// Options are options for ClusterOverview.
type Options struct {
	DashConfig config.Dash
}

// ClusterOverview is a module for the cluster overview.
type ClusterOverview struct {
	Options
	objectPath

	pathMatcher *describer.PathMatcher
}

var _ module.Module = (*ClusterOverview)(nil)

func New(ctx context.Context, options Options) *ClusterOverview {
	pm := describer.NewPathMatcher()
	for _, pf := range rootDescriber.PathFilters() {
		pm.Register(ctx, pf)
	}

	return &ClusterOverview{
		pathMatcher: pm,
		Options:     options,
	}
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
		LoadObject: loaderFactory.LoadObject,
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
	return []clustereye.Navigation{
		{
			Title: "Cluster Overview",
			Path:  path.Join("/content", co.ContentPath(), "/"),
			Children: []clustereye.Navigation{
				{
					Title: "RBAC",
					Path:  path.Join("/content", co.ContentPath(), "/rbac"),
					Children: []clustereye.Navigation{
						{
							Title: "Cluster Roles",
							Path:  path.Join("/content", co.ContentPath(), "/rbac", "cluster-roles"),
						},
						{
							Title: "Cluster Role Bindings",
							Path:  path.Join("/content", co.ContentPath(), "/rbac", "cluster-role-bindings"),
						},
					},
				},
				{
					Title: "Port Forwards",
					Path:  path.Join("/content", co.ContentPath(), "/port-forward"),
				},
			},
		},
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

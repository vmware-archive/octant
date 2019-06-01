package overview

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	kLabels "k8s.io/apimachinery/pkg/labels"

	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/config"
	"github.com/heptio/developer-dash/internal/describer"
	"github.com/heptio/developer-dash/internal/link"
	"github.com/heptio/developer-dash/internal/modules/overview/printer"
	"github.com/heptio/developer-dash/internal/modules/overview/resourceviewer"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/pkg/view/component"
)

type realGenerator struct {
	pathMatcher    *describer.PathMatcher
	componentCache resourceviewer.ComponentCache
	printer        printer.Printer
	dashConfig     config.Dash
}

// GeneratorOptions are additional options to pass a generator
type GeneratorOptions struct {
	LabelSet *kLabels.Set
}

// newGenerator creates a generator.
func newGenerator(pm *describer.PathMatcher, dashConfig config.Dash, componentCache resourceviewer.ComponentCache) (*realGenerator, error) {
	p := printer.NewResource(dashConfig)

	if err := printer.AddHandlers(p); err != nil {
		return nil, errors.Wrap(err, "add print handlers")
	}

	if pm == nil {
		return nil, errors.New("path matcher is nil")
	}

	return &realGenerator{
		pathMatcher:    pm,
		componentCache: componentCache,
		printer:        p,
		dashConfig:     dashConfig,
	}, nil
}

func (g *realGenerator) Generate(ctx context.Context, path, prefix, namespace string, opts GeneratorOptions) (component.ContentResponse, error) {
	ctx, span := trace.StartSpan(ctx, "Generate")
	defer span.End()

	pf, err := g.pathMatcher.Find(path)
	if err != nil {
		if err == describer.ErrPathNotFound {
			return emptyContentResponse, api.NewNotFoundError(path)
		}
		return emptyContentResponse, err
	}

	discoveryInterface, err := g.dashConfig.ClusterClient().DiscoveryClient()
	if err != nil {
		return emptyContentResponse, err
	}

	linkGenerator, err := link.NewFromDashConfig(g.dashConfig)
	if err != nil {
		return emptyContentResponse, err
	}

	q := queryer.New(g.dashConfig.ObjectStore(), discoveryInterface)
	g.componentCache.SetQueryer(q)

	loaderFactory := describer.NewObjectLoaderFactory(g.dashConfig)

	fields := pf.Fields(path)
	options := describer.Options{
		ComponentCache: g.componentCache,
		Queryer:        q,
		Fields:         fields,
		Printer:        g.printer,
		LabelSet:       opts.LabelSet,
		Dash:           g.dashConfig,
		Link:           linkGenerator,

		LoadObjects: loaderFactory.LoadObjects,
		LoadObject:  loaderFactory.LoadObject,
	}

	cResponse, err := pf.Describer.Describe(ctx, prefix, namespace, options)
	if err != nil {
		return emptyContentResponse, err
	}

	return cResponse, nil
}

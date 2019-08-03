/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	kLabels "k8s.io/apimachinery/pkg/labels"

	"github.com/vmware/octant/internal/api"
	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/describer"
	"github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/internal/modules/overview/printer"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/pkg/view/component"
)

type realGenerator struct {
	pathMatcher *describer.PathMatcher
	printer     printer.Printer
	dashConfig  config.Dash
}

// GeneratorOptions are additional options to pass a generator
type GeneratorOptions struct {
	LabelSet *kLabels.Set
}

// newGenerator creates a generator.
func newGenerator(pm *describer.PathMatcher, dashConfig config.Dash) (*realGenerator, error) {
	p := printer.NewResource(dashConfig)

	if err := printer.AddHandlers(p); err != nil {
		return nil, errors.Wrap(err, "add print handlers")
	}

	if pm == nil {
		return nil, errors.New("path matcher is nil")
	}

	return &realGenerator{
		pathMatcher: pm,
		printer:     p,
		dashConfig:  dashConfig,
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

	loaderFactory := describer.NewObjectLoaderFactory(g.dashConfig)

	fields := pf.Fields(path)
	options := describer.Options{
		Queryer:  q,
		Fields:   fields,
		Printer:  g.printer,
		LabelSet: opts.LabelSet,
		Dash:     g.dashConfig,
		Link:     linkGenerator,

		LoadObjects: loaderFactory.LoadObjects,
		LoadObject:  loaderFactory.LoadObject,
	}

	cResponse, err := pf.Describer.Describe(ctx, prefix, namespace, options)
	if err != nil {
		return emptyContentResponse, err
	}

	return cResponse, nil
}

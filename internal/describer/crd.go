/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/internal/modules/overview/yamlviewer"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/printer"
	"github.com/vmware-tanzu/octant/internal/queryer"
	"github.com/vmware-tanzu/octant/internal/resourceviewer"
	"github.com/vmware-tanzu/octant/pkg/icon"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type crdPrinter func(ctx context.Context, crd, object *unstructured.Unstructured, options printer.Options) (component.Component, error)
type metadataPrinter func(runtime.Object, link.Interface) (*component.FlexLayout, error)
type resourceViewerPrinter func(ctx context.Context, object *unstructured.Unstructured, dashConfig config.Dash, q queryer.Queryer) (component.Component, error)
type yamlPrinter func(runtime.Object) (*component.YAML, error)

type crdOption func(*crd)

type crd struct {
	base

	path                  string
	name                  string
	summaryPrinter        crdPrinter
	metadataPrinter       metadataPrinter
	resourceViewerPrinter resourceViewerPrinter
	yamlPrinter           yamlPrinter
}

var _ Describer = (*crd)(nil)

func newCRD(name, path string, options ...crdOption) *crd {
	d := &crd{
		path:                  path,
		name:                  name,
		summaryPrinter:        printer.CustomResourceHandler,
		metadataPrinter:       printer.MetadataHandler,
		resourceViewerPrinter: createCRDResourceViewer,
		yamlPrinter:           yamlviewer.ToComponent,
	}

	for _, option := range options {
		option(d)
	}

	return d
}

func (c *crd) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	objectStore := options.ObjectStore()
	crd, err := CustomResourceDefinition(ctx, c.name, objectStore)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	octantCRD, err := octant.NewCustomResourceDefinition(crd)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	crdVersions, err := octantCRD.Versions()
	if err != nil {
		return component.EmptyContentResponse, fmt.Errorf("get versions for crd %s: %w", crd.GetName(), err)
	} else if len(crdVersions) == 0 {
		return component.EmptyContentResponse, fmt.Errorf("crd %s has no no versions", crd.GetName())
	}

	crGVK, err := gvk.CustomResource(crd, crdVersions[0])
	if err != nil {
		return component.EmptyContentResponse, fmt.Errorf("get gvk for custom resource")
	}

	apiVersion, kind := crGVK.ToAPIVersionAndKind()

	key := store.Key{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       options.Fields["name"],
	}

	object, found, err := objectStore.Get(ctx, key)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	if !found {
		return component.EmptyContentResponse, err
	}

	title := component.Title(
		component.NewText("Custom Resources"),
		component.NewText(crd.GetName()),
		component.NewText(object.GroupVersionKind().Version),
		component.NewText(object.GetName()))

	iconName, iconSource := loadIcon(icon.CustomResourceDefinition)
	cr := component.NewContentResponse(title)
	cr.IconName = iconName
	cr.IconSource = iconSource

	linkGenerator, err := link.NewFromDashConfig(options)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	printOptions := printer.Options{
		DashConfig: options,
		Link:       linkGenerator,
	}

	summary, err := c.summaryPrinter(ctx, crd, object, printOptions)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	summary.SetAccessor("summary")

	cr.Add(summary)

	metadata, err := c.metadataPrinter(object, linkGenerator)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	metadata.SetAccessor("metadata")
	cr.Add(metadata)

	resourceViewerComponent, err := c.resourceViewerPrinter(ctx, object, options, options.Queryer)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	resourceViewerComponent.SetAccessor("resourceViewer")
	cr.Add(resourceViewerComponent)

	yvComponent, err := c.yamlPrinter(object)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	yvComponent.SetAccessor("yaml")
	cr.Add(yvComponent)

	pluginPrinter := options.PluginManager()
	tabs, err := pluginPrinter.Tabs(ctx, object)
	if err != nil {
		return component.EmptyContentResponse, errors.Wrap(err, "getting tabs from plugins")
	}

	for _, tab := range tabs {
		tab.Contents.SetAccessor(tab.Name)
		cr.Add(&tab.Contents)
	}

	return *cr, nil
}

func (c *crd) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(c.path, c),
	}
}

func createCRDResourceViewer(ctx context.Context, object *unstructured.Unstructured, dashConfig config.Dash, q queryer.Queryer) (component.Component, error) {
	rv, err := resourceviewer.New(dashConfig, resourceviewer.WithDefaultQueryer(dashConfig, q))
	if err != nil {
		return nil, err
	}

	handler, err := resourceviewer.NewHandler(dashConfig)
	if err != nil {
		return nil, err
	}

	if err := rv.Visit(ctx, object, handler); err != nil {
		return nil, err
	}

	return resourceviewer.GenerateComponent(ctx, handler, "")
}

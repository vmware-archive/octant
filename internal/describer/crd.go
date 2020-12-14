/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func defaultCustomResourceTabs(crdName string) []Tab {
	return []Tab{
		{Name: "Summary", Factory: CustomResourceSummaryTab(crdName)},
		{Name: "Metadata", Factory: MetadataTab},
		{Name: "Resource Viewer", Factory: ResourceViewerTab},
		{Name: "YAML", Factory: YAMLViewerTab},
	}
}

type crdOption func(*crd)

type crd struct {
	base

	path               string
	name               string
	version            string
	tabsGenerator      TabsGenerator
	tabFuncDescriptors []Tab
	resourceLoader     ResourceLoader
}

var _ Describer = (*crd)(nil)

func newCRD(name, path, version string, s store.Store, options ...crdOption) *crd {
	d := &crd{
		path:               path,
		name:               name,
		version:            version,
		tabsGenerator:      NewObjectTabsGenerator(),
		tabFuncDescriptors: defaultCustomResourceTabs(name),
		resourceLoader:     NewStoreResourceLoader(s),
	}

	for _, option := range options {
		option(d)
	}

	return d
}

func (c *crd) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	d := ResourceDescriptor{
		CustomResourceDefinitionName: c.name,
		Namespace:                    namespace,
		CustomResourceVersion:        c.version,
		CustomResourceName:           options.Fields["name"],
	}
	resp, err := c.resourceLoader.Load(ctx, d)
	if err != nil {
		return component.EmptyContentResponse, fmt.Errorf("load custom resource: %w", err)
	}

	crd := resp.CustomResourceDefinition
	object := resp.CustomResource

	name := fmt.Sprintf("%s (%s)", object.GetName(), object.GroupVersionKind().Version)
	title := getCrdTitle(namespace, crd, name)

	cr := component.NewContentResponse(title)

	generatorConfig := TabsGeneratorConfig{
		Object:      object,
		TabsFactory: objectTabsFactory(ctx, object, c.tabFuncDescriptors, options),
		Options:     options,
	}

	tabComponents, err := c.tabsGenerator.Generate(ctx, generatorConfig)
	if err != nil {
		return component.EmptyContentResponse, fmt.Errorf("generate tabs: %w", err)
	}

	cr.Add(tabComponents...)

	return *cr, nil
}

func (c *crd) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(c.path, c),
	}
}

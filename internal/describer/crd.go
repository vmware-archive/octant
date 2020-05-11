/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"fmt"
	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/octant"
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
	tabsGenerator      TabsGenerator
	tabFuncDescriptors []Tab
}

var _ Describer = (*crd)(nil)

func newCRD(name, path string, options ...crdOption) *crd {
	d := &crd{
		path:               path,
		name:               name,
		tabsGenerator:      NewObjectTabsGenerator(),
		tabFuncDescriptors: defaultCustomResourceTabs(name),
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

	object, err := objectStore.Get(ctx, key)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	if object == nil {
		return component.EmptyContentResponse, err
	}
	title := getCrdTitle(namespace, crd, object.GetName())

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

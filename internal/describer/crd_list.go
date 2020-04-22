/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/internal/modules/overview/yamlviewer"
	"github.com/vmware-tanzu/octant/internal/printer"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type crdListPrinter func(crdObject *unstructured.Unstructured, resources *unstructured.UnstructuredList, version string, linkGenerator link.Interface) (component.Component, error)

type crdListDescriptionOption func(*crdList)

type crdList struct {
	base

	name    string
	path    string
	printer crdListPrinter
}

var _ Describer = (*crdList)(nil)

func newCRDList(name, path string, options ...crdListDescriptionOption) *crdList {
	d := &crdList{
		name:    name,
		path:    path,
		printer: printer.CustomResourceListHandler,
	}

	for _, option := range options {
		option(d)
	}

	return d
}

func (cld *crdList) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	objectStore := options.ObjectStore()

	crd, err := CustomResourceDefinition(ctx, cld.name, objectStore)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	printOptions := printer.Options{
		DashConfig: options.Dash,
		Link:       options.Link,
	}

	view, err := printer.CustomResourceDefinitionHandler(ctx, crd, namespace, printOptions)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	view.SetAccessor("summary")

	title := getCrdTitle(namespace, crd, "")

	contentResponse := component.NewContentResponse(title)
	contentResponse.Add(view)

	metadata, err := printer.MetadataHandler(crd, options.Link)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	metadata.SetAccessor("metadata")
	contentResponse.Add(metadata)

	yamlView, err := yamlviewer.ToComponent(crd)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	yamlView.SetAccessor("yaml")

	contentResponse.Add(yamlView)

	return *contentResponse, nil
}

func (cld *crdList) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(cld.path, cld),
	}
}

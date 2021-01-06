/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/printer"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type crdListDescriptionOption func(*crdList)

type crdList struct {
	base

	name string
	path string
}

var _ Describer = (*crdList)(nil)

func newCRDList(name, path string, options ...crdListDescriptionOption) *crdList {
	d := &crdList{
		name: name,
		path: path,
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
		DashConfig:    options.Dash,
		Link:          options.Link,
		ObjectFactory: printer.NewDefaultObjectFactory(),
	}

	title := component.Title(component.NewText(""))
	contentResponse := component.NewContentResponse(title)

	view, err := printer.CustomResourceDefinitionVersionList(ctx, crd, namespace, printOptions)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	m := view.GetMetadata()
	m.Title = title
	view.SetMetadata(m)

	contentResponse.Add(view)

	return *contentResponse, nil
}

func (cld *crdList) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(cld.path, cld),
	}
}

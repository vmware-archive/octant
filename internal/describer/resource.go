/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"path"
	"reflect"
	"strings"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	ResourceNameRegex = "(?P<name>.*?)"
)

type ResourceTitle struct {
	List   string
	Object string
}

type ResourceLink struct {
	Title string
	Url   string
}

type ResourceOptions struct {
	Path                  string
	ObjectStoreKey        store.Key
	ListType              interface{}
	ObjectType            interface{}
	Titles                ResourceTitle
	DisableResourceViewer bool
	ClusterWide           bool
	IconName              string
	RootPath              ResourceLink
}

type Resource struct {
	base

	ResourceOptions
}

func NewResource(options ResourceOptions) *Resource {
	return &Resource{
		ResourceOptions: options,
	}
}

func (r *Resource) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	return r.List().Describe(ctx, namespace, options)
}

func (r *Resource) List() *List {

	return NewList(
		ListConfig{
			Path:     r.Path,
			Title:    r.Titles.List,
			StoreKey: r.ObjectStoreKey,
			ListType: func() interface{} {
				return reflect.New(reflect.ValueOf(r.ListType).Elem().Type()).Interface()
			},
			ObjectType: func() interface{} {
				return reflect.New(reflect.ValueOf(r.ObjectType).Elem().Type()).Interface()
			},
			IsClusterWide: r.ClusterWide,
			RootPath:      r.RootPath,
		},
	)
}

func (r *Resource) Object() *Object {

	return NewObject(
		ObjectConfig{
			Path:      path.Join(r.Path, ResourceNameRegex),
			BaseTitle: r.Titles.Object,
			StoreKey:  r.ObjectStoreKey,
			ObjectType: func() interface{} {
				return reflect.New(reflect.ValueOf(r.ObjectType).Elem().Type()).Interface()
			},
			RootPath:   r.RootPath,
		},
	)
}

func (r *Resource) PathFilters() []PathFilter {
	filters := []PathFilter{
		*NewPathFilter(r.Path, r.List()),
		*NewPathFilter(path.Join(r.Path, ResourceNameRegex), r.Object()),
	}

	return filters
}

func getBreadcrumb(rootPath ResourceLink, objectTitle string, objectUrl string, namespace string) []component.TitleComponent {
	var rootUrl = rootPath.Url
	if strings.Contains(rootPath.Url, "($NAMESPACE)") {
		rootUrl = strings.Replace(rootPath.Url, "($NAMESPACE)", namespace, 1)
	}
	var title []component.TitleComponent
	if len(rootUrl) > 0 {
		title = append(title, component.NewLink("", rootPath.Title, rootUrl))
	}
	title = append(title, component.NewLink("", objectTitle, objectUrl))
	return title
}

func getCrdTitle(namespace string, crd *unstructured.Unstructured, objectName string) []component.TitleComponent {
	var title []component.TitleComponent
	if namespace == "" {
		title= component.Title(component.NewLink("", "Cluster Overview", "/cluster-overview"),
			component.NewLink("", "Custom Resources", "/cluster-overview/custom-resources"))
	} else {
		title = component.Title(component.NewLink("", "Overview", "/overview/namespace/"+namespace),
			component.NewLink("", "Custom Resources", "/overview/namespace/"+namespace+"/custom-resources"))
	}

	if objectName == "" {
		title= append(title, component.NewText(crd.GetName()))
	} else {
		title= append(title, component.NewLink("", crd.GetName(), getCrdUrl(namespace, crd)))
		title= append(title, component.NewText(objectName))
	}
	return title
}

func getCrdUrl(namespace string, crd *unstructured.Unstructured) string {
	ref := path.Join("/overview/namespace", namespace, "custom-resources", crd.GetName())
	if namespace == "" {
		ref = path.Join("/cluster-overview/custom-resources", crd.GetName())
	}
	return ref
}

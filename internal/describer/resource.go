/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"path"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

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

type ResourceOptions struct {
	Path                  string
	ObjectStoreKey        store.Key
	ListType              interface{}
	ObjectType            interface{}
	Titles                ResourceTitle
	DisableResourceViewer bool
	ClusterWide           bool
	IconName              string
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

func getCrdUrl(namespace string, crd *unstructured.Unstructured) string {
	ref := path.Join("/overview/namespace", namespace, "custom-resources", crd.GetName())
	if namespace == "" {
		ref = path.Join("/cluster-overview/custom-resources", crd.GetName())
	}
	return ref
}

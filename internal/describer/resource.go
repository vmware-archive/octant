/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"path"
	"reflect"

	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
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
}

type Resource struct {
	ResourceOptions
}

func NewResource(options ResourceOptions) *Resource {
	return &Resource{
		ResourceOptions: options,
	}
}

func (r *Resource) Describe(ctx context.Context, prefix, namespace string, options Options) (component.ContentResponse, error) {
	return r.List().Describe(ctx, prefix, namespace, options)
}

func (r *Resource) List() *List {
	return NewList(
		r.Path,
		r.Titles.List,
		r.ObjectStoreKey,
		func() interface{} {
			return reflect.New(reflect.ValueOf(r.ListType).Elem().Type()).Interface()
		},
		func() interface{} {
			return reflect.New(reflect.ValueOf(r.ObjectType).Elem().Type()).Interface()
		},
		r.ClusterWide,
	)
}

func (r *Resource) Object() *Object {
	return NewObject(
		path.Join(r.Path, ResourceNameRegex),
		r.Titles.Object,
		r.ObjectStoreKey,
		func() interface{} {
			return reflect.New(reflect.ValueOf(r.ObjectType).Elem().Type()).Interface()
		},
		r.DisableResourceViewer,
	)
}

func (r *Resource) PathFilters() []PathFilter {
	filters := []PathFilter{
		*NewPathFilter(r.Path, r.List()),
		*NewPathFilter(path.Join(r.Path, ResourceNameRegex), r.Object()),
	}

	return filters
}

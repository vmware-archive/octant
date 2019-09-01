/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"fmt"
	"path"
	"reflect"

	"github.com/vmware/octant/pkg/icon"
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

func (r *Resource) Describe(ctx context.Context, prefix, namespace string, options Options) (component.ContentResponse, error) {
	return r.List().Describe(ctx, prefix, namespace, options)
}

func (r *Resource) List() *List {
	iconName, iconSource := loadIcon(r.IconName)

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
			IconName:      iconName,
			IconSource:    iconSource,
		},
	)
}

func (r *Resource) Object() *Object {
	iconName, iconSource := loadIcon(r.IconName)

	return NewObject(
		ObjectConfig{
			Path:      path.Join(r.Path, ResourceNameRegex),
			BaseTitle: r.Titles.Object,
			StoreKey:  r.ObjectStoreKey,
			ObjectType: func() interface{} {
				return reflect.New(reflect.ValueOf(r.ObjectType).Elem().Type()).Interface()
			},
			DisableResourceViewer: r.DisableResourceViewer,
			IconName:              iconName,
			IconSource:            iconSource,
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

func loadIcon(name string) (string, string) {
	source, err := icon.LoadIcon(name)
	if err != nil {
		return name, ""
	}

	internalName := fmt.Sprintf("internal:%s", name)

	return internalName, source
}

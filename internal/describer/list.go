/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

type ListConfig struct {
	Path          string
	Title         string
	StoreKey      store.Key
	ListType      func() interface{}
	ObjectType    func() interface{}
	IsClusterWide bool
	IconName      string
	IconSource    string
}

// List describes a list of objects.
type List struct {
	*base

	path           string
	title          string
	listType       func() interface{}
	objectType     func() interface{}
	objectStoreKey store.Key
	isClusterWide  bool
	iconName       string
	iconSource     string
}

// NewList creates an instance of List.
func NewList(c ListConfig) *List {
	return &List{
		path:           c.Path,
		title:          c.Title,
		base:           newBaseDescriber(),
		objectStoreKey: c.StoreKey,
		listType:       c.ListType,
		objectType:     c.ObjectType,
		isClusterWide:  c.IsClusterWide,
		iconName:       c.IconName,
		iconSource:     c.IconSource,
	}
}

// Describe creates content.
func (d *List) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	if options.Printer == nil {
		return component.EmptyContentResponse, errors.New("object list Describer requires a printer")
	}

	// Pass through selector if provided to filter objects
	var key = d.objectStoreKey // copy
	key.Selector = options.LabelSet

	if d.isClusterWide {
		namespace = ""
	}

	objectList, err := options.LoadObjects(ctx, namespace, options.Fields, []store.Key{key})
	if err != nil {
		logger := log.From(ctx)
		logger.WithErr(err).Warnf("error loading objects")
		objectList = &unstructured.UnstructuredList{}
	}

	list := component.NewList(d.title, nil)
	list.SetIcon(d.iconName, d.iconSource)

	listType := d.listType()

	v := reflect.ValueOf(listType)
	f := reflect.Indirect(v).FieldByName("Items")

	// Convert unstructured objects to typed runtime objects
	for i := range objectList.Items {
		item := d.objectType()
		if err := scheme.Scheme.Convert(&objectList.Items[i], item, nil); err != nil {
			return component.EmptyContentResponse, err
		}

		if err := copyObjectMeta(item, &objectList.Items[i]); err != nil {
			return component.EmptyContentResponse, err
		}

		newSlice := reflect.Append(f, reflect.ValueOf(item).Elem())
		f.Set(newSlice)
	}

	listObject, ok := listType.(runtime.Object)
	if !ok {
		return component.EmptyContentResponse, errors.Errorf("expected list to be a runtime object. It was a %T",
			listType)
	}

	viewComponent, err := options.Printer.Print(ctx, listObject, options.PluginManager())
	if err != nil {
		return component.EmptyContentResponse, err
	}

	if viewComponent != nil {
		if table, ok := viewComponent.(*component.Table); ok {
			list.Add(table)
		} else {
			list.Add(viewComponent)
		}
	}

	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

// PathFilters returns path filters for this Describer.
func (d *List) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(d.path, d),
	}
}

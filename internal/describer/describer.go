/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/internal/modules/overview/printer"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

// EmptyContentResponse is an empty content response.
var EmptyContentResponse = component.ContentResponse{}

type ObjectLoaderFactory struct {
	dashConfig config.Dash
}

func NewObjectLoaderFactory(dashConfig config.Dash) *ObjectLoaderFactory {
	return &ObjectLoaderFactory{
		dashConfig: dashConfig,
	}
}

func (f *ObjectLoaderFactory) LoadObject(ctx context.Context, namespace string, fields map[string]string, objectStoreKey store.Key) (*unstructured.Unstructured, error) {
	return LoadObject(ctx, f.dashConfig.ObjectStore(), namespace, fields, objectStoreKey)
}

func (f *ObjectLoaderFactory) LoadObjects(ctx context.Context, namespace string, fields map[string]string, objectStoreKeys []store.Key) (*unstructured.UnstructuredList, error) {
	return LoadObjects(ctx, f.dashConfig.ObjectStore(), namespace, fields, objectStoreKeys)
}

// loadObject loads a single object from the object store.
func LoadObject(ctx context.Context, objectStore store.Store, namespace string, fields map[string]string, objectStoreKey store.Key) (*unstructured.Unstructured, error) {
	objectStoreKey.Namespace = namespace

	if name, ok := fields["name"]; ok && name != "" {
		objectStoreKey.Name = name
	}

	object, err := objectStore.Get(ctx, objectStoreKey)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// loadObjects loads objects from the object store sorted by their name.
func LoadObjects(ctx context.Context, objectStore store.Store, namespace string, fields map[string]string, objectStoreKeys []store.Key) (*unstructured.UnstructuredList, error) {
	list := &unstructured.UnstructuredList{}

	for _, objectStoreKey := range objectStoreKeys {
		objectStoreKey.Namespace = namespace

		if name, ok := fields["name"]; ok && name != "" {
			objectStoreKey.Name = name
		}

		storedObjects, err := objectStore.List(ctx, objectStoreKey)
		if err != nil {
			return nil, err
		}

		list.Items = append(list.Items, storedObjects.Items...)
	}

	sort.SliceStable(list.Items, func(i, j int) bool {
		a, b := list.Items[i], list.Items[j]
		return a.GetName() < b.GetName()
	})

	return list, nil
}

// LoaderFunc loads an object from the object store.
type LoaderFunc func(ctx context.Context, o store.Store, namespace string, fields map[string]string) (*unstructured.Unstructured, error)

// Options provides options to describers
type Options struct {
	config.Dash

	Queryer  queryer.Queryer
	Fields   map[string]string
	Printer  printer.Printer
	LabelSet *kLabels.Set
	Link     link.Interface

	LoadObjects func(ctx context.Context, namespace string, fields map[string]string, objectStoreKeys []store.Key) (*unstructured.UnstructuredList, error)
	LoadObject  func(ctx context.Context, namespace string, fields map[string]string, objectStoreKey store.Key) (*unstructured.Unstructured, error)
}

// Describer creates content.
type Describer interface {
	Describe(ctx context.Context, prefix, namespace string, options Options) (component.ContentResponse, error)
	PathFilters() []PathFilter
}

type base struct{}

func newBaseDescriber() *base {
	return &base{}
}

func copyObjectMeta(to interface{}, from *unstructured.Unstructured) error {
	object, ok := to.(metav1.Object)
	if !ok {
		return errors.Errorf("%T is not an object", to)
	}

	t, err := meta.TypeAccessor(object)
	if err != nil {
		return errors.Wrapf(err, "accessing type meta")
	}
	t.SetAPIVersion(from.GetAPIVersion())
	t.SetKind(from.GetObjectKind().GroupVersionKind().Kind)

	object.SetNamespace(from.GetNamespace())
	object.SetName(from.GetName())
	object.SetGenerateName(from.GetGenerateName())
	object.SetUID(from.GetUID())
	object.SetResourceVersion(from.GetResourceVersion())
	object.SetGeneration(from.GetGeneration())
	object.SetSelfLink(from.GetSelfLink())
	object.SetCreationTimestamp(from.GetCreationTimestamp())
	object.SetDeletionTimestamp(from.GetDeletionTimestamp())
	object.SetDeletionGracePeriodSeconds(from.GetDeletionGracePeriodSeconds())
	object.SetLabels(from.GetLabels())
	object.SetAnnotations(from.GetAnnotations())
	object.SetInitializers(from.GetInitializers())
	object.SetOwnerReferences(from.GetOwnerReferences())
	object.SetClusterName(from.GetClusterName())
	object.SetFinalizers(from.GetFinalizers())

	return nil
}

func isPod(object runtime.Object) bool {
	gvk := object.GetObjectKind().GroupVersionKind()
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	return apiVersion == "v1" && kind == "Pod"
}

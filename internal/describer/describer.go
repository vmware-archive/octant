/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/config"
	oerrors "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/printer"
	"github.com/vmware-tanzu/octant/internal/queryer"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type ObjectLoaderFactory struct {
	dashConfig config.Dash
}

func NewObjectLoaderFactory(dashConfig config.Dash) *ObjectLoaderFactory {
	return &ObjectLoaderFactory{
		dashConfig: dashConfig,
	}
}

func (f *ObjectLoaderFactory) LoadObject(ctx context.Context, namespace string, fields map[string]string, objectStoreKey store.Key) (*unstructured.Unstructured, error) {
	return LoadObject(ctx, f.dashConfig.ObjectStore(), f.dashConfig.ErrorStore(), namespace, fields, objectStoreKey)
}

func (f *ObjectLoaderFactory) LoadObjects(ctx context.Context, namespace string, fields map[string]string, objectStoreKeys []store.Key) (*unstructured.UnstructuredList, error) {
	return LoadObjects(ctx, f.dashConfig.ObjectStore(), f.dashConfig.ErrorStore(), namespace, fields, objectStoreKeys)
}

// loadObject loads a single object from the object store.
func LoadObject(ctx context.Context, objectStore store.Store, errorStore oerrors.ErrorStore, namespace string, fields map[string]string, objectStoreKey store.Key) (*unstructured.Unstructured, error) {
	objectStoreKey.Namespace = namespace

	if name, ok := fields["name"]; ok && name != "" {
		objectStoreKey.Name = name
	}

	object, found, err := objectStore.Get(ctx, objectStoreKey)
	if err != nil {
		var ae *oerrors.AccessError
		if errors.As(err, &ae) {
			if ae.Name() == oerrors.OctantAccessError {
				found := errorStore.Add(ae)
				if !found {
					logger := log.From(ctx)
					logger.WithErr(ae).Errorf("loadObject")
				}
				return &unstructured.Unstructured{}, nil
			}
		}
		return nil, err
	}
	if !found {
		return nil, errors.New("object was not found")
	}

	return object, nil
}

// loadObjects loads objects from the object store sorted by their name.
func LoadObjects(ctx context.Context, objectStore store.Store, errorStore oerrors.ErrorStore, namespace string, fields map[string]string, objectStoreKeys []store.Key) (*unstructured.UnstructuredList, error) {
	list := &unstructured.UnstructuredList{}

	for _, objectStoreKey := range objectStoreKeys {
		objectStoreKey.Namespace = namespace

		if name, ok := fields["name"]; ok && name != "" {
			objectStoreKey.Name = name
		}

		storedObjects, _, err := objectStore.List(ctx, objectStoreKey)
		if err != nil {
			var ae *oerrors.AccessError
			if errors.As(err, &ae) {
				if ae.Name() == oerrors.OctantAccessError {
					logger := log.From(ctx)
					_ = objectStore.Unwatch(ctx, objectStoreKey.GroupVersionKind())
					found := errorStore.Add(ae)
					if !found {
						logger.WithErr(ae).Errorf("loadObjects")
					}
					return &unstructured.UnstructuredList{}, nil
				}
			}
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
	Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error)
	PathFilters() []PathFilter
	Reset(ctx context.Context) error
}

type base struct{}

func (b base) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	return component.EmptyContentResponse, nil
}

func (b base) PathFilters() []PathFilter {
	return nil
}

func (b base) Reset(ctx context.Context) error {
	return nil
}

var _ Describer = (*base)(nil)

func newBaseDescriber() *base {
	return &base{}
}

func copyObjectMeta(to interface{}, from *unstructured.Unstructured) error {
	object, ok := to.(metav1.Object)
	if !ok {
		return fmt.Errorf("%T is not an object", to)
	}

	t, err := meta.TypeAccessor(object)
	if err != nil {
		return fmt.Errorf("accessing type meta: %w", err)
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

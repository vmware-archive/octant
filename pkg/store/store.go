/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package store

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"

	"github.com/vmware/octant/internal/cluster"
)

//go:generate mockgen  -destination=./fake/mock_store.go -package=fake github.com/vmware/octant/pkg/store Store

// UpdateFn is a function that is called when
type UpdateFn func(store Store)

// Store stores Kubernetes objects.
type Store interface {
	List(ctx context.Context, key Key) (list *unstructured.UnstructuredList, loading bool, err error)
	Get(ctx context.Context, key Key) (object *unstructured.Unstructured, found bool, err error)
	Watch(ctx context.Context, key Key, handler cache.ResourceEventHandler) error
	HasAccess(context.Context, Key, string) error
	UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error
	RegisterOnUpdate(fn UpdateFn)
	Update(ctx context.Context, key Key, updater func(*unstructured.Unstructured) error) error
	IsLoading(ctx context.Context, key Key) bool
}

// Key is a key for the object store.
type Key struct {
	Namespace  string
	APIVersion string
	Kind       string
	Name       string
	Selector   *labels.Set
}

func (k Key) String() string {
	var sb strings.Builder

	sb.WriteString("CacheKey[")
	if k.Namespace != "" {
		sb.WriteString(fmt.Sprintf("Namespace='%s', ", k.Namespace))
	}
	sb.WriteString(fmt.Sprintf("APIVersion='%s', ", k.APIVersion))
	sb.WriteString(fmt.Sprintf("Kind='%s'", k.Kind))

	if k.Name != "" {
		sb.WriteString(fmt.Sprintf(", Name='%s'", k.Name))
	}

	if k.Selector != nil && k.Selector.String() != "" {
		sb.WriteString(fmt.Sprintf(", Selector='%s'", k.Selector.String()))
	}

	sb.WriteString("]")

	return sb.String()
}

// GroupVersionKind converts the Key to a GroupVersionKind.
func (k Key) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(k.APIVersion, k.Kind)
}

// KeyFromObject creates a key from a runtime object.
func KeyFromObject(object runtime.Object) (Key, error) {
	accessor := meta.NewAccessor()

	namespace, err := accessor.Namespace(object)
	if err != nil {
		return Key{}, err
	}

	apiVersion, err := accessor.APIVersion(object)
	if err != nil {
		return Key{}, err
	}

	kind, err := accessor.Kind(object)
	if err != nil {
		return Key{}, err
	}

	name, err := accessor.Name(object)
	if err != nil {
		return Key{}, err
	}

	return Key{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}, nil
}

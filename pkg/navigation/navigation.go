/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package navigation

import (
	"context"
	"fmt"
	"path"
	"sort"

	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/pkg/icon"
	"github.com/vmware/octant/pkg/store"
	octantunstructured "github.com/vmware/octant/thirdparty/unstructured"
)

// Option is an option for configuring navigation.
type Option func(*Navigation) error

// SetNavigationIcon sets the icon for the navigation entry.
func SetNavigationIcon(name string) Option {
	return func(n *Navigation) error {
		if name == "" {
			return nil
		}

		source, err := icon.LoadIcon(name)
		if err != nil {
			return err
		}

		n.IconName = fmt.Sprintf("internal:%s", name)
		n.IconSource = source

		return nil
	}
}

// Navigation is a set of navigation entries.
type Navigation struct {
	Title      string       `json:"title,omitempty"`
	Path       string       `json:"path,omitempty"`
	Children   []Navigation `json:"children,omitempty"`
	IconName   string       `json:"iconName,omitempty"`
	IconSource string       `json:"iconSource,omitempty"`
}

// New creates a Navigation.
func New(title, navigationPath string, options ...Option) (*Navigation, error) {
	navigation := &Navigation{Title: title, Path: navigationPath}

	for _, option := range options {
		if err := option(navigation); err != nil {
			return nil, err
		}
	}

	return navigation, nil
}

// CRDEntries generates navigation entries for CRDs.
func CRDEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, wantsClusterScoped bool) ([]Navigation, error) {
	var list []Navigation

	crds, err := CustomResourceDefinitions(ctx, objectStore)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving CRDs")
	}

	sort.Slice(crds, func(i, j int) bool {
		return crds[i].Name < crds[j].Name
	})

	for i := range crds {
		if wantsClusterScoped && crds[i].Spec.Scope != apiextv1beta1.ClusterScoped {
			continue
		} else if !wantsClusterScoped && crds[i].Spec.Scope != apiextv1beta1.NamespaceScoped {
			continue
		}

		objects, err := ListCustomResources(ctx, crds[i], namespace, objectStore, nil)
		if err != nil {
			return nil, err
		}

		if len(objects.Items) > 0 {
			navigation, err := New(crds[i].Name, path.Join(prefix, crds[i].Name), SetNavigationIcon(icon.CustomResourceDefinition))
			if err != nil {
				return nil, err
			}

			list = append(list, *navigation)
		}
	}

	return list, nil
}

func CustomResourceDefinitions(ctx context.Context, o store.Store) ([]*apiextv1beta1.CustomResourceDefinition, error) {
	key := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	rawList, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "listing CRDs")
	}

	logger := log.From(ctx)

	var list []*apiextv1beta1.CustomResourceDefinition
	for i := range rawList.Items {
		crd := &apiextv1beta1.CustomResourceDefinition{}

		// NOTE: (bryanl) vendored converter can't convert from int64 to float64. Watching
		// https://github.com/kubernetes-sigs/yaml/pull/14 to see when it gets pulled into
		// a release so Octant can switch back.
		if err := octantunstructured.DefaultUnstructuredConverter.FromUnstructured(rawList.Items[i].Object, crd); err != nil {
			logger.Errorf("%v", errors.Wrapf(errors.Wrapf(err, "converting unstructured object to custom resource definition"), rawList.Items[i].GetName()))
			continue
		}
		list = append(list, crd)
	}

	return list, nil
}

// CustomResourceDefinition retrieves a CRD.
func CustomResourceDefinition(ctx context.Context, name string, o store.Store) (*apiextv1beta1.CustomResourceDefinition, error) {
	key := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
		Name:       name,
	}

	crd := &apiextv1beta1.CustomResourceDefinition{}
	if err := store.GetObjectAs(ctx, o, key, crd); err != nil {
		return nil, errors.Wrap(err, "get CRD from object store")
	}

	return crd, nil
}

// ListCustomResources lists all custom resources given a CRD.
func ListCustomResources(
	ctx context.Context,
	crd *apiextv1beta1.CustomResourceDefinition,
	namespace string,
	o store.Store,
	selector *labels.Set) (*unstructured.UnstructuredList, error) {
	if crd == nil {
		return nil, errors.New("crd is nil")
	}
	gvk := schema.GroupVersionKind{
		Group:   crd.Spec.Group,
		Version: crd.Spec.Version,
		Kind:    crd.Spec.Names.Kind,
	}

	apiVersion, kind := gvk.ToAPIVersionAndKind()

	key := store.Key{
		APIVersion: apiVersion,
		Kind:       kind,
		Selector:   selector,
	}

	if crd.Spec.Scope == apiextv1beta1.NamespaceScoped {
		key.Namespace = namespace
	}

	objects, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "listing custom resources for %q", crd.Name)
	}

	return objects, nil
}

type navConfig struct {
	title    string
	suffix   string
	iconName string
}

// EntriesHelper generates navigation entries.
type EntriesHelper struct {
	navConfigs []navConfig
}

// Add adds an entry.
func (neh *EntriesHelper) Add(title, suffix, iconName string) {
	neh.navConfigs = append(neh.navConfigs, navConfig{
		title: title, suffix: suffix, iconName: iconName,
	})
}

// Generate generates navigation entries.
func (neh *EntriesHelper) Generate(prefix string) ([]Navigation, error) {
	var navigations []Navigation

	for _, nc := range neh.navConfigs {
		navigation, err := New(nc.title, path.Join(prefix, nc.suffix), SetNavigationIcon(nc.iconName))
		if err != nil {
			return nil, err
		}

		navigations = append(navigations, *navigation)
	}

	return navigations, nil
}

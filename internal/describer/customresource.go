/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"fmt"
	"path"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/pkg/store"
)

func CustomResourceDefinition(ctx context.Context, name string, o store.Store) (*unstructured.Unstructured, error) {
	key := store.KeyFromGroupVersionKind(gvk.CustomResourceDefinition)
	key.Name = name

	crd, err := o.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", key, err)
	}

	return crd, nil
}

func AddCRD(ctx context.Context, crd *unstructured.Unstructured, pm *PathMatcher, crdSection *CRDSection, m module.Module) {
	name := crd.GetName()

	logger := log.From(ctx).With("crd-name", name, "module", m.Name())
	logger.Debugf("adding CRD")

	cld := newCRDList(name, crdListPath(name))

	// TODO: this should add a list of custom resource definitions (GH#509)
	crdSection.Add(name, cld)

	for _, pf := range cld.PathFilters() {
		pm.Register(ctx, pf)
	}

	// TODO: there could be multiple paths here, so iterate through crd.spec.versions['name']

	versions, err := crdVersions(crd)
	if err != nil {
		logger.WithErr(err).Errorf("get crd versions: %w", err)
		return
	}

	for _, version := range versions {
		cd := newCRD(name, crdObjectPath(crd, version))
		for _, pf := range cd.PathFilters() {
			pm.Register(ctx, pf)
		}
	}

	if err := m.AddCRD(ctx, crd); err != nil {
		logger.WithErr(err).Errorf("unable to add CRD")
	}
}

func crdVersions(crd *unstructured.Unstructured) ([]string, error) {
	rawVersions, _, err := unstructured.NestedSlice(crd.Object, "spec", "versions")
	if err != nil {
		return nil, fmt.Errorf("unable to extract versions from crd: %w", err)
	}

	var list []string

	for _, rawVersion := range rawVersions {
		versionDesc, ok := rawVersion.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("crd version descriptor was not an object (it was %T)", rawVersion)
		}
		versionName, found, err := unstructured.NestedString(versionDesc, "name")
		if err != nil {
			return nil, fmt.Errorf("unable to extract version name from version descriptor: %w", err)
		}

		if found {
			list = append(list, versionName)
		}
	}

	return list, nil
}

func DeleteCRD(ctx context.Context, crd *unstructured.Unstructured, pm *PathMatcher, crdSection *CRDSection, m module.Module) {
	name := crd.GetName()

	logger := log.From(ctx).With("crd-name", name, "module", m.Name())
	logger.Debugf("deleting CRD")

	pm.Deregister(ctx, crdListPath(name))

	versions, err := crdVersions(crd)
	if err != nil {
		logger.WithErr(err).Errorf("get crd versions: %w", err)
		return
	}
	for _, version := range versions {
		pm.Deregister(ctx, crdObjectPath(crd, version))
	}

	crdSection.Remove(name)

	if err := m.RemoveCRD(ctx, crd); err != nil {
		logger.WithErr(err).Errorf("unable to remove CRD")
	}

}

func crdListPath(name string) string {
	return path.Join("/custom-resources", name)
}

func crdObjectPath(object *unstructured.Unstructured, version string) string {
	return path.Join(crdListPath(object.GetName()), version, ResourceNameRegex)
}

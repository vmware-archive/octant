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
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/store"
)

func CustomResourceDefinition(ctx context.Context, name string, o store.Store) (*unstructured.Unstructured, error) {
	key := store.KeyFromGroupVersionKind(gvk.CustomResourceDefinition)
	key.Name = name

	crd, _, err := o.Get(ctx, key)
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

	// TODO: this should add a list of custom resource definitions
	crdSection.Add(name, cld)

	for _, pf := range cld.PathFilters() {
		pm.Register(ctx, pf)
	}

	cd := newCRD(name, crdObjectPath(name))
	for _, pf := range cd.PathFilters() {
		pm.Register(ctx, pf)
	}

	if err := m.AddCRD(ctx, crd); err != nil {
		logger.WithErr(err).Errorf("unable to add CRD")
	}
}

func DeleteCRD(ctx context.Context, crd *unstructured.Unstructured, pm *PathMatcher, crdSection *CRDSection, m module.Module, s store.Store) {
	name := crd.GetName()

	logger := log.From(ctx).With("crd-name", name, "module", m.Name())
	logger.Debugf("deleting CRD")

	pm.Deregister(ctx, crdListPath(name))
	pm.Deregister(ctx, crdObjectPath(name))

	crdSection.Remove(name)

	if err := m.RemoveCRD(ctx, crd); err != nil {
		logger.WithErr(err).Errorf("unable to remove CRD")
	}

	list, err := kubernetes.CRDResources(crd)
	if err != nil {
		logger.WithErr(err).Errorf("unable to get group/version/kinds for CRD")

	}

	if err := s.Unwatch(ctx, list...); err != nil {
		logger.WithErr(err).Errorf("unable to unwatch CRD")
		return
	}
}

func crdListPath(name string) string {
	return path.Join("/custom-resources", name)
}

func crdObjectPath(name string) string {
	return path.Join(crdListPath(name), ResourceNameRegex)
}

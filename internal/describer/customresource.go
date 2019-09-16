/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"path"

	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/util/kubernetes"
	"github.com/vmware/octant/pkg/store"
)

func CustomResourceDefinition(ctx context.Context, name string, o store.Store) (*apiextv1beta1.CustomResourceDefinition, error) {
	key := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
		Name:       name,
	}

	crd := &apiextv1beta1.CustomResourceDefinition{}
	found, err := store.GetAs(ctx, o, key, crd)
	if err != nil {
		return nil, errors.Wrapf(err, "get object as custom resource definition %q from store", name)
	}
	if !found {
		return nil, errors.Errorf("custom resource definition %q was not found", name)
	}

	return crd, nil
}

func AddCRD(ctx context.Context, crd *unstructured.Unstructured, pm *PathMatcher, crdSection *CRDSection, m module.Module) {
	name := crd.GetName()

	logger := log.From(ctx).With("crd-name", name, "module", m.Name())
	logger.Debugf("adding CRD")

	cld := newCRDList(name, crdListPath(name))

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

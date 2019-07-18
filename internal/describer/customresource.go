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
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/pkg/store"
)

func CustomResourceDefinitionNames(ctx context.Context, o store.Store) ([]string, error) {
	key := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	if err := o.HasAccess(ctx, key, "list"); err != nil {
		return []string{}, nil
	}

	rawList, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "listing CRDs")
	}

	var list []string

	for _, object := range rawList {
		crd := &apiextv1beta1.CustomResourceDefinition{}

		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, crd); err != nil {
			return nil, errors.Wrap(err, "crd conversion failed")
		}

		list = append(list, crd.Name)
	}

	return list, nil
}

func CustomResourceDefinition(ctx context.Context, name string, o store.Store) (*apiextv1beta1.CustomResourceDefinition, error) {
	key := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
		Name:       name,
	}

	crd := &apiextv1beta1.CustomResourceDefinition{}
	if err := store.GetAs(ctx, o, key, crd); err != nil {
		return nil, errors.Wrap(err, "get CRD from object store")
	}

	return crd, nil
}

func AddCRD(ctx context.Context, crd *unstructured.Unstructured, pm *PathMatcher, crdSection *CRDSection, m module.Module) {
	name := crd.GetName()

	logger := log.From(ctx)
	logger.With("crd-name", name, "module", m.Name()).Debugf("adding CRD")

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
		logger.With("err", err).Errorf("unable to add CRD")
	}
}

func DeleteCRD(ctx context.Context, crd *unstructured.Unstructured, pm *PathMatcher, crdSection *CRDSection, m module.Module) {
	name := crd.GetName()

	logger := log.From(ctx)
	logger.With("crd-name", name).Debugf("deleting CRD")

	pm.Deregister(ctx, crdListPath(name))
	pm.Deregister(ctx, crdObjectPath(name))

	crdSection.Remove(name)

	if err := m.RemoveCRD(ctx, crd); err != nil {
		logger.With("err", err).Errorf("unable to remove CRD")
	}
}

func crdListPath(name string) string {
	return path.Join("/custom-resources", name)
}

func crdObjectPath(name string) string {
	return path.Join(crdListPath(name), ResourceNameRegex)
}

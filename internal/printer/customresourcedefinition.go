/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package printer

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func CustomResourceDefinitionHandler(ctx context.Context, crd *unstructured.Unstructured, namespace string, options Options) (component.Component, error) {
	octantCRD, err := octant.NewCustomResourceDefinition(crd)
	if err != nil {
		return nil, err
	}

	objectStore := options.DashConfig.ObjectStore()

	versions, err := octantCRD.Versions()
	if err != nil {
		return nil, err
	}

	list := component.NewList(nil, nil)

	for i := range versions {
		version := versions[i]

		crGVK, err := gvk.CustomResource(crd, version)
		if err != nil {
			return nil, err
		}

		key := store.KeyFromGroupVersionKind(crGVK)
		key.Namespace = namespace

		customResources, _, err := objectStore.List(ctx, key)
		if err != nil {
			return nil, err
		}

		view, err := CustomResourceListHandler(crd, customResources, version, options.Link)
		if err != nil {
			return nil, err
		}

		list.Add(view)
	}

	return list, nil
}

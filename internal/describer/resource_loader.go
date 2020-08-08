/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package describer

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// ResourceDescriptor describes a custom resource.
type ResourceDescriptor struct {
	CustomResourceDefinitionName string
	Namespace                    string
	CustomResourceVersion        string
	CustomResourceName           string
}

// ResourceLoadResponse is a response from a from ResourceLoader.
type ResourceLoadResponse struct {
	// CustomResource is the custom resource.
	CustomResource *unstructured.Unstructured
	// CustomResourceDefinition is the custom resource definition.
	CustomResourceDefinition *unstructured.Unstructured
}

// ResourceLoader is an interface which loads a custom resource.
type ResourceLoader interface {
	// Load loads a custom resource given a ResourceDescriptor.
	Load(ctx context.Context, d ResourceDescriptor) (*ResourceLoadResponse, error)
}

// StoreResourceLoader is an interface which loads a custom resource using Octant's store.
type StoreResourceLoader struct {
	store store.Store
}

var _ ResourceLoader = &StoreResourceLoader{}

// NewStoreResourceLoader creates an instance of StoreResourceLoader.
func NewStoreResourceLoader(s store.Store) *StoreResourceLoader {
	rl := &StoreResourceLoader{
		store: s,
	}
	return rl
}

// Load loads a custom resource.
func (rl *StoreResourceLoader) Load(ctx context.Context, d ResourceDescriptor) (*ResourceLoadResponse, error) {
	crd, err := CustomResourceDefinition(ctx, d.CustomResourceDefinitionName, rl.store)
	if err != nil {
		return nil, fmt.Errorf("find custom resource definition %q: %w", d.CustomResourceDefinitionName, err)
	}

	crGVK, err := gvk.CustomResource(crd, d.CustomResourceVersion)
	if err != nil {
		return nil, fmt.Errorf("get group/version/kind for custom resource: %w", err)
	}

	apiVersion, kind := crGVK.ToAPIVersionAndKind()

	key := store.Key{
		Namespace:  d.Namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       d.CustomResourceName,
	}

	object, err := rl.store.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if object == nil {
		return nil, errors.New("object is nil")
	}

	return &ResourceLoadResponse{
		CustomResource:           object,
		CustomResourceDefinition: crd,
	}, nil
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package module

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

//go:generate mockgen -destination=./fake/mock_module.go -package=fake github.com/vmware-tanzu/octant/internal/module Module

// ContentOptions are additional options for content generation
type ContentOptions struct {
	LabelSet *labels.Set
}

// Module is an octant plugin.
type Module interface {
	// Name is the name of the module.
	Name() string
	// ClientRequestHandlers are handlers for handling client requests.
	ClientRequestHandlers() []octant.ClientRequestHandler
	// Content generates content for a path.
	Content(ctx context.Context, contentPath string, opts ContentOptions) (component.ContentResponse, error)
	// ContentPath will be used to construct content paths.
	ContentPath() string
	// Navigation returns navigation entries for this module.
	Navigation(ctx context.Context, namespace, root string) ([]navigation.Navigation, error)
	// SetNamespace is called when the current namespace changes.
	SetNamespace(namespace string) error
	// Start starts the module.
	Start() error
	// Stop stops the module.
	Stop()

	// SetContext sets the current context name.
	SetContext(ctx context.Context, contextName string) error

	// Generators allow modules to send events to the frontend.
	Generators() []octant.Generator

	// SupportedGroupVersionKind returns a slice of supported GVKs it owns.
	SupportedGroupVersionKind() []schema.GroupVersionKind

	// GroupVersionKindPath returns the path for an object . It will
	// return an error if it is unable to generate a path
	GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error)

	// AddCRD adds a CRD this module is responsible for.
	AddCRD(ctx context.Context, crd *unstructured.Unstructured) error

	// RemoveCRD removes a CRD this module was responsible for.
	RemoveCRD(ctx context.Context, crd *unstructured.Unstructured) error

	// ResetCRDs removes all CRDs this module is responsible for.
	ResetCRDs(ctx context.Context) error
}

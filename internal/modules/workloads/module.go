/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package workloads

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/util/path_util"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/generator"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/icon"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// Options for configuring Module.
type Options struct {
	DashConfig config.Dash
}

// Module contains the implementation for the workloads module.
type Module struct {
	Options
	pathMatcher *describer.PathMatcher
}

var _ module.Module = (*Module)(nil)

// New creates an instance of Module.
func New(ctx context.Context, options Options) (*Module, error) {
	rootDescriber, err := NewHomeDescriber()
	if err != nil {
		return nil, fmt.Errorf("create home describe: %w", err)
	}

	pm := describer.NewPathMatcher("workloads")

	for _, pf := range rootDescriber.PathFilters() {
		pm.Register(ctx, pf)
	}

	detailDescriber, err := NewDetailDescriber()
	if err != nil {
		return nil, fmt.Errorf("create detail describer: %w", err)
	}

	for _, pf := range detailDescriber.PathFilters() {
		pm.Register(ctx, pf)
	}

	m := &Module{
		Options:     options,
		pathMatcher: pm,
	}

	return m, nil
}

// Name returns the module name.
func (m *Module) Name() string {
	return "workloads"
}

// Description returns the module description.
func (m *Module) Description() string {
	return "Application module displays all known applications and their status"
}

// ClientRequestHandlers returns nil.
func (m *Module) ClientRequestHandlers() []octant.ClientRequestHandler {
	return nil
}

// Content handles content for the module.
func (m *Module) Content(ctx context.Context, contentPath string, opts module.ContentOptions) (component.ContentResponse, error) {
	g, err := generator.NewGenerator(m.pathMatcher, m.DashConfig)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	return g.Generate(ctx, contentPath, generator.Options{})
}

// ContentPath returns the content path for this module.
func (m *Module) ContentPath() string {
	return m.Name()
}

// Navigation returns navigation entries for the module.
func (m *Module) Navigation(ctx context.Context, namespace, root string) ([]navigation.Navigation, error) {
	rootPath := path_util.NamespacedPath(m.ContentPath(), namespace)

	rootNav := navigation.Navigation{
		Title:    "Applications",
		Path:     rootPath,
		IconName: icon.Applications,
	}

	return []navigation.Navigation{rootNav}, nil
}

// SetNamespace is a no-op.
func (m Module) SetNamespace(namespace string) error {
	return nil
}

// Start is a no-op.
func (m Module) Start() error {
	return nil
}

// Stop is a no-op.
func (m Module) Stop() {
}

// SetContext is a no-op.
func (m Module) SetContext(ctx context.Context, contextName string) error {
	return nil
}

// Generators returns nil.
func (m Module) Generators() []octant.Generator {
	return nil
}

// SupportedGroupVersionKind returns nil.
func (m Module) SupportedGroupVersionKind() []schema.GroupVersionKind {
	return nil
}

// GroupVersionKindPath return return an error as this module does not support.
func (m Module) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	return "", fmt.Errorf("not supported")
}

// AddCRD is a no-op.
func (m Module) AddCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

// RemoveCRD is a no-op.
func (m Module) RemoveCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

// ResetCRDs is a no-op.
func (m Module) ResetCRDs(ctx context.Context) error {
	return nil
}

func (m Module) GvkFromPath(contentPath, namespace string) (schema.GroupVersionKind, error) {
	return schema.GroupVersionKind{}, errors.Errorf("not supported")
}

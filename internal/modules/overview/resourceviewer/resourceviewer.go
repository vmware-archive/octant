/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package resourceviewer

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/modules/overview/objectvisitor"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/internal/util/kubernetes"
	"github.com/vmware/octant/pkg/view/component"
)

// ViewerOpt is an option for ResourceViewer.
type ViewerOpt func(*ResourceViewer) error

// WithDefaultQueryer configures ResourceViewer with the default visitor.
func WithDefaultQueryer(dashConfig config.Dash, q queryer.Queryer) ViewerOpt {
	return func(rv *ResourceViewer) error {
		visitor, err := objectvisitor.NewDefaultVisitor(dashConfig, q)
		if err != nil {
			return err
		}

		rv.visitor = visitor
		return nil
	}
}

// ResourceViewer visits an object and creates a view component.
type ResourceViewer struct {
	dashConfig config.Dash
	visitor    objectvisitor.Visitor
}

// New creates an instance of ResourceViewer.
func New(dashConfig config.Dash, opts ...ViewerOpt) (*ResourceViewer, error) {
	rv := &ResourceViewer{
		dashConfig: dashConfig,
	}

	for _, opt := range opts {
		if err := opt(rv); err != nil {
			return nil, errors.Wrap(err, "invalid resource viewer option")
		}
	}

	if rv.visitor == nil {
		return nil, errors.New("resource viewer visitor is nil")
	}

	return rv, nil
}

// Visit visits an object and creates a view component.

func (rv *ResourceViewer) Visit(ctx context.Context, object runtime.Object) (*component.ResourceViewer, error) {
	ctx, span := trace.StartSpan(ctx, "resourceViewer")
	defer span.End()

	accessor, err := meta.Accessor(object)
	if err != nil {
		return nil, err
	}

	logger := log.From(ctx).With("object", kubernetes.PrintObject(object))
	logger.Debugf("starting resource viewer visit")

	now := time.Now()
	defer func() {
		logger.With("elapsed", time.Since(now)).Debugf("ending resource viewer visit")
	}()

	handler, err := NewHandler(rv.dashConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Create handler")
	}

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return nil, err
	}

	u := &unstructured.Unstructured{Object: m}

	if err := rv.visitor.Visit(ctx, u, handler, true); err != nil {
		return nil, errors.Wrapf(err, "error unable to visit object %s", kubernetes.PrintObject(object))
	}

	return GenerateComponent(ctx, handler, accessor.GetUID())
}

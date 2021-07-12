/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package resourceviewer

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/objectvisitor"
	"github.com/vmware-tanzu/octant/internal/queryer"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	visitMaxDuration = 3 * time.Second
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

// Create creates a resource viewer given a list objects.
func Create(ctx context.Context, dashConfig config.Dash, q queryer.Queryer, selection string, objects ...*unstructured.Unstructured) (*component.ResourceViewer, error) {
	rv, err := New(dashConfig, WithDefaultQueryer(dashConfig, q))
	if err != nil {
		return nil, fmt.Errorf("create resource viewer: %w", err)
	}

	handler, err := NewHandler(dashConfig)
	if err != nil {
		return nil, fmt.Errorf("create resource viewer handler: %w", err)
	}

	for _, object := range objects {
		if object == nil {
			continue
		}
		if err := rv.Visit(ctx, object, handler); err != nil {
			return nil, fmt.Errorf("unable to visit %s %s: %w",
				object.GroupVersionKind(),
				object.GetName(),
				err)
		}
	}

	c, err := GenerateComponent(ctx, handler, selection)
	if err != nil {
		return nil, fmt.Errorf("generate resource viewer component: %w", err)
	}

	return c, nil
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
func (rv *ResourceViewer) Visit(ctx context.Context, object runtime.Object, handler *Handler) error {
	ctx, span := trace.StartSpan(ctx, "resourceViewer")
	defer span.End()

	if handler == nil {
		return errors.New("handler is nil")
	}

	logger := log.From(ctx).With("object", kubernetes.PrintObject(object))

	now := time.Now()
	defer func() {
		elapsed := time.Since(now)
		if elapsed > visitMaxDuration {
			logger.With("elapsed", elapsed).Debugf("ending resource viewer visit")
		}
	}()

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return err
	}

	u := &unstructured.Unstructured{Object: m}

	level := 1
	if err := rv.visitor.Visit(ctx, u, handler, true, level); err != nil {
		return errors.Wrapf(err, "error unable to visit object %s", kubernetes.PrintObject(object))
	}

	sort.Slice(handler.edgeCache, func(i, j int) bool {
		first := handler.edgeCache[i]
		second := handler.edgeCache[j]

		// Sort edges first by depth level (to ensure good layout),
		// then by kind/name combination (to ensure repeatable layout)
		if first.level != second.level {
			return first.level < second.level
		} else {
			return fmt.Sprintf("%s(%s)-%s(%s)", first.from.GetKind(), first.from.GetName(), first.to.GetKind(), first.to.GetName()) <
				fmt.Sprintf("%s(%s)-%s(%s)", second.from.GetKind(), second.from.GetName(), second.to.GetKind(), second.to.GetName())
		}
	})

	for _, edge := range handler.edgeCache {
		handler.FinalizeEdge(ctx, edge.from, edge.to)
	}
	return nil
}

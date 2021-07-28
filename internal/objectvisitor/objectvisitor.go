/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectvisitor

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"

	"github.com/vmware-tanzu/octant/internal/queryer"
	"github.com/vmware-tanzu/octant/pkg/config"
)

//go:generate mockgen -destination=./fake/mock_object_handler.go -package=fake -self_package github.com/vmware-tanzu/octant/internal/objectvisitor/fake github.com/vmware-tanzu/octant/internal/objectvisitor ObjectHandler
//go:generate mockgen -destination=./fake/mock_default_typed_visitor.go -package=fake github.com/vmware-tanzu/octant/internal/objectvisitor DefaultTypedVisitor
//go:generate mockgen -destination=./fake/mock_typed_visitor.go -package=fake github.com/vmware-tanzu/octant/internal/objectvisitor TypedVisitor
//go:generate mockgen -destination=./fake/mock_visitor.go -package=fake github.com/vmware-tanzu/octant/internal/objectvisitor Visitor

// ObjectHandler performs actions on an object. Can be used to augment
// visitor actions with extra functionality.
type ObjectHandler interface {
	AddEdge(ctx context.Context, v1, v2 *unstructured.Unstructured, level int) error
	Process(ctx context.Context, object *unstructured.Unstructured) error
	SetLevel(objectKind string, level int) int
}

// DefaultTypedVisitor is the default typed visitors.
type DefaultTypedVisitor interface {
	Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitor Visitor, visitDescendants bool, level int) error
}

// TypedVisitor is a typed visitor for a specific gvk.
type TypedVisitor interface {
	DefaultTypedVisitor
	Supports() schema.GroupVersionKind
}

// Visitor is a visitor for cluster objects. It will visit an object and all of
// its ancestors and descendants.
type Visitor interface {
	Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitDescendants bool, level int) error
}

// DefaultVisitorOption is an option for configuring DefaultVisitor.
type DefaultVisitorOption func(*DefaultVisitor)

// SetDefaultHandler sets the default typed visitor for objects.
func SetDefaultHandler(dtv DefaultTypedVisitor) DefaultVisitorOption {
	return func(dv *DefaultVisitor) {
		dv.defaultHandler = dtv
	}
}

// SetTypedVisitors sets additional typed visitor for objects based on gvk.
func SetTypedVisitors(list []TypedVisitor) DefaultVisitorOption {
	return func(dv *DefaultVisitor) {
		dv.typedVisitors = list
	}
}

// DefaultVisitor is the default implementation of Visitor.
type DefaultVisitor struct {
	queryer   queryer.Queryer
	visited   map[types.UID]bool
	visitedMu sync.Mutex

	typedVisitors  []TypedVisitor
	defaultHandler DefaultTypedVisitor
}

var _ Visitor = (*DefaultVisitor)(nil)
var _ Visitor = (*DefaultVisitor)(nil)

// NewDefaultVisitor creates an instance of DefaultVisitor.
func NewDefaultVisitor(dashConfig config.Dash, q queryer.Queryer, options ...DefaultVisitorOption) (*DefaultVisitor, error) {
	dv := &DefaultVisitor{
		queryer: q,
		visited: make(map[types.UID]bool),
		typedVisitors: []TypedVisitor{
			NewIngress(q),
			NewPod(q),
			NewService(q),
			NewHorizontalPodAutoscaler(q),
			NewAPIService(dashConfig.ObjectStore()),
			NewMutatingWebhookConfiguration(dashConfig.ObjectStore()),
			NewValidatingWebhookConfiguration(dashConfig.ObjectStore()),
		},
		defaultHandler: NewObject(dashConfig, q),
	}

	for _, option := range options {
		option(dv)
	}

	return dv, nil
}

// hasVisited returns true if this object has already been visited. If the
// object has not been visited, it returns false, and sets the object
// visit status to true.
func (dv *DefaultVisitor) hasVisited(object runtime.Object) (bool, error) {
	if object == nil {
		return false, errors.Errorf("unable to check if nil object has been visited")
	}

	dv.visitedMu.Lock()
	defer dv.visitedMu.Unlock()

	accessor := meta.NewAccessor()
	uid, err := accessor.UID(object)
	if err != nil {
		return false, errors.Wrap(err, "get uid from object")
	}

	if _, ok := dv.visited[uid]; ok {
		return true, nil
	}

	dv.visited[uid] = true

	return false, nil
}

// Visit visits a runtime.Object.
func (dv *DefaultVisitor) Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitDescendants bool, level int) error {
	if ctx.Err() != nil {
		return nil
	}

	if object == nil {
		return errors.New("trying to visit a nil object")
	}

	if handler == nil {
		return errors.New("handler is nil")
	}

	hasVisited, err := dv.hasVisited(object)
	if err != nil {
		return errors.Wrapf(err, "check for visit object")
	}

	if hasVisited {
		return nil
	}

	return dv.visitObject(ctx, object, handler, visitDescendants, level)
}

// visitObject visits an object. If the object is a service, ingress, or pod, it
// also runs custom visitor code for them.
func (dv *DefaultVisitor) visitObject(ctx context.Context, object runtime.Object, handler ObjectHandler, visitDescendants bool, level int) error {
	ctx, span := trace.StartSpan(ctx, "visitObject")
	defer span.End()

	if object == nil {
		return errors.New("can't visit a nil object")
	}

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return err
	}

	u := &unstructured.Unstructured{Object: m}

	apiVersion := u.GetAPIVersion()
	kind := u.GetKind()

	objectGVK := schema.FromAPIVersionAndKind(apiVersion, kind)

	tvMap := make(map[schema.GroupVersionKind]TypedVisitor)
	for _, typedVisitor := range dv.typedVisitors {
		tvMap[typedVisitor.Supports()] = typedVisitor
	}

	tv, ok := tvMap[objectGVK]
	if ok {
		if err := tv.Visit(ctx, u, handler, dv, visitDescendants, level); err != nil {
			return err
		}
	}

	return dv.defaultHandler.Visit(ctx, u, handler, dv, visitDescendants, level)
}

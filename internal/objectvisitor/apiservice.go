/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectvisitor

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// APIService is a typed visitor for apiservices.
type APIService struct {
	objectStore store.Store
}

var _ TypedVisitor = (*APIService)(nil)

// NewAPIService creates an instance of APIService.
func NewAPIService(os store.Store) *APIService {
	return &APIService{
		objectStore: os,
	}
}

// Support returns the gvk this typed visitor supports.
func (p *APIService) Supports() schema.GroupVersionKind {
	return gvk.APIService
}

// Visit visits a apiservice. It looks for service accounts and services.
func (p *APIService) Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitor Visitor, visitDescendants bool) error {
	ctx, span := trace.StartSpan(ctx, "visitAPIService")
	defer span.End()

	if p.objectStore == nil {
		return errors.New("objectStore is nil")
	}

	apiservice := &apiregistrationv1.APIService{}
	if err := kubernetes.FromUnstructured(object, apiservice); err != nil {
		return err
	}

	var g errgroup.Group

	g.Go(func() error {
		if apiservice.Spec.Service == nil {
			return nil
		}

		key := store.KeyFromGroupVersionKind(gvk.Service)
		key.Namespace = apiservice.Spec.Service.Namespace
		key.Name = apiservice.Spec.Service.Name
		service, err := p.objectStore.Get(ctx, key)
		if err != nil {
			if kerrors.IsNotFound(err) {
				return nil
			}
			return err
		}

		if err := visitor.Visit(ctx, service, handler, true); err != nil {
			return errors.Wrapf(err, "apiservice %s visit service %s",
				kubernetes.PrintObject(apiservice), kubernetes.PrintObject(service))
		}

		return handler.AddEdge(ctx, object, service)
	})

	return g.Wait()
}

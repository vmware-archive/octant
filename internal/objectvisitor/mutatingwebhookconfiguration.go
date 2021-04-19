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
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// MutatingWebhookConfiguration is a typed visitor for mutatingwebhookconfigurations.
type MutatingWebhookConfiguration struct {
	objectStore store.Store
}

var _ TypedVisitor = (*MutatingWebhookConfiguration)(nil)

// NewMutatingWebhookConfiguration creates an instance of MutatingWebhookConfiguration.
func NewMutatingWebhookConfiguration(os store.Store) *MutatingWebhookConfiguration {
	return &MutatingWebhookConfiguration{
		objectStore: os,
	}
}

// Support returns the gvk this typed visitor supports.
func (p *MutatingWebhookConfiguration) Supports() schema.GroupVersionKind {
	return gvk.MutatingWebhookConfiguration
}

// Visit visits a mutatingwebhookconfiguration. It looks for service accounts and services.
func (p *MutatingWebhookConfiguration) Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitor Visitor, visitDescendants bool, level int) error {
	ctx, span := trace.StartSpan(ctx, "visitMutatingWebhookConfiguration")
	defer span.End()

	if p.objectStore == nil {
		return errors.New("objectStore is nil")
	}

	mutatingwebhookconfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{}
	if err := kubernetes.FromUnstructured(object, mutatingwebhookconfiguration); err != nil {
		return err
	}
	level = handler.SetLevel(mutatingwebhookconfiguration.Kind, level)

	var g errgroup.Group

	g.Go(func() error {
		for _, mutatingwebhook := range mutatingwebhookconfiguration.Webhooks {
			g.Go(func() error {
				if mutatingwebhook.ClientConfig.Service == nil {
					return nil
				}

				key := store.KeyFromGroupVersionKind(gvk.Service)
				key.Namespace = mutatingwebhook.ClientConfig.Service.Namespace
				key.Name = mutatingwebhook.ClientConfig.Service.Name
				service, err := p.objectStore.Get(ctx, key)
				if err != nil {
					if kerrors.IsNotFound(err) {
						return nil
					}
					return err
				}

				if visitDescendants {
					if err := visitor.Visit(ctx, service, handler, false, level); err != nil {
						return errors.Wrapf(err, "mutatingwebhookconfiguration %s visit service %s",
							kubernetes.PrintObject(mutatingwebhookconfiguration), kubernetes.PrintObject(service))
					}
				}

				return handler.AddEdge(ctx, object, service, level)
			})
		}

		return nil
	})

	return g.Wait()
}

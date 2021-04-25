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

// ValidatingWebhookConfiguration is a typed visitor for validatingwebhookconfigurations.
type ValidatingWebhookConfiguration struct {
	objectStore store.Store
}

var _ TypedVisitor = (*ValidatingWebhookConfiguration)(nil)

// NewValidatingWebhookConfiguration creates an instance of ValidatingWebhookConfiguration.
func NewValidatingWebhookConfiguration(os store.Store) *ValidatingWebhookConfiguration {
	return &ValidatingWebhookConfiguration{
		objectStore: os,
	}
}

// Support returns the gvk this typed visitor supports.
func (p *ValidatingWebhookConfiguration) Supports() schema.GroupVersionKind {
	return gvk.ValidatingWebhookConfiguration
}

// Visit visits a validatingwebhookconfiguration. It looks for service accounts and services.
func (p *ValidatingWebhookConfiguration) Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitor Visitor, visitDescendants bool, level int) error {
	ctx, span := trace.StartSpan(ctx, "visitValidatingWebhookConfiguration")
	defer span.End()

	if p.objectStore == nil {
		return errors.New("objectStore is nil")
	}

	validatingwebhookconfiguration := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	if err := kubernetes.FromUnstructured(object, validatingwebhookconfiguration); err != nil {
		return err
	}
	level = handler.SetLevel(validatingwebhookconfiguration.Kind, level)

	var g errgroup.Group

	g.Go(func() error {
		for _, validatingwebhook := range validatingwebhookconfiguration.Webhooks {
			g.Go(func() error {
				if validatingwebhook.ClientConfig.Service == nil {
					return nil
				}

				key := store.KeyFromGroupVersionKind(gvk.Service)
				key.Namespace = validatingwebhook.ClientConfig.Service.Namespace
				key.Name = validatingwebhook.ClientConfig.Service.Name
				service, err := p.objectStore.Get(ctx, key)
				if err != nil {
					if kerrors.IsNotFound(err) {
						return nil
					}
					return err
				}

				if visitDescendants {
					if err := visitor.Visit(ctx, service, handler, false, level); err != nil {
						return errors.Wrapf(err, "validatingwebhookconfiguration %s visit service %s",
							kubernetes.PrintObject(validatingwebhookconfiguration), kubernetes.PrintObject(service))
					}
				}

				return handler.AddEdge(ctx, object, service, level)
			})
		}

		return nil
	})

	return g.Wait()
}

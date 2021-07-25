/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/pkg/view/component"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
)

// MutatingWebhookConfigurationListHandler is a printFunc that prints mutating webhook configurations
func MutatingWebhookConfigurationListHandler(ctx context.Context, list *admissionregistrationv1.MutatingWebhookConfigurationList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("mutating webhook configuration list is nil")
	}

	cols := component.NewTableCols("Name", "Age")
	ot := NewObjectTable("Mutating Webhook Configurations", "We couldn't find any mutating webhook configurations!", cols, options.DashConfig.ObjectStore(), options.DashConfig.TerminateThreshold())
	ot.EnablePluginStatus(options.DashConfig.PluginManager())
	for _, mutatingWebhookConfiguration := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&mutatingWebhookConfiguration, mutatingWebhookConfiguration.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		ts := mutatingWebhookConfiguration.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		if err := ot.AddRowForObject(ctx, &mutatingWebhookConfiguration, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// MutatingWebhookConfigurationHandler is a printFunc that prints a mutating webhook configurations
func MutatingWebhookConfigurationHandler(ctx context.Context, mutatingWebhookConfiguration *admissionregistrationv1.MutatingWebhookConfiguration, options Options) (component.Component, error) {
	o := NewObject(mutatingWebhookConfiguration)

	ch, err := newMutatingWebhookConfigurationHandler(mutatingWebhookConfiguration, o)
	if err != nil {
		return nil, err
	}

	if err := ch.Webhooks(options); err != nil {
		return nil, errors.Wrap(err, "print mutatingwebhookconfiguration webhooks")
	}

	return o.ToComponent(ctx, options)
}

type mutatingWebhookConfigurationObject interface {
	Webhooks(options Options) error
}

type mutatingWebhookConfigurationHandler struct {
	mutatingWebhookConfiguration *admissionregistrationv1.MutatingWebhookConfiguration
	webhookFunc                  func(*admissionregistrationv1.MutatingWebhook, Options) (*component.Summary, error)
	object                       *Object
}

var _ mutatingWebhookConfigurationObject = (*mutatingWebhookConfigurationHandler)(nil)

func newMutatingWebhookConfigurationHandler(mutatingWebhookConfiguration *admissionregistrationv1.MutatingWebhookConfiguration, object *Object) (*mutatingWebhookConfigurationHandler, error) {
	if mutatingWebhookConfiguration == nil {
		return nil, errors.New("can't print a nil mutatingwebhookconfiguration")
	}

	if object == nil {
		return nil, errors.New("can't print a mutatingwebhookconfiguration using an nil object printer")
	}

	ch := &mutatingWebhookConfigurationHandler{
		mutatingWebhookConfiguration: mutatingWebhookConfiguration,
		webhookFunc:                  defaultMutatingWebhook,
		object:                       object,
	}
	return ch, nil
}

func (c *mutatingWebhookConfigurationHandler) Webhooks(options Options) error {
	for i := range c.mutatingWebhookConfiguration.Webhooks {
		webhook := &c.mutatingWebhookConfiguration.Webhooks[i]
		c.object.RegisterItems(
			ItemDescriptor{
				Width: component.WidthFull,
				Func: func() (component.Component, error) {
					return c.webhookFunc(webhook, options)
				},
			},
		)
	}
	return nil
}

func defaultMutatingWebhook(mutatingWebhook *admissionregistrationv1.MutatingWebhook, options Options) (*component.Summary, error) {
	return NewMutatingWebhook(mutatingWebhook).Create(options)
}

// MutatingWebhook generates a mutating webhook
type MutatingWebhook struct {
	mutatingWebhook *admissionregistrationv1.MutatingWebhook
}

// NewMutatingWebhook creates an instance of MutatingWebhook
func NewMutatingWebhook(mutatingWebhook *admissionregistrationv1.MutatingWebhook) *MutatingWebhook {
	return &MutatingWebhook{
		mutatingWebhook: mutatingWebhook,
	}
}

// Create creates a mutating webhook summary
func (c *MutatingWebhook) Create(options Options) (*component.Summary, error) {
	if c.mutatingWebhook == nil {
		return nil, errors.New("mutatingWebhook is nil")
	}

	var sections component.SummarySections

	client, err := admissionWebhookClientConfig(c.mutatingWebhook.ClientConfig, options)
	if err != nil {
		return nil, err
	}
	sections.Add("Client", client)
	rules, err := admissionWebhookRules(c.mutatingWebhook.Rules, options)
	if err != nil {
		return nil, err
	}
	sections.Add("Rules", rules)
	namespaceSelector := admissionWebhookLabelSelector(c.mutatingWebhook.NamespaceSelector)
	sections.Add("Namespace Selector", namespaceSelector)
	objectSelector := admissionWebhookLabelSelector(c.mutatingWebhook.ObjectSelector)
	sections.Add("Object Selector", objectSelector)
	reinvocationPolicy := admissionWebhookReinvocationPolicy(c.mutatingWebhook.ReinvocationPolicy)
	sections.Add("Reinvocation Policy", reinvocationPolicy)
	failurePolicy := admissionWebhookFailurePolicy(c.mutatingWebhook.FailurePolicy)
	sections.Add("Failure Policy", failurePolicy)
	matchPolicy := admissionWebhookMatchPolicy(c.mutatingWebhook.MatchPolicy)
	sections.Add("Match Policy", matchPolicy)
	sideEffects := admissionWebhookSideEffects(c.mutatingWebhook.SideEffects)
	sections.Add("Side Effects", sideEffects)
	timeout := admissionWebhookTimeout(c.mutatingWebhook.TimeoutSeconds)
	sections.Add("Timeout", timeout)
	admissionReviewVersions := admissionWebhookAdmissionReviewVersions(c.mutatingWebhook.AdmissionReviewVersions)
	sections.Add("Admission Review Versions", admissionReviewVersions)

	summary := component.NewSummary(c.mutatingWebhook.Name, sections...)

	return summary, nil
}

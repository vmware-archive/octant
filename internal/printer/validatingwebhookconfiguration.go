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

// ValidatingWebhookConfigurationListHandler is a printFunc that prints validating webhook configurations
func ValidatingWebhookConfigurationListHandler(ctx context.Context, list *admissionregistrationv1.ValidatingWebhookConfigurationList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("validating webhook configuration list is nil")
	}

	cols := component.NewTableCols("Name", "Age")
	ot := NewObjectTable("Validating Webhook Configurations", "We couldn't find any validating webhook configurations!", cols, options.DashConfig.ObjectStore(), options.DashConfig.TerminateThreshold())
	ot.EnablePluginStatus(options.DashConfig.PluginManager())
	for _, validatingWebhookConfiguration := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&validatingWebhookConfiguration, validatingWebhookConfiguration.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		ts := validatingWebhookConfiguration.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		if err := ot.AddRowForObject(ctx, &validatingWebhookConfiguration, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// ValidatingWebhookConfigurationHandler is a printFunc that prints a validating webhook configurations
func ValidatingWebhookConfigurationHandler(ctx context.Context, validatingWebhookConfiguration *admissionregistrationv1.ValidatingWebhookConfiguration, options Options) (component.Component, error) {
	o := NewObject(validatingWebhookConfiguration)

	ch, err := newValidatingWebhookConfigurationHandler(validatingWebhookConfiguration, o)
	if err != nil {
		return nil, err
	}

	if err := ch.Webhooks(options); err != nil {
		return nil, errors.Wrap(err, "print validatingwebhookconfiguration webhooks")
	}

	return o.ToComponent(ctx, options)
}

type validatingWebhookConfigurationObject interface {
	Webhooks(options Options) error
}

type validatingWebhookConfigurationHandler struct {
	validatingWebhookConfiguration *admissionregistrationv1.ValidatingWebhookConfiguration
	webhookFunc                    func(*admissionregistrationv1.ValidatingWebhook, Options) (*component.Summary, error)
	object                         *Object
}

var _ validatingWebhookConfigurationObject = (*validatingWebhookConfigurationHandler)(nil)

func newValidatingWebhookConfigurationHandler(validatingWebhookConfiguration *admissionregistrationv1.ValidatingWebhookConfiguration, object *Object) (*validatingWebhookConfigurationHandler, error) {
	if validatingWebhookConfiguration == nil {
		return nil, errors.New("can't print a nil validatingwebhookconfiguration")
	}

	if object == nil {
		return nil, errors.New("can't print a validatingwebhookconfiguration using an nil object printer")
	}

	ch := &validatingWebhookConfigurationHandler{
		validatingWebhookConfiguration: validatingWebhookConfiguration,
		webhookFunc:                    defaultValidatingWebhook,
		object:                         object,
	}
	return ch, nil
}

func (c *validatingWebhookConfigurationHandler) Webhooks(options Options) error {
	for i := range c.validatingWebhookConfiguration.Webhooks {
		webhook := &c.validatingWebhookConfiguration.Webhooks[i]
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

func defaultValidatingWebhook(validatingWebhook *admissionregistrationv1.ValidatingWebhook, options Options) (*component.Summary, error) {
	return NewValidatingWebhook(validatingWebhook).Create(options)
}

// ValidatingWebhook generates a validating webhook
type ValidatingWebhook struct {
	validatingWebhook *admissionregistrationv1.ValidatingWebhook
}

// NewValidatingWebhook creates an instance of ValidatingWebhook
func NewValidatingWebhook(validatingWebhook *admissionregistrationv1.ValidatingWebhook) *ValidatingWebhook {
	return &ValidatingWebhook{
		validatingWebhook: validatingWebhook,
	}
}

// Create creates a validating webhook summary
func (c *ValidatingWebhook) Create(options Options) (*component.Summary, error) {
	if c.validatingWebhook == nil {
		return nil, errors.New("validatingWebhook is nil")
	}

	var sections component.SummarySections

	client, err := admissionWebhookClientConfig(c.validatingWebhook.ClientConfig, options)
	if err != nil {
		return nil, err
	}
	sections.Add("Client", client)
	rules, err := admissionWebhookRules(c.validatingWebhook.Rules, options)
	if err != nil {
		return nil, err
	}
	sections.Add("Rules", rules)
	namespaceSelector := admissionWebhookLabelSelector(c.validatingWebhook.NamespaceSelector)
	sections.Add("Namespace Selector", namespaceSelector)
	objectSelector := admissionWebhookLabelSelector(c.validatingWebhook.ObjectSelector)
	sections.Add("Object Selector", objectSelector)
	failurePolicy := admissionWebhookFailurePolicy(c.validatingWebhook.FailurePolicy)
	sections.Add("Failure Policy", failurePolicy)
	matchPolicy := admissionWebhookMatchPolicy(c.validatingWebhook.MatchPolicy)
	sections.Add("Match Policy", matchPolicy)
	sideEffects := admissionWebhookSideEffects(c.validatingWebhook.SideEffects)
	sections.Add("Side Effects", sideEffects)
	timeout := admissionWebhookTimeout(c.validatingWebhook.TimeoutSeconds)
	sections.Add("Timeout", timeout)
	admissionReviewVersions := admissionWebhookAdmissionReviewVersions(c.validatingWebhook.AdmissionReviewVersions)
	sections.Add("Admission Review Versions", admissionReviewVersions)

	summary := component.NewSummary(c.validatingWebhook.Name, sections...)

	return summary, nil
}

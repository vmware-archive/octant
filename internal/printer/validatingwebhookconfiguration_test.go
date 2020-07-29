/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_ValidatingWebhookConfigurationListHandler(t *testing.T) {
	object := testutil.CreateValidatingWebhookConfiguration("validatingWebhookConfiguration")
	object.CreationTimestamp = *testutil.CreateTimestamp()

	list := &admissionregistrationv1beta1.ValidatingWebhookConfigurationList{
		Items: []admissionregistrationv1beta1.ValidatingWebhookConfiguration{*object},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	tpo.PathForObject(object, object.Name, "/path")

	now := testutil.Time()

	ctx := context.Background()
	got, err := ValidatingWebhookConfigurationListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Age")
	expected := component.NewTable("Validating Webhook Configurations", "We couldn't find any validating webhook configurations!", cols)

	expected.Add(component.TableRow{
		"Name": component.NewLink("", object.Name, "/path",
			genObjectStatus(component.TextStatusOK, []string{
				"admissionregistration.k8s.io/v1beta1 ValidatingWebhookConfiguration is OK",
			})),
		"Age": component.NewTimestamp(now),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, object),
		}),
	})

	testutil.AssertJSONEqual(t, expected, got)
}

func Test_NewValidatingWebhook(t *testing.T) {
	namespacedScope := admissionregistrationv1beta1.NamespacedScope
	ignore := admissionregistrationv1beta1.Ignore
	exact := admissionregistrationv1beta1.Exact
	sideEffectClassNone := admissionregistrationv1beta1.SideEffectClassNone

	cases := []struct {
		name              string
		validatingWebhook *admissionregistrationv1beta1.ValidatingWebhook
		isErr             bool
		expected          *component.Summary
	}{
		{
			name: "normal",
			validatingWebhook: &admissionregistrationv1beta1.ValidatingWebhook{
				Name: "test-webhook",
				ClientConfig: admissionregistrationv1beta1.WebhookClientConfig{
					Service: &admissionregistrationv1beta1.ServiceReference{
						Namespace: "default",
						Name:      "service",
					},
				},
				Rules: []admissionregistrationv1beta1.RuleWithOperations{
					{
						Rule: admissionregistrationv1beta1.Rule{
							APIGroups:   []string{"apps"},
							APIVersions: []string{"v1"},
							Resources:   []string{"*"},
							Scope:       &namespacedScope,
						},
						Operations: []admissionregistrationv1beta1.OperationType{
							admissionregistrationv1beta1.Create,
							admissionregistrationv1beta1.Update,
						},
					},
				},
				NamespaceSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"intercept": "true"},
				},
				ObjectSelector:          &metav1.LabelSelector{},
				FailurePolicy:           &ignore,
				MatchPolicy:             &exact,
				SideEffects:             &sideEffectClassNone,
				AdmissionReviewVersions: []string{"v1beta1"},
			},
			expected: component.NewSummary("test-webhook", []component.SummarySection{
				{
					Header:  "Client",
					Content: component.NewLink("", "default/service", "/path"),
				},
				{
					Header: "Rules",
					Content: testWebhookRulesTable(
						component.TableRow{
							"API Groups":   component.NewText("apps"),
							"API Versions": component.NewText("v1"),
							"Resources":    component.NewText("*"),
							"Operations":   component.NewMarkdownText("- CREATE\n- UPDATE\n"),
							"Scope":        component.NewText("Namespaced"),
						},
					),
				},
				{
					Header:  "Namespace Selector",
					Content: component.NewText("intercept:true"),
				},
				{
					Header:  "Object Selector",
					Content: component.NewText("*"),
				},
				{
					Header:  "Failure Policy",
					Content: component.NewText("Ignore"),
				},
				{
					Header:  "Match Policy",
					Content: component.NewText("Exact"),
				},
				{
					Header:  "Side Effects",
					Content: component.NewText("None"),
				},
				{
					Header:  "Timeout",
					Content: component.NewText("10s"),
				},
				{
					Header:  "Admission Review Versions",
					Content: component.NewText("v1beta1"),
				},
			}...),
		},
		{
			name: "default",
			validatingWebhook: &admissionregistrationv1beta1.ValidatingWebhook{
				Name: "test-webhook",
			},
			expected: component.NewSummary("test-webhook", []component.SummarySection{
				{
					Header:  "Client",
					Content: component.NewText("unknown"),
				},
				{
					Header:  "Rules",
					Content: testWebhookRulesTable(),
				},
				{
					Header:  "Namespace Selector",
					Content: component.NewText("*"),
				},
				{
					Header:  "Object Selector",
					Content: component.NewText("*"),
				},
				{
					Header:  "Failure Policy",
					Content: component.NewText("Fail"),
				},
				{
					Header:  "Match Policy",
					Content: component.NewText("Equivalent"),
				},
				{
					Header:  "Side Effects",
					Content: component.NewText("Unknown"),
				},
				{
					Header:  "Timeout",
					Content: component.NewText("10s"),
				},
				{
					Header:  "Admission Review Versions",
					Content: component.NewMarkdownText("null\n"),
				},
			}...),
		},
		{
			name:              "validatingwebhook is nil",
			validatingWebhook: nil,
			isErr:             true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			tpo.link.EXPECT().
				ForGVK("default", "v1", "Service", "service", "default/service").
				Return(component.NewLink("", "default/service", "/path"), nil).
				AnyTimes()

			cc := NewValidatingWebhook(test.validatingWebhook)

			summary, err := cc.Create(printOptions)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, test.expected, summary)
		})
	}
}

func init() {
	admissionregistrationv1beta1.AddToScheme(scheme.Scheme)
}

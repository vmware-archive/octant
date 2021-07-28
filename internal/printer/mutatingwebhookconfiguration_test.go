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
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_MutatingWebhookConfigurationListHandler(t *testing.T) {
	object := testutil.CreateMutatingWebhookConfiguration("mutatingWebhookConfiguration")
	object.CreationTimestamp = *testutil.CreateTimestamp()

	list := &admissionregistrationv1.MutatingWebhookConfigurationList{
		Items: []admissionregistrationv1.MutatingWebhookConfiguration{*object},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	tpo.PathForObject(object, object.Name, "/path")

	now := testutil.Time()

	ctx := context.Background()
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, object)
	got, err := MutatingWebhookConfigurationListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Age")
	expected := component.NewTable("Mutating Webhook Configurations", "We couldn't find any mutating webhook configurations!", cols)

	expected.Add(component.TableRow{
		"Name": component.NewLink("", object.Name, "/path",
			genObjectStatus(component.TextStatusOK, []string{
				"admissionregistration.k8s.io/v1 MutatingWebhookConfiguration is OK",
			})),
		"Age": component.NewTimestamp(now),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, object),
		}),
	})

	testutil.AssertJSONEqual(t, expected, got)
}

func Test_NewMutatingWebhook(t *testing.T) {
	namespacedScope := admissionregistrationv1.NamespacedScope
	ifNeededReinvocationPolicy := admissionregistrationv1.IfNeededReinvocationPolicy
	ignore := admissionregistrationv1.Ignore
	exact := admissionregistrationv1.Exact
	sideEffectClassNone := admissionregistrationv1.SideEffectClassNone

	cases := []struct {
		name            string
		mutatingWebhook *admissionregistrationv1.MutatingWebhook
		isErr           bool
		expected        *component.Summary
	}{
		{
			name: "normal",
			mutatingWebhook: &admissionregistrationv1.MutatingWebhook{
				Name: "test-webhook",
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					Service: &admissionregistrationv1.ServiceReference{
						Namespace: "default",
						Name:      "service",
					},
				},
				Rules: []admissionregistrationv1.RuleWithOperations{
					{
						Rule: admissionregistrationv1.Rule{
							APIGroups:   []string{"apps"},
							APIVersions: []string{"v1"},
							Resources:   []string{"*"},
							Scope:       &namespacedScope,
						},
						Operations: []admissionregistrationv1.OperationType{
							admissionregistrationv1.Create,
							admissionregistrationv1.Update,
						},
					},
				},
				NamespaceSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"intercept": "true"},
				},
				ObjectSelector:          &metav1.LabelSelector{},
				ReinvocationPolicy:      &ifNeededReinvocationPolicy,
				FailurePolicy:           &ignore,
				MatchPolicy:             &exact,
				SideEffects:             &sideEffectClassNone,
				AdmissionReviewVersions: []string{"v1"},
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
					Header:  "Reinvocation Policy",
					Content: component.NewText("IfNeeded"),
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
					Content: component.NewText("v1"),
				},
			}...),
		},
		{
			name: "default",
			mutatingWebhook: &admissionregistrationv1.MutatingWebhook{
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
					Header:  "Reinvocation Policy",
					Content: component.NewText("Never"),
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
			name:            "mutatingwebhook is nil",
			mutatingWebhook: nil,
			isErr:           true,
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

			cc := NewMutatingWebhook(test.mutatingWebhook)

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
	admissionregistrationv1.AddToScheme(scheme.Scheme)
}

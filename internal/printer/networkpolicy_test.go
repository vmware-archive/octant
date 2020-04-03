/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_NetworkPolicyListHandler(t *testing.T) {
	cols := component.NewTableCols("Name", "Labels", "Age")
	now := testutil.Time()

	labels := map[string]string{
		"app": "testing",
	}

	object := testutil.CreateNetworkPolicy("networkPolicy")
	object.Labels = labels
	object.CreationTimestamp = metav1.Time{Time: now}

	list := &networkingv1.NetworkPolicyList{
		Items: []networkingv1.NetworkPolicy{*object},
	}

	cases := []struct {
		name     string
		list     *networkingv1.NetworkPolicyList
		expected *component.Table
		isErr    bool
	}{
		{
			name: "in general",
			list: list,
			expected: component.NewTableWithRows("Network Policies", "We couldn't find any network policies!", cols,
				[]component.TableRow{
					{
						"Name":   component.NewLink("", "networkPolicy", "/networkPolicy"),
						"Labels": component.NewLabels(labels),
						"Age":    component.NewTimestamp(now),
					},
				}),
		},
		{
			name:  "list is nil",
			list:  nil,
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			ctx := context.Background()

			if tc.list != nil {
				tpo.PathForObject(&tc.list.Items[0], tc.list.Items[0].Name, "/networkPolicy")
			}

			got, err := NetworkPolicyListHandler(ctx, tc.list, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, got)
		})
	}
}

func Test_NetworkPolicyConfiguration(t *testing.T) {
	selector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": "my_app",
		},
	}

	networkPolicy := testutil.CreateNetworkPolicy("networkPolicy")
	networkPolicy.Spec.PodSelector = selector
	networkPolicy.Spec.PolicyTypes = []networkingv1.PolicyType{
		networkingv1.PolicyTypeIngress,
		networkingv1.PolicyTypeEgress,
	}

	cases := []struct {
		name          string
		networkPolicy *networkingv1.NetworkPolicy
		expected      component.Component
		isErr         bool
	}{
		{
			name:          "general",
			networkPolicy: networkPolicy,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Policy Types",
					Content: component.NewText("Ingress, Egress"),
				},
				{
					Header:  "Selectors",
					Content: component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "my_app")}),
				},
			}...),
		},
		{
			name:          "nil networkPolicy",
			networkPolicy: nil,
			isErr:         true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			npc := NewNetworkPolicyConfiguration(tc.networkPolicy)
			summary, err := npc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_NetworkPolicySummaryStatus(t *testing.T) {
	networkPolicy := testutil.CreateNetworkPolicy("networkPolicy")

	got, err := createNetworkPolicySummaryStatus(networkPolicy)
	require.NoError(t, err)

	ingressTable, err := createIngressRules(networkPolicy.Spec.Ingress)
	require.NoError(t, err)

	egressTable, err := createEgressRules(networkPolicy.Spec.Egress)
	require.NoError(t, err)

	sections := component.SummarySections{
		{Header: "Allowing ingress traffic", Content: ingressTable},
		{Header: "Allowing egress traffic", Content: egressTable},
	}

	expected := component.NewSummary("Status", sections...)
	assert.Equal(t, expected, got)
}

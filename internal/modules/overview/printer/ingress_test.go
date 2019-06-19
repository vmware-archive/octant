/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func Test_IngressListHandler(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := testutil.CreateIngress("ingress")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Labels = labels

	list := &extv1beta1.IngressList{
		Items: []extv1beta1.Ingress{*object},
	}

	tlsObject := testutil.CreateIngress("ingress")
	tlsObject.CreationTimestamp = metav1.Time{Time: now}
	tlsObject.Labels = labels
	tlsObject.Spec.TLS = []extv1beta1.IngressTLS{{}}

	tlsList := &extv1beta1.IngressList{
		Items: []extv1beta1.Ingress{*tlsObject},
	}

	cols := component.NewTableCols("Name", "Labels", "Hosts", "Address", "Ports", "Age")

	cases := []struct {
		name     string
		list     *extv1beta1.IngressList
		expected *component.Table
		isErr    bool
	}{
		{
			name: "in general",
			list: list,
			expected: component.NewTableWithRows("Ingresses", cols,
				[]component.TableRow{
					{
						"Name":    component.NewLink("", "ingress", "/ingress"),
						"Labels":  component.NewLabels(labels),
						"Age":     component.NewTimestamp(now),
						"Hosts":   component.NewText("*"),
						"Address": component.NewText(""),
						"Ports":   component.NewText("80"),
					},
				}),
		},
		{
			name: "with TLS",
			list: tlsList,
			expected: component.NewTableWithRows("Ingresses", cols,
				[]component.TableRow{
					{
						"Name":    component.NewLink("", "ingress", "/ingress"),
						"Labels":  component.NewLabels(labels),
						"Age":     component.NewTimestamp(now),
						"Hosts":   component.NewText("*"),
						"Address": component.NewText(""),
						"Ports":   component.NewText("80, 443"),
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

			if tc.list != nil {
				tpo.PathForObject(&tc.list.Items[0], tc.list.Items[0].Name, "/ingress")

			}

			ctx := context.Background()
			got, err := IngressListHandler(ctx, tc.list, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)

		})
	}

}

func Test_printIngressConfig(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := testutil.CreateIngress("ingress")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Labels = labels

	objectNoBackend := testutil.CreateIngress("ingress")
	objectNoBackend.CreationTimestamp = metav1.Time{Time: now}
	objectNoBackend.Labels = labels
	objectNoBackend.Spec.Backend = nil

	cases := []struct {
		name     string
		object   *extv1beta1.Ingress
		expected component.Component
		isErr    bool
	}{
		{
			name:   "in general",
			object: object,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Default Backend",
					Content: component.NewLink("", "service", "/service"),
				},
			}...),
		},
		{
			name:   "no default backend",
			object: objectNoBackend,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Default Backend",
					Content: component.NewText("Default is not configured"),
				},
			}...),
		},
		{
			name:   "nil ingress",
			object: nil,
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			if tc.object != nil {
				stubIngressBackendLinks(tpo)
			}

			got, err := printIngressConfig(tc.object, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertComponentEqual(t, tc.expected, got)
		})
	}
}

func Test_printIngressHosts(t *testing.T) {
	object := testutil.CreateIngress("ingress")

	objectWithRules := testutil.CreateIngress("ingress")
	objectWithRules.Spec.Rules = []extv1beta1.IngressRule{
		{

			Host: "",
			IngressRuleValue: extv1beta1.IngressRuleValue{
				HTTP: &extv1beta1.HTTPIngressRuleValue{
					Paths: []extv1beta1.HTTPIngressPath{
						{
							Path: "/",
							Backend: extv1beta1.IngressBackend{
								ServiceName: "b1",
								ServicePort: intstr.FromInt(80),
							},
						},
					},
				},
			},
		},
	}

	cols := component.NewTableCols("Host", "Path", "Backends")

	cases := []struct {
		name     string
		object   *extv1beta1.Ingress
		expected component.Component
		isErr    bool
	}{
		{
			name:   "in general",
			object: object,
			expected: component.NewTableWithRows("Rules", cols, []component.TableRow{
				{
					"Backends": component.NewLink("", "service", "/service"),
					"Host":     component.NewText("*"),
					"Path":     component.NewText("*"),
				},
			}),
		},
		{
			name:   "with rules",
			object: objectWithRules,
			expected: component.NewTableWithRows("Rules", cols, []component.TableRow{
				{
					"Backends": component.NewLink("", "service", "/service"),
					"Host":     component.NewText("*"),
					"Path":     component.NewText("/"),
				},
			}),
		},
		{
			name:   "nil ingress",
			object: nil,
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			if tc.object != nil {
				stubIngressBackendLinks(tpo)
			}

			got, err := printRulesForIngress(tc.object, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertComponentEqual(t, tc.expected, got)
		})
	}
}

func stubIngressBackendLinks(tpo *testPrinterOptions) {
	serviceLink := component.NewLink("", "service", "/service")
	tpo.link.EXPECT().
		ForGVK(gomock.Any(), "v1", "Service", gomock.Any(), gomock.Any()).
		Return(serviceLink, nil).
		AnyTimes()
}

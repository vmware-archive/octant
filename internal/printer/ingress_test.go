/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_IngressListHandler(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	object := testutil.CreateIngress("ingress")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Labels = labels

	list := &networkingv1.IngressList{
		Items: []networkingv1.Ingress{*object},
	}

	tlsObject := testutil.CreateIngress("ingress")
	tlsObject.CreationTimestamp = metav1.Time{Time: now}
	tlsObject.Labels = labels
	tlsObject.Spec.TLS = []networkingv1.IngressTLS{{}}

	hostTest1 := testutil.CreateIngress("ingress")
	hostTest1.CreationTimestamp = metav1.Time{Time: now}
	hostTest1.Labels = labels
	hostTest1.Spec.TLS = []networkingv1.IngressTLS{}
	hostTest1.Spec.Rules = []networkingv1.IngressRule{
		{
			Host: "hello-world.info",
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path: "/v2",
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: "app",
									Port: networkingv1.ServiceBackendPort{
										Number: 80,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	hostTest2 := testutil.CreateIngress("ingress")
	hostTest2.CreationTimestamp = metav1.Time{Time: now}
	hostTest2.Labels = labels
	hostTest2.Spec.TLS = []networkingv1.IngressTLS{
		{
			SecretName: "secret",
			Hosts:      []string{"echo1.example.com", "echo2.example.com"},
		},
	}
	hostTest2.Spec.Rules = []networkingv1.IngressRule{
		{
			Host: "echo1.example.com",
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path: "path1/example.com",
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: "app",
									Port: networkingv1.ServiceBackendPort{
										Number: 8080,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Host: "echo2.example.com",
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path: "path2/example.com",
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: "app",
									Port: networkingv1.ServiceBackendPort{
										Number: 8080,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	tlsList := &networkingv1.IngressList{
		Items: []networkingv1.Ingress{*tlsObject},
	}

	hostTest1List := &networkingv1.IngressList{
		Items: []networkingv1.Ingress{*hostTest1},
	}

	hostTest2List := &networkingv1.IngressList{
		Items: []networkingv1.Ingress{*hostTest2},
	}

	cols := component.NewTableCols("Name", "Labels", "Hosts", "Address", "Ports", "Age")

	service := testutil.ToUnstructured(t, testutil.CreateService("service"))
	secret := testutil.ToUnstructured(t, testutil.CreateSecret("secret"))

	cases := []struct {
		name     string
		list     *networkingv1.IngressList
		expected *component.Table
		isErr    bool
	}{
		{
			name: "in general",
			list: list,
			expected: component.NewTableWithRows("Ingresses", "We couldn't find any ingresses!", cols,
				[]component.TableRow{
					{
						"Name": component.NewLink("", "ingress", "/ingress",
							genObjectStatus(component.TextStatusError, []string{
								`Backend for service "app" specifies an invalid port`,
							})),
						"Labels":  component.NewLabels(labels),
						"Age":     component.NewTimestamp(now),
						"Hosts":   component.NewText("*"),
						"Address": component.NewText(""),
						"Ports":   component.NewText("80"),
						component.GridActionKey: gridActionsFactory([]component.GridAction{
							buildObjectDeleteAction(t, object),
						}),
					},
				}),
		},
		{
			name: "with TLS",
			list: tlsList,
			expected: component.NewTableWithRows("Ingresses", "We couldn't find any ingresses!", cols,
				[]component.TableRow{
					{
						"Name": component.NewLink("", "ingress", "/ingress",
							genObjectStatus(component.TextStatusError, []string{
								`Backend for service "app" specifies an invalid port`,
								"TLS configuration did not define a secret name",
							})),
						"Labels":  component.NewLabels(labels),
						"Age":     component.NewTimestamp(now),
						"Hosts":   component.NewText("*"),
						"Address": component.NewText(""),
						"Ports":   component.NewText("80, 443"),
						component.GridActionKey: gridActionsFactory([]component.GridAction{
							buildObjectDeleteAction(t, object),
						}),
					},
				}),
		},
		{
			name: "host URL",
			list: hostTest1List,
			expected: component.NewTableWithRows("Ingresses", "We couldn't find any ingresses!", cols,
				[]component.TableRow{
					{
						"Name": component.NewLink("", "ingress", "/ingress",
							genObjectStatus(component.TextStatusError, []string{
								`Backend for service "app" specifies an invalid port`,
								`Backend for service "app" specifies an invalid port`,
							})),
						"Labels":  component.NewLabels(labels),
						"Age":     component.NewTimestamp(now),
						"Hosts":   component.NewLink("", "hello-world.info", "http://hello-world.info"),
						"Address": component.NewText(""),
						"Ports":   component.NewText("80"),
						component.GridActionKey: gridActionsFactory([]component.GridAction{
							buildObjectDeleteAction(t, object),
						}),
					},
				}),
		},
		{
			name: "multiple host TLS URLs",
			list: hostTest2List,
			expected: component.NewTableWithRows("Ingresses", "We couldn't find any ingresses!", cols,
				[]component.TableRow{
					{
						"Name": component.NewLink("", "ingress", "/ingress",
							genObjectStatus(component.TextStatusError, []string{
								`Backend for service "app" specifies an invalid port`,
								`Backend for service "app" specifies an invalid port`,
								`Backend for service "app" specifies an invalid port`,
							})),
						"Labels":  component.NewLabels(labels),
						"Age":     component.NewTimestamp(now),
						"Hosts":   component.NewMarkdownText("[echo1.example.com](https://echo1.example.com), [echo2.example.com](https://echo2.example.com)"),
						"Address": component.NewText(""),
						"Ports":   component.NewText("80, 443"),
						component.GridActionKey: gridActionsFactory([]component.GridAction{
							buildObjectDeleteAction(t, object),
						}),
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

			ctx := context.Background()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			if tc.list != nil {
				tpo.PathForObject(&tc.list.Items[0], tc.list.Items[0].Name, "/ingress")
				tpo.pluginManager.EXPECT().ObjectStatus(ctx, &tc.list.Items[0])
			}

			tpo.objectStore.EXPECT().
				Get(gomock.Any(), store.Key{
					Namespace:  "namespace",
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "app"}).
				Return(service, nil).
				AnyTimes()
			tpo.objectStore.EXPECT().
				Get(gomock.Any(), store.Key{
					Namespace:  "namespace",
					APIVersion: "v1",
					Kind:       "Secret",
					Name:       "secret"}).
				Return(secret, nil).
				AnyTimes()

			got, err := IngressListHandler(ctx, tc.list, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, got)

		})
	}

}

func Test_IngressConfiguration(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	ingress := testutil.CreateIngress("ingress")
	ingress.CreationTimestamp = metav1.Time{Time: now}
	ingress.Labels = labels

	ingressNoBackend := testutil.CreateIngress("ingress")
	ingressNoBackend.CreationTimestamp = metav1.Time{Time: now}
	ingressNoBackend.Labels = labels
	ingressNoBackend.Spec.DefaultBackend = nil

	ingressALB := testutil.CreateIngress("ingress")
	ingressALB.Annotations = map[string]string{
		"alb.ingress.kubernetes.io/actions.ssl-redirect": `{"Type": "redirect", "RedirectConfig": { "Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}`,
	}
	ingressALB.Spec.DefaultBackend = nil
	ingressALB.Spec.Rules = []networkingv1.IngressRule{
		{
			Host: "",
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path: "/",
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: "ssl-redirect",
									Port: networkingv1.ServiceBackendPort{
										Name: "use-annotation",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	cases := []struct {
		name     string
		ingress  *networkingv1.Ingress
		expected component.Component
		isErr    bool
	}{
		{
			name:    "in general",
			ingress: ingress,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Default Backend",
					Content: component.NewLink("", "service", "/service"),
				},
			}...),
		},
		{
			name:    "no default backend",
			ingress: ingressNoBackend,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Default Backend",
					Content: component.NewText("Default is not configured"),
				},
			}...),
		},
		{
			name:    "alb ingress controller",
			ingress: ingressALB,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Default Backend",
					Content: component.NewText("Default is not configured"),
				},
				{
					Header:  "Action: ssl-redirect",
					Content: component.NewText(`{"Type": "redirect", "RedirectConfig": { "Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}`),
				},
			}...)},
		{
			name:    "nil ingress",
			ingress: nil,
			isErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			if tc.ingress != nil {
				stubIngressBackendLinks(tpo)
			}

			ic := NewIngressConfiguration(tc.ingress)

			summary, err := ic.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createIngressRules(t *testing.T) {
	ingress := testutil.CreateIngress("ingress")

	ingressWithRules := testutil.CreateIngress("ingress")
	ingressWithRules.Spec.Rules = []networkingv1.IngressRule{
		{

			Host: "",
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path: "/",
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: "b1",
									Port: networkingv1.ServiceBackendPort{
										Number: 80,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Host: "",
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path: "/aws",
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: "ssl-redirect",
									Port: networkingv1.ServiceBackendPort{
										Name: "use-annotation",
									},
								},
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
		ingress  *networkingv1.Ingress
		expected *component.Table
		isErr    bool
	}{
		{
			name:    "in general",
			ingress: ingress,
			expected: component.NewTableWithRows("Rules", "There are no rules defined!", cols, []component.TableRow{
				{
					"Backends": component.NewLink("", "service", "/service"),
					"Host":     component.NewText("*"),
					"Path":     component.NewText("*"),
				},
			}),
		},
		{
			name:    "with rules",
			ingress: ingressWithRules,
			expected: component.NewTableWithRows("Rules", "There are no rules defined!", cols, []component.TableRow{
				{
					"Backends": component.NewLink("", "service", "/service"),
					"Host":     component.NewText("*"),
					"Path":     component.NewText("/"),
				},
				{
					"Backends": component.NewMarkdownText("*defined via use-annotation*"),
					"Host":     component.NewText("*"),
					"Path":     component.NewText("/aws"),
				},
			}),
		},
		{
			name:    "nil ingress",
			ingress: nil,
			isErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			if tc.ingress != nil {
				stubIngressBackendLinks(tpo)
			}

			got, err := createIngressRulesView(tc.ingress, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, got)
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

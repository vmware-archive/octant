/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_ServiceListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	object := &corev1.ServiceList{
		Items: []corev1.Service{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Service",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service",
					Namespace: "default",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
					Labels: labels,
				},
				Spec: corev1.ServiceSpec{
					Selector: map[string]string{
						"app": "myapp",
					},
					Type:        corev1.ServiceTypeClusterIP,
					ClusterIP:   "1.2.3.4",
					ExternalIPs: []string{"8.8.8.8", "8.8.4.4"},
					Ports: []corev1.ServicePort{
						{
							Port:     8000,
							Protocol: corev1.ProtocolTCP,
							TargetPort: intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: 8181,
							},
						},
						{
							Port:     8888,
							Protocol: corev1.ProtocolUDP,
						},
					},
				},
			},
		},
	}

	tpo.PathForObject(&object.Items[0], object.Items[0].Name, "/service")

	ctx := context.Background()
	got, err := ServiceListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Type", "Cluster IP", "External IP", "Ports", "Age", "Selector")
	expected := component.NewTable("Services", "We couldn't find any services!", cols)
	expected.Add(component.TableRow{
		"Name":        component.NewLink("", "service", "/service"),
		"Labels":      component.NewLabels(labels),
		"Type":        component.NewText("ClusterIP"),
		"Cluster IP":  component.NewText("1.2.3.4"),
		"External IP": component.NewText("8.8.8.8, 8.8.4.4"),
		"Ports":       component.NewText("8000/TCP, 8888/UDP"),
		"Age":         component.NewTimestamp(now),
		"Selector":    component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
	})

	component.AssertEqual(t, expected, got)
}

func Test_ServiceConfiguration(t *testing.T) {
	validService := &corev1.Service{
		Spec: corev1.ServiceSpec{
			ExternalTrafficPolicy:    corev1.ServiceExternalTrafficPolicyTypeCluster,
			HealthCheckNodePort:      31311,
			LoadBalancerSourceRanges: []string{"range1", "range2"},
			Ports: []corev1.ServicePort{
				{Name: "http", Port: 8080},
			},
			Selector:        map[string]string{"app": "app1"},
			SessionAffinity: corev1.ServiceAffinityNone,
			Type:            corev1.ServiceTypeClusterIP,
		},
	}
	validService.Namespace = "default"

	cases := []struct {
		name     string
		service  *corev1.Service
		isErr    bool
		expected *component.Summary
	}{
		{
			name:    "general",
			service: validService,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Selectors",
					Content: component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "app1")}),
				},
				{
					Header:  "Type",
					Content: component.NewText("ClusterIP"),
				},
				{
					Header:  "Ports",
					Content: component.NewText("http 8080/TCP"),
				},
				{
					Header:  "Session Affinity",
					Content: component.NewText("None"),
				},
				{
					Header:  "External Traffic Policy",
					Content: component.NewText("Cluster"),
				},
				{
					Header:  "Health Check Node Port",
					Content: component.NewText("31311"),
				},
				{
					Header:  "Load Balancer Source Ranges",
					Content: component.NewText("range1, range2"),
				},
			}...),
		},
		{
			name:    "service is nil",
			service: nil,
			isErr:   true,
		},
	}
	for _, tc := range cases {
		controller := gomock.NewController(t)
		defer controller.Finish()

		tpo := newTestPrinterOptions(controller)
		printOptions := tpo.ToOptions()

		ctx := context.Background()

		podKey := store.Key{
			Namespace:  "default",
			APIVersion: "v1",
			Kind:       "Pod",
		}

		pods := testutil.ToUnstructuredList(t, testutil.CreatePod("pod"))

		tpo.objectStore.EXPECT().
			List(gomock.Any(), podKey).
			Return(pods, false, nil).AnyTimes()

		sc := NewServiceConfiguration(tc.service)

		summary, err := sc.Create(ctx, printOptions)
		if tc.isErr {
			require.Error(t, err)
			return
		}
		require.NoError(t, err)

		editAction, err := editServiceAction(ctx, tc.service, printOptions)
		require.NoError(t, err)
		tc.expected.AddAction(editAction)

		component.AssertEqual(t, tc.expected, summary)
	}
}

func Test_createServiceSummaryStatus(t *testing.T) {
	cases := []struct {
		name     string
		service  *corev1.Service
		sections []component.SummarySection
	}{
		{
			name: "from spec",
			service: &corev1.Service{
				Spec: corev1.ServiceSpec{
					ClusterIP:      "10.5.5.5",
					ExternalIPs:    []string{"10.20.1.5", "10.21.1.6"},
					ExternalName:   "my-service",
					LoadBalancerIP: "10.100.1.32",
				},
			},
			sections: []component.SummarySection{
				{
					Header:  "Cluster IP",
					Content: component.NewText("10.5.5.5"),
				},
				{
					Header:  "External IPs",
					Content: component.NewText("10.20.1.5, 10.21.1.6"),
				},
				{
					Header:  "Load Balancer IP",
					Content: component.NewText("10.100.1.32"),
				},
				{
					Header:  "External Name",
					Content: component.NewText("my-service"),
				},
			},
		},
		{
			name: "from ingress",
			service: &corev1.Service{
				Spec: corev1.ServiceSpec{
					ClusterIP:      "10.5.5.5",
					ExternalName:   "my-service",
					LoadBalancerIP: "10.100.1.32",
				},
				Status: corev1.ServiceStatus{
					LoadBalancer: corev1.LoadBalancerStatus{
						Ingress: []corev1.LoadBalancerIngress{
							{
								Hostname: "example.com",
							},
						},
					},
				},
			},
			sections: []component.SummarySection{
				{
					Header:  "Cluster IP",
					Content: component.NewText("10.5.5.5"),
				},
				{
					Header:  "External IPs",
					Content: component.NewText("example.com"),
				},
				{
					Header:  "Load Balancer IP",
					Content: component.NewText("10.100.1.32"),
				},
				{
					Header:  "External Name",
					Content: component.NewText("my-service"),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := createServiceSummaryStatus(tc.service)
			require.NoError(t, err)

			expected := component.NewSummary("Status", tc.sections...)
			component.AssertEqual(t, expected, got)
		})
	}
}

func Test_createServiceEndpointsView(t *testing.T) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "service",
		},
	}

	nodeName := "node"
	endpoints := &corev1.Endpoints{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Endpoints"},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						TargetRef: &corev1.ObjectReference{
							Kind:      "Pod",
							Name:      "pod-1",
							Namespace: "default",
						},
						NodeName: &nodeName,
						IP:       "10.1.1.1",
					},
				},
			},
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	key := store.Key{Namespace: "default", APIVersion: "v1", Kind: "Endpoints", Name: "service"}
	tpo.objectStore.EXPECT().
		Get(gomock.Any(), gomock.Eq(key)).
		Return(toUnstructured(t, endpoints), true, nil)

	podLink := component.NewLink("", "pod", "/pod")
	tpo.link.EXPECT().
		ForGVK(gomock.Any(), "v1", "Pod", gomock.Any(), gomock.Any()).
		Return(podLink, nil).
		AnyTimes()

	ctx := context.Background()
	got, err := createServiceEndpointsView(ctx, service, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Target", "IP", "Node Name")
	expected := component.NewTable("Endpoints", "There are no endpoints!", cols)
	expected.Add(component.TableRow{
		"Target":    component.NewLink("", "pod", "/pod"),
		"IP":        component.NewText("10.1.1.1"),
		"Node Name": component.NewText("node"),
	})

	component.AssertEqual(t, expected, got)
}

func Test_describePortShort(t *testing.T) {
	port := corev1.ServicePort{
		Port:       8080,
		TargetPort: intstr.FromInt(80),
		Protocol:   corev1.ProtocolTCP,
	}

	got := describePortShort(port)
	expected := "8080/TCP"
	assert.Equal(t, expected, got)
}

func Test_describePort(t *testing.T) {
	cases := []struct {
		name     string
		port     corev1.ServicePort
		expected string
	}{
		{
			name: "port",
			port: corev1.ServicePort{
				Name: "http",
				Port: 80,
			},
			expected: "http 80/TCP",
		},
		{
			name: "port is not named",
			port: corev1.ServicePort{
				Port: 80,
			},
			expected: "80/TCP",
		},
		{
			name: "has node port",
			port: corev1.ServicePort{
				Name:     "http",
				NodePort: 31000,
				Port:     80,
			},
			expected: "http 80:31000/TCP",
		},
		{
			name: "port has target port",
			port: corev1.ServicePort{
				Name:     "http",
				NodePort: 31000,
				TargetPort: intstr.IntOrString{
					Type:   intstr.String,
					StrVal: "http",
				},
				Port: 80,
			},
			expected: "http 80:31000/TCP -> http",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := describePort(tc.port)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func toUnstructured(t *testing.T, object runtime.Object) *unstructured.Unstructured {
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	require.NoError(t, err)

	return &unstructured.Unstructured{Object: m}
}

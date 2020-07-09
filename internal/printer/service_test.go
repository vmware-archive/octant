/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"github.com/vmware-tanzu/octant/internal/portforward"
	"github.com/vmware-tanzu/octant/pkg/action"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

	endpoints := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Endpoints",
		},
	}

	tpo.objectStore.EXPECT().
		Get(gomock.Any(), store.Key{
			Namespace:  "default",
			APIVersion: "v1",
			Kind:       "Endpoints",
			Name:       "service",
		}).Return(endpoints, nil).AnyTimes()

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
		"Name": component.NewLink("", "service", "/service",
			genObjectStatus(component.TextStatusWarning, []string{
				"Service has no endpoint addresses",
			})),
		"Labels":      component.NewLabels(labels),
		"Type":        component.NewText("ClusterIP"),
		"Cluster IP":  component.NewText("1.2.3.4"),
		"External IP": component.NewText("8.8.8.8, 8.8.4.4"),
		"Ports":       component.NewText("8000/TCP, 8888/UDP"),
		"Age":         component.NewTimestamp(now),
		"Selector":    component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, &object.Items[0]),
		}),
	})

	component.AssertEqual(t, expected, got)
}

func createServiceWithPort(port corev1.ServicePort) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "default"},
		TypeMeta:   metav1.TypeMeta{APIVersion: "api/v1", Kind: "Service"},
		Spec: corev1.ServiceSpec{
			ExternalTrafficPolicy:    corev1.ServiceExternalTrafficPolicyTypeCluster,
			HealthCheckNodePort:      31311,
			LoadBalancerSourceRanges: []string{"range1", "range2"},
			Ports: []corev1.ServicePort{
				port,
			},
			Selector:        map[string]string{"app": "app1"},
			SessionAffinity: corev1.ServiceAffinityNone,
			Type:            corev1.ServiceTypeClusterIP,
		},
	}
}

func createExpected(targetPortName string, state component.PortForwardState, button *component.ButtonGroup) *component.Summary {
	return component.NewSummary("Configuration", []component.SummarySection{
		{
			Header:  "Selectors",
			Content: component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "app1")}),
		},
		{
			Header:  "Type",
			Content: component.NewText("ClusterIP"),
		},
		{
			Header: "Ports",
			Content: component.NewPorts([]component.Port{
				{
					Config: component.PortConfig{
						Port:           8181,
						Protocol:       "TCP",
						TargetPort:     8080,
						TargetPortName: targetPortName,
						State:          state,
						Button:         button,
					},
				},
			}),
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
	}...)
}

func Test_ServiceConfiguration(t *testing.T) {
	validService := createServiceWithPort(corev1.ServicePort{Name: "http", Port: 8181, TargetPort: intstr.FromInt(8080), Protocol: corev1.ProtocolTCP})
	validServiceWithNamedPort := createServiceWithPort(corev1.ServicePort{Name: "http", Port: 8181, TargetPort: intstr.FromString("pod-port"), Protocol: corev1.ProtocolTCP})

	startPortForwardingButtonGroup := component.NewButtonGroup()
	startPortForwardingButtonGroup.Config = component.ButtonGroupConfig{Buttons: []component.Button{
		component.NewButton("Start port forward", action.Payload{
			"action":     "overview/startPortForward",
			"apiVersion": validService.APIVersion,
			"kind":       validService.Kind,
			"name":       validService.Name,
			"namespace":  validService.Namespace,
			"port":       validService.Spec.Ports[0].TargetPort.IntVal,
		}),
	}}

	stopPortForwardingButtonGroup := component.NewButtonGroup()
	stopPortForwardingButtonGroup.Config = component.ButtonGroupConfig{Buttons: []component.Button{
		component.NewButton("Stop port forward", action.Payload{
			"action": "overview/stopPortForward",
			"id":     "an-id",
		}),
	}}
	cases := []struct {
		name     string
		service  *corev1.Service
		isErr    bool
		expected *component.Summary
		states   []portforward.State
	}{
		{
			name:    "port-forwarding-already-running",
			service: &validService,
			expected: createExpected("", component.PortForwardState{
				IsForwardable: true,
				IsForwarded:   true,
				Port:          45275,
				ID:            "an-id",
			}, stopPortForwardingButtonGroup),
			states: []portforward.State{
				{
					ID:        "an-id",
					CreatedAt: testutil.Time(),
					Ports: []portforward.ForwardedPort{
						{
							Local:  uint16(45275),
							Remote: uint16(8080),
						},
					},
					Pod: portforward.Target{
						GVK:       schema.GroupVersionKind{Group: "", Version: "api/v1", Kind: "Pod"},
						Namespace: "namespace",
						Name:      "pod",
					},
					Target: portforward.Target{
						GVK:       schema.GroupVersionKind{Group: "", Version: "api/v1", Kind: "Service"},
						Namespace: "default",
						Name:      "service",
					},
				},
			},
		},
		{
			name:    "port-forwarding-not-running",
			service: &validService,
			expected: createExpected("", component.PortForwardState{
				IsForwardable: true,
			}, startPortForwardingButtonGroup),
			states: []portforward.State{},
		},
		{
			name:    "port-forwarding-not-running-named-port",
			service: &validServiceWithNamedPort,
			expected: createExpected("pod-port", component.PortForwardState{
				IsForwardable: true,
			}, startPortForwardingButtonGroup),
			states: []portforward.State{},
		},
		{
			name:    "service is nil",
			service: nil,
			isErr:   true,
			states:  []portforward.State{},
		},
	}
	for _, tc := range cases {
		func() {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			pf := tpo.portForwarder
			printOptions := tpo.ToOptions()

			ctx := context.Background()

			pods := testutil.ToUnstructuredList(t, testutil.CreatePod("pod", func(pod *corev1.Pod) {
				pod.Spec.Containers = []corev1.Container{
					{
						Ports: []corev1.ContainerPort{
							{
								Name:          "pod-port",
								ContainerPort: 8080,
							},
						},
					},
				}
			}))

			labelSet := labels.Set(map[string]string{"app": "app1"})
			podSelectorKey := store.Key{
				Namespace:  "default",
				APIVersion: "api/v1",
				Kind:       "Pod",
				Selector:   &labelSet,
			}

			tpo.objectStore.EXPECT().
				List(gomock.Any(), podSelectorKey).
				Return(pods, false, nil).AnyTimes()

			podKey := store.Key{
				Namespace:  "default",
				APIVersion: "v1",
				Kind:       "Pod",
			}

			tpo.objectStore.EXPECT().
				List(gomock.Any(), podKey).
				Return(pods, false, nil).AnyTimes()

			pf.EXPECT().FindTarget("default", gomock.Any(), "service").
				Return(tc.states, nil).AnyTimes()

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
		}()
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
	cols := component.NewTableCols("Target", "IP", "Node Name")

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

	cases := []struct {
		name    string
		service *corev1.Service
		table   *component.Table
		rows    component.TableRow
	}{
		{
			name: "endpoint",
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "service",
				},
			},
			table: component.NewTable("Endpoints", "There are no endpoints!", cols),
			rows: component.TableRow{
				"Target":    component.NewLink("", "pod", "/pod"),
				"IP":        component.NewText("10.1.1.1"),
				"Node Name": component.NewText("node"),
			},
		},
		{
			name: "externalName",
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "service",
				},
				Spec: corev1.ServiceSpec{
					ExternalName: "test",
				},
			},
			table: component.NewTable("Endpoints", "There are no endpoints!", cols),
		},
	}

	for _, tc := range cases {
		controller := gomock.NewController(t)
		defer controller.Finish()

		tpo := newTestPrinterOptions(controller)
		printOptions := tpo.ToOptions()

		if tc.service.Spec.ExternalName == "" {
			key := store.Key{Namespace: "default", APIVersion: "v1", Kind: "Endpoints", Name: "service"}
			tpo.objectStore.EXPECT().
				Get(gomock.Any(), gomock.Eq(key)).
				Return(toUnstructured(t, endpoints), nil)

			podLink := component.NewLink("", "pod", "/pod")
			tpo.link.EXPECT().
				ForGVK(gomock.Any(), "v1", "Pod", gomock.Any(), gomock.Any()).
				Return(podLink, nil).
				AnyTimes()
		}

		ctx := context.Background()
		got, err := createServiceEndpointsView(ctx, tc.service, printOptions)
		require.NoError(t, err)

		if tc.rows != nil {
			tc.table.Add(tc.rows)
		}

		component.AssertEqual(t, tc.table, got)
	}
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

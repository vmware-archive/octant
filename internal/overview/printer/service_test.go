package printer

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/cache"
	fakecache "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Test_ServiceListHandler(t *testing.T) {
	printOptions := Options{
		Cache: cache.NewMemoryCache(),
	}

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := &corev1.ServiceList{
		Items: []corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "service",
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
						corev1.ServicePort{
							Port:     8000,
							Protocol: corev1.ProtocolTCP,
							TargetPort: intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: 8181,
							},
						},
						corev1.ServicePort{
							Port:     8888,
							Protocol: corev1.ProtocolUDP,
						},
					},
				},
			},
		},
	}

	got, err := ServiceListHandler(object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Type", "Cluster IP", "External IP", "Target Ports", "Age", "Selector")
	expected := component.NewTable("Services", cols)
	expected.Add(component.TableRow{
		"Name":         component.NewLink("", "service", "/content/overview/discovery-and-load-balancing/services/service"),
		"Labels":       component.NewLabels(labels),
		"Type":         component.NewText("", "ClusterIP"),
		"Cluster IP":   component.NewText("", "1.2.3.4"),
		"External IP":  component.NewText("", "8.8.8.8,8.8.4.4"),
		"Target Ports": component.NewText("", "8181/TCP, 8888/UDP"),
		"Age":          component.NewTimestamp(now),
		"Selector":     component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
	})

	assert.Equal(t, expected, got)
}

func Test_describeServiceConfiguration(t *testing.T) {
	service := &corev1.Service{
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

	got, err := serviceConfiguration(service)
	require.NoError(t, err)

	sections := []component.SummarySection{
		{
			Header:  "Selectors",
			Content: component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "app1")}),
		},
		{
			Header:  "Type",
			Content: component.NewText("", "ClusterIP"),
		},
		{
			Header:  "Ports",
			Content: component.NewText("", "http 8080/TCP"),
		},
		{
			Header:  "Session Affinity",
			Content: component.NewText("", "None"),
		},
		{
			Header:  "External Traffic Policy",
			Content: component.NewText("", "Cluster"),
		},
		{
			Header:  "Health Check Node Port",
			Content: component.NewText("", "31311"),
		},
		{
			Header:  "Load Balancer Source Ranges",
			Content: component.NewText("", "range1, range2"),
		},
	}

	expected := component.NewSummary("Configuration", sections...)
	assert.Equal(t, expected, got)
}

func Test_serviceSummary(t *testing.T) {
	service := &corev1.Service{
		Spec: corev1.ServiceSpec{
			ClusterIP:      "10.5.5.5",
			ExternalIPs:    []string{"10.20.1.5", "10.21.1.6"},
			ExternalName:   "my-service",
			LoadBalancerIP: "10.100.1.32",
		},
	}

	got, err := serviceSummary(service)
	require.NoError(t, err)

	sections := []component.SummarySection{
		{
			Header:  "Cluster IP",
			Content: component.NewText("", "10.5.5.5"),
		},
		{
			Header:  "External IPs",
			Content: component.NewText("", "10.20.1.5, 10.21.1.6"),
		},
		{
			Header:  "Load Balancer IP",
			Content: component.NewText("", "10.100.1.32"),
		},
		{
			Header:  "External Name",
			Content: component.NewText("", "my-service"),
		},
	}

	expected := component.NewSummary("Status", sections...)
	assert.Equal(t, expected, got)
}

func Test_serviceEndpoints(t *testing.T) {
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
							APIVersion: "v1",
							Kind:       "Pod",
							Name:       "pod-1",
						},
						NodeName: &nodeName,
						IP:       "10.1.1.1",
					},
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := fakecache.NewMockCache(ctrl)

	key := cache.Key{Namespace: "default", APIVersion: "v1", Kind: "Endpoints", Name: "service"}
	c.EXPECT().
		Get(gomock.Eq(key)).
		Return(toUnstructured(t, endpoints), nil)

	got, err := serviceEndpoints(c, service)
	require.NoError(t, err)

	cols := component.NewTableCols("Target", "IP", "Node Name")
	expected := component.NewTable("Endpoints", cols)
	expected.Add(component.TableRow{
		"Target":    component.NewLink("", "pod-1", "/content/overview/workloads/pods/pod-1"),
		"IP":        component.NewText("", "10.1.1.1"),
		"Node Name": component.NewText("", "node"),
	})

	assert.Equal(t, expected, got)
}

func Test_describeTargetPort(t *testing.T) {
	port := corev1.ServicePort{
		Port:     8080,
		Protocol: corev1.ProtocolTCP,
	}

	got := describeTargetPort(port)
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

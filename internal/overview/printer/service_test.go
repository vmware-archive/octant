package printer_test

import (
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Test_ServiceListHandler(t *testing.T) {
	printOptions := printer.Options{
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
							Protocol: corev1.ProtocolTCP,
							TargetPort: intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: 8000,
							},
						},
						corev1.ServicePort{
							Protocol: corev1.ProtocolUDP,
							TargetPort: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: "8888",
							},
						},
					},
				},
			},
		},
	}

	got, err := printer.ServiceListHandler(object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Type", "Cluster IP", "External IP", "Ports", "Age", "Selector")
	expected := component.NewTable("Services", cols)
	expected.Add(component.TableRow{
		"Name":        component.NewLink("", "service", "/content/overview/discovery-and-load-balancing/services/service"),
		"Labels":      component.NewLabels(labels),
		"Type":        component.NewText("", "ClusterIP"),
		"Cluster IP":  component.NewText("", "1.2.3.4"),
		"External IP": component.NewText("", "8.8.8.8,8.8.4.4"),
		"Ports":       component.NewText("", "8000/TCP,8888/UDP"),
		"Age":         component.NewTimestamp(now),
		"Selector":    component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
	})

	assert.Equal(t, expected, got)
}

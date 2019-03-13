package printer

import (
	"testing"

	"github.com/heptio/developer-dash/internal/portforward"

	"github.com/golang/mock/gomock"

	pffake "github.com/heptio/developer-dash/internal/portforward/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	propagation    = corev1.MountPropagationHostToContainer
	validContainer = &corev1.Container{
		Name:  "nginx",
		Image: "nginx:1.15",
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:     "http",
				HostPort: 80,
				Protocol: corev1.ProtocolTCP,
			},
			corev1.ContainerPort{
				Name:     "metrics",
				HostPort: 8080,
				Protocol: corev1.ProtocolTCP,
			},
			corev1.ContainerPort{
				Name:          "tls",
				ContainerPort: 443,
				Protocol:      corev1.ProtocolTCP,
			},
			corev1.ContainerPort{
				Name:          "dtls",
				ContainerPort: 443,
				Protocol:      corev1.ProtocolUDP,
			},
		},
		Command: []string{"/usr/bin/nginx"},
		Args:    []string{"-v", "-p", "80"},

		VolumeMounts: []corev1.VolumeMount{
			corev1.VolumeMount{
				Name:      "config",
				ReadOnly:  true,
				MountPath: "/etc/nginx",
			},
			corev1.VolumeMount{
				Name:             "data",
				MountPath:        "/var/www",
				SubPath:          "/content",
				MountPropagation: &propagation,
			},
		},
		Env: []corev1.EnvVar{
			corev1.EnvVar{
				Name:  "tier",
				Value: "prod",
			},
			corev1.EnvVar{
				Name: "fieldref",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.name"},
				},
			},
			corev1.EnvVar{
				Name: "resourcefieldref",
				ValueFrom: &corev1.EnvVarSource{
					ResourceFieldRef: &corev1.ResourceFieldSelector{
						Resource: "requests.cpu",
					},
				},
			},
			corev1.EnvVar{
				Name: "configmapref",
				ValueFrom: &corev1.EnvVarSource{
					ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{Name: "myconfig"},
						Key:                  "somekey",
					},
				},
			},
			corev1.EnvVar{
				Name: "secretref",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{Name: "mysecret"},
						Key:                  "somesecretkey",
					},
				},
			},
		},
		EnvFrom: []corev1.EnvFromSource{
			corev1.EnvFromSource{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: "fromconfig"},
				},
			},
			corev1.EnvFromSource{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: "fromsecret"},
				},
			},
		},
	}
)

func Test_ContainerConfiguration(t *testing.T) {
	envTable := component.NewTable("Environment",
		component.NewTableCols("Name", "Value", "Source"))
	envTable.Add(
		component.TableRow{
			"Name":   component.NewText("tier"),
			"Value":  component.NewText("prod"),
			"Source": component.NewText(""),
		},
		component.TableRow{
			"Name":   component.NewText("fieldref"),
			"Value":  component.NewText(""),
			"Source": component.NewText("metadata.name"),
		},
		component.TableRow{
			"Name":   component.NewText("resourcefieldref"),
			"Value":  component.NewText(""),
			"Source": component.NewText("requests.cpu"),
		},
		component.TableRow{
			"Name":   component.NewText("configmapref"),
			"Value":  component.NewText(""),
			"Source": component.NewLink("", "myconfig:somekey", "/content/overview/namespace/namespace/config-and-storage/config-maps/myconfig"),
		},
		component.TableRow{
			"Name":   component.NewText("secretref"),
			"Value":  component.NewText(""),
			"Source": component.NewLink("", "mysecret:somesecretkey", "/content/overview/namespace/namespace/config-and-storage/secrets/mysecret"),
		},
		// EnvFromSource
		component.TableRow{
			"Source": component.NewLink("", "fromconfig", "/content/overview/namespace/namespace/config-and-storage/config-maps/fromconfig"),
		},
		component.TableRow{
			"Source": component.NewLink("", "fromsecret", "/content/overview/namespace/namespace/config-and-storage/secrets/fromsecret"),
		},
	)

	volTable := component.NewTable("Volume Mounts",
		component.NewTableCols("Name", "Mount Path", "Propagation"))
	volTable.Add(
		component.TableRow{
			"Name":        component.NewText("config"),
			"Mount Path":  component.NewText("/etc/nginx (ro)"),
			"Propagation": component.NewText(""),
		},
		component.TableRow{
			"Name":        component.NewText("data"),
			"Mount Path":  component.NewText("/var/www/content (rw)"),
			"Propagation": component.NewText("HostToContainer"),
		},
	)

	cases := []struct {
		name      string
		container *corev1.Container
		isErr     bool
		expected  *component.Summary
	}{
		{
			name:      "general",
			container: validContainer,
			expected: component.NewSummary("Container nginx", []component.SummarySection{
				{
					Header:  "Image",
					Content: component.NewText("nginx:1.15"),
				},
				{
					Header:  "Host Ports",
					Content: component.NewText("80/TCP, 8080/TCP"),
				},
				{
					Header: "Container Ports",
					Content: component.NewPorts([]component.Port{
						*component.NewPort("namespace", "v1", "Pod", "pod", 443, "TCP", component.PortForwardState{IsForwardable: true, IsForwarded: true}),
						*component.NewPort("namespace", "v1", "Pod", "pod", 443, "UDP", component.PortForwardState{IsForwardable: false, IsForwarded: false}),
					}),
				},
				{
					Header:  "Environment",
					Content: envTable,
				},
				{
					Header:  "Command",
					Content: component.NewText("['/usr/bin/nginx']"),
				},
				{
					Header:  "Args",
					Content: component.NewText("['-v', '-p', '80']"),
				},
				{
					Header:  "Volume Mounts",
					Content: volTable,
				},
			}...),
		},
		{
			name:      "container is nil",
			container: nil,
			isErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			pf := pffake.NewMockPortForwarder(controller)
			gvk := schema.GroupVersionKind{Version: "v1", Kind: "Pod"}

			state := portforward.State{}
			pf.EXPECT().Find("namespace", gomock.Eq(gvk), "pod").Return(state, nil).AnyTimes()

			parentPod := testutil.CreatePod("pod")
			parentPod.Namespace = "namespace"

			cc := NewContainerConfiguration(parentPod, tc.container, pf, false)
			summary, err := cc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, summary)
		})
	}
}

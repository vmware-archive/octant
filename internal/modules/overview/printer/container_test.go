package printer

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/heptio/developer-dash/internal/portforward"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	pffake "github.com/heptio/developer-dash/internal/portforward/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
)

var (
	propagation    = corev1.MountPropagationHostToContainer
	validContainer = &corev1.Container{
		Name:  "nginx",
		Image: "nginx:1.15",
		Ports: []corev1.ContainerPort{
			{
				Name:     "http",
				HostPort: 80,
				Protocol: corev1.ProtocolTCP,
			},
			{
				Name:     "metrics",
				HostPort: 8080,
				Protocol: corev1.ProtocolTCP,
			},
			{
				Name:          "tls",
				ContainerPort: 443,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				Name:          "dtls",
				ContainerPort: 443,
				Protocol:      corev1.ProtocolUDP,
			},
		},
		Command: []string{"/usr/bin/nginx"},
		Args:    []string{"-v", "-p", "80"},

		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "config",
				ReadOnly:  true,
				MountPath: "/etc/nginx",
			},
			{
				Name:             "data",
				MountPath:        "/var/www",
				SubPath:          "/content",
				MountPropagation: &propagation,
			},
		},
		Env: []corev1.EnvVar{
			{
				Name:  "tier",
				Value: "prod",
			},
			{
				Name: "fieldref",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.name"},
				},
			},
			{
				Name: "resourcefieldref",
				ValueFrom: &corev1.EnvVarSource{
					ResourceFieldRef: &corev1.ResourceFieldSelector{
						Resource: "requests.cpu",
					},
				},
			},
			{
				Name: "configmapref",
				ValueFrom: &corev1.EnvVarSource{
					ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{Name: "myconfig"},
						Key:                  "somekey",
					},
				},
			},
			{
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
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: "fromconfig"},
				},
			},
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: "fromsecret"},
				},
			},
		},
	}
)

func Test_ContainerConfiguration(t *testing.T) {
	now := time.Now()

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
			"Source": component.NewLink("", "myconfig:somekey", "/configMap"),
		},
		component.TableRow{
			"Name":   component.NewText("secretref"),
			"Value":  component.NewText(""),
			"Source": component.NewLink("", "mysecret:somesecretkey", "/secret"),
		},
		// EnvFromSource
		component.TableRow{
			"Source": component.NewLink("", "fromconfig", "/fromConfig"),
		},
		component.TableRow{
			"Source": component.NewLink("", "fromsecret", "/fromSecret"),
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
			name:      "in general",
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
					Header:  "Last State",
					Content: component.NewText(fmt.Sprintf("terminated with 255 at %s: reason", now)),
				},
				{
					Header:  "Current State",
					Content: component.NewText(fmt.Sprintf("started at %s", now)),
				},
				{
					Header:  "Ready",
					Content: component.NewText("true"),
				},
				{
					Header:  "Restart Count",
					Content: component.NewText("2"),
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

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			pf := pffake.NewMockPortForwarder(controller)
			gvk := schema.GroupVersionKind{Version: "v1", Kind: "Pod"}

			state := portforward.State{}
			pf.EXPECT().Find("namespace", gomock.Eq(gvk), "pod").Return(state, nil).AnyTimes()

			tpo.PathForGVK("namespace", "v1", "Secret", "mysecret", "mysecret:somesecretkey", "/secret")
			tpo.PathForGVK("namespace", "v1", "ConfigMap", "myconfig", "myconfig:somekey", "/configMap")
			tpo.PathForGVK("namespace", "v1", "Secret", "fromsecret", "fromsecret", "/fromSecret")
			tpo.PathForGVK("namespace", "v1", "ConfigMap", "fromconfig", "fromconfig", "/fromConfig")

			parentPod := testutil.CreatePod("pod")
			parentPod.Namespace = "namespace"
			parentPod.Status = corev1.PodStatus{
				ContainerStatuses: []corev1.ContainerStatus{
					{
						Name:         "nginx",
						Ready:        true,
						RestartCount: 2,
						State: corev1.ContainerState{
							Running: &corev1.ContainerStateRunning{
								StartedAt: metav1.Time{Time: now},
							},
						},
						LastTerminationState: corev1.ContainerState{
							Terminated: &corev1.ContainerStateTerminated{
								FinishedAt: metav1.Time{Time: now},
								Reason:     "reason",
								ExitCode:   255,
							},
						},
					},
				},
			}

			cc := NewContainerConfiguration(parentPod, tc.container, pf, false, printOptions)
			summary, err := cc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertComponentEqual(t, tc.expected, summary)
		})
	}
}

func Test_containerNotFoundError(t *testing.T) {
	e := containerNotFoundError{name: "name"}

	expected := fmt.Sprintf("container %q not found", "name")
	assert.Equal(t, expected, e.Error())

	assert.False(t, e.isContainerFound())
}

func Test_findContainerStatus(t *testing.T) {
	tests := []struct {
		name            string
		podFactory      func(statusName string) *corev1.Pod
		statusName      string
		isInit          bool
		isErr           bool
		expectedErrType reflect.Type
		expectedStatus  *corev1.ContainerStatus
	}{
		{
			name:       "container with status",
			statusName: "name",
			podFactory: func(statusName string) *corev1.Pod {
				pod := testutil.CreatePod("pod")

				pod.Status.ContainerStatuses = append(
					pod.Status.ContainerStatuses,
					corev1.ContainerStatus{Name: statusName})

				return pod
			},
			expectedStatus: &corev1.ContainerStatus{Name: "name"},
		},
		{
			name:       "init container with status",
			isInit:     true,
			statusName: "name",
			podFactory: func(statusName string) *corev1.Pod {
				pod := testutil.CreatePod("pod")

				pod.Status.InitContainerStatuses = append(
					pod.Status.ContainerStatuses,
					corev1.ContainerStatus{Name: statusName})

				return pod
			},
			expectedStatus: &corev1.ContainerStatus{Name: "name"},
		},
		{
			name: "no containers",
			podFactory: func(statusName string) *corev1.Pod {
				pod := testutil.CreatePod("pod")
				return pod
			},
			isErr:           true,
			expectedErrType: reflect.TypeOf(&containerNotFoundError{}),
		},
		{
			name: "pod is nil",
			podFactory: func(statusName string) *corev1.Pod {
				return nil
			},
			isErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pod := test.podFactory(test.statusName)

			status, err := findContainerStatus(pod, test.statusName, test.isInit)
			if test.isErr {
				require.Error(t, err)

				if test.expectedErrType != nil {
					errType := reflect.TypeOf(err)
					require.Equal(t, test.expectedErrType, errType)
				}

				return
			}
			require.NoError(t, err)

			require.Equal(t, test.expectedStatus, status)
		})
	}
}

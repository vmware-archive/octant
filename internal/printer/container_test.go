/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vmware-tanzu/octant/internal/portforward"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	pffake "github.com/vmware-tanzu/octant/internal/portforward/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_ContainerConfiguration(t *testing.T) {
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
					Name:          "application",
					ContainerPort: 8443,
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
		validInitContainer = &corev1.Container{
			Name:    "busybox",
			Image:   "busybox:1.28",
			Command: []string{"sh"},
			Args:    []string{"-c", "until nslookup mydb; do echo waiting for mydb; sleep 2; done;"},
		}
	)

	now := time.Now()

	envTable := component.NewTable("Environment", "There are no defined environment variables!",
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

	volTable := component.NewTable("Volume Mounts", "There are no volume mounts!",
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
		isInit    bool
		action    component.Action
		expected  *component.Summary
	}{
		{
			name:      "in general",
			container: validContainer,
			action: component.Action{
				Name:  "Execute Command",
				Title: "Execute nginx Command",
				Form: component.Form{
					Fields: []component.FormField{
						component.NewFormFieldHidden("containerName", "nginx"),
						component.NewFormFieldText("Command", "containerCommand", ""),
						component.NewFormFieldHidden("apiVersion", "v1"),
						component.NewFormFieldHidden("kind", "Pod"),
						component.NewFormFieldHidden("name", "pod"),
						component.NewFormFieldHidden("namespace", "namespace"),
						component.NewFormFieldHidden("action", "overview/commandExec"),
					},
				},
			},
			expected: component.NewSummary("Container nginx", []component.SummarySection{
				{
					Header:  "Image",
					Content: component.NewText("nginx:1.15"),
				},
				{
					Header:  "Image ID",
					Content: component.NewText("nginx-image-id"),
				},
				{
					Header:  "Host Ports",
					Content: component.NewText("80/TCP, 8080/TCP"),
				},
				{
					Header: "Container Ports",
					Content: component.NewPorts([]component.Port{
						*component.NewPort("namespace", "v1", "Pod", "pod", 8443, "TCP", component.PortForwardState{IsForwardable: true, IsForwarded: true}),
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
			name:      "init containers",
			container: validInitContainer,
			isInit:    true,
			action: component.Action{
				Name:  "Execute Command",
				Title: "Execute busybox Command",
				Form: component.Form{
					Fields: []component.FormField{
						component.NewFormFieldHidden("containerName", "busybox"),
						component.NewFormFieldText("Command", "containerCommand", ""),
						component.NewFormFieldHidden("apiVersion", "v1"),
						component.NewFormFieldHidden("kind", "Pod"),
						component.NewFormFieldHidden("name", "pod"),
						component.NewFormFieldHidden("namespace", "namespace"),
						component.NewFormFieldHidden("action", "overview/commandExec"),
					},
				},
			},
			expected: component.NewSummary("Init Container busybox", []component.SummarySection{
				{
					Header:  "Image",
					Content: component.NewText("busybox:1.28"),
				},
				{
					Header:  "Image ID",
					Content: component.NewText("busybox-image-id"),
				},
				{
					Header:  "Command",
					Content: component.NewText("['sh']"),
				},
				{
					Header:  "Args",
					Content: component.NewText("['-c', 'until nslookup mydb; do echo waiting for mydb; sleep 2; done;']"),
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

			ctx := context.Background()
			configMap := testutil.CreateConfigMap("myconfig")

			pf := pffake.NewMockPortForwarder(controller)
			gvk := schema.GroupVersionKind{Version: "v1", Kind: "Pod"}

			states := []portforward.State{
				{
					CreatedAt: testutil.Time(),
					Ports: []portforward.ForwardedPort{
						{
							Local:  uint16(45275),
							Remote: uint16(8443),
						},
					},
					Pod: portforward.Target{
						GVK:       gvk,
						Namespace: "namespace",
						Name:      "pod",
					},
				},
			}

			state := createPortForwardState("stateid", "namespace", "pod", gvk)

			pf.EXPECT().Find("namespace", gomock.Eq(gvk), "pod").Return(states, nil).AnyTimes()
			pf.EXPECT().Get(gomock.Any()).Return(state, true).AnyTimes()

			tpo.PathForGVK("namespace", "v1", "Secret", "mysecret", "mysecret:somesecretkey", "/secret")
			tpo.PathForGVK("namespace", "v1", "ConfigMap", "myconfig", "myconfig:somekey", "/configMap")
			tpo.PathForGVK("namespace", "v1", "Secret", "fromsecret", "fromsecret", "/fromSecret")
			tpo.PathForGVK("namespace", "v1", "ConfigMap", "fromconfig", "fromconfig", "/fromConfig")

			parentPod := testutil.CreatePod("pod")
			parentPod.Namespace = "namespace"
			parentPod.Status = corev1.PodStatus{
				InitContainerStatuses: []corev1.ContainerStatus{
					{
						Name:    "busybox",
						ImageID: "busybox-image-id",
					},
				},
				ContainerStatuses: []corev1.ContainerStatus{
					{
						Name:         "nginx",
						ImageID:      "nginx-image-id",
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

			if tc.container != nil {
				key := store.Key{
					Kind:       configMap.Kind,
					APIVersion: configMap.APIVersion,
					Name:       configMap.Name,
					Namespace:  configMap.Namespace,
				}
				tpo.objectStore.EXPECT().Get(ctx, gomock.Eq(key)).Return(testutil.ToUnstructured(t, configMap), true, nil).AnyTimes()
			}

			cc := NewContainerConfiguration(ctx, parentPod, tc.container, pf, tc.isInit, printOptions)
			summary, err := cc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			tc.expected.AddAction(tc.action)

			component.AssertEqual(t, tc.expected, summary)
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

func Test_editContainerAction(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment", testutil.WithGenericDeployment())
	container := deployment.Spec.Template.Spec.Containers[0]

	got, err := editContainerAction(deployment, &container)
	require.NoError(t, err)

	form, err := component.CreateFormForObject("overview/containerEditor", deployment,
		component.NewFormFieldText("Image", "containerImage", container.Image),
		component.NewFormFieldHidden("containersPath", `["spec","template","spec","containers"]`),
		component.NewFormFieldHidden("containerName", container.Name),
	)
	require.NoError(t, err)

	expected := component.Action{
		Name:  "Edit",
		Title: "Container container-name Editor",
		Form:  form,
	}
	require.Equal(t, expected, got)
}

func createPortForwardState(id, namespace, targetName string, gvk schema.GroupVersionKind) portforward.State {
	return portforward.State{
		ID:        id,
		CreatedAt: testutil.Time(),
		Pod: portforward.Target{
			GVK:       gvk,
			Namespace: namespace,
			Name:      targetName,
		},
		Target: portforward.Target{
			GVK:       gvk,
			Namespace: namespace,
			Name:      targetName,
		}}
}

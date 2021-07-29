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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/conversion"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_PodListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	nodeLink := component.NewLink("", "node", "/node")
	tpo.link.EXPECT().
		ForGVK("", "v1", "Node", "node", "node").
		Return(nodeLink, nil)
	printOptions := tpo.ToOptions()

	now := testutil.Time()

	labels := map[string]string{
		"app": "testing",
	}

	pod := testutil.CreatePod("pod")
	pod.CreationTimestamp = metav1.Time{Time: now}
	pod.Labels = labels
	pod.Spec.Containers = []corev1.Container{
		{
			Name:  "nginx",
			Image: "nginx:1.15",
		},
		{
			Name:  "kuard",
			Image: "gcr.io/kuar-demo/kuard-amd64:1",
		},
	}
	pod.Spec.NodeName = "node"
	pod.Status = corev1.PodStatus{
		Phase: "Pending",
		ContainerStatuses: []corev1.ContainerStatus{
			{
				Name:         "nginx",
				Image:        "nginx:1.15",
				RestartCount: 0,
				Ready:        true,
				State: corev1.ContainerState{
					Waiting: &corev1.ContainerStateWaiting{
						Reason:  "ContainerCreating",
						Message: "",
					},
					Running:    nil,
					Terminated: nil,
				},
			},
			{
				Name:         "kuard",
				Image:        "gcr.io/kuar-demo/kuard-amd64:1",
				RestartCount: 0,
				Ready:        false,
				State: corev1.ContainerState{
					Waiting: &corev1.ContainerStateWaiting{
						Reason:  "ContainerCreating",
						Message: "",
					},
					Running:    nil,
					Terminated: nil,
				},
			},
		},
	}

	object := &corev1.PodList{
		Items: []corev1.Pod{*pod},
	}

	tpo.PathForObject(pod, pod.Name, "/pod")

	ctx := context.Background()
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, pod)
	got, err := PodListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Ready", "Phase", "Status", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", "We couldn't find any pods!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "pod", "/pod",
			genObjectStatus(component.TextStatusWarning, []string{
				"Pod may require additional action",
			})),
		"Labels":   component.NewLabels(labels),
		"Ready":    component.NewText("1/2"),
		"Phase":    component.NewText("Pending"),
		"Status":   component.NewText("ContainerCreating"),
		"Restarts": component.NewText("0"),
		"Age":      component.NewTimestamp(now),
		"Node":     nodeLink,
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, pod),
		}),
	})
	addPodTableFilters(expected)

	component.AssertEqual(t, expected, got)
}

func Test_PodListHandlerNoLabel(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	nodeLink := component.NewLink("", "node", "/node")
	tpo.link.EXPECT().
		ForGVK("", "v1", "Node", "node", "node").
		Return(nodeLink, nil)
	printOptions := tpo.ToOptions()

	printOptions.DisableLabels = true
	now := testutil.Time()

	pod := testutil.CreatePod("pi-7xpxr")
	pod.CreationTimestamp = metav1.Time{Time: now}
	pod.Spec.Containers = []corev1.Container{
		{
			Name:  "pi",
			Image: "perl",
		},
	}
	pod.Spec.NodeName = "node"
	pod.Status = corev1.PodStatus{
		Phase: "Succeeded",
		ContainerStatuses: []corev1.ContainerStatus{
			{
				Name:         "pi",
				Image:        "perl",
				RestartCount: 0,
				Ready:        false,
				State: corev1.ContainerState{
					Waiting:    nil,
					Running:    &corev1.ContainerStateRunning{StartedAt: metav1.Time{Time: now}},
					Terminated: nil,
				},
			},
		},
	}

	object := &corev1.PodList{
		Items: []corev1.Pod{*pod},
	}

	tpo.PathForObject(pod, pod.Name, "/pi-7xpxr")

	ctx := context.Background()
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, pod)
	got, err := PodListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Ready", "Phase", "Status", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", "We couldn't find any pods!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "pi-7xpxr", "/pi-7xpxr",
			genObjectStatus(component.TextStatusWarning, []string{"Pod may require additional action"})),
		"Ready":    component.NewText("0/1"),
		"Phase":    component.NewText("Succeeded"),
		"Status":   component.NewText("Running"),
		"Restarts": component.NewText("0"),
		"Age":      component.NewTimestamp(now),
		"Node":     nodeLink,
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, pod),
		}),
	})
	addPodTableFilters(expected)

	component.AssertEqual(t, expected, got)
}

func Test_PodListHandlerTerminating(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	nodeLink := component.NewLink("", "node", "/node")
	tpo.link.EXPECT().
		ForGVK("", "v1", "Node", "node", "node").
		Return(nodeLink, nil)
	printOptions := tpo.ToOptions()

	now := testutil.Time()

	labels := map[string]string{
		"app": "testing",
	}

	pod := testutil.CreatePod("pi-7xpxr")
	pod.CreationTimestamp = metav1.Time{Time: now}
	pod.Labels = labels
	pod.Spec.Containers = []corev1.Container{
		{
			Name:  "pi",
			Image: "perl",
		},
	}
	pod.Spec.NodeName = "node"
	pod.Status = corev1.PodStatus{
		Phase: "Running",
		ContainerStatuses: []corev1.ContainerStatus{
			{
				Name:         "pi",
				Image:        "perl",
				RestartCount: 0,
				Ready:        false,
				State: corev1.ContainerState{
					Waiting:    nil,
					Running:    nil,
					Terminated: nil,
				},
			},
		},
	}
	pod.DeletionTimestamp = &metav1.Time{Time: now}

	object := &corev1.PodList{
		Items: []corev1.Pod{*pod},
	}

	tpo.PathForObject(pod, pod.Name, "/pi-7xpxr")

	ctx := context.Background()
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, pod)
	got, err := PodListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Ready", "Phase", "Status", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", "We couldn't find any pods!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "pi-7xpxr", "/pi-7xpxr",
			genObjectStatus(component.TextStatusWarning, []string{
				"Pod is being deleted",
			}),
		),
		"Labels":     component.NewLabels(labels),
		"Ready":      component.NewText("0/1"),
		"Phase":      component.NewText("Running"),
		"Status":     component.NewText("Terminating"),
		"Restarts":   component.NewText("0"),
		"Age":        component.NewTimestamp(now),
		"_isDeleted": component.NewText("deleted"),
		"Node":       nodeLink,
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, pod),
		}),
	})
	addPodTableFilters(expected)

	component.AssertEqual(t, expected, got)
}

func TestPodListHandler_sorted(t *testing.T) {
	pod1 := testutil.CreatePod("pod1")
	pod2 := testutil.CreatePod("pod2")

	list := &corev1.PodList{
		Items: []corev1.Pod{
			*pod2,
			*pod1,
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	tpo.PathForObject(pod1, pod1.Name, "/pod1")
	tpo.PathForObject(pod2, pod2.Name, "/pod2")

	ctx := context.Background()
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, pod1)
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, pod2)

	got, err := PodListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Ready", "Phase", "Status", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", "We couldn't find any pods!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "pod1", "/pod1",
			genObjectStatus(component.TextStatusWarning, []string{"Pod may require additional action"})),
		"Labels":   component.NewLabels(make(map[string]string)),
		"Ready":    component.NewText("0/0"),
		"Phase":    component.NewText(""),
		"Status":   component.NewText(""),
		"Restarts": component.NewText("0"),
		"Age":      component.NewTimestamp(pod1.CreationTimestamp.Time),
		"Node":     component.NewText("<not scheduled>"),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, pod1),
		}),
	})
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "pod2", "/pod2",
			genObjectStatus(component.TextStatusWarning, []string{"Pod may require additional action"})),
		"Labels":   component.NewLabels(make(map[string]string)),
		"Ready":    component.NewText("0/0"),
		"Phase":    component.NewText(""),
		"Status":   component.NewText(""),
		"Restarts": component.NewText("0"),
		"Age":      component.NewTimestamp(pod1.CreationTimestamp.Time),
		"Node":     component.NewText("<not scheduled>"),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, pod2),
		}),
	})
	addPodTableFilters(expected)

	component.AssertEqual(t, expected, got)
}

func Test_PodConfiguration(t *testing.T) {
	now := testutil.Time()
	validPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "ReplicationController",
					Name:       "myreplicationcontroller",
					Controller: conversion.PtrBool(true),
				},
			},
			CreationTimestamp: metav1.Time{
				Time: now,
			},
			DeletionTimestamp: &metav1.Time{
				Time: now,
			},
			DeletionGracePeriodSeconds: conversion.PtrInt64(30),
		},
		Spec: corev1.PodSpec{
			NodeName:           "node",
			Priority:           conversion.PtrInt32(1000000),
			PriorityClassName:  "high-priority",
			ServiceAccountName: "default",
		},
		Status: corev1.PodStatus{
			StartTime:         &metav1.Time{Time: now},
			Phase:             corev1.PodRunning,
			Reason:            "SleepExpired",
			Message:           "Sleep expired",
			NominatedNodeName: "mynode",
			QOSClass:          "Guaranteed",
		},
	}

	nodeLink := component.NewLink("", "node", "/node")

	cases := []struct {
		name     string
		pod      *corev1.Pod
		isErr    bool
		expected *component.Summary
	}{
		{
			name: "general",
			pod:  validPod,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Priority",
					Content: component.NewText("1000000"),
				},
				{
					Header:  "PriorityClassName",
					Content: component.NewText("high-priority"),
				},
				{
					Header:  "Node",
					Content: nodeLink,
				},
				{
					Header:  "Service Account",
					Content: component.NewLink("", "serviceAccount", "/service-account"),
				},
			}...),
		},
		{
			name:  "pod is nil",
			pod:   nil,
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			tpo.link.EXPECT().
				ForGVK("", "v1", "Node", "node", "node").
				Return(nodeLink, nil).AnyTimes()

			printOptions := tpo.ToOptions()

			if tc.pod != nil {
				tpo.PathForObject(tc.pod, tc.pod.Name, "/pod")

				serviceAccountLink := component.NewLink("", "serviceAccount", "/service-account")
				tpo.link.EXPECT().
					ForGVK(gomock.Any(), "v1", "ServiceAccount", gomock.Any(), gomock.Any()).
					Return(serviceAccountLink, nil).
					AnyTimes()
			}

			cc := NewPodConfiguration(tc.pod)

			summary, err := cc.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createPodSummaryStatus(t *testing.T) {
	pod := testutil.CreatePod("pod")
	pod.Status.QOSClass = corev1.PodQOSBestEffort
	pod.Status.Phase = corev1.PodRunning
	pod.Status.PodIP = "10.1.1.1"
	pod.Status.HostIP = "10.2.1.1"

	got, err := createPodSummaryStatus(pod)
	require.NoError(t, err)

	sections := component.SummarySections{
		{Header: "QoS", Content: component.NewText("BestEffort")},
		{Header: "Phase", Content: component.NewText("Running")},
		{Header: "Pod IP", Content: component.NewText("10.1.1.1")},
		{Header: "Host IP", Content: component.NewText("10.2.1.1")},
	}
	expected := component.NewSummary("Status", sections...)

	assert.Equal(t, expected, got)
}

func createPodWithPhase(name string, podLabels map[string]string, phase corev1.PodPhase, owner *metav1.OwnerReference) *corev1.Pod {
	pod := testutil.CreatePod(name)
	pod.Namespace = "testing"
	pod.Labels = podLabels
	pod.Status.Phase = phase

	if owner != nil {
		pod.SetOwnerReferences([]metav1.OwnerReference{*owner})
	}
	return pod
}

func Test_printPodResources(t *testing.T) {
	pod := testutil.CreatePod("pod")
	pod.Spec.Containers = []corev1.Container{
		{
			Name: "container-a",
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("1Mi"),
					corev1.ResourceCPU:    resource.MustParse("2Mi"),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("3Mi"),
					corev1.ResourceCPU:    resource.MustParse("4Mi"),
				},
			},
		},
	}

	got, err := printPodResources(pod.Spec)
	require.NoError(t, err)

	expected := component.NewTable("Resources", "Pod has no resource needs", podResourceCols)
	expected.Add(component.TableRow{
		"Container":       component.NewText("container-a"),
		"Request: Memory": component.NewText("1Mi"),
		"Request: CPU":    component.NewText("2Mi"),
		"Limit: Memory":   component.NewText("3Mi"),
		"Limit: CPU":      component.NewText("4Mi"),
	})

	assert.Equal(t, expected, got)
}

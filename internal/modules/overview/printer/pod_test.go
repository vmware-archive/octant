/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware/octant/internal/conversion"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_PodListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	now := testutil.Time()

	labels := map[string]string{
		"app": "testing",
	}

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod",
			Namespace: "default",
			CreationTimestamp: metav1.Time{
				Time: now,
			},
			Labels: map[string]string{
				"app": "testing",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:1.15",
				},
				{
					Name:  "kuard",
					Image: "gcr.io/kuar-demo/kuard-amd64:1",
				},
			},
			NodeName: "node",
		},
		Status: corev1.PodStatus{
			Phase: "Pending",
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:         "nginx",
					Image:        "nginx:1.15",
					RestartCount: 0,
					Ready:        true,
				},
				{
					Name:         "kuard",
					Image:        "gcr.io/kuar-demo/kuard-amd64:1",
					RestartCount: 0,
					Ready:        false,
				},
			},
		},
	}

	object := &corev1.PodList{
		Items: []corev1.Pod{*pod},
	}

	tpo.PathForObject(pod, pod.Name, "/pod")

	ctx := context.Background()
	got, err := PodListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Ready", "Phase", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", cols)
	expected.Add(component.TableRow{
		"Name":     component.NewLink("", "pod", "/pod"),
		"Labels":   component.NewLabels(labels),
		"Ready":    component.NewText("1/2"),
		"Phase":    component.NewText("Pending"),
		"Restarts": component.NewText("0"),
		"Age":      component.NewTimestamp(now),
		"Node":     component.NewText("node"),
	})

	component.AssertEqual(t, expected, got)
}

func Test_PodHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	now := testutil.Time()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"app": "testing",
	}

	sidecar := testutil.CreatePod("pod")
	sidecar.ObjectMeta.CreationTimestamp = *testutil.CreateTimestamp()
	sidecar.ObjectMeta.Labels = labels
	sidecar.Spec.Containers = []corev1.Container{
		{
			Name:  "nginx",
			Image: "nginx:1.15",
		},
		{
			Name:  "kuard",
			Image: "gcr.io/kuar-demo/kuard-amd64:1",
		},
	}

	tpo.PathForObject(sidecar, sidecar.Name, "/pod")

	serviceAccountLink := component.NewLink("", "serviceAccount", "/service-account")
	tpo.link.EXPECT().
		ForGVK(gomock.Any(), "v1", "ServiceAccount", gomock.Any(), gomock.Any()).
		Return(serviceAccountLink, nil).
		AnyTimes()

	printResponse := &plugin.PrintResponse{}
	tpo.pluginManager.EXPECT().
		Print(gomock.Any()).Return(printResponse, nil)

	key := store.Key{
		Namespace:  "namespace",
		APIVersion: "v1",
		Kind:       "Event",
	}
	eventList := []*unstructured.Unstructured{}
	tpo.objectStore.EXPECT().List(gomock.Any(), gomock.Eq(key)).Return(eventList, nil)

	ctx := context.Background()
	got, err := PodHandler(ctx, sidecar, printOptions)
	require.NoError(t, err)

	configSection := component.SummarySections{
		{Header: "Service Account", Content: component.NewLink("", "serviceAccount", "/service-account")},
	}
	configSummary := component.NewSummary("Configuration", configSection...)

	statusSections := component.SummarySections{
		{Header: "QoS", Content: component.NewText("")},
		{Header: "Phase", Content: component.NewText("")},
		{Header: "Pod IP", Content: component.NewText("")},
		{Header: "Host IP", Content: component.NewText("")},
	}
	statusSummary := component.NewSummary("Status", statusSections...)

	metadataSections := component.SummarySections{
		{Header: "Age", Content: component.NewTimestamp(now)},
		{Header: "Labels", Content: component.NewLabels(labels)},
	}
	metadataSummary := component.NewSummary("Metadata", metadataSections...)

	conditionsCols := component.NewTableCols("Type", "Last Transition Time", "Message", "Reason")
	conditionTable := component.NewTable("Pod Conditions", conditionsCols)

	container1Sections := component.SummarySections{
		{Header: "Image", Content: component.NewText("nginx:1.15")},
	}
	container1Summary := component.NewSummary("Container nginx", container1Sections...)

	container2Sections := component.SummarySections{
		{
			Header:  "Image",
			Content: component.NewText("gcr.io/kuar-demo/kuard-amd64:1"),
		},
	}
	container2Summary := component.NewSummary("Container kuard", container2Sections...)

	volumeCols := component.NewTableCols("Name", "Kind", "Description")
	volumeTable := component.NewTable("Volumes", volumeCols)

	taintCols := component.NewTableCols("Description")
	taintTable := component.NewTable("Taints and Tolerations", taintCols)

	affinityCols := component.NewTableCols("Type", "Description")
	affinityTable := component.NewTable("Affinities and Anti-Affinities", affinityCols)

	expected := component.NewFlexLayout("Summary")
	expected.AddSections(
		component.FlexLayoutSection{
			{
				Width: component.WidthHalf,
				View:  configSummary,
			},
			{
				Width: component.WidthHalf,
				View:  statusSummary,
			},
		},
		component.FlexLayoutSection{
			{
				Width: component.WidthFull,
				View:  metadataSummary,
			},
		},
		component.FlexLayoutSection{
			{
				Width: component.WidthFull,
				View:  conditionTable,
			},
		},
		component.FlexLayoutSection{},
		component.FlexLayoutSection{
			{
				Width: component.WidthHalf,
				View:  container1Summary,
			},
			{
				Width: component.WidthHalf,
				View:  container2Summary,
			},
		},
		component.FlexLayoutSection{
			{
				Width: component.WidthHalf,
				View:  volumeTable,
			},
			{
				Width: component.WidthHalf,
				View:  taintTable,
			},
			{
				Width: component.WidthHalf,
				View:  affinityTable,
			},
		},
	)

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
	got, err := PodListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Ready", "Phase", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", cols)
	expected.Add(component.TableRow{
		"Name":     component.NewLink("", "pod1", "/pod1"),
		"Labels":   component.NewLabels(make(map[string]string)),
		"Ready":    component.NewText("0/0"),
		"Phase":    component.NewText(""),
		"Restarts": component.NewText("0"),
		"Age":      component.NewTimestamp(pod1.CreationTimestamp.Time),
		"Node":     component.NewText(""),
	})
	expected.Add(component.TableRow{
		"Name":     component.NewLink("", "pod2", "/pod2"),
		"Labels":   component.NewLabels(make(map[string]string)),
		"Ready":    component.NewText("0/0"),
		"Phase":    component.NewText(""),
		"Restarts": component.NewText("0"),
		"Age":      component.NewTimestamp(pod1.CreationTimestamp.Time),
		"Node":     component.NewText(""),
	})

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

func Test_createPodConditionsView(t *testing.T) {
	now := metav1.Time{Time: time.Now()}

	pod := testutil.CreatePod("pod")
	pod.Status.Conditions = []corev1.PodCondition{
		{
			Type:               corev1.PodInitialized,
			LastTransitionTime: now,
			Message:            "message",
			Reason:             "reason",
		},
	}

	got, err := createPodConditionsView(pod)
	require.NoError(t, err)

	cols := component.NewTableCols("Type", "Last Transition Time", "Message", "Reason")
	expected := component.NewTable("Pod Conditions", cols)
	expected.Add([]component.TableRow{
		{
			"Type":                 component.NewText("Initialized"),
			"Last Transition Time": component.NewTimestamp(now.Time),
			"Message":              component.NewText("message"),
			"Reason":               component.NewText("reason"),
		},
	}...)

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

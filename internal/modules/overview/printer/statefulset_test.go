/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware/octant/internal/conversion"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_StatefulSetListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	statefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "web",
			CreationTimestamp: metav1.Time{
				Time: now,
			},
			Labels: labels,
		},
		Status: appsv1.StatefulSetStatus{
			Replicas: 1,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "myapp",
				},
			},
			Replicas: conversion.PtrInt32(3),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "k8s.gcr.io/nginx-slim:0.8",
						},
					},
				},
			},
		},
	}

	tpo.PathForObject(statefulSet, "web", "/path")

	object := &appsv1.StatefulSetList{
		Items: []appsv1.StatefulSet{*statefulSet},
	}

	ctx := context.Background()
	got, err := StatefulSetListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Desired", "Current", "Age", "Selector")
	expected := component.NewTable("StatefulSets", "We couldn't find any stateful sets!", cols)
	expected.Add(component.TableRow{
		"Name":     component.NewLink("", "web", "/path"),
		"Labels":   component.NewLabels(labels),
		"Desired":  component.NewText("3"),
		"Current":  component.NewText("1"),
		"Age":      component.NewTimestamp(now),
		"Selector": component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
	})

	component.AssertEqual(t, expected, got)
}

func Test_StatefulSetStatus(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"app": "myapp",
	}

	sts := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "statefulset",
			Namespace: "testing",
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "myapp",
				},
			},
		},
	}

	pods := &corev1.PodList{
		Items: []corev1.Pod{
			*createPodWithPhase("web-0", labels, corev1.PodRunning, metav1.NewControllerRef(sts, sts.GroupVersionKind())),
			*createPodWithPhase("web-1", labels, corev1.PodRunning, metav1.NewControllerRef(sts, sts.GroupVersionKind())),
			*createPodWithPhase("web-2", labels, corev1.PodPending, metav1.NewControllerRef(sts, sts.GroupVersionKind())),
			*createPodWithPhase("random-pod", nil, corev1.PodRunning, nil),
		},
	}

	podList := &unstructured.UnstructuredList{}
	for _, p := range pods.Items {
		podList.Items = append(podList.Items, *testutil.ToUnstructured(t, &p))
	}

	key := store.Key{
		Namespace:  "testing",
		APIVersion: "v1",
		Kind:       "Pod",
	}

	tpo.objectStore.EXPECT().List(gomock.Any(), gomock.Eq(key)).Return(podList, false, nil)

	ctx := context.Background()

	stsc := NewStatefulSetStatus(ctx, sts, printOptions)
	got, err := stsc.Create()
	require.NoError(t, err)

	expected := component.NewQuadrant("Status")
	require.NoError(t, expected.Set(component.QuadNW, "Running", "2"))
	require.NoError(t, expected.Set(component.QuadNE, "Waiting", "1"))
	require.NoError(t, expected.Set(component.QuadSW, "Succeeded", "0"))
	require.NoError(t, expected.Set(component.QuadSE, "Failed", "0"))

	assert.Equal(t, expected, got)
}

func Test_StatefulSetConfiguration(t *testing.T) {
	now := testutil.Time()
	validStatefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "web",
			CreationTimestamp: metav1.Time{
				Time: now,
			},
			Labels: map[string]string{
				"foo": "bar",
			},
		},
		Status: appsv1.StatefulSetStatus{
			Replicas: 1,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "myapp",
				},
			},
			Replicas: conversion.PtrInt32(3),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "k8s.gcr.io/nginx-slim:0.8",
						},
					},
				},
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			PodManagementPolicy: appsv1.OrderedReadyPodManagement,
		},
	}

	cases := []struct {
		name        string
		statefulSet *appsv1.StatefulSet
		isErr       bool
		expected    *component.Summary
	}{
		{
			name:        "default",
			statefulSet: validStatefulSet,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Update Strategy",
					Content: component.NewText("RollingUpdate"),
				},
				{
					Header:  "Selectors",
					Content: component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
				},
				{
					Header:  "Replicas",
					Content: component.NewText("3 Desired / 1 Total"),
				},
				{
					Header:  "Pod Management Policy",
					Content: component.NewText("OrderedReady"),
				},
			}...),
		},
		{
			name:        "statefulset is nil",
			statefulSet: nil,
			isErr:       true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			sc := NewStatefulSetConfiguration(tc.statefulSet)

			summary, err := sc.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_StatefulSetPods(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	nodeLink := component.NewLink("", "node", "/node")
	tpo.link.EXPECT().
		ForGVK("", "v1", "Node", "node", "node").
		Return(nodeLink, nil).AnyTimes()

	ctx := context.Background()

	now := testutil.Time()

	labels := map[string]string{
		"app": "testing",
	}

	statefulSet := testutil.CreateStatefulSet("statefulset")
	statefulSet.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: labels,
	}

	pod := testutil.CreatePod("web-0")
	pod.SetOwnerReferences(testutil.ToOwnerReferences(t, statefulSet))
	pod.CreationTimestamp = metav1.Time{Time: now}
	pod.Labels = labels
	pod.Spec.Containers = []corev1.Container{
		{
			Name:  "nginx",
			Image: "nginx:1.15",
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
			},
		},
	}

	pods := &corev1.PodList{
		Items: []corev1.Pod{*pod},
	}

	tpo.PathForObject(pod, pod.Name, "/pod")

	podList := &unstructured.UnstructuredList{}
	for _, p := range pods.Items {
		podList.Items = append(podList.Items, *testutil.ToUnstructured(t, &p))
	}
	key := store.Key{
		Namespace:  "namespace",
		APIVersion: "v1",
		Kind:       "Pod",
	}

	tpo.objectStore.EXPECT().List(gomock.Any(), gomock.Eq(key)).Return(podList, false, nil)

	printOptions := tpo.ToOptions()
	printOptions.DisableLabels = false

	got, err := createPodListView(ctx, statefulSet, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Ready", "Phase", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", "We couldn't find any pods!", cols)
	expected.Add(component.TableRow{
		"Name":     component.NewLink("", "web-0", "/pod"),
		"Ready":    component.NewText("1/1"),
		"Phase":    component.NewText("Pending"),
		"Restarts": component.NewText("0"),
		"Node":     nodeLink,
		"Age":      component.NewTimestamp(now),
	})

	component.AssertEqual(t, expected, got)
}

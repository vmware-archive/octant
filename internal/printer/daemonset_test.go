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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_DaemonSetListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	object := testutil.CreateDaemonSet("ds")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Labels = labels

	tpo.PathForObject(object, object.Name, "/path")

	list := &appsv1.DaemonSetList{
		Items: []appsv1.DaemonSet{*object},
	}

	ctx := context.Background()
	got, err := DaemonSetListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Desired", "Current", "Ready",
		"Up-To-Date", "Age", "Node Selector")
	expected := component.NewTable("Daemon Sets", "We couldn't find any daemon sets!", cols)
	expected.Add(component.TableRow{
		"Name":          component.NewLink("", object.Name, "/path"),
		"Labels":        component.NewLabels(labels),
		"Age":           component.NewTimestamp(now),
		"Desired":       component.NewText("1"),
		"Current":       component.NewText("1"),
		"Ready":         component.NewText("1"),
		"Up-To-Date":    component.NewText("1"),
		"Node Selector": component.NewSelectors(nil),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, object),
		}),
	})

	component.AssertEqual(t, expected, got)
}

func Test_DaemonSetConfiguration(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	ds := testutil.CreateDaemonSet("ds")
	ds.CreationTimestamp = metav1.Time{Time: now}
	ds.Labels = labels
	ds.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: labels,
	}
	ds.Spec.Template.Spec.NodeSelector = labels

	cases := []struct {
		name      string
		daemonSet *appsv1.DaemonSet
		isErr     bool
		expected  *component.Summary
	}{
		{
			name:      "daemonset",
			daemonSet: ds,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Update Strategy",
					Content: component.NewText("Max Unavailable 1"),
				},
				{
					Header:  "Revision History Limit",
					Content: component.NewText("10"),
				},
				{
					Header:  "Selectors",
					Content: printSelectorMap(labels),
				},
				{
					Header:  "Node Selectors",
					Content: printSelectorMap(labels),
				},
			}...),
		},
		{
			name:      "daemonset is nil",
			daemonSet: nil,
			isErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dc := NewDaemonSetConfiguration(tc.daemonSet)

			summary, err := dc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createDaemonSetSummaryStatus(t *testing.T) {
	ds := testutil.CreateDaemonSet("ds")

	got, err := createDaemonSetSummaryStatus(ds)
	require.NoError(t, err)

	sections := component.SummarySections{
		{Header: "Current Number Scheduled", Content: component.NewText("1")},
		{Header: "Desired Number Scheduled", Content: component.NewText("1")},
		{Header: "Number Available", Content: component.NewText("1")},
		{Header: "Number Mis-scheduled", Content: component.NewText("0")},
		{Header: "Number Ready", Content: component.NewText("1")},
		{Header: "Updated Number Scheduled", Content: component.NewText("1")},
	}
	expected := component.NewSummary("Status", sections...)

	assert.Equal(t, expected, got)
}

func Test_DaemonSetPods(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	ctx := context.Background()

	now := testutil.Time()

	nodeLink := component.NewLink("", "node", "/node")
	tpo.link.EXPECT().
		ForGVK("", "v1", "Node", "node", "node").
		Return(nodeLink, nil).AnyTimes()

	daemonSet := testutil.CreateDaemonSet("daemonset")

	pod := testutil.CreatePod("fluentd-elasticsearch-dvskv")
	pod.SetOwnerReferences(testutil.ToOwnerReferences(t, daemonSet))
	pod.CreationTimestamp = metav1.Time{Time: now}
	pod.Spec.Containers = []corev1.Container{
		{
			Name:  "fluentd-elasticsearch",
			Image: "fluentd:1.7",
		},
	}
	pod.Spec.NodeName = "node"
	pod.Status = corev1.PodStatus{
		Phase: "Pending",
		ContainerStatuses: []corev1.ContainerStatus{
			{
				Name:         "fluentd-elasticsearch",
				Image:        "fluentd:1.7",
				RestartCount: 0,
				Ready:        false,
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

	got, err := createPodListView(ctx, daemonSet, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Ready", "Phase", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", "We couldn't find any pods!", cols)
	expected.Add(component.TableRow{
		"Name":     component.NewLink("", "fluentd-elasticsearch-dvskv", "/pod"),
		"Ready":    component.NewText("0/1"),
		"Phase":    component.NewText("Pending"),
		"Restarts": component.NewText("0"),
		"Node":     nodeLink,
		"Age":      component.NewTimestamp(now),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, pod),
		}),
	})
	addPodTableFilters(expected)

	component.AssertEqual(t, expected, got)
}

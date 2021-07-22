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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_ReplicationControllerListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	validReplicationControllerLabels := map[string]string{
		"foo": "bar",
	}

	validReplicationControllerCreationTime := testutil.Time()

	validReplicationController := &corev1.ReplicationController{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ReplicationController",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rc-test",
			Namespace: "default",
			CreationTimestamp: metav1.Time{
				Time: validReplicationControllerCreationTime,
			},
			Labels: validReplicationControllerLabels,
		},
		Status: corev1.ReplicationControllerStatus{
			Replicas:          3,
			AvailableReplicas: 0,
		},
		Spec: corev1.ReplicationControllerSpec{
			Selector: map[string]string{
				"app": "myapp",
			},
			Template: &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.15",
						},
					},
				},
			},
		},
	}

	validReplicationControllerList := &corev1.ReplicationControllerList{
		Items: []corev1.ReplicationController{
			*validReplicationController,
		},
	}

	tpo.PathForObject(validReplicationController, validReplicationController.Name, "/rc")

	ctx := context.Background()
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, validReplicationController)
	got, err := ReplicationControllerListHandler(ctx, validReplicationControllerList, printOptions)
	require.NoError(t, err)

	containers := component.NewContainers()
	containers.Add("nginx", "nginx:1.15")

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	expected := component.NewTable("ReplicationControllers", "We couldn't find any replication controllers!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "rc-test", "/rc",
			genObjectStatus(component.TextStatusWarning, []string{
				"Replication Controller pods are not ready",
			})),
		"Labels":     component.NewLabels(validReplicationControllerLabels),
		"Status":     component.NewText("0/3"),
		"Age":        component.NewTimestamp(validReplicationControllerCreationTime),
		"Containers": containers,
		"Selector":   component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, validReplicationController),
		}),
	})

	component.AssertEqual(t, expected, got)
}

func Test_ReplicationControllerConfiguration(t *testing.T) {
	var replicas int32 = 3

	rc := testutil.CreateReplicationController("rc")
	rc.Spec.Replicas = &replicas
	rc.Status = corev1.ReplicationControllerStatus{
		ReadyReplicas: 3,
		Replicas:      3,
	}

	cases := []struct {
		name                  string
		replicationController *corev1.ReplicationController
		isErr                 bool
		expected              *component.Summary
	}{
		{
			name:                  "replicationcontroller",
			replicationController: rc,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Replica Status",
					Content: component.NewText("Current 3 / Desired 3"),
				},
				{
					Header:  "Replicas",
					Content: component.NewText("3"),
				},
			}...),
		},
		{
			name:                  "replicationcontroller is nil",
			replicationController: nil,
			isErr:                 true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			rcc := NewReplicationControllerConfiguration(tc.replicationController)

			summary, err := rcc.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func TestReplicationControllerStatus(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	replicationController := testutil.CreateReplicationController("rc")
	replicationController.Labels = map[string]string{
		"foo": "bar",
	}
	replicationController.Spec.Selector = map[string]string{
		"foo": "bar",
	}
	replicationController.Namespace = "testing"

	p1 := *createPodWithPhase(
		"nginx-g7f72",
		replicationController.Labels,
		corev1.PodRunning,
		metav1.NewControllerRef(replicationController, replicationController.GroupVersionKind()))

	p2 := *createPodWithPhase(
		"nginx-p64jr",
		replicationController.Labels,
		corev1.PodRunning,
		metav1.NewControllerRef(replicationController, replicationController.GroupVersionKind()))

	p3 := *createPodWithPhase(
		"nginx-x8nrk",
		replicationController.Labels,
		corev1.PodRunning,
		metav1.NewControllerRef(replicationController, replicationController.GroupVersionKind()))

	pods := &corev1.PodList{
		Items: []corev1.Pod{p1, p2, p3},
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
	rcs := NewReplicationControllerStatus(ctx, replicationController, printOptions)
	got, err := rcs.Create()
	require.NoError(t, err)

	expected := component.NewQuadrant("Status")
	require.NoError(t, expected.Set(component.QuadNW, "Running", "3"))
	require.NoError(t, expected.Set(component.QuadNE, "Waiting", "0"))
	require.NoError(t, expected.Set(component.QuadSW, "Succeeded", "0"))
	require.NoError(t, expected.Set(component.QuadSE, "Failed", "0"))

	assert.Equal(t, expected, got)
}

func Test_ReplicationControllerPods(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	ctx := context.Background()

	now := testutil.Time()

	nodeLink := component.NewLink("", "node", "/node")
	tpo.link.EXPECT().
		ForGVK("", "v1", "Node", "node", "node").
		Return(nodeLink, nil).AnyTimes()

	rc := testutil.CreateReplicationController("replicationcontroller")

	pod := testutil.CreatePod("nginx-hv4qs")
	pod.SetOwnerReferences(testutil.ToOwnerReferences(t, rc))
	pod.CreationTimestamp = metav1.Time{Time: now}
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
				Ready:        false,
			},
		},
	}

	pods := &corev1.PodList{
		Items: []corev1.Pod{*pod},
	}

	tpo.PathForObject(pod, pod.Name, "/pod")
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, pod)

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

	got, err := createPodListView(ctx, rc, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Ready", "Phase", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", "We couldn't find any pods!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "nginx-hv4qs", "/pod",
			genObjectStatus(component.TextStatusWarning, []string{"Pod may require additional action"})),
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

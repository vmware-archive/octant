/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/store"
	storefake "github.com/vmware/octant/pkg/store/fake"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_ReplicaSetListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	object := &appsv1.ReplicaSetList{
		Items: []appsv1.ReplicaSet{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps/v1",
					Kind:       "ReplicaSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "replicaset-test",
					Namespace: "default",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
					Labels: labels,
				},
				Status: appsv1.ReplicaSetStatus{
					Replicas:          3,
					AvailableReplicas: 2,
				},
				Spec: appsv1.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "myapp",
						},
					},
					Template: corev1.PodTemplateSpec{
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
						},
					},
				},
			},
		},
	}

	tpo.PathForObject(&object.Items[0], object.Items[0].Name, "/replica-set")

	ctx := context.Background()
	got, err := ReplicaSetListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	containers := component.NewContainers()
	containers.Add("nginx", "nginx:1.15")
	containers.Add("kuard", "gcr.io/kuar-demo/kuard-amd64:1")

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	expected := component.NewTable("ReplicaSets", "We couldn't find any replica sets!", cols)
	expected.Add(component.TableRow{
		"Name":       component.NewLink("", "replicaset-test", "/replica-set"),
		"Labels":     component.NewLabels(labels),
		"Age":        component.NewTimestamp(now),
		"Selector":   component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
		"Status":     component.NewText("2/3"),
		"Containers": containers,
	})

	component.AssertEqual(t, expected, got)
}

func TestReplicaSetConfiguration(t *testing.T) {

	var replicas int32 = 3
	isController := true

	rs := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rs-frontend",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Controller: &isController,
					Name:       "replicaset-controller",
					Kind:       "ReplicationController",
				},
			},
		},
		Spec: appsv1.ReplicaSetSpec{
			Replicas: &replicas,
		},
		Status: appsv1.ReplicaSetStatus{
			ReadyReplicas: 3,
			Replicas:      3,
		},
	}

	cases := []struct {
		name       string
		replicaset *appsv1.ReplicaSet
		isErr      bool
		expected   *component.Summary
	}{
		{
			name:       "replicaset",
			replicaset: rs,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Controlled By",
					Content: component.NewLink("", "replicaset-controller", "/owner"),
				},
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
			name:       "replicaset is nil",
			replicaset: nil,
			isErr:      true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			rc := NewReplicaSetConfiguration(tc.replicaset)

			if tc.replicaset != nil && len(tc.replicaset.OwnerReferences) > 0 {
				tpo.PathForOwner(tc.replicaset, &tc.replicaset.OwnerReferences[0], "/owner")
			}

			summary, err := rc.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func TestReplicaSetStatus(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockStore(controller)

	labels := map[string]string{
		"app": "myapp",
	}

	rs := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rs-frontend",
			Namespace: "testing",
		},
		Spec: appsv1.ReplicaSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "myapp",
				},
			},
		},
	}

	pods := &corev1.PodList{
		Items: []corev1.Pod{
			*createPodWithPhase("frontend-l82ph", labels, corev1.PodRunning, metav1.NewControllerRef(rs, rs.GroupVersionKind())),
			*createPodWithPhase("frontend-rs95v", labels, corev1.PodRunning, metav1.NewControllerRef(rs, rs.GroupVersionKind())),
			*createPodWithPhase("frontend-sl8sv", labels, corev1.PodRunning, metav1.NewControllerRef(rs, rs.GroupVersionKind())),
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

	o.EXPECT().List(gomock.Any(), gomock.Eq(key)).Return(podList, false, nil).AnyTimes()

	ctx := context.Background()
	rsc := NewReplicaSetStatus(rs)
	got, err := rsc.Create(ctx, o)
	require.NoError(t, err)

	filteredPods, err := listPods(ctx, rs.Namespace, rs.Spec.Selector, rs.UID, o)
	require.NoError(t, err)
	assert.Equal(t, 3, len(filteredPods))

	expected := component.NewQuadrant("Status")
	require.NoError(t, expected.Set(component.QuadNW, "Running", "3"))
	require.NoError(t, expected.Set(component.QuadNE, "Waiting", "0"))
	require.NoError(t, expected.Set(component.QuadSW, "Succeeded", "0"))
	require.NoError(t, expected.Set(component.QuadSE, "Failed", "0"))

	assert.Equal(t, expected, got)
}

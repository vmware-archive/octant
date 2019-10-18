/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_PersistentVolumeListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()
	object := testutil.CreatePersistentVolumeClaim("pvc")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Labels = labels

	tpo.PathForObject(object, object.Name, "/pvc")

	list := &corev1.PersistentVolumeClaimList{
		Items: []corev1.PersistentVolumeClaim{*object},
	}

	ctx := context.Background()
	got, err := PersistentVolumeClaimListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Status", "Volume", "Capacity", "Access Modes",
		"Storage Class", "Age")
	expected := component.NewTable("Persistent Volume Claims", "We couldn't find any persistent volume claims!", cols)
	expected.Add(component.TableRow{
		"Name":          component.NewLink("", object.Name, "/pvc"),
		"Status":        component.NewText("Bound"),
		"Volume":        component.NewText("task-pv-volume"),
		"Capacity":      component.NewText("10Gi"),
		"Access Modes":  component.NewText("RWO"),
		"Storage Class": component.NewText("manual"),
		"Age":           component.NewTimestamp(now),
	})

	component.AssertEqual(t, expected, got)
}

func Test_PersistentVolumeClaimConfiguration(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	pvc := testutil.CreatePersistentVolumeClaim("pvc")
	pvc.CreationTimestamp = metav1.Time{Time: now}
	pvc.Finalizers = []string{"kubernetes.io/pvc-protection"}
	pvc.Labels = labels
	pvc.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: labels,
	}

	cases := []struct {
		name                  string
		persistentVolumeClaim *corev1.PersistentVolumeClaim
		isErr                 bool
		expected              *component.Summary
	}{
		{
			name:                  "general",
			persistentVolumeClaim: pvc,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Volume Mode",
					Content: component.NewText("Filesystem"),
				},
				{
					Header:  "Access Modes",
					Content: component.NewText("RWO"),
				},
				{
					Header:  "Finalizers",
					Content: component.NewText("[kubernetes.io/pvc-protection]"),
				},
				{
					Header:  "Storage Class Name",
					Content: component.NewText("manual"),
				},
				{
					Header:  "Labels",
					Content: component.NewLabels(labels),
				},
				{
					Header:  "Selectors",
					Content: printSelectorMap(labels),
				},
			}...),
		},
		{
			name:                  "pvc is nil",
			persistentVolumeClaim: nil,
			isErr:                 true,
		},
	}

	for _, tc := range cases {
		controller := gomock.NewController(t)
		defer controller.Finish()

		tpo := newTestPrinterOptions(controller)
		printOptions := tpo.ToOptions()

		pc := NewPersistentVolumeClaimConfiguration(tc.persistentVolumeClaim)

		summary, err := pc.Create(printOptions)
		if tc.isErr {
			require.Error(t, err)
			return
		}
		require.NoError(t, err)

		component.AssertEqual(t, tc.expected, summary)
	}
}

func Test_createPersistentVolumeClaimStatusView(t *testing.T) {
	object := testutil.CreatePersistentVolumeClaim("pvc")

	got, err := createPersistentVolumeClaimStatusView(object)
	require.NoError(t, err)

	sections := component.SummarySections{}
	sections.AddText("Claim Status", "Bound")
	sections.AddText("Storage Requested", "3Gi")
	sections.AddText("Bound Volume", "task-pv-volume")
	sections.AddText("Total Volume Capacity", "10Gi")
	expected := component.NewSummary("Status", sections...)

	component.AssertEqual(t, expected, got)
}

func Test_PersistentVolumeClaimMountedPodsList(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	ctx := context.Background()

	now := testutil.Time()

	nodeLink := component.NewLink("", "node", "/node")
	tpo.link.EXPECT().
		ForGVK("", "v1", "Node", "node", "node").
		Return(nodeLink, nil).AnyTimes()

	pvc := testutil.CreatePersistentVolumeClaim("mysql-pv-claim")

	pod := testutil.CreatePod("wordpress-mysql-67565bd57-8fzbh")
	pod.CreationTimestamp = metav1.Time{Time: now}
	pod.Spec.Containers = []corev1.Container{
		{
			Name:  "mysql",
			Image: "mysql:5.6",
		},
	}
	pod.Spec.NodeName = "node"
	pod.Status = corev1.PodStatus{
		Phase: "Running",
		ContainerStatuses: []corev1.ContainerStatus{
			{
				Name:         "mysql",
				Image:        "mysql:5.6",
				RestartCount: 0,
				Ready:        true,
			},
		},
	}
	pod.Spec.Volumes = []corev1.Volume{
		{
			Name: "mysql-persistent-storage",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "mysql-pv-claim",
				},
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

	got, err := createMountedPodListView(ctx, pvc.Namespace, pvc.Name, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Ready", "Phase", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", "We couldn't find any pods!", cols)
	expected.Add(component.TableRow{
		"Name":     component.NewLink("", "wordpress-mysql-67565bd57-8fzbh", "/pod"),
		"Ready":    component.NewText("1/1"),
		"Phase":    component.NewText("Running"),
		"Restarts": component.NewText("0"),
		"Node":     nodeLink,
		"Age":      component.NewTimestamp(now),
	})
	addPodTableFilters(expected)

	component.AssertEqual(t, expected, got)
}

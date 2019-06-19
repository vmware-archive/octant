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
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_PersistentVolumeListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

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
	expected := component.NewTable("Persistent Volume Claims", cols)
	expected.Add(component.TableRow{
		"Name":          component.NewLink("", object.Name, "/pvc"),
		"Status":        component.NewText("Bound"),
		"Volume":        component.NewText("task-pv-volume"),
		"Capacity":      component.NewText("10Gi"),
		"Access Modes":  component.NewText("RWO"),
		"Storage Class": component.NewText("manual"),
		"Age":           component.NewTimestamp(now),
	})

	assert.Equal(t, expected, got)
}

func Test_printPersistentVolumeClaimConfig(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := testutil.CreatePersistentVolumeClaim("pvc")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Finalizers = []string{"kubernetes.io/pvc-protection"}
	object.Labels = labels
	object.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: labels,
	}

	got, err := printPersistentVolumeClaimConfig(object)
	require.NoError(t, err)

	sections := component.SummarySections{}
	sections.AddText("Volume Mode", "Filesystem")
	sections.AddText("Access Modes", "RWO")
	sections.AddText("Finalizers", "[kubernetes.io/pvc-protection]")
	sections.AddText("Storage Class Name", "manual")
	sections.Add("Labels", component.NewLabels(labels))
	sections.Add("Selectors", printSelectorMap(labels))
	expected := component.NewSummary("Configuration", sections...)

	assert.Equal(t, expected, got)
}

func Test_printPersistentVolumeClaimStatus(t *testing.T) {
	object := testutil.CreatePersistentVolumeClaim("pvc")

	got, err := printPersistentVolumeClaimStatus(object)
	require.NoError(t, err)

	sections := component.SummarySections{}
	sections.AddText("Claim Status", "Bound")
	sections.AddText("Storage Requested", "3Gi")
	sections.AddText("Bound Volume", "task-pv-volume")
	sections.AddText("Total Volume Capacity", "10Gi")
	expected := component.NewSummary("Status", sections...)

	assert.Equal(t, expected, got)
}

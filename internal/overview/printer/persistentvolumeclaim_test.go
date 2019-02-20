package printer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/heptio/developer-dash/internal/overview/link"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_PersistentVolumeListHandler(t *testing.T) {
	printOptions := Options{
		Cache: cache.NewMemoryCache(),
	}

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := testutil.CreatePersistentVolumeClaim("pvc")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Labels = labels

	list := &corev1.PersistentVolumeClaimList{
		Items: []corev1.PersistentVolumeClaim{*object},
	}

	got, err := PersistentVolumeClaimListHandler(list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Status", "Volume", "Capacity", "Access Modes",
		"Storage Class", "Age")
	expected := component.NewTable("Persistent Volume Claims", cols)
	expected.Add(component.TableRow{
		"Name":          link.ForObject(object, object.Name),
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

	var sections component.SummarySections
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

	var sections component.SummarySections
	sections.AddText("Claim Status", "Bound")
	sections.AddText("Storage Requested", "3Gi")
	sections.AddText("Bound Volume", "task-pv-volume")
	sections.AddText("Total Volume Capacity", "10Gi")
	expected := component.NewSummary("Status", sections...)

	assert.Equal(t, expected, got)
}

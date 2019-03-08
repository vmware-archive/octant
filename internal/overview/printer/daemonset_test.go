package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_DaemonSetListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		Cache: cachefake.NewMockCache(controller),
	}

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := testutil.CreateDaemonSet("ds")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Labels = labels

	list := &appsv1.DaemonSetList{
		Items: []appsv1.DaemonSet{*object},
	}

	ctx := context.Background()
	got, err := DaemonSetListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Desired", "Current", "Ready",
		"Up-To-Date", "Age", "Node Selector")
	expected := component.NewTable("Daemon Sets", cols)
	expected.Add(component.TableRow{
		"Name":          link.ForObject(object, object.Name),
		"Labels":        component.NewLabels(labels),
		"Age":           component.NewTimestamp(now),
		"Desired":       component.NewText("1"),
		"Current":       component.NewText("1"),
		"Ready":         component.NewText("1"),
		"Up-To-Date":    component.NewText("1"),
		"Node Selector": component.NewSelectors(nil),
	})

	assert.Equal(t, expected, got)
}

func Test_printDaemonSetConfig(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := testutil.CreateDaemonSet("ds")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Labels = labels
	object.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: labels,
	}
	object.Spec.Template.Spec.NodeSelector = labels

	got, err := printDaemonSetConfig(object)
	require.NoError(t, err)

	var sections component.SummarySections
	sections.AddText("Update Strategy", "Max Unavailable 1")
	sections.AddText("Revision History Limit", "10")
	sections.Add("Selectors", printSelectorMap(labels))
	sections.Add("Node Selectors", printSelectorMap(labels))
	expected := component.NewSummary("Configuration", sections...)

	assert.Equal(t, expected, got)
}

func Test_printDaemonSetSummary(t *testing.T) {

	object := testutil.CreateDaemonSet("ds")

	got, err := printDaemonSetStatus(object)
	require.NoError(t, err)

	var sections component.SummarySections
	sections.AddText("Current Number Scheduled", "1")
	sections.AddText("Desired Number Scheduled", "1")
	sections.AddText("Number Available", "1")
	sections.AddText("Number Mis-scheduled", "0")
	sections.AddText("Number Ready", "1")
	sections.AddText("Updated Number Scheduled", "1")
	expected := component.NewSummary("Status", sections...)

	assert.Equal(t, expected, got)
}

package view

import (
	"context"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/core"
)

func TestPodCondition_InvalidObject(t *testing.T) {
	pc := NewPodCondition()
	ctx := context.Background()

	object := &unstructured.Unstructured{}

	_, err := pc.Content(ctx, object, nil)
	require.Error(t, err)
}

func TestPodCondition(t *testing.T) {
	pc := NewPodCondition()

	lastProbeTime := metav1.Time{
		Time: time.Unix(1539603521, 0),
	}

	lastTransitionTime := metav1.Time{
		Time: time.Unix(1539603521, 0),
	}

	ctx := context.Background()
	object := &core.Pod{
		Status: core.PodStatus{
			Conditions: []core.PodCondition{
				{
					Type:               core.PodScheduled,
					Status:             core.ConditionTrue,
					LastProbeTime:      lastProbeTime,
					LastTransitionTime: lastTransitionTime,
					Reason:             "reason",
					Message:            "message",
				},
			},
		},
	}

	contents, err := pc.Content(ctx, object, nil)
	require.NoError(t, err)

	require.Len(t, contents, 1)

	table, ok := contents[0].(*content.Table)
	require.True(t, ok)
	require.Len(t, table.Rows, 1)

	expectedColumns := []string{"Type", "Status", "Last probe time",
		"Last transition time", "Reason", "Message"}
	assert.Equal(t, expectedColumns, table.ColumnNames())

	expectedRow := content.TableRow{
		"Type":                 content.NewStringText("PodScheduled"),
		"Status":               content.NewStringText("True"),
		"Last probe time":      content.NewStringText("2018-10-15T11:38:41Z"),
		"Last transition time": content.NewStringText("2018-10-15T11:38:41Z"),
		"Reason":               content.NewStringText("reason"),
		"Message":              content.NewStringText("message"),
	}
	assert.Equal(t, expectedRow, table.Rows[0])
}

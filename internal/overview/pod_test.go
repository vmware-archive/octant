package overview

import (
	"context"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/content"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func TestPodList_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewPodList("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestPodList(t *testing.T) {
	pl := NewPodList("prefix", "ns", clock.NewFakeClock(time.Now()))

	cache := NewMemoryCache()

	rs := loadFromFile(t, "replicaset-1.yaml")
	rs = convertToInternal(t, rs)

	storeFromFile(t, "rs-pod-1.yaml", cache)

	ctx := context.Background()

	contents, err := pl.Content(ctx, rs, cache)
	require.NoError(t, err)

	podColumns := []content.TableColumn{
		tableCol("Name"),
		tableCol("Ready"),
		tableCol("Status"),
		tableCol("Restarts"),
		tableCol("Age"),
		tableCol("IP"),
		tableCol("Node"),
		tableCol("Nominated Node"),
		tableCol("Labels"),
	}

	listTable := content.NewTable("Pods", "No pods were found")
	listTable.Columns = podColumns
	listTable.AddRow(
		content.TableRow{
			"Name":           content.NewLinkText("rs1-s8mj8", "/content/overview/workloads/pods/rs1-s8mj8"),
			"Ready":          content.NewStringText("1/1"),
			"Status":         content.NewStringText("Running"),
			"Restarts":       content.NewStringText("0"),
			"Age":            content.NewStringText("1d"),
			"IP":             content.NewStringText("10.1.114.100"),
			"Node":           content.NewStringText("node1"),
			"Nominated Node": content.NewStringText("<none>"),
			"Labels":         content.NewStringText("app=myapp,pod-template-hash=2350241137"),
		},
	)

	expected := []content.Content{
		&listTable,
	}

	assert.Equal(t, expected, contents)
}

func TestPodCondition_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewPodCondition("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestPodCondition(t *testing.T) {
	pc := NewPodCondition("prefix", "ns", clock.NewFakeClock(time.Now()))

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
		"Last probe time":      content.NewTimeText("2018-10-15T11:38:41Z"),
		"Last transition time": content.NewTimeText("2018-10-15T11:38:41Z"),
		"Reason":               content.NewStringText("reason"),
		"Message":              content.NewStringText("message"),
	}
	assert.Equal(t, expectedRow, table.Rows[0])
}

func Test_createPodStatus(t *testing.T) {
	pods := []*core.Pod{
		{Status: core.PodStatus{Phase: core.PodRunning}},
		{Status: core.PodStatus{Phase: core.PodPending}},
		{Status: core.PodStatus{Phase: core.PodSucceeded}},
		{Status: core.PodStatus{Phase: core.PodFailed}},
	}

	ps := createPodStatus(pods)

	expected := podStatus{
		Running:   1,
		Waiting:   1,
		Succeeded: 1,
		Failed:    1,
	}

	assert.Equal(t, expected, ps)
}

func Test_getSelector(t *testing.T) {
	labels := map[string]string{
		"app": "app",
	}

	cases := []struct {
		name     string
		object   runtime.Object
		expected *metav1.LabelSelector
		isErr    bool
	}{
		{
			name: "daemon set",
			object: &extensions.DaemonSet{
				Spec: extensions.DaemonSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
				},
			},
			expected: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
		{
			name: "stateful set",
			object: &apps.StatefulSet{
				Spec: apps.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
				},
			},
			expected: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
		{
			name:     "job",
			object:   &batch.Job{},
			expected: nil,
		},
		{
			name: "replication controller",
			object: &core.ReplicationController{
				Spec: core.ReplicationControllerSpec{
					Selector: labels,
				},
			},
			expected: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
		{
			name: "replica set",
			object: &extensions.ReplicaSet{
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
				},
			},
			expected: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
		{
			name:     "empty replica set",
			object:   &extensions.ReplicaSet{},
			expected: nil,
		},
		{
			name: "deployment",
			object: &extensions.Deployment{
				Spec: extensions.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
				},
			},
			expected: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
		{
			name:     "empty deployment",
			object:   &extensions.Deployment{},
			isErr:    false,
			expected: nil,
		},
		{
			name: "service",
			object: &core.Service{
				Spec: core.ServiceSpec{
					Selector: labels,
				},
			},
			expected: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			selector, err := getSelector(tc.object)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			assert.Equal(t, tc.expected, selector)
		})
	}
}

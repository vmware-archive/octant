package overview

import (
	"context"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view"

	"github.com/heptio/developer-dash/internal/content"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestPodList_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewPodList("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestPodList(t *testing.T) {
	pl := NewPodList("prefix", "ns", clock.NewFakeClock(time.Now()))

	c := cache.NewMemoryCache()

	rs := loadFromFile(t, "replicaset-1.yaml")

	storeFromFile(t, "rs-pod-1.yaml", c)

	ctx := context.Background()

	contents, err := pl.Content(ctx, rs, c)
	require.NoError(t, err)

	podColumns := []content.TableColumn{
		view.TableCol("Name"),
		view.TableCol("Ready"),
		view.TableCol("Status"),
		view.TableCol("Restarts"),
		view.TableCol("Age"),
		view.TableCol("IP"),
		view.TableCol("Node"),
		view.TableCol("Nominated Node"),
		view.TableCol("Readiness Gates"),
		view.TableCol("Labels"),
	}

	listTable := content.NewTable("Pods", "No pods were found")
	listTable.Columns = podColumns
	listTable.AddRow(
		content.TableRow{
			"Name":            content.NewLinkText("rs1-s8mj8", "/content/overview/workloads/pods/rs1-s8mj8"),
			"Ready":           content.NewStringText("1/1"),
			"Status":          content.NewStringText("Running"),
			"Restarts":        content.NewStringText("0"),
			"Age":             content.NewStringText("24h"),
			"IP":              content.NewStringText("10.1.114.100"),
			"Node":            content.NewStringText("node1"),
			"Nominated Node":  content.NewStringText("<none>"),
			"Readiness Gates": content.NewStringText("<none>"),
			"Labels":          content.NewStringText("app=myapp,pod-template-hash=2350241137"),
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
	object := &corev1.Pod{
		Status: corev1.PodStatus{
			Conditions: []corev1.PodCondition{
				{
					Type:               corev1.PodScheduled,
					Status:             corev1.ConditionTrue,
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
	pods := []*corev1.Pod{
		{Status: corev1.PodStatus{Phase: corev1.PodRunning}},
		{Status: corev1.PodStatus{Phase: corev1.PodPending}},
		{Status: corev1.PodStatus{Phase: corev1.PodSucceeded}},
		{Status: corev1.PodStatus{Phase: corev1.PodFailed}},
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
			object: &appsv1.DaemonSet{
				Spec: appsv1.DaemonSetSpec{
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
			object: &appsv1.StatefulSet{
				Spec: appsv1.StatefulSetSpec{
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
			object:   &batchv1beta1.CronJob{},
			expected: nil,
		},
		{
			name: "replication controller",
			object: &corev1.ReplicationController{
				Spec: corev1.ReplicationControllerSpec{
					Selector: labels,
				},
			},
			expected: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
		{
			name: "replica set",
			object: &appsv1.ReplicaSet{
				Spec: appsv1.ReplicaSetSpec{
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
			object:   &appsv1.ReplicaSet{},
			expected: nil,
		},
		{
			name: "deployment",
			object: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
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
			object:   &appsv1.Deployment{},
			isErr:    false,
			expected: nil,
		},
		{
			name: "service",
			object: &corev1.Service{
				Spec: corev1.ServiceSpec{
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

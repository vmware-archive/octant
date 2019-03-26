package printer

import (
	"context"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	cacheutil "github.com/heptio/developer-dash/internal/cache/util"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_ReplicaSetListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		Cache: cachefake.NewMockCache(controller),
	}

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

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
								corev1.Container{
									Name:  "nginx",
									Image: "nginx:1.15",
								},
								corev1.Container{
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

	ctx := context.Background()
	got, err := ReplicaSetListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	containers := component.NewContainers()
	containers.Add("nginx", "nginx:1.15")
	containers.Add("kuard", "gcr.io/kuar-demo/kuard-amd64:1")

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	expected := component.NewTable("ReplicaSets", cols)
	expected.Add(component.TableRow{
		"Name":       component.NewLink("", "replicaset-test", "/content/overview/namespace/default/workloads/replica-sets/replicaset-test"),
		"Labels":     component.NewLabels(labels),
		"Age":        component.NewTimestamp(now),
		"Selector":   component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
		"Status":     component.NewText("2/3"),
		"Containers": containers,
	})

	assert.Equal(t, expected, got)
}

func TestReplicaSetConfiguration(t *testing.T) {

	var replicas int32 = 3
	isController := true

	rs := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rs-frontend",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				metav1.OwnerReference{
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
					Content: component.NewLink("", "replicaset-controller", "/content/overview/namespace/default/workloads/replication-controllers/replicaset-controller"),
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
			rc := NewReplicaSetConfiguration(tc.replicaset)

			summary, err := rc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, summary)
		})
	}
}

func TestReplicaSetStatus(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	c := cachefake.NewMockCache(controller)

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

	var podList []*unstructured.Unstructured
	for _, p := range pods.Items {
		u := testutil.ToUnstructured(t, &p)
		podList = append(podList, u)
	}
	key := cacheutil.Key{
		Namespace:  "testing",
		APIVersion: "v1",
		Kind:       "Pod",
	}

	c.EXPECT().List(gomock.Any(), gomock.Eq(key)).Return(podList, nil)

	ctx := context.Background()
	rsc := NewReplicaSetStatus(rs)
	got, err := rsc.Create(ctx, c)
	require.NoError(t, err)

	expected := component.NewQuadrant("Status")
	require.NoError(t, expected.Set(component.QuadNW, "Running", "3"))
	require.NoError(t, expected.Set(component.QuadNE, "Waiting", "0"))
	require.NoError(t, expected.Set(component.QuadSW, "Succeeded", "0"))
	require.NoError(t, expected.Set(component.QuadSE, "Failed", "0"))

	assert.Equal(t, expected, got)
}

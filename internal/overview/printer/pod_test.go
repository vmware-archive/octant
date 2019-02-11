package printer_test

import (
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
)

func Test_PodListHandler(t *testing.T) {
	printOptions := printer.Options{
		Cache: cache.NewMemoryCache(),
	}

	labels := map[string]string{
		"app": "testing",
	}

	now := time.Unix(1547211430, 0)

	object := &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
					Labels: labels,
				},
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
				Status: corev1.PodStatus{
					Phase: "Pending",
					ContainerStatuses: []corev1.ContainerStatus{
						corev1.ContainerStatus{
							Name:         "nginx",
							Image:        "nginx:1.15",
							RestartCount: 0,
							Ready:        true,
						},
						corev1.ContainerStatus{
							Name:         "kuard",
							Image:        "gcr.io/kuar-demo/kuard-amd64:1",
							RestartCount: 0,
							Ready:        false,
						},
					},
				},
			},
		},
	}

	got, err := printer.PodListHandler(object, printOptions)
	require.NoError(t, err)

	containers := component.NewContainers()
	containers.Add("nginx", "nginx:1.15")

	cols := component.NewTableCols("Name", "Labels", "Ready", "Status", "Restarts", "Age")
	expected := component.NewTable("Pods", cols)
	expected.Add(component.TableRow{
		"Name":     component.NewText("pod"),
		"Labels":   component.NewLabels(labels),
		"Ready":    component.NewText("1/2"),
		"Status":   component.NewText("Pending"),
		"Restarts": component.NewText("0"),
		"Age":      component.NewTimestamp(now),
	})

	assert.Equal(t, expected, got)
}

var (
	now      = time.Unix(1547211430, 0)
	validPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod",
			OwnerReferences: []metav1.OwnerReference{
				metav1.OwnerReference{
					APIVersion: "v1",
					Kind:       "ReplicationController",
					Name:       "myreplicationcontroller",
					Controller: ptrbool(true),
				},
			},
			CreationTimestamp: metav1.Time{
				Time: now,
			},
			DeletionTimestamp: &metav1.Time{
				Time: now,
			},
			DeletionGracePeriodSeconds: ptr64(30),
		},
		Spec: corev1.PodSpec{
			Priority:           ptrInt32(1000000),
			PriorityClassName:  "high-priority",
			ServiceAccountName: "default",
		},
		Status: corev1.PodStatus{
			StartTime:         &metav1.Time{Time: now},
			Phase:             corev1.PodRunning,
			Reason:            "SleepExpired",
			Message:           "Sleep expired",
			NominatedNodeName: "mynode",
			QOSClass:          "Guaranteed",
		},
	}
)

func Test_PodConfiguration(t *testing.T) {
	cases := []struct {
		name     string
		pod      *corev1.Pod
		isErr    bool
		expected *component.Summary
	}{
		{
			name: "general",
			pod:  validPod,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Priority",
					Content: component.NewText("1000000"),
				},
				{
					Header:  "PriorityClassName",
					Content: component.NewText("high-priority"),
				},
				{
					Header:  "Start Time",
					Content: component.NewTimestamp(now),
				},
				{
					Header:  "Status: Terminating",
					Content: component.NewTimestamp(now),
				},
				{
					Header:  "Termination Grace Period",
					Content: component.NewText("30s"),
				},
				{
					Header:  "Reason",
					Content: component.NewText("SleepExpired"),
				},
				{
					Header:  "Message",
					Content: component.NewText("Sleep expired"),
				},
				{
					Header:  "Controlled By",
					Content: component.NewLink("", "myreplicationcontroller", "/content/overview/workloads/replication-controllers/myreplicationcontroller"),
				},
				{
					Header:  "NominatedNodeName",
					Content: component.NewText("mynode"),
				},
				{
					Header:  "QoS Class",
					Content: component.NewText("Guaranteed"),
				},
				{
					Header:  "Service Account",
					Content: component.NewLink("", "default", "/content/overview/config-and-storage/service-accounts/default"),
				},
			}...),
		},
		{
			name:  "pod is nil",
			pod:   nil,
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cc := printer.NewPodConfiguration(tc.pod)
			summary, err := cc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, summary)
		})
	}
}

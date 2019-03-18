package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/conversion"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_PodListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		Cache: cachefake.NewMockCache(controller),
	}

	labels := map[string]string{
		"app": "testing",
	}

	now := time.Unix(1547211430, 0)

	object := &corev1.PodList{
		Items: []corev1.Pod{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod",
					Namespace: "default",
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
					NodeName: "node",
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

	ctx := context.Background()
	got, err := PodListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	containers := component.NewContainers()
	containers.Add("nginx", "nginx:1.15")

	cols := component.NewTableCols("Name", "Labels", "Ready", "Status", "Restarts", "Node", "Age")
	expected := component.NewTable("Pods", cols)
	expected.Add(component.TableRow{
		"Name":     component.NewLink("", "pod", "/content/overview/namespace/default/workloads/pods/pod"),
		"Labels":   component.NewLabels(labels),
		"Ready":    component.NewText("1/2"),
		"Status":   component.NewText("Pending"),
		"Restarts": component.NewText("0"),
		"Age":      component.NewTimestamp(now),
		"Node":     component.NewText("node"),
	})

	assert.Equal(t, expected, got)
}

var (
	now      = time.Unix(1547211430, 0)
	validPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				metav1.OwnerReference{
					APIVersion: "v1",
					Kind:       "ReplicationController",
					Name:       "myreplicationcontroller",
					Controller: conversion.PtrBool(true),
				},
			},
			CreationTimestamp: metav1.Time{
				Time: now,
			},
			DeletionTimestamp: &metav1.Time{
				Time: now,
			},
			DeletionGracePeriodSeconds: conversion.PtrInt64(30),
		},
		Spec: corev1.PodSpec{
			Priority:           conversion.PtrInt32(1000000),
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
					Header:  "Service Account",
					Content: component.NewLink("", "default", "/content/overview/namespace/default/config-and-storage/service-accounts/default"),
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
			cc := NewPodConfiguration(tc.pod)
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

func Test_createPodSummaryStatus(t *testing.T) {
	pod := testutil.CreatePod("pod")
	pod.Status.QOSClass = corev1.PodQOSBestEffort
	pod.Status.Phase = corev1.PodRunning
	pod.Status.PodIP = "10.1.1.1"
	pod.Status.HostIP = "10.2.1.1"

	got, err := createPodSummaryStatus(pod)
	require.NoError(t, err)

	sections := component.SummarySections{
		{Header: "QoS", Content: component.NewText("BestEffort")},
		{Header: "Status", Content: component.NewText("Running")},
		{Header: "Pod IP", Content: component.NewText("10.1.1.1")},
		{Header: "Host IP", Content: component.NewText("10.2.1.1")},
	}
	expected := component.NewSummary("Status", sections...)

	assert.Equal(t, expected, got)
}

func Test_createPodConditionsView(t *testing.T) {
	now := metav1.Time{Time: time.Now()}

	pod := testutil.CreatePod("pod")
	pod.Status.Conditions = []corev1.PodCondition{
		{
			Type:               corev1.PodInitialized,
			LastTransitionTime: now,
			Message:            "message",
			Reason:             "reason",
		},
	}

	got, err := createPodConditionsView(pod)
	require.NoError(t, err)

	cols := component.NewTableCols("Type", "Last Transition Time", "Message", "Reason")
	expected := component.NewTable("Pod Conditions", cols)
	expected.Add([]component.TableRow{
		{
			"Type":                 component.NewText("Initialized"),
			"Last Transition Time": component.NewTimestamp(now.Time),
			"Message":              component.NewText("message"),
			"Reason":               component.NewText("reason"),
		},
	}...)

	assert.Equal(t, expected, got)
}

func createPodWithPhase(name string, podLabels map[string]string, phase corev1.PodPhase, owner *metav1.OwnerReference) *corev1.Pod {
	pod := testutil.CreatePod(name)
	pod.Namespace = "testing"
	pod.Labels = podLabels
	pod.Status.Phase = phase

	if owner != nil {
		pod.SetOwnerReferences([]metav1.OwnerReference{*owner})
	}
	return pod
}

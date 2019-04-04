package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/pkg/cacheutil"
	"github.com/heptio/developer-dash/internal/conversion"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_StatefulSetListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		Cache: cachefake.NewMockCache(controller),
	}
	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := &appsv1.StatefulSetList{
		Items: []appsv1.StatefulSet{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps/v1",
					Kind:       "StatefulSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "web",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
					Labels: labels,
				},
				Status: appsv1.StatefulSetStatus{
					Replicas: 1,
				},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "myapp",
						},
					},
					Replicas: conversion.PtrInt32(3),
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								corev1.Container{
									Name:  "nginx",
									Image: "k8s.gcr.io/nginx-slim:0.8",
								},
							},
						},
					},
				},
			},
		},
	}

	ctx := context.Background()
	got, err := StatefulSetListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Desired", "Current", "Age", "Selector")
	expected := component.NewTable("StatefulSets", cols)
	expected.Add(component.TableRow{
		"Name":     component.NewLink("", "web", "/content/overview/workloads/stateful-sets/web"),
		"Labels":   component.NewLabels(labels),
		"Desired":  component.NewText("3"),
		"Current":  component.NewText("1"),
		"Age":      component.NewTimestamp(now),
		"Selector": component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
	})

	assert.Equal(t, expected, got)
}

func Test_StatefulSetStatus(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	c := cachefake.NewMockCache(controller)
	printOptions := Options{
		Cache: c,
	}

	labels := map[string]string{
		"app": "myapp",
	}

	sts := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "statefulset",
			Namespace: "testing",
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "myapp",
				},
			},
		},
	}

	pods := &corev1.PodList{
		Items: []corev1.Pod{
			*createPodWithPhase("web-0", labels, corev1.PodRunning, metav1.NewControllerRef(sts, sts.GroupVersionKind())),
			*createPodWithPhase("web-1", labels, corev1.PodRunning, metav1.NewControllerRef(sts, sts.GroupVersionKind())),
			*createPodWithPhase("web-2", labels, corev1.PodPending, metav1.NewControllerRef(sts, sts.GroupVersionKind())),
			*createPodWithPhase("random-pod", nil, corev1.PodRunning, nil),
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

	stsc := NewStatefulSetStatus(sts)
	ctx := context.Background()
	got, err := stsc.Create(ctx, printOptions.Cache)
	require.NoError(t, err)

	expected := component.NewQuadrant("Status")
	require.NoError(t, expected.Set(component.QuadNW, "Running", "2"))
	require.NoError(t, expected.Set(component.QuadNE, "Waiting", "1"))
	require.NoError(t, expected.Set(component.QuadSW, "Succeeded", "0"))
	require.NoError(t, expected.Set(component.QuadSE, "Failed", "0"))

	assert.Equal(t, expected, got)
}

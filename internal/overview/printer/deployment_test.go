package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/conversion"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func Test_DeploymentListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		Cache: cachefake.NewMockCache(controller),
	}

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	object := &appsv1.DeploymentList{
		Items: []appsv1.Deployment{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
					Labels: labels,
				},
				Status: appsv1.DeploymentStatus{
					Replicas:            3,
					AvailableReplicas:   2,
					UnavailableReplicas: 1,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: conversion.PtrInt32(3),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "my_app",
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
	got, err := DeploymentListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	containers := component.NewContainers()
	containers.Add("nginx", "nginx:1.15")
	containers.Add("kuard", "gcr.io/kuar-demo/kuard-amd64:1")

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	expected := component.NewTable("Deployments", cols)
	expected.Add(component.TableRow{
		"Name":       component.NewLink("", "deployment", "/content/overview/namespace/default/workloads/deployments/deployment"),
		"Labels":     component.NewLabels(labels),
		"Age":        component.NewTimestamp(now),
		"Selector":   component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "my_app")}),
		"Status":     component.NewText("2/3"),
		"Containers": containers,
	})

	assert.Equal(t, expected, got)
}

func Test_deploymentConfiguration(t *testing.T) {
	cases := []struct {
		name       string
		deployment *appsv1.Deployment
		isErr      bool
		expected   *component.Summary
	}{
		{
			name:       "rolling update",
			deployment: validDeployment,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Deployment Strategy",
					Content: component.NewText("RollingUpdate"),
				},
				{
					Header:  "Rolling Update Strategy",
					Content: component.NewText("Max Surge 25%, Max Unavailable 25%"),
				},
				{
					Header: "Selectors",
					Content: component.NewSelectors(
						[]component.Selector{
							component.NewExpressionSelector("key", component.OperatorIn, []string{"value1", "value2"}),
							component.NewLabelSelector("app", "my_app"),
						},
					),
				},
				{
					Header:  "Min Ready Seconds",
					Content: component.NewText("0"),
				},
				{
					Header:  "Revision History Limit",
					Content: component.NewText("5"),
				},
				{
					Header:  "Replicas",
					Content: component.NewText("3"),
				},
			}...),
		},
		{
			name:       "deployment is nil",
			deployment: nil,
			isErr:      true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dc := NewDeploymentConfiguration(tc.deployment)
			summary, err := dc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, summary)
		})
	}
}

var (
	rhl             int32 = 5
	validDeployment       = &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment",
			CreationTimestamp: metav1.Time{
				Time: time.Unix(1548377609, 0),
			},
			Labels: map[string]string{
				"app": "app",
			},
		},
		Status: appsv1.DeploymentStatus{
			Replicas:            3,
			AvailableReplicas:   2,
			UnavailableReplicas: 1,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:             conversion.PtrInt32(3),
			RevisionHistoryLimit: &rhl,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "my_app",
				},
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "key",
						Operator: "In",
						Values:   []string{"value1", "value2"},
					},
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
					MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
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
	}
)

func TestDeploymentStatus(t *testing.T) {
	d := &appsv1.Deployment{
		Status: appsv1.DeploymentStatus{
			UpdatedReplicas:     1,
			Replicas:            2,
			UnavailableReplicas: 3,
			AvailableReplicas:   4,
		},
	}

	ds := NewDeploymentStatus(d)
	got, err := ds.Create()
	require.NoError(t, err)

	expected := component.NewQuadrant("Status")
	require.NoError(t, expected.Set(component.QuadNW, "Updated", "1"))
	require.NoError(t, expected.Set(component.QuadNE, "Total", "2"))
	require.NoError(t, expected.Set(component.QuadSW, "Unavailable", "3"))
	require.NoError(t, expected.Set(component.QuadSE, "Available", "4"))

	assert.Equal(t, expected, got)
}

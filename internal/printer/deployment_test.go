/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/vmware-tanzu/octant/internal/octant"

	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/conversion"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_DeploymentListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	objectLabels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	object := &appsv1.Deployment{
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
			Labels: objectLabels,
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
						{
							Name:  "nginx",
							Image: "nginx:1.15",
						},
						{
							Name:  "kuard",
							Image: "gcr.io/kuar-demo/kuard-amd64:1",
						},
					},
				},
			},
		},
	}

	tpo.PathForObject(object, object.Name, "/path")

	list := &appsv1.DeploymentList{
		Items: []appsv1.Deployment{*object},
	}

	ctx := context.Background()
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, object)
	got, err := DeploymentListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	containers := component.NewContainers()
	containers.Add("nginx", "nginx:1.15")
	containers.Add("kuard", "gcr.io/kuar-demo/kuard-amd64:1")

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	expected := component.NewTable("Deployments", "We couldn't find any deployments!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "deployment", "/path",
			genObjectStatus(component.TextStatusWarning, []string{
				"Expected 3 replicas, but 2 are available"})),
		"Labels":     component.NewLabels(objectLabels),
		"Age":        component.NewTimestamp(now),
		"Selector":   component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "my_app")}),
		"Status":     component.NewText("2/3"),
		"Containers": containers,
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, object),
		}),
	})

	component.AssertEqual(t, expected, got)
}

func Test_deploymentConfiguration(t *testing.T) {
	var rhl int32 = 5
	validDeployment := testutil.CreateDeployment("deployment")
	validDeployment.Spec = appsv1.DeploymentSpec{
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
					{
						Name:  "nginx",
						Image: "nginx:1.15",
					},
					{
						Name:  "kuard",
						Image: "gcr.io/kuar-demo/kuard-amd64:1",
					},
				},
			},
		},
	}

	cases := []struct {
		name       string
		deployment *appsv1.Deployment
		isErr      bool
		expected   *component.Summary
	}{
		{
			name:       "deployment",
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
			dc.actionGenerators = []actionGeneratorFunction{}

			summary, err := dc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createDeploymentSummaryStatus(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment")
	deployment.Status.AvailableReplicas = 1
	deployment.Status.ReadyReplicas = 2
	deployment.Status.Replicas = 3
	deployment.Status.UnavailableReplicas = 4
	deployment.Status.UpdatedReplicas = 5

	got, err := createDeploymentSummaryStatus(deployment)
	require.NoError(t, err)

	sections := component.SummarySections{
		{Header: "Available Replicas", Content: component.NewText("1")},
		{Header: "Ready Replicas", Content: component.NewText("2")},
		{Header: "Total Replicas", Content: component.NewText("3")},
		{Header: "Unavailable Replicas", Content: component.NewText("4")},
		{Header: "Updated Replicas", Content: component.NewText("5")},
	}
	expected := component.NewSummary("Status", sections...)

	assert.Equal(t, expected, got)
}

func Test_DeploymentPods(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	var replicas int32 = 3

	deployment := testutil.CreateDeployment("deployment")

	replicaSet := testutil.CreateAppReplicaSet("replicaset")
	replicaSet.Spec.Replicas = &replicas
	replicaSet.SetOwnerReferences(testutil.ToOwnerReferences(t, deployment))

	now := testutil.Time()
	pod := testutil.CreatePod("pod")
	pod.SetOwnerReferences(testutil.ToOwnerReferences(t, replicaSet))
	pod.ObjectMeta.CreationTimestamp = metav1.Time{Time: now}
	pod.Status.Phase = corev1.PodRunning
	pod.Spec.Containers = []corev1.Container{
		{
			Name:  "nginx",
			Image: "nginx:1.15",
		},
	}
	pod.Status.ContainerStatuses = []corev1.ContainerStatus{
		{
			Name:         "nginx",
			Image:        "nginx:1.15",
			RestartCount: 0,
			Ready:        true,
		},
	}

	tpo.PathForObject(pod, pod.Name, "/pod")

	podKey := store.Key{
		Namespace:  "namespace",
		APIVersion: "v1",
		Kind:       "Pod",
	}

	replicaSetKey := store.Key{
		Namespace:  "namespace",
		APIVersion: "apps/v1",
		Kind:       "ReplicaSet",
	}

	tpo.objectStore.EXPECT().List(gomock.Any(), replicaSetKey).
		Return(testutil.ToUnstructuredList(t, replicaSet), false, nil)

	tpo.objectStore.EXPECT().List(gomock.Any(), podKey).
		Return(testutil.ToUnstructuredList(t, pod), false, nil)

	ctx := context.Background()
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, pod)

	replicaSets, err := listReplicaSetsAsObjects(ctx, deployment, printOptions)
	require.NoError(t, err)

	got, err := createRollingPodListView(ctx, replicaSets, printOptions)
	require.NoError(t, err)

	expected := component.NewTableWithRows("Pods", "We couldn't find any pods!", podColsWithOutLabels, []component.TableRow{
		{
			"Name": component.NewLink("", pod.Name, "/pod",
				genObjectStatus(component.TextStatusOK, []string{"Pod is OK"})),
			"Age":      component.NewTimestamp(now),
			"Ready":    component.NewText("1/1"),
			"Restarts": component.NewText("0"),
			"Phase":    component.NewText("Running"),
			"Node":     component.NewText("<not scheduled>"),
			component.GridActionKey: gridActionsFactory([]component.GridAction{
				buildObjectDeleteAction(t, pod),
			}),
		},
	})
	addPodTableFilters(expected)

	component.AssertEqual(t, expected, got)
}

func Test_editDeploymentAction(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment")
	deployment.Spec.Replicas = pointer.Int32Ptr(3)

	actions, err := editDeploymentAction(deployment)
	require.NoError(t, err)
	assert.Len(t, actions, 1)

	got := actions[0]

	apiVersion, kind := deployment.GroupVersionKind().ToAPIVersionAndKind()

	expected := component.Action{
		Name:  "Edit",
		Title: "Deployment Editor",
		Form: component.Form{
			Fields: []component.FormField{
				component.NewFormFieldNumber("Replicas", "replicas", "3"),
				component.NewFormFieldHidden("apiVersion", apiVersion),
				component.NewFormFieldHidden("kind", kind),
				component.NewFormFieldHidden("name", deployment.Name),
				component.NewFormFieldHidden("namespace", deployment.Namespace),
				component.NewFormFieldHidden("action", octant.ActionDeploymentConfiguration),
			},
		},
	}

	assert.Equal(t, expected, got)
}

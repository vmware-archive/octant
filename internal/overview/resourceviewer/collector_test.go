package resourceviewer

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/testutil"

	"github.com/heptio/developer-dash/internal/conversion"

	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func Test_Collector(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment")
	deployment.Status = appsv1.DeploymentStatus{
		Replicas:          1,
		AvailableReplicas: 1,
	}

	replicaSet1 := testutil.CreateReplicaSet("replicaSet1")
	replicaSet1.Spec = appsv1.ReplicaSetSpec{
		Replicas: conversion.PtrInt32(1),
	}
	replicaSet1.Status = appsv1.ReplicaSetStatus{
		Replicas:          1,
		AvailableReplicas: 1,
	}

	replicaSet2 := testutil.CreateReplicaSet("replicaSet2")

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod",
			UID:  types.UID("pod"),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(replicaSet1,
					schema.FromAPIVersionAndKind(replicaSet1.APIVersion,
						replicaSet1.Kind)),
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()
	o := storefake.NewMockObjectStore(controller)
	c := NewCollector(o)

	ctx := context.Background()

	err := c.Process(ctx, deployment)
	require.NoError(t, err)

	err = c.Process(ctx, replicaSet1)
	require.NoError(t, err)

	err = c.Process(ctx, replicaSet2)
	require.NoError(t, err)

	err = c.Process(ctx, pod)
	require.NoError(t, err)

	err = c.AddChild(deployment, replicaSet1, replicaSet2)
	require.NoError(t, err)

	err = c.AddChild(replicaSet1, pod)
	require.NoError(t, err)

	got, err := c.Component("deployment")
	require.NoError(t, err)

	q := url.Values{}

	expected := component.NewResourceViewer("Resource Viewer")
	expected.AddEdge("deployment", "replicaSet1", component.EdgeTypeExplicit)
	expected.AddEdge("replicaSet1", "pods-replicaSet1", component.EdgeTypeExplicit)
	expected.AddNode("deployment", component.Node{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Name:       "deployment",
		Status:     "ok",
		Details:    []component.Component{component.NewText("Deployment is OK")},
		Path:       link.ForObjectWithQuery(deployment, deployment.Name, q),
	})

	expected.AddNode("replicaSet1", component.Node{
		APIVersion: "extensions/v1beta1",
		Kind:       "ReplicaSet",
		Name:       "replicaSet1",
		Status:     "ok",
		Details:    []component.Component{component.NewText("Replica Set is OK")},
		Path:       link.ForObjectWithQuery(replicaSet1, replicaSet1.Name, q),
	})

	podStatus := component.NewPodStatus()
	details := []component.Component{
		component.NewText(""),
	}
	podStatus.AddSummary("pod", details, component.NodeStatusOK)

	expected.AddNode("pods-replicaSet1", component.Node{
		APIVersion: "v1",
		Kind:       "Pod",
		Name:       "replicaSet1 pods",
		Status:     "ok",
		Details: []component.Component{
			podStatus,
			component.NewText("Pod count: 1"),
		},
	})
	expected.Select("deployment")

	assertComponentEqual(t, expected, got)

	got, err = c.Component("pod")
	require.NoError(t, err)

	expected.Select("pods-replicaSet1")

	assertComponentEqual(t, expected, got)
}

func assertComponentEqual(t *testing.T, expected, got component.Component) {
	transformer := func(in component.Component) string {
		b, err := json.MarshalIndent(in, "  ", "  ")
		require.NoError(t, err)

		return string(b)

	}

	expectedString := transformer(expected)
	gotString := transformer(got)

	assert.Equal(t, expectedString, gotString)
}

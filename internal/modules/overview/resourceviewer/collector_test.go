package resourceviewer

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	configFake "github.com/heptio/developer-dash/internal/config/fake"
	"github.com/heptio/developer-dash/internal/conversion"
	linkFake "github.com/heptio/developer-dash/internal/link/fake"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
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

	setOwner(t, replicaSet1, deployment)

	replicaSet2 := testutil.CreateReplicaSet("replicaSet2")

	serviceAccount := testutil.CreateServiceAccount("service-account")
	serviceAccount.UID = types.UID("service-account")

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod",
			UID:  types.UID("pod"),
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}
	pod.Spec.ServiceAccountName = serviceAccount.Name

	setOwner(t, pod, replicaSet1)

	controller := gomock.NewController(t)
	defer controller.Finish()

	objectStore := storefake.NewMockObjectStore(controller)

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

	linkGenerator := linkFake.NewMockInterface(controller)

	q := url.Values{}
	deploymentLink := component.NewLink("", "deployment", "/deployment")
	linkGenerator.EXPECT().
		ForObjectWithQuery(deployment, deployment.Name, q).
		Return(deploymentLink, nil)

	replicaSetLink := component.NewLink("", "replicaSet1", "/replica-set")
	linkGenerator.EXPECT().
		ForObjectWithQuery(replicaSet1, replicaSet1.Name, q).
		Return(replicaSetLink, nil)

	serviceAccountLink := component.NewLink("", "service-account", "/service-account")
	linkGenerator.EXPECT().
		ForObjectWithQuery(serviceAccount, serviceAccount.Name, q).
		Return(serviceAccountLink, nil)

	mockLink := func(c *Collector) {
		c.link = linkGenerator
	}

	c, err := NewCollector(dashConfig, mockLink)
	require.NoError(t, err)

	ctx := context.Background()

	err = c.Process(ctx, deployment)
	require.NoError(t, err)

	err = c.Process(ctx, replicaSet1)
	require.NoError(t, err)

	err = c.Process(ctx, replicaSet2)
	require.NoError(t, err)

	err = c.Process(ctx, serviceAccount)
	require.NoError(t, err)

	err = c.Process(ctx, pod)
	require.NoError(t, err)

	err = c.AddChild(deployment, replicaSet1, replicaSet2)
	require.NoError(t, err)

	err = c.AddChild(replicaSet1, pod)
	require.NoError(t, err)

	err = c.AddChild(pod, serviceAccount)
	require.NoError(t, err)

	got, err := c.Component("deployment")
	require.NoError(t, err)

	expected := component.NewResourceViewer("Resource Viewer")
	expected.AddNode("apps/v1-Deployment-deployment", component.Node{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Name:       "deployment",
		Status:     "ok",
		Details:    []component.Component{component.NewText("Deployment is OK")},
		Path:       deploymentLink,
	})

	expected.AddNode("apps/v1-ReplicaSet-replicaSet1", component.Node{
		APIVersion: "extensions/v1beta1",
		Kind:       "ReplicaSet",
		Name:       "replicaSet1",
		Status:     "ok",
		Details:    []component.Component{component.NewText("Replica Set is OK")},
		Path:       replicaSetLink,
	})

	podStatus := component.NewPodStatus()
	details := []component.Component{
		component.NewText(""),
	}
	podStatus.AddSummary("pod", details, component.NodeStatusOK)

	expected.AddNode("pods-apps/v1-ReplicaSet-replicaSet1", component.Node{
		APIVersion: "v1",
		Kind:       "Pod",
		Name:       "replicaSet1 pods",
		Status:     "ok",
		Details: []component.Component{
			podStatus,
			component.NewText("Pod count: 1"),
		},
	})

	expected.AddNode("v1-ServiceAccount-service-account", component.Node{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
		Name:       "service-account",
		Status:     "ok",
		Details:    []component.Component{component.NewText("v1 ServiceAccount is OK")},
		Path:       serviceAccountLink,
	})

	require.NoError(t, expected.AddEdge("apps/v1-Deployment-deployment", "apps/v1-ReplicaSet-replicaSet1", component.EdgeTypeExplicit))
	require.NoError(t, expected.AddEdge("apps/v1-ReplicaSet-replicaSet1", "pods-apps/v1-ReplicaSet-replicaSet1", component.EdgeTypeExplicit))
	require.NoError(t, expected.AddEdge("v1-ServiceAccount-service-account", "pods-apps/v1-ReplicaSet-replicaSet1", component.EdgeTypeExplicit))

	expected.Select("deployment")

	assertComponentEqual(t, expected, got)

	got, err = c.Component("pod")
	require.NoError(t, err)

	expected.Select("pods-apps/v1-ReplicaSet-replicaSet1")

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

func setOwner(t *testing.T, object metav1.Object, owner runtime.Object) {
	require.NotNil(t, object)
	require.NotNil(t, owner)

	accessor := meta.NewAccessor()
	apiVersion, err := accessor.APIVersion(owner)
	require.NoError(t, err)

	kind, err := accessor.Kind(owner)
	require.NoError(t, err)

	name, err := accessor.Name(owner)
	require.NoError(t, err)

	uid, err := accessor.UID(owner)
	require.NoError(t, err)

	boolTrue := true

	object.SetOwnerReferences([]metav1.OwnerReference{
		{
			APIVersion:         apiVersion,
			Kind:               kind,
			Name:               name,
			UID:                uid,
			Controller:         &boolTrue,
			BlockOwnerDeletion: &boolTrue,
		},
	})
}

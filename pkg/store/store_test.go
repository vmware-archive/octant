package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/action"
)

func TestKey_ToActionPayload(t *testing.T) {
	pod := testutil.CreatePod("pod")
	key := Key{
		Namespace:  pod.Namespace,
		APIVersion: pod.APIVersion,
		Kind:       pod.Kind,
		Name:       pod.Name,
	}

	got := key.ToActionPayload()

	expected := action.Payload{
		"namespace":  pod.Namespace,
		"apiVersion": pod.APIVersion,
		"kind":       pod.Kind,
		"name":       pod.Name,
	}

	assert.Equal(t, expected, got)
}

func TestKey_GroupVersionKind(t *testing.T) {
	pod := testutil.CreatePod("pod")
	key := Key{
		Namespace:  pod.Namespace,
		APIVersion: pod.APIVersion,
		Kind:       pod.Kind,
		Name:       pod.Name,
	}

	got := key.GroupVersionKind()
	assert.Equal(t, gvk.Pod, got)
}

func TestKeyFromObject(t *testing.T) {
	pod := testutil.CreatePod("pod")

	got, err := KeyFromObject(pod)
	require.NoError(t, err)

	expected := Key{
		Namespace:  pod.Namespace,
		APIVersion: pod.APIVersion,
		Kind:       pod.Kind,
		Name:       pod.Name,
	}

	assert.Equal(t, expected, got)
}

func TestKeyFromPayload(t *testing.T) {
	pod := testutil.CreatePod("pod")

	payload := action.Payload{
		"namespace":  pod.Namespace,
		"apiVersion": pod.APIVersion,
		"kind":       pod.Kind,
		"name":       pod.Name,
	}

	got, err := KeyFromPayload(payload)
	require.NoError(t, err)

	expected := Key{
		Namespace:  pod.Namespace,
		APIVersion: pod.APIVersion,
		Kind:       pod.Kind,
		Name:       pod.Name,
	}
	assert.Equal(t, expected, got)
}

func TestKeyFromGroupVersionKind(t *testing.T) {
	actual := KeyFromGroupVersionKind(gvk.Pod)
	expected := Key{
		APIVersion: "v1",
		Kind:       "Pod",
	}
	require.Equal(t, expected, actual)
}

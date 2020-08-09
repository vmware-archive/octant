package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/action"
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

func TestKey_Validate(t *testing.T) {
	tests := []struct {
		name    string
		key     Key
		wantErr bool
	}{
		{
			name: "valid with namespace",
			key: Key{
				Namespace:  "test",
				APIVersion: "apiVersion",
				Kind:       "kind",
				Name:       "name",
			},
		},
		{
			name: "valid without namespace",
			key: Key{
				APIVersion: "apiVersion",
				Kind:       "kind",
				Name:       "name",
			},
		},
		{
			name: "missing api version",
			key: Key{
				Kind: "kind",
				Name: "name",
			},
			wantErr: true,
		},
		{
			name: "missing kind",
			key: Key{
				APIVersion: "apiVersion",
				Name:       "name",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.key.Validate()
			testutil.RequireErrorOrNot(t, tt.wantErr, err)
		})
	}
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

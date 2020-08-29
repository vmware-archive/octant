package store

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

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
		{
			name: "Exists operator with values",
			key: Key{
				APIVersion: "apiVersion",
				Name:       "name",
				Kind:       "kind",
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{Key: "foo", Operator: metav1.LabelSelectorOpExists, Values: []string{"foo1"}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "DoesNotExist operator with values",
			key: Key{
				APIVersion: "apiVersion",
				Name:       "name",
				Kind:       "kind",
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{Key: "foo", Operator: metav1.LabelSelectorOpDoesNotExist, Values: []string{"foo1"}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "In operator without values",
			key: Key{
				APIVersion: "apiVersion",
				Name:       "name",
				Kind:       "kind",
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{Key: "foo", Operator: metav1.LabelSelectorOpIn},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "NotIn operator without values",
			key: Key{
				APIVersion: "apiVersion",
				Name:       "name",
				Kind:       "kind",
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{Key: "foo", Operator: metav1.LabelSelectorOpNotIn},
					},
				},
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
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Key
	}{
		{
			name: "in general",
			args: args{map[string]interface{}{
				"namespace":  "namespace",
				"apiVersion": "apiVersion",
				"kind":       "kind",
				"name":       "name",
			}},
			want: Key{
				Namespace:  "namespace",
				APIVersion: "apiVersion",
				Kind:       "kind",
				Name:       "name",
			},
		},
		{
			name: "with label selector",
			args: args{map[string]interface{}{
				"namespace":  "namespace",
				"apiVersion": "apiVersion",
				"kind":       "kind",
				"labelSelector": map[string]interface{}{
					"matchLabels": map[string]string{
						"foo": "bar",
					},
				},
			}},
			want: Key{
				Namespace:  "namespace",
				APIVersion: "apiVersion",
				Kind:       "kind",
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"foo": "bar",
					},
				},
			},
		},
		{
			name: "with label set",
			args: args{map[string]interface{}{
				"namespace":  "namespace",
				"apiVersion": "apiVersion",
				"kind":       "kind",
				"selector": map[string]string{
					"foo": "bar",
				},
			}},
			want: Key{
				Namespace:  "namespace",
				APIVersion: "apiVersion",
				Kind:       "kind",
				Selector: &labels.Set{
					"foo": "bar",
				},
			},
		},
		{
			name: "missing required field",
			args: args{map[string]interface{}{
				// apiVersion is required
				"namespace": "namespace",
				"kind":      "kind",
				"name":      "name",
			}},
			wantErr: true,
		},
		{
			name: "invalid type for field",
			args: args{map[string]interface{}{
				// namespace should be a string.
				"namespace":  1,
				"apiVersion": "apiVersion",
				"kind":       "kind",
				"name":       "name",
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KeyFromPayload(tt.args.m)
			testutil.RequireErrorOrNot(t, tt.wantErr, err, func() {
				require.Equal(t, tt.want, got)
			})
		})
	}
}

func TestKeyFromGroupVersionKind(t *testing.T) {
	actual := KeyFromGroupVersionKind(gvk.Pod)
	expected := Key{
		APIVersion: "v1",
		Kind:       "Pod",
	}
	require.Equal(t, expected, actual)
}

package applications

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu/octant/internal/testutil"
)

func Test_sortedDescriptorList(t *testing.T) {
	descriptors := map[descriptor]metadata{
		descriptor{name: "b", instance: "b", version: "a"}: {},
		descriptor{name: "b", instance: "c", version: "a"}: {},
		descriptor{name: "a", instance: "b", version: "a"}: {},
		descriptor{name: "b", instance: "b", version: "d"}: {},
		descriptor{name: "a", instance: "a", version: "a"}: {},
	}

	actual := sortedDescriptorList(descriptors)

	expected := []descriptor{
		{name: "a", instance: "a", version: "a"},
		{name: "a", instance: "b", version: "a"},
		{name: "b", instance: "b", version: "a"},
		{name: "b", instance: "b", version: "d"},
		{name: "b", instance: "c", version: "a"},
	}

	require.Equal(t, expected, actual)
}

func Test_podBelongsToApplication(t *testing.T) {
	tests := []struct {
		name         string
		labels       map[string]string
		shouldBelong bool
		expected     descriptor
		isErr        bool
	}{
		{
			name: "pod with home labels",
			labels: map[string]string{
				appLabelName:     "name",
				appLabelInstance: "instance",
				appLabelVersion:  "version",
			},
			shouldBelong: true,
			expected: descriptor{
				name:     "name",
				instance: "instance",
				version:  "version",
			},
		},
		{
			name:         "pod with out labels",
			labels:       map[string]string{},
			shouldBelong: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			object := testutil.ToUnstructured(
				t,
				testutil.CreatePod("pod", withPodLabels(test.labels)))

			actualDescriptor, actualBelongs, err := podBelongsToApplication(object)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.shouldBelong, actualBelongs)
			assert.Equal(t, test.expected, actualDescriptor)
		})
	}

}

func Test_getLabel(t *testing.T) {
	tests := []struct {
		name     string
		labels   map[string]string
		key      string
		expected string
		isErr    bool
	}{
		{
			name:     "label exists",
			labels:   map[string]string{"key": "value"},
			key:      "key",
			expected: "value",
		},
		{
			name:   "label does not exist",
			labels: map[string]string{"key": "value"},
			key:    "other",
			isErr:  true,
		},
		{
			name:   "key is blank",
			labels: map[string]string{"key": "value"},
			key:    "",
			isErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			object := testutil.ToUnstructured(
				t,
				testutil.CreatePod("pod", withPodLabels(test.labels)))

			got, err := getLabel(object, test.key)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.expected, got)
		})
	}
}

func withPodLabels(labels map[string]string) testutil.PodOption {
	return func(pod *corev1.Pod) {
		pod.SetLabels(labels)
	}
}

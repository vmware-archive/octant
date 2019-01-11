package component

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Containers_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        ViewComponent
		expectedPath string
		isErr        bool
	}{
		{
			name: "general",
			input: &Containers{
				Config: ContainersConfig{
					Containers: []ContainerDef{
						ContainerDef{
							Name:  "nginx",
							Image: "nginx:1.15",
						},
						ContainerDef{
							Name:  "kuard",
							Image: "gcr.io/kuar-demo/kuard-amd64:1",
						},
					},
				},
			},
			expectedPath: "container.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := (err != nil)
			if isErr != tc.isErr {
				t.Fatalf("Unexepected error: %v", err)
			}

			expected, err := ioutil.ReadFile(filepath.Join("testdata", tc.expectedPath))
			require.NoError(t, err, "reading test fixtures")
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

func Test_Containers_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    ViewComponent
		expected bool
	}{
		{
			name: "general",
			input: &Containers{
				Config: ContainersConfig{
					Containers: []ContainerDef{
						ContainerDef{},
					},
				},
			},
			expected: false,
		},
		{
			name:     "empty",
			input:    &Containers{},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.IsEmpty(), "IsEmpty mismatch")
		})
	}
}

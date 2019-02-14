package component

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func Test_YAML_Marshal(t *testing.T) {
	cases := []struct {
		name         string
		input        *YAML
		expectedPath string
		isErr        bool
	}{
		{
			name: "in general",
			input: &YAML{
				Config: YAMLConfig{
					Data: "---\nfoo: bar",
				},
				Metadata: Metadata{},
			},
			expectedPath: "yaml1.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := (err != nil)
			if isErr != tc.isErr {
				t.Fatalf("Unexpected error: %v", err)
			}

			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err, "reading test fixtures")
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

func Test_YAML_Data(t *testing.T) {
	y := NewYAML(Title(NewText("Title")))

	pod := &corev1.Pod{}
	require.NoError(t, y.Data(pod))

	got := y.Config.Data
	expected := "---\nmetadata:\n  creationTimestamp: null\nspec:\n  containers: null\nstatus: {}\n"

	assert.Equal(t, expected, got)
}

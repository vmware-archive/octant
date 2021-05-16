package component

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SelectFile_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        Component
		expectedPath string
		isErr        bool
	}{
		{
			name:         "in general",
			expectedPath: "select_file.json",
			input: &SelectFile{
				Base: newBase(TypeSelectFile, nil),
				Config: SelectFileConfig{
					Label:         "Open File",
					Multiple:      false,
					Status:        FileStatusSuccess,
					StatusMessage: "Success",
					Layout:        LayoutCompact,
					Action:        "action.octant.dev/SelectFileAction",
				}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := err != nil
			if isErr != tc.isErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err)
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

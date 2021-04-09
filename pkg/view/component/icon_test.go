package component

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

func Test_Icon_Marshal(t *testing.T) {
	test := []struct {
		name         string
		input        Component
		expectedPath string
		isErr        bool
	}{
		{
			name: "in general",
			input: &Icon{
				Base: newBase(TypeIcon, nil),
				Config: IconConfig{
					Shape:      "user",
					Size:       "16",
					Direction:  DirectionDown,
					Flip:       FlipHorizontal,
					Solid:      true,
					Status:     StatusDanger,
					Inverse:    false,
					Badge:      BadgeDanger,
					Color:      "#add8e6",
					BadgeColor: "purple",
					Label:      "example icon",
				},
			},
			expectedPath: "icon.json",
			isErr:        false,
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := err != nil
			if isErr != tc.isErr {
				t.Fatalf("UnExpected error: %v", err)
			}

			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err)
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

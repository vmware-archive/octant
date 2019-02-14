package component

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Grid_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        ViewComponent
		expectedPath string
		isErr        bool
	}{
		{
			name: "general",
			input: &Grid{
				Config: GridConfig{
					Panels: []Panel{
						Panel{
							Config: PanelConfig{
								Position: PanelPosition{X: 0, Y: 0, W: 12, H: 7},
								Content: &Text{
									Config: TextConfig{
										Text: "Panel contents",
									},
								},
							},
						},
					},
				},
			},
			expectedPath: "grid.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := (err != nil)
			if isErr != tc.isErr {
				t.Fatalf("Unexepected error: %v", err)
			}

			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err, "reading test fixtures")
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

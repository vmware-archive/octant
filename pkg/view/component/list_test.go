package component

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_List_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        Component
		expectedPath string
		isErr        bool
	}{
		{
			name: "general",
			input: &List{
				base: newBase(typeList, TitleFromString("mylist")),
				Config: ListConfig{
					Items: []Component{
						&Link{
							Config: LinkConfig{
								Text: "nginx-deployment",
								Ref:  "/overview/deployments/nginx-deployment",
							},
						},
						&Labels{
							Config: LabelsConfig{
								Labels: map[string]string{
									"app": "nginx",
								},
							},
						},
					},
				},
			},
			expectedPath: "list.json",
		},
	}

	for _, tc := range tests {
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

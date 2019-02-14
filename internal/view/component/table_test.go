package component

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Table_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        ViewComponent
		expectedPath string
		isErr        bool
	}{
		{
			name: "general",
			input: &Table{
				Metadata: Metadata{
					Title: []TitleViewComponent{NewText("my table")},
				},
				Config: TableConfig{
					Columns: []TableCol{
						TableCol{Name: "Name", Accessor: "Name"},
						TableCol{Name: "Description", Accessor: "Description"},
					},
					Rows: []TableRow{
						TableRow{
							"Name": &Text{
								Config: TextConfig{
									Text: "First",
								},
							},
							"Description": &Text{
								Config: TextConfig{
									Text: "The first row",
								},
							},
						},
						TableRow{
							"Name": &Text{
								Config: TextConfig{
									Text: "Last",
								},
							},
							"Description": &Text{
								Config: TextConfig{
									Text: "The last row",
								},
							},
						},
					},
					EmptyContent: "",
				},
			},
			expectedPath: "table.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := (err != nil)
			if isErr != tc.isErr {
				t.Fatalf("Unexpected error: %v", err)
			}

			expected, err := ioutil.ReadFile(filepath.Join("testdata", tc.expectedPath))
			require.NoError(t, err, "reading test fixtures")
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

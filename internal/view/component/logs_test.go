package component

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Logs_Marshal(t *testing.T) {
	cases := []struct {
		name         string
		input        *Logs
		expectedPath string
		isErr        bool
	}{
		{
			name: "in general",
			input: &Logs{
				Metadata: Metadata{
					Type:  "logs",
					Title: Title(NewText("Logs")),
				},
				Config: LogsConfig{
					Containers: []string{"one", "two"},
				},
			},
			expectedPath: "logs.json",
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

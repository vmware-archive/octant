package component

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ButtonGroup_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        func() *ButtonGroup
		expectedFile string
		isErr        bool
	}{
		{
			name: "empty button group",
			input: func() *ButtonGroup {
				return NewButtonGroup()
			},
			expectedFile: "button_group_empty.json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := json.Marshal(test.input())
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expected, err := ioutil.ReadFile(filepath.Join("testdata", test.expectedFile))
			require.NoError(t, err)

			assert.JSONEq(t, string(expected), string(got))

		})
	}
}

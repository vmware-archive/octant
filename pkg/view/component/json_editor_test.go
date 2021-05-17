/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

func TestJSONEditor_MarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		input        func() *JSONEditor
		expectedFile string
		isError      bool
	}{
		{
			name: "in general",
			input: func() *JSONEditor {
				return NewJSONEditor("{ \"hello\": 123 }")
			},
			expectedFile: "json_editor.json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := json.Marshal(test.input())
			if test.isError {
				require.NoError(t, err)
				return
			}
			require.NoError(t, err)

			expected, err := ioutil.ReadFile(filepath.Join("testdata", test.expectedFile))
			require.NoError(t, err)

			assert.JSONEq(t, string(expected), string(got))
		})
	}
}

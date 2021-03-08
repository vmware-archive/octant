/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Code_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    Component
		expected string
		isErr    bool
	}{
		{
			name: "general",
			input: &Code{
				Config: CodeConfig{
					Code: "hello world\nthis is a newline",
				},
			},
			expected: `
			{
				"metadata": {
								"type": "codeBlock"
				},
				"config": {
					"value": "hello world\nthis is a newline"
				}
			}
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := (err != nil)
			if isErr != tc.isErr {
				t.Fatalf("Unexpected error: %v", err)
			}

			assert.JSONEq(t, tc.expected, string(actual))
		})
	}
}

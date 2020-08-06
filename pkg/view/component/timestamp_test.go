/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Timestamp_Marshal(t *testing.T) {
	ts, err := time.Parse(time.RFC3339, "1969-07-21T02:56:00+00:00")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    Component
		expected string
		isErr    bool
	}{
		{
			name: "general",
			input: &Timestamp{
				Config: TimestampConfig{
					Timestamp: ts.Unix(),
				},
			},
			expected: `
            {
                "metadata": {
                  "type": "timestamp"
                },
                "config": {
                  "timestamp": -14159040
                }
            }
`,
		},
		{
			name: "with title",
			input: &Timestamp{
				Base: newBase(TypeTimestamp, TitleFromString("LandedOn")),
				Config: TimestampConfig{
					Timestamp: ts.Unix(),
				},
			},
			expected: `
            {
                "metadata": {
									"type": "timestamp",
									"title": [
										{
											"config": { "value": "LandedOn" },
											"metadata": { "type": "text" }
										}
									]
                },
                "config": {
                  "timestamp": -14159040
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

func Test_Timestamp_LessThan(t *testing.T) {
	cases := []struct {
		name     string
		ts       Timestamp
		other    Component
		expected bool
	}{
		{
			name:     "is less",
			ts:       *NewTimestamp(time.Unix(5, 0)),
			other:    NewTimestamp(time.Unix(6, 0)),
			expected: true,
		},
		{
			name:     "is not less",
			ts:       *NewTimestamp(time.Unix(5, 0)),
			other:    NewTimestamp(time.Unix(4, 0)),
			expected: false,
		},
		{
			name:     "other is not a timestamp",
			ts:       *NewTimestamp(time.Unix(5, 0)),
			other:    nil,
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.ts.LessThan(tc.other)
			assert.Equal(t, tc.expected, got)
		})
	}
}

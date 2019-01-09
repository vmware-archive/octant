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
		input    ViewComponent
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
				Metadata: Metadata{
					Title: "LandedOn",
				},
				Config: TimestampConfig{
					Timestamp: ts.Unix(),
				},
			},
			expected: `
            {
                "metadata": {
                  "type": "timestamp",
                  "title": "LandedOn"
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
				t.Fatalf("Unexepected error: %v", err)
			}

			assert.JSONEq(t, tc.expected, string(actual))
		})
	}
}

func Test_Timestamp_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    ViewComponent
		expected bool
	}{
		{
			name: "general",
			input: &Timestamp{
				Config: TimestampConfig{
					Timestamp: -14159040,
				},
			},
			expected: false,
		},
		{
			name:     "empty",
			input:    &Timestamp{},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.IsEmpty(), "IsEmpty mismatch")
		})
	}
}

package component

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Text_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    ViewComponent
		expected string
		isErr    bool
	}{
		{
			name: "general",
			input: &Text{
				Config: TextConfig{
					Text: "lorem ipsum",
				},
			},
			expected: `
            {
                "metadata": {
                  "type": "text"
                },
                "config": {
                  "value": "lorem ipsum"
                }
            }
`,
		},
		{
			name: "with title",
			input: &Text{
				Metadata: Metadata{
					Title: []TitleViewComponent{NewText("image")},
				},
				Config: TextConfig{
					Text: "nginx:latest",
				},
			},
			expected: `
            {
                "metadata": {
									"type": "text",
									"title": [
										{
											"config": { "value": "image" },
											"metadata": { "type": "text" }
										}
									]
                },
                "config": {
                  "value": "nginx:latest"
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

func Test_Text_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    ViewComponent
		expected bool
	}{
		{
			name: "general",
			input: &Text{
				Config: TextConfig{
					Text: "lorem ipsum",
				},
			},
			expected: false,
		},
		{
			name:     "empty",
			input:    &Text{},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.IsEmpty(), "IsEmpty mismatch")
		})
	}
}

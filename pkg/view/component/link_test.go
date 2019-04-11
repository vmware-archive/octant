package component

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Link_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    Component
		expected string
		isErr    bool
	}{
		{
			name: "general",
			input: &Link{
				Config: LinkConfig{
					Text: "nginx-deployment",
					Ref:  "/overview/deployments/nginx-deployment",
				},
			},
			expected: `
            {
                "metadata": {
                  "type": "link"
                },
                "config": {
                  "value": "nginx-deployment",
                  "ref": "/overview/deployments/nginx-deployment"
                }
            }
`,
		},
		{
			name: "with title",
			input: &Link{
				base: newBase(typeLink, TitleFromString("Name")),
				Config: LinkConfig{
					Text: "nginx-deployment",
					Ref:  "/overview/deployments/nginx-deployment",
				},
			},
			expected: `
            {
                "metadata": {
                  "type": "link",
                  "title": [
									  {
											"metadata": { "type": "text" },
											"config": { "value": "Name" }
										}
									]
                },
                "config": {
                  "value": "nginx-deployment",
                  "ref": "/overview/deployments/nginx-deployment"
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

func Test_Link_String(t *testing.T) {
	c := NewLink("title", "string", "/path")
	assert.Equal(t, "string", c.String())
}

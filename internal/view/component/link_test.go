package component

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Link_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    ViewComponent
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
				Metadata: Metadata{
					Title: []TitleViewComponent{NewText("Name")},
				},
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

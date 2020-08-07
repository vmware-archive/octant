/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

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
				Base: newBase(TypeLink, TitleFromString("Name")),
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
			isErr := err != nil
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

func Test_Link_LessThan(t *testing.T) {
	tests := []struct {
		name     string
		link     Link
		other    Component
		expected bool
	}{
		{
			name:     "is less",
			link:     *NewLink("", "a", "/a"),
			other:    NewLink("", "b", "/b"),
			expected: true,
		},
		{
			name:     "is not less",
			link:     *NewLink("", "b", "/b"),
			other:    NewLink("", "a", "/a"),
			expected: false,
		},
		{
			name:     "other is not link",
			link:     *NewLink("", "b", "/b"),
			other:    nil,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.link.LessThan(test.other)
			assert.Equal(t, test.expected, got)
		})
	}
}

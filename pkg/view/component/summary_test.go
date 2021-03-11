/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummary_Add(t *testing.T) {
	tc := []struct {
		name      string
		summary   *Summary
		additions SummarySections
		expected  []SummarySection
	}{
		{
			name:     "empty with no additions",
			summary:  NewSummary("title"),
			expected: nil,
		},
		{
			name:    "empty with additions",
			summary: NewSummary("title"),
			additions: SummarySections{
				{Header: "a", Content: NewText("a")},
			},
			expected: SummarySections{
				{Header: "a", Content: NewText("a")},
			},
		},
		{
			name: "existing with additions",
			summary: NewSummary("title", SummarySection{
				Header:  "a",
				Content: NewText("a"),
			}),
			additions: SummarySections{
				{Header: "b", Content: NewText("b")},
			},
			expected: SummarySections{
				{Header: "a", Content: NewText("a")},
				{Header: "b", Content: NewText("b")},
			},
		},
		{
			name: "existing with updates",
			summary: NewSummary("title", []SummarySection{
				{Header: "a", Content: NewText("a")},
				{Header: "b", Content: NewText("b")},
			}...),
			additions: SummarySections{
				{Header: "a", Content: NewText("updated")},
			},
			expected: SummarySections{
				{Header: "a", Content: NewText("updated")},
				{Header: "b", Content: NewText("b")},
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			tt.summary.Add(tt.additions...)
			got := tt.summary.Sections()

			require.Equal(t, tt.expected, got)
		})
	}
}

func Test_Summary_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        Component
		expectedPath string
		isErr        bool
	}{
		{
			name: "general",
			input: &Summary{
				Base: newBase(TypeSummary, TitleFromString("my summary")),
				Config: SummaryConfig{
					Alert: &Alert{
						Type:    AlertTypeInfo,
						Message: "info",
					},
					Sections: []SummarySection{
						{
							Header: "Containers",
							Content: &List{
								Base: newBase(TypeList, TitleFromString("nginx")),
								Config: ListConfig{
									Items: []Component{
										&Text{
											Base: newBase(TypeText, TitleFromString("Image")),
											Config: TextConfig{
												Text: "nginx:latest",
											},
										},
										&Text{
											Base: newBase(TypeText, TitleFromString("Port")),
											Config: TextConfig{
												Text: "80/TCP",
											},
										},
									},
								},
							},
						},
						{
							Header: "Empty Section",
							Content: &Text{
								Config: TextConfig{
									Text: "Nothing to see here",
								},
							},
						},
					},
				},
			},
			expectedPath: "summary.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := err != nil
			if isErr != tc.isErr {
				t.Fatalf("Unexepected error: %v", err)
			}

			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err, "reading test fixtures")
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

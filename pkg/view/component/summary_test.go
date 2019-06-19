/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				base: newBase(typeSummary, TitleFromString("mysummary")),
				Config: SummaryConfig{
					Sections: []SummarySection{
						SummarySection{
							Header: "Containers",
							Content: &List{
								base: newBase(typeList, TitleFromString("nginx")),
								Config: ListConfig{
									Items: []Component{
										&Text{
											base: newBase(typeText, TitleFromString("Image")),
											Config: TextConfig{
												Text: "nginx:latest",
											},
										},
										&Text{
											base: newBase(typeText, TitleFromString("Port")),
											Config: TextConfig{
												Text: "80/TCP",
											},
										},
									},
								},
							},
						},
						SummarySection{
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
			isErr := (err != nil)
			if isErr != tc.isErr {
				t.Fatalf("Unexepected error: %v", err)
			}

			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err, "reading test fixtures")
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

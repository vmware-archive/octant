/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

func TestAccordion_Add(t *testing.T) {
	accordion := NewAccordion("Accordion", nil)
	accordion.Add(AccordionRow{
		Title:   "row title",
		Content: NewText("row content"),
	})

	expected := []AccordionRow{
		{
			Title:   "row title",
			Content: NewText("row content"),
		},
	}

	require.Equal(t, accordion.Config.Rows, expected)
}

func TestAccordion_AllowMultipleExpanded(t *testing.T) {
	accordion := NewAccordion("accordion", nil)
	accordion.AllowMultipleExpanded()
	require.Equal(t, accordion.Config.AllowMultipleExpanded, true)
}

func TestAccordion_Marshal(t *testing.T) {
	cases := []struct {
		name         string
		input        *Accordion
		expectedPath string
		isErr        bool
	}{
		{
			name: "in general",
			input: &Accordion{
				Base: newBase(TypeAccordion, TitleFromString("my accordion")),
				Config: AccordionConfig{
					Rows: []AccordionRow{
						{
							Title:   "Row title",
							Content: NewText("row content"),
						},
					},
					AllowMultipleExpanded: true,
				},
			},
			expectedPath: "accordion.json",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := err != nil
			if isErr != tc.isErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err)
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

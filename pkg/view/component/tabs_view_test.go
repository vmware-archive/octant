/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/action"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

func Test_Tabs_Marshal(t *testing.T) {
	fl := NewFlexLayout("title")
	button := NewButton("test", action.Payload{"foo": "bar"})
	section := FlexLayoutSection{
		{
			Width: WidthFull,
			View:  button,
		},
	}
	fl.AddSections(section)

	test := []struct {
		name         string
		input        Component
		expectedPath string
		isErr        bool
	}{
		{
			name: "in general",
			input: &TabsView{
				Base: newBase(TypeTabsView, nil),
				Config: TabsViewConfig{
					Tabs: []SingleTab{
						{
							Name:     "title",
							Contents: *fl,
						},
						{
							Name:     "title 2",
							Contents: *fl,
						},
					},
				},
			},
			expectedPath: "tabs.json",
			isErr:        false,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := err != nil
			if isErr != tc.isErr {
				t.Fatalf("UnExpected error: %v", err)
			}

			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err)
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

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

func Test_Signpost_Marshal(t *testing.T) {
	test := []struct {
		name         string
		input        Component
		expectedPath string
		isErr        bool
	}{
		{
			name: "in general",
			input: &Signpost{
				Base: newBase(TypeSignpost, nil),
				Config: SignpostConfig{
					Trigger:  NewIcon("user"),
					Message:  "Message",
					Position: PositionTopLeft,
				},
			},
			expectedPath: "signpost.json",
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

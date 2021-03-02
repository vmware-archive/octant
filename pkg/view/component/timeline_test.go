/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path"
	"testing"
)

func Test_Timeline_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        Component
		expectedPath string
		isErr        bool
	}{
		{
			name: "in general",
			input: &Timeline{
				Base: newBase(TypeTimeline, nil),
				Config: TimelineConfig{
					Steps: []TimelineStep{
						{
							State:       TimelineStepSuccess,
							Title:       "Step 1",
							Header:      "success header",
							Description: "this is a success",
						},
						{
							State:       TimelineStepError,
							Title:       "Step 2",
							Header:      "error header",
							Description: "this is an error",
						},
						{
							State:       TimelineStepCurrent,
							Title:       "Step 3",
							Header:      "current header",
							Description: "this is a current step",
						},
						{
							State:       TimelineStepProcessing,
							Title:       "Step 4",
							Header:      "processing header",
							Description: "this is processing",
						},
						{
							State:       TimelineStepNotStarted,
							Title:       "Step 5",
							Header:      "not started header",
							Description: "this has not started",
						},
					},
					Vertical: false,
				},
			},
			expectedPath: "timeline.json",
		},
	}
	for _, tc := range tests {
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

func Test_Timeline_Add(t *testing.T) {
	step := TimelineStep{
		State:       TimelineStepCurrent,
		Title:       "Title",
		Header:      "Header",
		Description: "Description",
	}
	timeline := NewTimeline([]TimelineStep{}, true)
	timeline.Add(step)

	expected := []TimelineStep{
		step,
	}
	assert.Equal(t, expected, timeline.Config.Steps)
}

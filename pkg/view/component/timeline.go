/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"sync"
)

// Timeline is a component for timeline
// +octant:component
type Timeline struct {
	Base
	Config TimelineConfig `json:"config"`

	mu sync.Mutex
}

// TimelineConfig is the contents of Timeline
type TimelineConfig struct {
	Steps    []TimelineStep `json:"steps"`
	Vertical bool           `json:"vertical"`
}

// TimelineStep is the data for each timeline step
type TimelineStep struct {
	State       TimelineState `json:"state"`
	Header      string        `json:"header"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
}

// TimelineState is the state of a timeline step
type TimelineState string

const (
	TimelineStepNotStarted TimelineState = "not-started"
	TimelineStepCurrent    TimelineState = "current"
	TimelineStepProcessing TimelineState = "processing"
	TimelineStepSuccess    TimelineState = "success"
	TimelineStepError      TimelineState = "error"
)

// NewTimeline creates a timeline component
func NewTimeline(steps []TimelineStep, vertical bool) *Timeline {
	return &Timeline{
		Base: newBase(TypeTimeline, nil),
		Config: TimelineConfig{
			Steps:    steps,
			Vertical: vertical,
		},
	}
}

// Add adds an additional step to the timeline
func (t *Timeline) Add(steps ...TimelineStep) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Config.Steps = append(t.Config.Steps, steps...)
}

type timelineMarshal Timeline

func (t *Timeline) MarshalJSON() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	m := timelineMarshal{
		Base:   t.Base,
		Config: t.Config,
	}
	m.Metadata.Type = TypeTimeline
	return json.Marshal(&m)
}

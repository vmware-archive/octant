/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/util/json"
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
	ButtonGroup *ButtonGroup  `json:"buttonGroup,omitempty"`
}

func (t *TimelineStep) UnmarshalJSON(data []byte) error {
	x := struct {
		State       TimelineState `json:"state"`
		Header      string        `json:"header"`
		Title       string        `json:"title"`
		Description string        `json:"description"`
		ButtonGroup *TypedObject  `json:"buttonGroup,omitempty"`
	}{}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x.ButtonGroup != nil {
		component, err := x.ButtonGroup.ToComponent()
		if err != nil {
			return err
		}
		buttonGroup, ok := component.(*ButtonGroup)
		if !ok {
			return errors.New("item was not a buttonGroup")
		}
		t.ButtonGroup = buttonGroup
	}
	t.State = x.State
	t.Title = x.Title
	t.Header = x.Header
	t.Description = x.Description

	return nil
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

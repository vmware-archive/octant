package component

import (
	"encoding/json"
	"sync"
)

type Timeline struct {
	Base
	Config TimelineConfig `json:"config"`

	mu sync.Mutex
}

type TimelineConfig struct {
	Steps    []TimelineStep `json:"steps"`
	Vertical bool           `json:"vertical"`
}

type TimelineStep struct {
	State       TimelineState `json:"state"`
	Header      string        `json:"header"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
}

type TimelineState string

const (
	TimelineStepNotStarted TimelineState = "not-started"
	TimelineStepCurrent    TimelineState = "current"
	TimelineStepProcessing TimelineState = "processing"
	TimelineStepSuccess    TimelineState = "success"
	TimelineStepError      TimelineState = "error"
)

func NewTimeline(title string, steps []TimelineStep, vertical bool) *Timeline {
	return &Timeline{
		Base: newBase(TypeTimeline, TitleFromString(title)),
		Config: TimelineConfig{
			Steps:    steps,
			Vertical: vertical,
		},
	}
}

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

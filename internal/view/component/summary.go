package component

import "encoding/json"

// Summary contains other ViewComponents
type Summary struct {
	Metadata Metadata      `json:"metadata"`
	Config   SummaryConfig `json:"config"`
}

// SummaryConfig is the contents of a Summary
type SummaryConfig struct {
	Sections []SummarySection `json:"sections"`
}

// SummarySection is a section within a summary
type SummarySection struct {
	Header  string        `json:"header"`
	Content ViewComponent `json:"content"`
}

// NewSummary creates a summary component
func NewSummary(title string, sections ...SummarySection) *Summary {
	s := append([]SummarySection(nil), sections...) // Make a copy
	return &Summary{
		Metadata: Metadata{
			Type:  "summary",
			Title: title,
		},
		Config: SummaryConfig{
			Sections: s,
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Summary) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifes whether the component is considered empty. Implements ViewComponent.
func (t *Summary) IsEmpty() bool {
	return len(t.Config.Sections) == 0
}

// Add adds additional items to the tail of the summary.
func (t *Summary) Add(sections ...SummarySection) {
	t.Config.Sections = append(t.Config.Sections, sections...)
}

type summaryMarshal Summary

// MarshalJSON implements json.Marshaler
func (t *Summary) MarshalJSON() ([]byte, error) {
	m := summaryMarshal(*t)
	m.Metadata.Type = "summary"
	return json.Marshal(&m)
}

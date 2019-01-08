package component

import (
	"encoding/json"
	"time"
)

// Timestamp is a component representing a point in time
type Timestamp struct {
	Metadata Metadata        `json:"metadata"`
	Config   TimestampConfig `json:"config"`
}

// TimestampConfig is the contents of Timestamp
type TimestampConfig struct {
	Timestamp int64 `json:"timestamp"`
}

// NewTimestamp creates a timestamp component
func NewTimestamp(t time.Time) *Timestamp {
	return &Timestamp{
		Metadata: Metadata{
			Type: "timestamp",
		},
		Config: TimestampConfig{
			Timestamp: t.Unix(),
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Timestamp) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifes whether the component is considered empty. Implements ViewComponent.
func (t *Timestamp) IsEmpty() bool {
	return t.Config.Timestamp == time.Time{}.Unix()
}

type timestampMarshal Timestamp

// MarshalJSON implements json.Marshaler
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	m := timestampMarshal(*t)
	m.Metadata.Type = "timestamp"
	return json.Marshal(&m)
}

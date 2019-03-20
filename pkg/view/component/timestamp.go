package component

import (
	"encoding/json"
	"time"
)

// Timestamp is a component representing a point in time
type Timestamp struct {
	base
	Config TimestampConfig `json:"config"`
}

// TimestampConfig is the contents of Timestamp
type TimestampConfig struct {
	Timestamp int64 `json:"timestamp"`
}

// NewTimestamp creates a timestamp component
func NewTimestamp(t time.Time) *Timestamp {
	return &Timestamp{
		base: newBase(typeTimestamp, nil),
		Config: TimestampConfig{
			Timestamp: t.Unix(),
		},
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (t *Timestamp) GetMetadata() Metadata {
	return t.Metadata
}

type timestampMarshal Timestamp

// MarshalJSON implements json.Marshaler
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	m := timestampMarshal(*t)
	m.Metadata.Type = typeTimestamp
	return json.Marshal(&m)
}

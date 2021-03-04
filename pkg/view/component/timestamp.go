/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"time"
)

// Timestamp is a component representing a point in time
//
// +octant:component
type Timestamp struct {
	Base
	Config TimestampConfig `json:"config"`
}

var _ Component = (*Timestamp)(nil)

// TimestampConfig is the contents of Timestamp
type TimestampConfig struct {
	Timestamp int64 `json:"timestamp"`
}

// NewTimestamp creates a timestamp component
func NewTimestamp(t time.Time) *Timestamp {
	return &Timestamp{
		Base: newBase(TypeTimestamp, nil),
		Config: TimestampConfig{
			Timestamp: t.Unix(),
		},
	}
}

type timestampMarshal Timestamp

// MarshalJSON implements json.Marshaler
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	m := timestampMarshal(*t)
	m.Metadata.Type = TypeTimestamp
	return json.Marshal(&m)
}

// LessThan returns true if this component's value is less than the argument supplied.
func (t *Timestamp) LessThan(i interface{}) bool {
	v, ok := i.(*Timestamp)
	if !ok {
		return false
	}

	return t.Config.Timestamp < v.Config.Timestamp

}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// Annotations is a component representing key/value based annotations
type Annotations struct {
	base
	Config AnnotationsConfig `json:"config"`
}

// AnnotationsConfig is the contents of Annotations
type AnnotationsConfig struct {
	Annotations map[string]string `json:"annotations"`
}

// NewAnnotations creates a annotations component
func NewAnnotations(annotations map[string]string) *Annotations {
	return &Annotations{
		base: newBase(typeAnnotations, nil),
		Config: AnnotationsConfig{
			Annotations: annotations,
		},
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (t *Annotations) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifies whether the component is considered empty. Implements Component.
func (t *Annotations) IsEmpty() bool {
	return len(t.Config.Annotations) == 0
}

type annotationsMarshal Annotations

// MarshalJSON implements json.Marshaler.
func (t *Annotations) MarshalJSON() ([]byte, error) {
	m := annotationsMarshal(*t)
	m.Metadata.Type = "annotations"
	m.Metadata.Title = t.Metadata.Title
	return json.Marshal(&m)
}

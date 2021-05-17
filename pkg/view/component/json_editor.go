/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"github.com/vmware-tanzu/octant/internal/util/json"
)

type JSONEditorMode string

const (
	// ViewMode can show JSON up to 500MiB
	ViewMode JSONEditorMode = "view"
	TextMode JSONEditorMode = "text"
)

// JSONEditor is an JSON editor component.
//
// +octant:component
type JSONEditor struct {
	Base
	Config JSONEditorConfig `json:"config"`
}

type JSONEditorConfig struct {
	Mode    JSONEditorMode `json:"mode"`
	Content string         `json:"content"`
}

// NewJSONEditor creates an instance of JSONEditor
func NewJSONEditor(content string) *JSONEditor {
	return &JSONEditor{
		Base: newBase(TypeJSONEditor, nil),
		Config: JSONEditorConfig{
			Mode:    ViewMode,
			Content: content,
		},
	}
}

// GetMetadata returns the component's metadata.
func (j *JSONEditor) GetMetadata() Metadata {
	return j.Metadata
}

type jsonEditorMarshal JSONEditor

// MarshalJSON implements json.Marshaler
func (j *JSONEditor) MarshalJSON() ([]byte, error) {
	m := jsonEditorMarshal(*j)
	m.Metadata.Type = TypeJSONEditor
	return json.Marshal(&m)
}

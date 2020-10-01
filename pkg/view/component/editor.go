/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// Editor is an editor component.
//
// +octant:component
type Editor struct {
	Base
	Config EditorConfig `json:"config"`
}

// EditorConfig is configuration for Editor.
type EditorConfig struct {
	Value        string            `json:"value"`
	Language     string            `json:"language"`
	ReadOnly     bool              `json:"readOnly"`
	Metadata     map[string]string `json:"metadata"`
	SubmitAction string            `json:"submitAction,omitempty"`
	SubmitLabel  string            `json:"submitLabel,omitempty"`
}

// NewEditor creates an instance of an editor component.
func NewEditor(title []TitleComponent, value string, readOnly bool) *Editor {
	return &Editor{
		Base: newBase(TypeEditor, title),
		Config: EditorConfig{
			Value:    value,
			ReadOnly: readOnly,
		},
	}
}

func (e *Editor) SetValueFromObject(object runtime.Object) error {
	s, err := kubernetes.SerializeToString(object)
	if err != nil {
		return fmt.Errorf("serialize object: %w", err)
	}

	e.Config.Value = s

	key, err := store.KeyFromObject(object)
	if err != nil {
		return fmt.Errorf("create key from object: %w", err)
	}

	e.Config.Metadata = map[string]string{
		"namespace":  key.Namespace,
		"apiVersion": key.APIVersion,
		"kind":       key.Kind,
		"name":       key.Name,
	}

	return nil
}

// GetMetadata returns the component's metadata.
func (e *Editor) GetMetadata() Metadata {
	return e.Metadata
}

type editorMarshal Editor

// MarshalJSON implements json.Marshaler
func (e *Editor) MarshalJSON() ([]byte, error) {
	m := editorMarshal(*e)
	m.Metadata.Type = TypeEditor
	return json.Marshal(&m)
}

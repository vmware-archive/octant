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

// Value is a component for code
type Editor struct {
	base
	Config EditorConfig `json:"config"`
}

// CodeConfig is the contents of Value
type EditorConfig struct {
	Value    string            `json:"value"`
	ReadOnly bool              `json:"readOnly"`
	Metadata map[string]string `json:"metadata"`
}

// NewCodeBlock creates a code component
func NewEditor(title []TitleComponent, value string, readOnly bool) *Editor {
	return &Editor{
		base: newBase(typeEditor, title),
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
	m.Metadata.Type = typeEditor
	return json.Marshal(&m)
}

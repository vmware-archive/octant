/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"fmt"
)

// Error is a component for freetext
type Error struct {
	base
	Config ErrorConfig `json:"config"`
}

// ErrorConfig is the contents of Text
type ErrorConfig struct {
	Data string `json:"data,omitempty"`
}

// NewError creates a text component
func NewError(title []TitleComponent, err error) *Error {
	return &Error{
		base: newBase(typeError, title),
		Config: ErrorConfig{
			Data: fmt.Sprintf("%+v", err),
		},
	}
}

// SupportsTitle denotes this is a TextComponent.
func (t *Error) SupportsTitle() {}

type errorMarshal Error

// MarshalJSON implements json.Marshaler
func (t *Error) MarshalJSON() ([]byte, error) {
	m := errorMarshal(*t)
	m.Metadata.Type = typeError
	return json.Marshal(&m)
}

// String returns the text content of the component.
func (t *Error) String() string {
	return t.Config.Data
}

// LessThan returns true if this component's value is less than the argument supplied.
func (t *Error) LessThan(i interface{}) bool {
	v, ok := i.(*Error)
	if !ok {
		return false
	}

	return t.Config.Data < v.Config.Data

}

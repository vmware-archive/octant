/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "github.com/vmware-tanzu/octant/internal/util/json"

// Value is a component for code
//
// +octant:component
type Code struct {
	Base
	Config CodeConfig `json:"config"`
}

// CodeConfig is the contents of Value
type CodeConfig struct {
	Code string `json:"value"`
}

// NewCodeBlock creates a code component
func NewCodeBlock(s string) *Code {
	return &Code{
		Base: newBase(TypeCode, nil),
		Config: CodeConfig{
			Code: s,
		},
	}
}

type codeMarshal Code

// MarshalJSON implements json.Marshaler
func (c *Code) MarshalJSON() ([]byte, error) {
	m := codeMarshal(*c)
	m.Metadata.Type = TypeCode
	return json.Marshal(&m)
}

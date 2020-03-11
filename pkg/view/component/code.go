/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// Value is a component for code
type Code struct {
	base
	Config CodeConfig `json:"config"`
}

// CodeConfig is the contents of Value
type CodeConfig struct {
	Code string `json:"value"`
}

// NewCodeBlock creates a code component
func NewCodeBlock(s string) *Code {
	return &Code{
		base: newBase(typeCodeBlock, nil),
		Config: CodeConfig{
			Code: s,
		},
	}
}

type codeMarshal Code

// MarshalJSON implements json.Marshaler
func (c *Code) MarshalJSON() ([]byte, error) {
	m := codeMarshal(*c)
	m.Metadata.Type = typeCodeBlock
	return json.Marshal(&m)
}

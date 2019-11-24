/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// IFrame is a component for displaying content in an iframe
type IFrame struct {
	base
	Config IFrameConfig `json:"config"`
}

// IFrameConfig is the title and url of the iframe
type IFrameConfig struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

// NewIFrame creates an iframe component
func NewIFrame(url string, title string) *IFrame {
	return &IFrame{
		base: newBase(typeText, nil),
		Config: IFrameConfig{
			Url:   url,
			Title: title,
		},
	}
}

type IFrameMarshal IFrame

// MarshalJSON implements json.Marshaler
func (t *IFrame) MarshalJSON() ([]byte, error) {
	m := IFrameMarshal(*t)
	m.Metadata.Type = typeIFrame
	return json.Marshal(&m)
}

// String returns the url content of the component.
func (t *IFrame) String() string {
	return t.Config.Url
}

// LessThan returns true if this component's url is lexically less than the argument supplied.
func (t *IFrame) LessThan(i interface{}) bool {
	v, ok := i.(*IFrame)
	if !ok {
		return false
	}

	return t.Config.Url < v.Config.Url

}

/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"github.com/vmware-tanzu/octant/internal/util/json"
)

// Signpost is a component for a signpost
// +octant:component
type Signpost struct {
	Base
	Config SignpostConfig `json:"config"`
}

// SignpostConfig is the contents of a signpost
type SignpostConfig struct {
	// Trigger is the component that will trigger the signpost.
	Trigger Component `json:"trigger"`
	// Message is the text that will be displayed.
	Message string `json:"message"`
	// Position of Signpost
	Position Position `json:"position"`
}

type Position string

const (
	PositionTopLeft      Position = "top-left"
	PositionTopMiddle    Position = "top-middle"
	PositionTopRight     Position = "top-right"
	PositionRightTop     Position = "right-top"
	PositionRightMiddle  Position = "right-middle"
	PositionRightBottom  Position = "right-bottom"
	PositionBottomRight  Position = "bottom-right"
	PositionBottomMiddle Position = "bottom-middle"
	PositionBottomLeft   Position = "bottom-left"
	PositionLeftBottom   Position = "left-bottom"
	PositionLeftMiddle   Position = "left-middle"
	PositionLeftTop      Position = "left-top"
)

func NewSignpost(t Component, m string) *Signpost {
	so := &Signpost{
		Base: newBase(TypeSignpost, nil),
		Config: SignpostConfig{
			Trigger: t,
			Message: m,
		},
	}

	return so
}

// SetStatus sets the status of the text component.
func (t *Signpost) SetPosition(position Position) {
	t.Config.Position = position
}

type signpostMarshal Signpost

// MarshalJSON implements json.Marshaller
func (t *Signpost) MarshalJSON() ([]byte, error) {
	m := signpostMarshal(*t)
	m.Metadata.Type = TypeSignpost
	return json.Marshal(&m)
}

func (c *SignpostConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Trigger  *TypedObject `json:"trigger"`
		Message  string       `json:"message,omitempty"`
		Position Position     `json:"position"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x.Trigger != nil {
		trigger, err := x.Trigger.ToComponent()
		if err != nil {
			return err
		}
		c.Trigger = trigger
	}

	c.Message = x.Message
	c.Position = x.Position

	return nil
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package flexlayout

import (
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	// maxWidth is the maximum width of an item
	maxWidth = 24
)

type SectionMember struct {
	View  component.Component
	Width int
}

type Section struct {
	Members []SectionMember
}

func NewSection() *Section {
	return &Section{}
}

func (s *Section) Add(view component.Component, width int) error {
	if width > maxWidth {
		return errors.Errorf("component width %d; max width %d", width, maxWidth)
	}
	member := SectionMember{View: view, Width: width}
	s.Members = append(s.Members, member)

	return nil
}

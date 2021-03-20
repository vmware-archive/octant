/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"sync"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

// Accordion is a component for accordion
//
// +octant:component
type Accordion struct {
	Base
	Config AccordionConfig `json:"config"`

	mu sync.Mutex
}

// Accordion is the contents of Accordion
type AccordionConfig struct {
	Rows                  []AccordionRow `json:"rows"`
	AllowMultipleExpanded bool           `json:"allowMultipleExpanded"`
}

// AccordionRow is a row of an accordion
type AccordionRow struct {
	Title   string    `json:"title"`
	Content Component `json:"content"`
}

// Unmarshal unmarshals an accordion config from JSON.
func (a *AccordionRow) UnmarshalJSON(data []byte) error {
	x := struct {
		Title   string       `json:"title"`
		Content *TypedObject `json:"content"`
	}{}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x.Content != nil {
		var err error
		a.Content, err = x.Content.ToComponent()
		if err != nil {
			return err
		}
	}
	a.Title = x.Title
	return nil
}

// NewAccordion creates a new accordion component
func NewAccordion(title string, rows []AccordionRow, options ...func(*Accordion)) *Accordion {
	a := &Accordion{
		Base: newBase(TypeAccordion, TitleFromString(title)),
		Config: AccordionConfig{
			Rows: rows,
		},
	}

	for _, option := range options {
		option(a)
	}
	return a
}

// AllowMultipleExpanded sets an accordion to allow expanding multiple rows.
func (a *Accordion) AllowMultipleExpanded() {
	a.Config.AllowMultipleExpanded = true
}

// AddRow adds one or more rows to an accordion.
func (a *Accordion) Add(row ...AccordionRow) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Config.Rows = append(a.Config.Rows, row...)
}

type accordionMarshal Accordion

// MarshalJSON implements json.Marshaler
func (a *Accordion) MarshalJSON() ([]byte, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	m := accordionMarshal{
		Base:   a.Base,
		Config: a.Config,
	}

	m.Metadata.Type = TypeAccordion
	return json.Marshal(&m)
}

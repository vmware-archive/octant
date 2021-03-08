/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package component

type SingleStatValue struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}

type SingleStatConfig struct {
	Title string          `json:"title"`
	Value SingleStatValue `json:"value"`
}

// Single stat shows a single statistic.
//
// +octant:component
type SingleStat struct {
	Base
	Config SingleStatConfig `json:"config"`
}

var _ Component = (*SingleStat)(nil)

func NewSingleStat(title, valueText, color string) *SingleStat {
	return &SingleStat{
		Base: newBase(TypeSingleStat, nil),
		Config: SingleStatConfig{
			Title: title,
			Value: SingleStatValue{
				Text:  valueText,
				Color: color,
			},
		},
	}
}

type singleStatMarshal SingleStat

func (ss *SingleStat) MarshalJSON() ([]byte, error) {
	m := singleStatMarshal(*ss)
	m.Metadata.Type = TypeSingleStat
	return json.Marshal(&m)
}

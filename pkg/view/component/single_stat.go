/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package component

import "encoding/json"

type SingleStatValue struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}

type SingleStateConfig struct {
	Title string          `json:"title"`
	Value SingleStatValue `json:"value"`
}

type SingleStat struct {
	base
	Config SingleStateConfig `json:"config"`
}

var _ Component = (*SingleStat)(nil)

func NewSingleStat(title, valueText, color string) *SingleStat {
	return &SingleStat{
		base: newBase(typeSingleStat, nil),
		Config: SingleStateConfig{
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
	m.Metadata.Type = typeSingleStat
	return json.Marshal(&m)
}

/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package component

import "encoding/json"

type ChartColor string

const (
	ChartColorError   ChartColor = "#f52f22"
	ChartColorWarning ChartColor = "#fac400"
	ChartColorOK      ChartColor = "#60b515"
)

type BulletBand struct {
	Min   int        `json:"min"`
	Max   int        `json:"max"`
	Color ChartColor `json:"color"`
	Label string     `json:"label"`
}

type VerticalBulletChartConfig struct {
	Bands        []BulletBand `json:"bands"`
	Measure      int          `json:"measure"`
	MeasureLabel string       `json:"measureLabel"`
	Label        string       `json:"label"`
}

type VerticalBulletChart struct {
	base
	Config VerticalBulletChartConfig `json:"config"`
}

var _ Component = (*VerticalBulletChart)(nil)

func NewVerticalBulletChart(label string) *VerticalBulletChart {
	return &VerticalBulletChart{
		base: newBase(typeVerticalBulletChart, nil),
		Config: VerticalBulletChartConfig{
			Label: label,
		},
	}
}

func (vbc *VerticalBulletChart) SetBands(bands []BulletBand) {
	vbc.Config.Bands = bands
}

func (vbc *VerticalBulletChart) SetMeasure(label string, val int) {
	vbc.Config.MeasureLabel = label
	vbc.Config.Measure = val
}

type verticalBulletChartMarshal VerticalBulletChart

func (vbc *VerticalBulletChart) MarshalJSON() ([]byte, error) {
	m := verticalBulletChartMarshal(*vbc)
	m.Metadata.Type = typeVerticalBulletChart
	return json.Marshal(&m)
}

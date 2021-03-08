/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package component

type DonutChartSize int

const (
	DonutChartSizeSmall  DonutChartSize = 50
	DonutChartSizeMedium DonutChartSize = 100
)

type DonutChartLabels struct {
	Plural   string `json:"plural"`
	Singular string `json:"singular"`
}

type DonutSegment struct {
	Count       int        `json:"count"`
	Status      NodeStatus `json:"status"`
	Color       string     `json:"color,omitempty"`
	Description string     `json:"description,omitempty"`
	Thickness   int        `json:"thickness,omitempty"`
}

type DonutChartConfig struct {
	Segments  []DonutSegment   `json:"segments"`
	Labels    DonutChartLabels `json:"labels"`
	Size      DonutChartSize   `json:"size"`
	Thickness int              `json:"thickness,omitempty"`
}

// +octant:component
type DonutChart struct {
	Base
	Config DonutChartConfig `json:"config"`
}

var _ Component = (*DonutChart)(nil)

func NewDonutChart() *DonutChart {
	dc := &DonutChart{
		Base: newBase(TypeDonutChart, nil),
	}

	return dc
}

type donutChartMarshal DonutChart

func (dc *DonutChart) SetSegments(segments []DonutSegment) {
	dc.Config.Segments = segments
}

func (dc *DonutChart) SetLabels(plural string, singular string) {
	dc.Config.Labels = DonutChartLabels{
		Plural:   plural,
		Singular: singular,
	}
}

func (dc *DonutChart) SetSize(size DonutChartSize) {
	dc.Config.Size = size
}

// Set donut chart thickness - trimmed to be inside [2-100] interval
// where 2 is barely visible and 100 turns it to a pie chart
func (dc *DonutChart) SetThickness(thickness int) {
	dc.Config.Thickness = thickness
}

func (dc *DonutChart) MarshalJSON() ([]byte, error) {
	m := donutChartMarshal(*dc)
	m.Metadata.Type = TypeDonutChart
	return json.Marshal(&m)
}

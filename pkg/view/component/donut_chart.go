/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package component

import "encoding/json"

type DonutChartLabels struct {
	Plural   string `json:"plural"`
	Singular string `json:"singular"`
}

type DonutSegment struct {
	Count  int        `json:"count"`
	Status NodeStatus `json:"status"`
}

type DonutChartConfig struct {
	Segments []DonutSegment   `json:"segments"`
	Labels   DonutChartLabels `json:"labels"`
}

type DonutChart struct {
	base
	Config DonutChartConfig `json:"config"`
}

var _ Component = (*DonutChart)(nil)

func NewDonutChart() *DonutChart {
	dc := &DonutChart{
		base: newBase(typeDonutChart, nil),
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

func (dc *DonutChart) MarshalJSON() ([]byte, error) {
	m := donutChartMarshal(*dc)
	m.Metadata.Type = typeDonutChart
	return json.Marshal(&m)
}

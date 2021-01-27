/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_Donut_Legacy_Chart(t *testing.T) {
	donut := component.NewDonutChart()
	donut.SetSegments([]component.DonutSegment{
		{
			Count:  3,
			Status: "ok",
		},
		{
			Count:  1,
			Status: "error",
		},
	})
	donut.SetLabels("pods", "pod")
	donut.SetSize(component.DonutChartSizeMedium)

	got, err := donut.MarshalJSON()
	require.NoError(t, err)

	expected, err := ioutil.ReadFile(filepath.Join("testdata", "donut.json"))
	require.NoError(t, err)

	assert.JSONEq(t, string(expected), string(got))
}

func Test_Donut_Chart_With_Colors(t *testing.T) {
	donut := component.NewDonutChart()
	donut.SetSegments([]component.DonutSegment{
		{
			Count:       1,
			Status:      "critical",
			Color:       "purple",
			Description: "Critical vulnerability issue",
		},
		{
			Count:       3,
			Status:      "high",
			Color:       "red",
			Description: "High vulnerability issue",
		},
		{
			Count:       7,
			Status:      "medium",
			Color:       "orange",
			Description: "Medium vulnerability issue",
		},
		{
			Count:       11,
			Status:      "low",
			Color:       "green",
			Description: "Low vulnerability issue",
		},
	})
	donut.SetLabels("Vulnerabilities", "Vulnerability")
	donut.SetSize(component.DonutChartSizeMedium)

	got, err := donut.MarshalJSON()
	require.NoError(t, err)

	expected, err := ioutil.ReadFile(filepath.Join("testdata", "donut_with_colors.json"))
	require.NoError(t, err)

	assert.JSONEq(t, string(expected), string(got))
}

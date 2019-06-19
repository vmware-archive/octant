/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
	"github.com/vmware/octant/pkg/view/flexlayout"
)

func Test_Metadata(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	fl := flexlayout.New()

	deployment := testutil.CreateDeployment("deployment")
	metadata, err := NewMetadata(deployment, tpo.link)
	require.NoError(t, err)

	require.NoError(t, metadata.AddToFlexLayout(fl))

	got := fl.ToComponent("Summary")

	expected := component.NewFlexLayout("Summary")
	expected.AddSections([]component.FlexLayoutSection{
		{
			{
				Width: component.WidthFull,
				View: component.NewSummary("Metadata", component.SummarySections{
					{
						Header:  "Age",
						Content: component.NewTimestamp(deployment.CreationTimestamp.Time),
					},
				}...),
			},
		},
	}...)

	assert.Equal(t, expected, got)
}

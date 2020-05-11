/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package applications_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/modules/applications"
	"github.com/vmware-tanzu/octant/internal/modules/applications/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_homeDescriber_Describe(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	table := component.NewTable("table", "table", component.NewTableCols("col"))

	s := fake.NewMockSummarizer(controller)
	s.EXPECT().
		Summarize(gomock.Any(), "default", gomock.Any()).
		Return(table, nil)

	dashConfig := configFake.NewMockDash(controller)

	d := applications.NewHomeDescriber(applications.WithHomeDescriberSummarizer(s))

	ctx := context.Background()
	options := describer.Options{
		Dash: dashConfig,
	}
	actual, err := d.Describe(ctx, "default", options)
	require.NoError(t, err)

	expected := component.ContentResponse{
		Title:      component.TitleFromString("Applications"),
		Components: []component.Component{table},
	}
	require.Equal(t, expected, actual)
}

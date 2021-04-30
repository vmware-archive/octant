/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/vmware-tanzu/octant/pkg/view/flexlayout"
)

func Test_Metadata(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	fl := flexlayout.New()

	deployment := testutil.CreateDeployment("deployment")
	metadata, err := NewMetadata(deployment, tpo.link)
	require.NoError(t, err)

	fieldEntry := metav1.ManagedFieldsEntry{
		Manager:    "octant",
		Operation:  metav1.ManagedFieldsOperationUpdate,
		Time:       testutil.CreateTimestamp(),
		FieldsType: "FieldsV1",
		FieldsV1: &metav1.FieldsV1{
			Raw: []byte(`{"hello": "world"}`),
		},
	}
	deployment.ManagedFields = []metav1.ManagedFieldsEntry{
		fieldEntry,
	}

	require.NoError(t, metadata.AddToFlexLayout(fl))

	got := fl.ToComponent("Summary")

	fieldJSONData, err := convertFieldsToFormattedString(fieldEntry.FieldsV1)
	require.NoError(t, err)

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
			{
				Width: component.WidthFull,
				View: component.NewTableWithRows("Managed Fields", "There are no managed fields!", component.NewTableCols("Manager", "Operation", "Time", "Fields"), []component.TableRow{
					{
						"Manager":   component.NewText("octant"),
						"Operation": component.NewText(string(metav1.ManagedFieldsOperationUpdate)),
						"Time":      component.NewTimestamp(testutil.Time()),
						"Fields":    component.NewCodeBlock(fieldJSONData),
					},
				}),
			},
		},
	}...)

	component.AssertEqual(t, expected, got)
}

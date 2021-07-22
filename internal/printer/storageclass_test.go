/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_StorageClassListHandler(t *testing.T) {
	cols := component.NewTableCols("Name", "Provisioner", "Age")
	now := testutil.Time()

	object := testutil.CreateStorageClass("sc")

	object.CreationTimestamp = metav1.Time{Time: now}

	list := &storagev1.StorageClassList{
		Items: []storagev1.StorageClass{*object},
	}

	cases := []struct {
		name     string
		list     *storagev1.StorageClassList
		expected *component.Table
		isErr    bool
	}{
		{
			name: "in general",
			list: list,
			expected: component.NewTableWithRows("Storage Class", "We couldn't find any storage class!", cols,
				[]component.TableRow{
					{
						"Name": component.NewLink("", "sc", "/sc",
							genObjectStatus(component.TextStatusOK, []string{
								"storage.k8s.io/v1 StorageClass is OK",
							})),
						"Provisioner": component.NewText("manual"),
						"Age":         component.NewTimestamp(now),
						component.GridActionKey: gridActionsFactory([]component.GridAction{
							buildObjectDeleteAction(t, object),
						}),
					},
				}),
		},
		{
			name:  "list is nil",
			list:  nil,
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			ctx := context.Background()

			if tc.list != nil {
				tpo.PathForObject(&tc.list.Items[0], tc.list.Items[0].Name, "/"+tc.list.Items[0].Name)

				tpo.objectStore.EXPECT().Get(ctx, store.Key{
					Namespace:  object.Namespace,
					APIVersion: object.APIVersion,
					Kind:       object.Kind,
					Name:       object.Name,
				}).Return(testutil.ToUnstructured(t, object), nil).AnyTimes()
				tpo.pluginManager.EXPECT().ObjectStatus(ctx, object)
			}
			got, err := StorageClassListHandler(ctx, tc.list, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			component.AssertEqual(t, tc.expected, got)
		})
	}
}

func Test_StorageClassConfiguration(t *testing.T) {
	storageClass := testutil.CreateStorageClass("storageClass")

	cases := []struct {
		name         string
		storageClass *storagev1.StorageClass
		expected     component.Component
		isErr        bool
	}{
		{
			name:         "local",
			storageClass: storageClass,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Provisioner",
					Content: component.NewText("manual"),
				},
			}...),
		},
		{
			name:         "nil storageClass",
			storageClass: nil,
			isErr:        true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			sc := NewStorageClassConfiguration(tc.storageClass)

			summary, err := sc.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createStorageClassParametersView(t *testing.T) {
	now := testutil.Time()

	parameters := map[string]string{
		"Key": "Value",
	}

	sc := testutil.CreateStorageClass("storageClass")
	sc.CreationTimestamp = metav1.Time{Time: now}
	sc.Parameters = parameters

	observed, err := createStorageClassParameterView(sc)
	require.NoError(t, err)

	columns := component.NewTableCols("Key", "Value")
	expected := component.NewTable("Parameters", "There are no parameters!", columns)

	row := component.TableRow{}
	row["Key"] = component.NewText("Key")
	row["Value"] = component.NewText("Value")

	expected.Add(row)

	component.AssertEqual(t, expected, observed)
}

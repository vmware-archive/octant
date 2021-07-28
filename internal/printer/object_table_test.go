/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package printer

import (
	"context"
	"fmt"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/plugin"
	pluginFake "github.com/vmware-tanzu/octant/pkg/plugin/fake"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestObjectTable(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	cols := component.NewTableCols("A", "B")

	genDeleteGA := func(object runtime.Object) *component.GridActions {
		ga := component.NewGridActions()
		ga.AddGridAction(buildObjectDeleteAction(t, object))
		return ga
	}

	pod1 := testutil.CreatePod("pod1")
	pod1A := component.NewLink("", "pod1", "/pod1", func(l *component.Link) {
		list := component.NewList(nil, []component.Component{
			component.NewText("Pod may require additional action"),
		})
		l.SetStatus(component.TextStatusWarning, list)
	})
	pod2 := testutil.CreatePod("pod2")
	pod2A := component.NewLink("", "pod2", "/pod2", func(l *component.Link) {
		list := component.NewList(nil, []component.Component{
			component.NewText("Pod may require additional action"),
		})
		l.SetStatus(component.TextStatusWarning, list)
	})

	pluginManager := pluginFake.NewMockManagerInterface(controller)

	tests := []struct {
		name                 string
		mutateFn             func(*ObjectTable)
		wanted               func() *component.Table
		enablePluginResponse bool
	}{
		{
			name: "no mutations",
			mutateFn: func(table *ObjectTable) {

			},
			wanted: func() *component.Table {

				return component.NewTableWithRows("table", "placeholder", cols, []component.TableRow{
					{
						"A":                     pod1A,
						"B":                     component.NewText("0"),
						component.GridActionKey: genDeleteGA(pod1),
					},
					{
						"A":                     pod2A,
						"B":                     component.NewText("1"),
						component.GridActionKey: genDeleteGA(pod2),
					},
				})
			},
		},
		{
			name: "set sort order",
			mutateFn: func(table *ObjectTable) {
				table.SetSortOrder("A", true)
			},
			wanted: func() *component.Table {
				return component.NewTableWithRows("table", "placeholder", cols, []component.TableRow{
					{
						"A":                     pod2A,
						"B":                     component.NewText("1"),
						component.GridActionKey: genDeleteGA(pod2),
					},
					{
						"A":                     pod1A,
						"B":                     component.NewText("0"),
						component.GridActionKey: genDeleteGA(pod1),
					},
				})
			},
		},
		{
			name: "add column filters",
			mutateFn: func(table *ObjectTable) {
				table.AddFilters(map[string]component.TableFilter{
					"A": {
						Values:   []string{"pod1", "pod2"},
						Selected: []string{"pod1"},
					},
				})
			},
			wanted: func() *component.Table {
				table := component.NewTableWithRows("table", "placeholder", cols, []component.TableRow{
					{
						"A":                     pod1A,
						"B":                     component.NewText("0"),
						component.GridActionKey: genDeleteGA(pod1),
					},
					{
						"A":                     pod2A,
						"B":                     component.NewText("1"),
						component.GridActionKey: genDeleteGA(pod2),
					},
				})

				table.AddFilter("A", component.TableFilter{
					Values:   []string{"pod1", "pod2"},
					Selected: []string{"pod1"},
				})

				return table
			},
		},
		{
			name:                 "additional plugin status",
			enablePluginResponse: true,
			mutateFn:             func(table *ObjectTable) {},
			wanted: func() *component.Table {
				pod1Status := component.NewLink("", "pod1", "/pod1", func(l *component.Link) {
					list := component.NewList(nil, []component.Component{
						component.NewText("Pod may require additional action"),
						component.NewText("detail"),
					})
					l.SetStatus(component.TextStatusError, list)
				})
				pod2Status := component.NewLink("", "pod2", "/pod2", func(l *component.Link) {
					list := component.NewList(nil, []component.Component{
						component.NewText("Pod may require additional action"),
						component.NewText("detail"),
					})
					l.SetStatus(component.TextStatusError, list)
				})
				table := component.NewTableWithRows("table", "placeholder", cols, []component.TableRow{
					{
						"A":                     pod1Status,
						"B":                     component.NewText("0"),
						component.GridActionKey: genDeleteGA(pod1),
					},
					{
						"A":                     pod2Status,
						"B":                     component.NewText("1"),
						component.GridActionKey: genDeleteGA(pod2),
					},
				})
				return table
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			objectStore := fake.NewMockStore(ctrl)

			ot := NewObjectTable("table", "placeholder", cols, objectStore)
			if test.enablePluginResponse {
				pluginResponse := plugin.ObjectStatusResponse{
					ObjectStatus: component.PodSummary{
						Status:  component.NodeStatusError,
						Details: []component.Component{component.NewText("detail")},
					},
				}

				pluginManager.EXPECT().ObjectStatus(context.Background(), pod1).Return(&pluginResponse, nil)
				pluginManager.EXPECT().ObjectStatus(context.Background(), pod2).Return(&pluginResponse, nil)
				ot.EnablePluginStatus(pluginManager)
			}

			for i, pod := range []*corev1.Pod{pod1, pod2} {
				err := ot.AddRowForObject(ctx, pod, component.TableRow{
					"A": component.NewLink("", pod.Name, "/"+pod.Name),
					"B": component.NewText(fmt.Sprintf("%d", i)),
				})
				require.NoError(t, err)
			}
			test.mutateFn(ot)

			actual, err := ot.ToComponent()
			require.NoError(t, err)
			testutil.AssertJSONEqual(t, test.wanted(), actual)
		})
	}
}

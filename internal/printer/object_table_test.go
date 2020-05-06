/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package printer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestObjectTable(t *testing.T) {
	cols := component.NewTableCols("A", "B")

	genDeleteGA := func(object runtime.Object) *component.GridActions {
		ga := component.NewGridActions()
		ga.AddGridAction(buildObjectDeleteAction(t, object))
		return ga
	}

	pod1 := testutil.CreatePod("pod1")
	pod2 := testutil.CreatePod("pod2")

	tests := []struct {
		name     string
		mutateFn func(*ObjectTable)
		wanted   func() *component.Table
	}{
		{
			name: "no mutations",
			mutateFn: func(table *ObjectTable) {

			},
			wanted: func() *component.Table {
				return component.NewTableWithRows("table", "placeholder", cols, []component.TableRow{
					{
						"A":                     component.NewText("pod1"),
						"B":                     component.NewText("0"),
						component.GridActionKey: genDeleteGA(pod1),
					},
					{
						"A":                     component.NewText("pod2"),
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
						"A":                     component.NewText("pod2"),
						"B":                     component.NewText("1"),
						component.GridActionKey: genDeleteGA(pod2),
					},
					{
						"A":                     component.NewText("pod1"),
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
						"A":                     component.NewText("pod1"),
						"B":                     component.NewText("0"),
						component.GridActionKey: genDeleteGA(pod1),
					},
					{
						"A":                     component.NewText("pod2"),
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ot := NewObjectTable("table", "placeholder", cols)

			for i, pod := range []*corev1.Pod{pod1, pod2} {
				err := ot.AddRowForObject(pod, component.TableRow{
					"A": component.NewText(pod.Name),
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

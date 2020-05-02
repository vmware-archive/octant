/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package printer

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// ObjectTable is a helper for creating a table containing a list of objects.
type ObjectTable struct {
	cols        []component.TableCol
	title       string
	placeholder string
	rows        []component.TableRow
	filters     map[string]component.TableFilter
	sortOrder   *tableSetOrder
}

// NewObjectTable creates an instance of ObjectTable.
func NewObjectTable(title, placeholder string, cols []component.TableCol) *ObjectTable {
	ol := ObjectTable{
		cols:        cols,
		title:       title,
		placeholder: placeholder,
		filters:     map[string]component.TableFilter{},
	}

	return &ol
}

// AddFilters adds filters to a set of table columns.
func (ol *ObjectTable) AddFilters(filters map[string]component.TableFilter) {
	for k, v := range filters {
		ol.filters[k] = v
	}
}

// AddRowForObject adds a row for an object to the table.
func (ol *ObjectTable) AddRowForObject(object runtime.Object, row component.TableRow) error {
	gridAction, err := objectDeleteAction(object)
	if err != nil {
		return fmt.Errorf("create object delete action: %w", err)
	}

	row.AddAction(gridAction)

	ol.rows = append(ol.rows, row)

	return nil
}

func objectDeleteAction(object runtime.Object) (component.GridAction, error) {
	key, err := store.KeyFromObject(object)
	if err != nil {
		return component.GridAction{}, fmt.Errorf("create key from object: %w", err)
	}

	payload := key.ToActionPayload()

	deleteConfirmation, err := octant.DeleteObjectConfirmation(object)
	if err != nil {
		return component.GridAction{}, fmt.Errorf("create delete object confirmation: %w", err)
	}

	return component.GridAction{
		Name:         "Delete",
		ActionPath:   octant.ActionDeleteObject,
		Payload:      payload,
		Confirmation: deleteConfirmation,
		Type:         component.GridActionDanger,
	}, nil

}

type tableSetOrder struct {
	name    string
	reverse bool
}

// SetSortOrder sets the sort order for the table.
func (ol *ObjectTable) SetSortOrder(name string, reverse bool) {
	ol.sortOrder = &tableSetOrder{
		name:    name,
		reverse: reverse,
	}
}

// ToComponent converts the ObjectTable instance to a component.
func (ol *ObjectTable) ToComponent() (component.Component, error) {
	table := component.NewTableWithRows(ol.title, ol.placeholder, ol.cols, ol.rows)

	for name, filter := range ol.filters {
		table.AddFilter(name, filter)
	}

	if so := ol.sortOrder; so != nil {
		table.Sort(so.name, so.reverse)
	}

	return table, nil
}

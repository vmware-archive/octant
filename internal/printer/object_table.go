/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package printer

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/octant/pkg/plugin"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/objectstatus"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// ObjectTable is a helper for creating a table containing a list of objects.
type ObjectTable struct {
	cols          []component.TableCol
	title         string
	placeholder   string
	rows          []component.TableRow
	filters       map[string]component.TableFilter
	sortOrder     *tableSetOrder
	store         store.Store
	pluginManager plugin.ManagerInterface
}

// NewObjectTable creates an instance of ObjectTable.
func NewObjectTable(title, placeholder string, cols []component.TableCol, objectStore store.Store) *ObjectTable {
	ol := ObjectTable{
		cols:        cols,
		title:       title,
		placeholder: placeholder,
		filters:     map[string]component.TableFilter{},
		store:       objectStore,
	}

	return &ol
}

// AddFilters adds filters to a set of table columns.
func (ol *ObjectTable) AddFilters(filters map[string]component.TableFilter) {
	for k, v := range filters {
		ol.filters[k] = v
	}
}

func (ol *ObjectTable) EnablePluginStatus(pluginManager plugin.ManagerInterface) {
	ol.pluginManager = pluginManager
}

type componentStatus interface {
	SetStatus(status component.TextStatus, detail component.Component)
}

// AddRowForObject adds a row for an object to the table.
func (ol *ObjectTable) AddRowForObject(ctx context.Context, object runtime.Object, row component.TableRow) error {
	gridAction, err := objectDeleteAction(object)
	if err != nil {
		return fmt.Errorf("create object delete action: %w", err)
	}

	accessor, err := meta.Accessor(object)
	if err != nil {
		return fmt.Errorf("get accessor for object: %w", err)
	}

	if accessor.GetDeletionTimestamp() != nil {
		row["_isDeleted"] = component.NewText("deleted")
	}

	status, err := objectstatus.Status(ctx, object, ol.store, nil)
	if err != nil {
		return fmt.Errorf("get status for object: %w", err)
	}

	var pluginStatus *plugin.ObjectStatusResponse
	var detailComponent *component.List
	if ol.pluginManager != nil {
		pluginStatus, err = ol.pluginManager.ObjectStatus(ctx, object)
		if err != nil {
			return err
		}
	}

	if len(ol.cols) > 0 {
		firstRow := row[ol.cols[0].Name]
		details := status.Details
		cs, ok := firstRow.(componentStatus)
		if ok {
			detailComponent = component.NewList(nil, details)
			cs.SetStatus(convertNodeStatusToTextStatus(status.Status()), detailComponent)
		}
		if pluginStatus != nil {
			details = append(details, pluginStatus.ObjectStatus.Details...)
			detailComponent = component.NewList(nil, details)

			if pluginStatus.ObjectStatus.Status != "" {
				cs.SetStatus(convertNodeStatusToTextStatus(pluginStatus.ObjectStatus.Status), detailComponent)
			}
		}
	}

	row.AddAction(gridAction)

	ol.rows = append(ol.rows, row)

	return nil
}

func convertNodeStatusToTextStatus(nodeStatus component.NodeStatus) component.TextStatus {
	switch nodeStatus {
	case component.NodeStatusOK:
		return component.TextStatusOK
	case component.NodeStatusWarning:
		return component.TextStatusWarning
	case component.NodeStatusError:
		return component.TextStatusError
	default:
		return 0
	}
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
		table.Sort(so.name)
		if so.reverse {
			table.Reverse()
		}
	}

	return table, nil
}

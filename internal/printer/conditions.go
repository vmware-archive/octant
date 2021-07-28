/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/vmware-tanzu/octant/pkg/view/flexlayout"
)

var (
	conditionType           = "Type"
	conditionReason         = "Reason"
	conditionStatus         = "Status"
	conditionMessage        = "Message"
	conditionLastUpdate     = "Last Update"
	conditionLastTransition = "Last Transition"

	conditionMap = [][]string{
		{conditionType, "type"},
		{conditionReason, "reason"},
		{conditionStatus, "status"},
		{conditionMessage, "message"},
		{conditionLastUpdate, "lastUpdateTime"},
		{conditionLastTransition, "lastTransitionTime"},
	}
)

// parseConditions returns an error if no status is found or conditions fail to parse in
// to the expected `NestedSlice` format.
func parseConditions(u unstructured.Unstructured) ([]interface{}, error) {
	status, ok, err := unstructured.NestedMap(u.Object, "status")
	if err != nil {
		return nil, err
	}
	// No status found
	if !ok {
		return nil, fmt.Errorf("no status found for object")
	}

	conditions, ok, err := unstructured.NestedSlice(status, "conditions")
	if err != nil {
		return nil, err
	}
	// No conditions found
	if !ok {
		return make([]interface{}, 0), nil
	}

	return conditions, nil
}

// createConditionsTable returns a component.Table, if conditions exist on the object, and any err.
// For objects with empty conditions, an empty table with placeholder text is used.
func createConditionsTable(conditions []interface{}, sortKey string, customConditionMap [][]string) *component.Table {
	localConditionMap := conditionMap
	if customConditionMap != nil {
		localConditionMap = customConditionMap
	}

	table := component.NewTable("Conditions", "There are no conditions!", nil)

	columnSet := map[string]bool{}

	for _, c := range conditions {
		row := component.TableRow{}
		cm, _ := c.(map[string]interface{})

		for _, pair := range localConditionMap {
			if len(pair) != 2 {
				continue
			}
			columnKey, jsonKey := pair[0], pair[1]
			conditionValue, ok := cm[jsonKey]
			if !columnSet[columnKey] {
				table.AddColumn(columnKey)
				columnSet[columnKey] = true
			}
			if ok && conditionValue != nil {
				if strings.Contains(jsonKey, "Time") {
					t, err := time.Parse(time.RFC3339, conditionValue.(string))
					if err == nil {
						row[columnKey] = component.NewTimestamp(t)
						continue
					}
				}
				row[columnKey] = component.NewText(conditionValue.(string))
			}
		}
		table.Add(row)
	}

	table.Sort(sortKey)
	return table
}

func createConditionsForObject(ctx context.Context, fl *flexlayout.FlexLayout, object runtime.Object, sortKey string, columns [][]string, mapFn mapGenFn) error {
	var obj map[string]interface{}
	var err error

	if mapFn != nil {
		obj, err = mapFn(object)
	} else {
		obj, err = runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	}

	if err != nil {
		return nil
	}
	conditions, err := parseConditions(unstructured.Unstructured{Object: obj})
	if err != nil {
		return err
	}

	if sortKey == "" {
		sortKey = conditionType
	}

	conditionsTable := createConditionsTable(conditions, sortKey, columns)
	conditionsSection := fl.AddSection()
	if err := conditionsSection.Add(conditionsTable, component.WidthFull); err != nil {
		return fmt.Errorf("add conditions table to layout: %w", err)
	}
	return nil
}

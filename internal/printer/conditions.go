/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/pkg/view/component"
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

// createConditionsTable returns a component.Table, if conditions exist on the object, and any err.
// For objects with empty conditions, an empty table with placeholder text is used.
func createConditionsTable(u *unstructured.Unstructured, sortKey string, customConditionMap [][]string) (*component.Table, bool, error) {
	if u == nil {
		return nil, false, errors.New("object is nil")
	}

	localConditionMap := conditionMap
	if customConditionMap != nil {
		localConditionMap = customConditionMap
	}

	table := component.NewTable("Conditions", "There are no conditions!", nil)
	status, ok, err := unstructured.NestedMap(u.Object, "status")
	if err != nil {
		return nil, false, err
	}
	// No status found
	if !ok {
		return table, false, nil
	}

	conditions, ok, err := unstructured.NestedSlice(status, "conditions")
	if err != nil {
		return nil, false, err
	}
	// No conditions found
	if !ok {
		return table, false, nil
	}

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
			if ok {
				if !columnSet[columnKey] {
					table.AddColumn(columnKey)
					columnSet[columnKey] = true
				}
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
	return table, false, nil
}

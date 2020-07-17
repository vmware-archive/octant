/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func testWebhookRulesTable(rows ...component.TableRow) *component.Table {
	columns := component.NewTableCols("API Groups", "API Versions", "Resources", "Operations", "Scope")
	table := component.NewTable("Rules", "There are no webhook rules!", columns)
	for _, row := range rows {
		table.Add(row)
	}
	return table
}

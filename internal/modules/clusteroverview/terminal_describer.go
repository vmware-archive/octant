/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// TerminalListDescriber describes a list of terminals.
type TerminalListDescriber struct {
}

func NewTerminalListDescriber() *TerminalListDescriber {
	return &TerminalListDescriber{}
}

var _ describer.Describer = (*TerminalListDescriber)(nil)

// Describe describes a list of port forwards as content
func (d *TerminalListDescriber) Describe(ctx context.Context, namespace string, options describer.Options) (component.ContentResponse, error) {
	tm := options.TerminalManager()

	list := component.NewList("Terminals", nil)

	tblCols := component.NewTableCols("Container", "Command", "ID", "Age")
	tbl := component.NewTable("Terminals", "There are no terminals!", tblCols)
	list.Add(tbl)

	for _, t := range tm.List(ctx) {
		tRow := component.TableRow{
			"Container": component.NewText(t.Container()),
			"Command":   component.NewText(t.Command()),
			"ID":        component.NewText(t.ID()),
			"Age":       component.NewTimestamp(t.CreatedAt()),
		}
		tbl.Add(tRow)
	}

	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

func (d *TerminalListDescriber) PathFilters() []describer.PathFilter {
	filter := describer.NewPathFilter("/terminal", d)
	return []describer.PathFilter{*filter}
}

func (d *TerminalListDescriber) Reset(ctx context.Context) error {
	return nil
}

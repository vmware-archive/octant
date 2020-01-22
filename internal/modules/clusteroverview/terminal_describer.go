/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/pkg/action"
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

	tblCols := component.NewTableCols("Container", "Command", "Active", "ID", "Age", "")
	tbl := component.NewTable("Terminals", "There are no terminals!", tblCols)
	list.Add(tbl)

	for _, t := range tm.List(namespace) {
		nameLink, err := options.Link.ForGVK(t.Key().Namespace, t.Key().APIVersion, t.Key().Kind, t.Key().Name, t.Key().Name)
		if err != nil {
			return component.EmptyContentResponse, err
		}

		nameLink.Config.Text = t.Container()

		buttonGroup := component.NewButtonGroup()
		buttonGroup.AddButton(
			component.NewButton("Delete",
				action.CreatePayload("overview/deleteTerminal", action.Payload{"terminalID": t.ID()}),
				component.WithButtonConfirmation(
					"Delete Terminal",
					fmt.Sprintf("Are you sure you want to delete *Terminal* **%s**? You will lose access to the scrollback buffer.", t.ID()),
				)))

		tRow := component.TableRow{
			"Container": nameLink,
			"Command":   component.NewText(t.Command()),
			"Active":    component.NewText(strconv.FormatBool(t.Active())),
			"ID":        component.NewText(t.ID()),
			"Age":       component.NewTimestamp(t.CreatedAt()),
			"":          buttonGroup,
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

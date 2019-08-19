/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/vmware/octant/internal/describer"
	"github.com/vmware/octant/pkg/view/component"
)

// PluginListDescriber describes a list of plugins
type PluginListDescriber struct {
}

var _ describer.Describer = (*PluginListDescriber)(nil)

// Describe describes a list of plugins
func (d *PluginListDescriber) Describe(ctx context.Context, prefix, namespace string, options describer.Options) (component.ContentResponse, error) {
	pluginStore := options.PluginManager().Store()

	list := component.NewList("Plugins", nil)
	tableCols := component.NewTableCols("Name", "Description", "Capabilities")
	tbl := component.NewTable("Plugins", "There are no plugins!", tableCols)
	list.Add(tbl)

	for _, n := range pluginStore.ClientNames() {
		metadata, err := pluginStore.GetMetadata(n)
		if err != nil {
			return describer.EmptyContentResponse, errors.New("metadata is nil")
		}

		capability, err := json.Marshal(metadata.Capabilities)
		if err != nil {
			return describer.EmptyContentResponse, err
		}

		row := component.TableRow{
			"Name":        component.NewText(metadata.Name),
			"Description": component.NewText(metadata.Description),
			"Capability":  component.NewText(string(capability)),
		}
		tbl.Add(row)
	}

	tbl.Sort("Name", false)

	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

func (d *PluginListDescriber) PathFilters() []describer.PathFilter {
	filter := describer.NewPathFilter("/plugins", d)
	return []describer.PathFilter{*filter}
}

func (d *PluginListDescriber) Reset(ctx context.Context) error {
	return nil
}

func NewPluginListDescriber() *PluginListDescriber {
	return &PluginListDescriber{}
}

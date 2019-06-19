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

// Describe describes a list of plugins
func (d *PluginListDescriber) Describe(ctx context.Context, prefix, namespace string, options describer.Options) (component.ContentResponse, error) {
	pluginStore := options.PluginManager().Store()

	list := component.NewList("Plugins", nil)
	tableCols := component.NewTableCols("Name", "Description", "Capabilities")
	tbl := component.NewTable("Plugins", tableCols)
	list.Add(tbl)

	for _, n := range pluginStore.ClientNames() {
		metadata, err := pluginStore.GetMetadata(n)
		if err != nil {
			return component.ContentResponse{}, errors.New("metadata is nil")
		}

		capability, _ := json.Marshal(metadata.Capabilities)

		row := component.TableRow{
			"Name":        component.NewText(metadata.Name),
			"Description": component.NewText(metadata.Description),
			"Capability":  component.NewText(string(capability)),
		}
		tbl.Add(row)
	}

	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

func (d *PluginListDescriber) PathFilters() []describer.PathFilter {
	filter := describer.NewPathFilter("/plugins", d)
	return []describer.PathFilter{*filter}
}

func NewPluginListDescriber() *PluginListDescriber {
	return &PluginListDescriber{}
}

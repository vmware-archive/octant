/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// PluginListDescriber describes a list of plugins
type PluginListDescriber struct {
}

var _ describer.Describer = (*PluginListDescriber)(nil)

// Describe describes a list of plugins
func (d *PluginListDescriber) Describe(ctx context.Context, namespace string, options describer.Options) (component.ContentResponse, error) {
	pluginStore := options.PluginManager().Store()
	title := append([]component.TitleComponent{}, component.NewText("Plugins"))
	list := component.NewList(title, nil)
	tableCols := component.NewTableCols("Name", "Description", "Capabilities")
	tbl := component.NewTable("Plugins", "There are no plugins!", tableCols)
	list.Add(tbl)

	for _, n := range pluginStore.ClientNames() {
		var metadata *plugin.Metadata
		if plugin.IsJavaScriptPlugin(n) {
			jsPlugin, ok := pluginStore.GetJS(n)
			if !ok {
				return component.EmptyContentResponse, fmt.Errorf("plugin %s not found", n)
			}
			metadata = jsPlugin.Metadata()
		} else {
			var err error
			metadata, err = pluginStore.GetMetadata(n)
			if err != nil {
				return component.EmptyContentResponse, fmt.Errorf("metadata is nil")
			}
		}

		var summaryItems []string
		if metadata.Capabilities.IsModule {
			summaryItems = append(summaryItems, "Module")
		}

		if actionNames := metadata.Capabilities.ActionNames; len(actionNames) > 0 {
			summaryItems = append(summaryItems, fmt.Sprintf("Actions: %s",
				strings.Join(actionNames, ", ")))
		}

		in := []struct {
			name string
			list []schema.GroupVersionKind
		}{
			{name: "Object Status", list: metadata.Capabilities.SupportsObjectStatus},
			{name: "Printer Config", list: metadata.Capabilities.SupportsPrinterConfig},
			{name: "Printer Items", list: metadata.Capabilities.SupportsPrinterItems},
			{name: "Printer Status", list: metadata.Capabilities.SupportsPrinterStatus},
			{name: "Tab", list: metadata.Capabilities.SupportsTab},
		}

		for _, item := range in {
			support, ok := summarizeSupports(item.name, item.list)
			if ok {
				summaryItems = append(summaryItems, support)
			}
		}

		var sb strings.Builder
		for i := range summaryItems {
			sb.WriteString(fmt.Sprintf("[%s]", summaryItems[i]))
			if i < len(summaryItems)-1 {
				sb.WriteString(", ")
			}
		}

		row := component.TableRow{
			"Name":         component.NewText(metadata.Name),
			"Description":  component.NewText(metadata.Description),
			"Capabilities": component.NewText(sb.String()),
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

func summarizeSupports(name string, list []schema.GroupVersionKind) (string, bool) {
	if len(list) < 1 {
		return "", false
	}

	var items []string
	for _, groupVersionKind := range list {
		apiVersion, kind := groupVersionKind.ToAPIVersionAndKind()
		items = append(items, fmt.Sprintf("%s %s", apiVersion, kind))
	}

	return fmt.Sprintf("%s: %s",
		name, strings.Join(items, ", "),
	), true
}

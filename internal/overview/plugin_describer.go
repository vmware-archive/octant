package overview

import (
	"context"
	"encoding/json"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
)

// PluginListDescriber describes a list of plugins
type PluginListDescriber struct {
}

// Describe describes a list of plugins
func (d *PluginListDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (component.ContentResponse, error) {
	if options.PluginManagerStore == nil {
		return component.ContentResponse{}, errors.New("plugin store is nil")
	}

	list := component.NewList("Plugins", nil)
	tableCols := component.NewTableCols("Name", "Description", "Capabilities")
	tbl := component.NewTable("Plugins", tableCols)
	list.Add(tbl)

	for _, n := range options.PluginManagerStore.ClientNames() {
		metadata, err := options.PluginManagerStore.GetMetadata(n)
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

func (d *PluginListDescriber) PathFilters() []pathFilter {
	filter := newPathFilter("/plugins", d)
	return []pathFilter{*filter}
}

func NewPluginListDescriber() *PluginListDescriber {
	return &PluginListDescriber{}
}

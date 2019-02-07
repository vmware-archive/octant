package printer

import (
	"fmt"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
)

// ConfigMapListHandler is a printFunc that prints ConfigMaps
func ConfigMapListHandler(list *corev1.ConfigMapList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("list is nil")
	}

	// Data column
	cols := component.NewTableCols("Name", "Labels", "Data", "Age")
	tbl := component.NewTable("ConfigMaps", cols)

	for _, c := range list.Items {
		row := component.TableRow{}
		configmapPath := gvkPath(c.TypeMeta.APIVersion, c.TypeMeta.Kind, c.Name)
		row["Name"] = component.NewLink("", c.Name, configmapPath)
		row["Labels"] = component.NewLabels(c.Labels)

		data := fmt.Sprintf("%d", len(c.Data))
		row["Data"] = component.NewText("", data)

		ts := c.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		tbl.Add(row)
	}

	return tbl, nil
}

// ConfigMapHandler is a printFunc that prints a ConfigMap
func ConfigMapHandler(configmap *corev1.ConfigMap, options Options) (component.ViewComponent, error) {
	grid := component.NewGrid("Summary")

	detailsSummary := component.NewSummary("Details")

	detailsPanel := component.NewPanel("", detailsSummary)
	grid.Add(*detailsPanel)

	list := component.NewList("", []component.ViewComponent{grid})

	return list, nil
}

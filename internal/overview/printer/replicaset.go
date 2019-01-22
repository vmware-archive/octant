package printer

import (
	"errors"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/heptio/developer-dash/internal/view/component"
)

// ReplicaSetListHandler is a printFunc that lists deployments
func ReplicaSetListHandler(list *appsv1.ReplicaSetList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	tbl := component.NewTable("ReplicaSets", cols)

	for _, d := range list.Items {
		row := component.TableRow{}
		row["Name"] = component.NewText("", d.Name)
		row["Labels"] = component.NewLabels(d.Labels)

		status := fmt.Sprintf("%d/%d", d.Status.AvailableReplicas, d.Status.Replicas)
		row["Status"] = component.NewText("", status)

		ts := d.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		containers := component.NewContainers()
		for _, c := range d.Spec.Template.Spec.Containers {
			containers.Add(c.Name, c.Image)
		}
		row["Containers"] = containers
		row["Selector"] = printSelector(d.Spec.Selector)

		tbl.Add(row)
	}
	return tbl, nil
}

// ReplicaSetHandler is a printFunc that prints a ReplicaSets.
// TODO: This handler is incomplete.
func ReplicaSetHandler(deployment *appsv1.ReplicaSet, options Options) (component.ViewComponent, error) {
	grid := component.NewGrid("Summary")

	detailsSummary := component.NewSummary("Details")

	detailsPanel := component.NewPanel("", detailsSummary)
	grid.Add(*detailsPanel)

	list := component.NewList("", []component.ViewComponent{grid})

	return list, nil
}

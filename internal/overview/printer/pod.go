package printer

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/view/component"
	corev1 "k8s.io/api/core/v1"
)

// PodListHandler is a printFunc that prints pods
func PodListHandler(list *corev1.PodList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Ready", "Status", "Restarts", "Age")
	tbl := component.NewTable("Pods", cols)

	for _, d := range list.Items {
		row := component.TableRow{}
		row["Name"] = component.NewText("", d.Name)
		row["Labels"] = component.NewLabels(d.Labels)

		readyCounter := 0
		for _, c := range d.Status.ContainerStatuses {
			if c.Ready {
				readyCounter++
			}
		}
		ready := fmt.Sprintf("%d/%d", readyCounter, len(d.Status.ContainerStatuses))
		row["Ready"] = component.NewText("", ready)

		status := fmt.Sprintf("%s", d.Status.Phase)
		row["Status"] = component.NewText("", status)

		restartCounter := 0
		for _, c := range d.Status.ContainerStatuses {
			restartCounter += int(c.RestartCount)
		}
		restarts := fmt.Sprintf("%d", restartCounter)
		row["Restarts"] = component.NewText("", restarts)

		ts := d.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		tbl.Add(row)
	}

	return tbl, nil
}

// PodHandler is a printFunc that prints Pods
// TODO: This handler is incomplete
func PodHandler(p *corev1.Pod, opts Options) (component.ViewComponent, error) {
	grid := component.NewGrid("Summary")
	summary := component.NewSummary("Details")
	panel := component.NewPanel("", summary)

	grid.Add(*panel)

	list := component.NewList("", []component.ViewComponent{grid})

	return list, nil
}

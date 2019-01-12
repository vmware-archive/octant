package printer

import (
	"errors"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/view/component"
)

// DeploymentListHandler is a printFunc that lists deployments
func DeploymentListHandler(list *appsv1.DeploymentList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	tbl := component.NewTable("Deployments", cols)

	for _, d := range list.Items {
		row := component.TableRow{}
		row["Name"] = component.NewText("", d.Name)
		row["Labels"] = component.NewLabels(d.Labels)

		status := fmt.Sprintf("%d/%d", d.Status.AvailableReplicas, d.Status.AvailableReplicas+d.Status.UnavailableReplicas)
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

// DeploymentHandler is a printFunc that printers a Deployments.
// This handler is incomplete.
func DeploymentHandler(deployment *appsv1.Deployment, options Options) (component.ViewComponent, error) {
	grid := component.NewGrid("Summary")

	detailsSummary := component.NewSummary("Details")

	detailsPanel := component.NewPanel("", detailsSummary)
	grid.Add(*detailsPanel)

	list := component.NewList("", []component.ViewComponent{grid})

	return list, nil
}

func printSelector(selector *metav1.LabelSelector) *component.Selectors {
	s := component.NewSelectors(nil)
	if selector == nil {
		return s
	}

	for k, v := range selector.MatchLabels {
		s.Add(component.NewLabelSelector(k, v))
	}

	for _, e := range selector.MatchExpressions {
		es := component.NewExpressionSelector(e.Key, component.Operator(e.Operator), e.Values)
		s.Add(es)
	}

	return s
}

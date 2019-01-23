package printer

import (
	"errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/heptio/developer-dash/internal/view/component"
)

// ServiceListHandler is a printFunc that lists deployments
func ServiceListHandler(list *corev1.ServiceList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Type", "Cluster IP", "External IP", "Ports", "Age", "Selector")
	tbl := component.NewTable("Services", cols)

	for _, d := range list.Items {
		row := component.TableRow{}
		row["Name"] = component.NewText("", d.Name)
		row["Labels"] = component.NewLabels(d.Labels)
		row["Type"] = component.NewText("", string(d.Spec.Type))
		row["Cluster IP"] = component.NewText("", d.Spec.ClusterIP)
		row["External IP"] = component.NewText("", strings.Join(d.Spec.ExternalIPs, ","))
		row["Ports"] = printServicePorts(d.Spec.Ports)

		ts := d.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		row["Selector"] = printSelectorMap(d.Spec.Selector)

		tbl.Add(row)
	}
	return tbl, nil
}

// ServiceHandler is a printFunc that prints a Services.
// TODO: This handler is incomplete.
func ServiceHandler(deployment *corev1.Service, options Options) (component.ViewComponent, error) {
	grid := component.NewGrid("Summary")

	detailsSummary := component.NewSummary("Details")

	detailsPanel := component.NewPanel("", detailsSummary)
	grid.Add(*detailsPanel)

	list := component.NewList("", []component.ViewComponent{grid})

	return list, nil
}

func printServicePorts(ports []corev1.ServicePort) component.ViewComponent {
	out := make([]string, len(ports))
	for i, port := range ports {
		if port.TargetPort.Type == intstr.Type(intstr.Int) {
			out[i] = fmt.Sprintf("%d/%s", port.TargetPort.IntVal, port.Protocol)
		} else {
			out[i] = fmt.Sprintf("%s/%s", port.TargetPort.StrVal, port.Protocol)
		}
		if port.NodePort != 0 {
			out[i] = fmt.Sprintf("%d/%s", port.NodePort, port.Protocol)
		}
	}

	return component.NewText("", strings.Join(out, ","))
}

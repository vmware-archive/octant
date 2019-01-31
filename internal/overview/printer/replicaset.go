package printer

import (
	"errors"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"

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
func ReplicaSetHandler(rs *appsv1.ReplicaSet, options Options) (component.ViewComponent, error) {
	grid := component.NewGrid("Summary")

	detailsSummary := component.NewSummary("Details")

	detailsPanel := component.NewPanel("", detailsSummary)
	grid.Add(*detailsPanel)

	list := component.NewList("", []component.ViewComponent{grid})

	return list, nil
}

// ReplicaSet generates a replicaset satatus
type ReplicaSetConfiguration struct {
	replicaset *appsv1.ReplicaSet
	client     clientset.Interface
}

// NewReplicaSetConfiguration ...
func NewReplicaSetConfiguration(rs *appsv1.ReplicaSet, c clientset.Interface) *ReplicaSetConfiguration {
	return &ReplicaSetConfiguration{
		replicaset: rs,
		client:     c,
	}
}

// Create generates a replicaset status
func (rc *ReplicaSetConfiguration) Create() (*component.Quadrant, error) {
	if rc.replicaset == nil {
		return nil, errors.New("replicaset is nil")
	}
	pods := rc.client.Core().Pods(rc.replicaset.Namespace)

	selector, _ := metav1.LabelSelectorAsSelector(rc.replicaset.Spec.Selector)

	ps := CreatePodStatus(pods, selector)

	quadrant := component.NewQuadrant()
	if err := quadrant.Set(component.QuadNW, "Running", fmt.Sprintf("%d", ps.Running)); err != nil {
		return nil, errors.New("unable to set quadrant nw")
	}
	if err := quadrant.Set(component.QuadNE, "Waiting", fmt.Sprintf("%d", ps.Waiting)); err != nil {
		return nil, errors.New("unable to set quadrant ne")
	}
	if err := quadrant.Set(component.QuadSW, "Succeeded", fmt.Sprintf("%d", ps.Succeeded)); err != nil {
		return nil, errors.New("unable to set quadrant sw")
	}
	if err := quadrant.Set(component.QuadSE, "Failed", fmt.Sprintf("%d", ps.Failed)); err != nil {
		return nil, errors.New("unable to set quadrant se")
	}

	return quadrant, nil
}

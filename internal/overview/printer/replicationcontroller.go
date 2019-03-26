package printer

import (
	"context"
	"fmt"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplicationControllerListHandler is a printFunc that lists ReplicationControllers
func ReplicationControllerListHandler(ctx context.Context, list *corev1.ReplicationControllerList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	tbl := component.NewTable("ReplicationControllers", cols)

	for _, rc := range list.Items {
		row := component.TableRow{}

		row["Name"] = link.ForObject(&rc, rc.Name)
		row["Labels"] = component.NewLabels(rc.Labels)

		status := fmt.Sprintf("%d/%d", rc.Status.AvailableReplicas, rc.Status.Replicas)
		row["Status"] = component.NewText(status)

		ts := rc.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		containers := component.NewContainers()
		for _, c := range rc.Spec.Template.Spec.Containers {
			containers.Add(c.Name, c.Image)
		}
		row["Containers"] = containers

		row["Selector"] = printSelectorMap(rc.Spec.Selector)

		tbl.Add(row)
	}
	return tbl, nil
}

// ReplicationControllerHandler is a printFunc that prints a ReplicationController
func ReplicationControllerHandler(ctx context.Context, rc *corev1.ReplicationController, options Options) (component.Component, error) {
	o := NewObject(rc)

	rcConfigGen := NewReplicationControllerConfiguration(rc)
	o.RegisterConfig(func() (component.Component, error) {
		return rcConfigGen.Create()
	}, 16)

	rcSummaryGen := NewReplicationControllerStatus(rc)
	o.RegisterSummary(func() (component.Component, error) {
		return rcSummaryGen.Create(ctx, options.Cache)
	}, 8)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return createPodListView(ctx, rc, options)
		},
		Width: component.WidthFull,
	})

	o.EnablePodTemplate(*rc.Spec.Template)

	o.EnableEvents()

	return o.ToComponent(ctx, options)
}

// ReplicationControllerConfiguration generates a replicationcontroller configuration
type ReplicationControllerConfiguration struct {
	replicationcontroller *corev1.ReplicationController
}

// NewReplicationControllerConfiguration creates an instance of ReplicationControllerConfiguration
func NewReplicationControllerConfiguration(rc *corev1.ReplicationController) *ReplicationControllerConfiguration {
	return &ReplicationControllerConfiguration{
		replicationcontroller: rc,
	}
}

// Create generates a replicationcontroller configuration summary
func (rcc *ReplicationControllerConfiguration) Create() (*component.Summary, error) {
	if rcc == nil || rcc.replicationcontroller == nil {
		return nil, errors.New("replicationcontroller is nil")
	}

	replicationController := rcc.replicationcontroller

	sections := component.SummarySections{}

	if controllerRef := metav1.GetControllerOf(replicationController); controllerRef != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Controlled By",
			Content: link.ForOwner(replicationController, controllerRef),
		})
	}

	current := fmt.Sprintf("%d", replicationController.Status.ReadyReplicas)

	if desired := replicationController.Spec.Replicas; desired != nil {
		desiredReplicas := fmt.Sprintf("%d", *desired)
		status := fmt.Sprintf("Current %s / Desired %s", current, desiredReplicas)
		sections.AddText("Replica Status", status)
	}

	replicas := fmt.Sprintf("%d", replicationController.Status.Replicas)
	sections.AddText("Replicas", replicas)

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

// ReplicationControllerStatus generates a replication controller status
type ReplicationControllerStatus struct {
	replicationcontroller *corev1.ReplicationController
}

// NewReplicationControllerStatus creates an instance of ReplicationControllerStatus
func NewReplicationControllerStatus(replicationController *corev1.ReplicationController) *ReplicationControllerStatus {
	return &ReplicationControllerStatus{
		replicationcontroller: replicationController,
	}
}

// Create generates a replicaset status quadrant
func (replicationController *ReplicationControllerStatus) Create(ctx context.Context, c cache.Cache) (*component.Quadrant, error) {
	if replicationController.replicationcontroller == nil {
		return nil, errors.New("replicationcontroller is nil")
	}

	selectors := metav1.LabelSelector{
		MatchLabels: replicationController.replicationcontroller.Spec.Selector,
	}

	pods, err := listPods(ctx, replicationController.replicationcontroller.Namespace, &selectors, replicationController.replicationcontroller.GetUID(), c)
	if err != nil {
		return nil, err
	}

	ps := createPodStatus(pods)

	quadrant := component.NewQuadrant("Status")
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

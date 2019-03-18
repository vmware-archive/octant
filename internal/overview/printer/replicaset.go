package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
)

// ReplicaSetListHandler is a printFunc that lists deployments
func ReplicaSetListHandler(ctx context.Context, list *appsv1.ReplicaSetList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	tbl := component.NewTable("ReplicaSets", cols)

	for _, rs := range list.Items {
		row := component.TableRow{}
		row["Name"] = link.ForObject(&rs, rs.Name)
		row["Labels"] = component.NewLabels(rs.Labels)

		status := fmt.Sprintf("%d/%d", rs.Status.AvailableReplicas, rs.Status.Replicas)
		row["Status"] = component.NewText(status)

		ts := rs.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		containers := component.NewContainers()
		for _, c := range rs.Spec.Template.Spec.Containers {
			containers.Add(c.Name, c.Image)
		}
		row["Containers"] = containers
		row["Selector"] = printSelector(rs.Spec.Selector)

		tbl.Add(row)
	}
	return tbl, nil
}

// ReplicaSetHandler is a printFunc that prints a ReplicaSets.
func ReplicaSetHandler(ctx context.Context, rs *appsv1.ReplicaSet, options Options) (component.ViewComponent, error) {
	o := NewObject(rs)

	replicaSetConfigGen := NewReplicaSetConfiguration(rs)
	o.RegisterConfig(func() (component.ViewComponent, error) {
		return replicaSetConfigGen.Create()
	}, 16)

	replicaSetStatusGen := NewReplicaSetStatus(rs)
	o.RegisterSummary(func() (component.ViewComponent, error) {
		return replicaSetStatusGen.Create(ctx, options.Cache)
	}, 8)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.ViewComponent, error) {
			return createPodListView(ctx, rs, options)
		},
		Width: component.WidthFull,
	})

	o.EnablePodTemplate(rs.Spec.Template)

	o.EnableEvents()

	return o.ToComponent(ctx, options)
}

// ReplicaSetConfiguration generates a replicaset configuration
type ReplicaSetConfiguration struct {
	replicaset *appsv1.ReplicaSet
}

// NewReplicaSetConfiguration creates an instance of ReplicaSetConfiguration
func NewReplicaSetConfiguration(rs *appsv1.ReplicaSet) *ReplicaSetConfiguration {
	return &ReplicaSetConfiguration{
		replicaset: rs,
	}
}

// Create generates a replicaset configuration summary
func (rc *ReplicaSetConfiguration) Create() (*component.Summary, error) {
	if rc == nil || rc.replicaset == nil {
		return nil, errors.New("replicaset is nil")
	}

	rs := rc.replicaset

	sections := component.SummarySections{}

	if controllerRef := metav1.GetControllerOf(rs); controllerRef != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Controlled By",
			Content: link.ForOwner(rs, controllerRef),
		})
	}

	current := fmt.Sprintf("%d", rs.Status.ReadyReplicas)

	if desired := rs.Spec.Replicas; desired != nil {
		desiredReplicas := fmt.Sprintf("%d", *desired)
		status := fmt.Sprintf("Current %s / Desired %s", current, desiredReplicas)
		sections.AddText("Replica Status", status)
	}

	replicas := fmt.Sprintf("%d", rs.Status.Replicas)
	sections.AddText("Replicas", replicas)

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

// ReplicaSetStatus generates a replicaset status
type ReplicaSetStatus struct {
	replicaset *appsv1.ReplicaSet
}

// NewReplicaSetStatus creates an instance of ReplicaSetStatus
func NewReplicaSetStatus(rs *appsv1.ReplicaSet) *ReplicaSetStatus {
	return &ReplicaSetStatus{
		replicaset: rs,
	}
}

// Create generates a replicaset status quadrant
func (rs *ReplicaSetStatus) Create(ctx context.Context, c cache.Cache) (*component.Quadrant, error) {
	if rs == nil || rs.replicaset == nil {
		return nil, errors.New("replicaset is nil")
	}
	pods, err := listPods(ctx, rs.replicaset.Namespace, rs.replicaset.Spec.Selector, rs.replicaset.GetUID(), c)
	if err != nil {
		return nil, err
	}

	ps := createPodStatus(pods)

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

package printer

import (
	"fmt"

	"github.com/heptio/developer-dash/internal/overview/link"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
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
		row["Name"] = link.ForObject(&d, d.Name)
		row["Labels"] = component.NewLabels(d.Labels)

		status := fmt.Sprintf("%d/%d", d.Status.AvailableReplicas, d.Status.AvailableReplicas+d.Status.UnavailableReplicas)
		row["Status"] = component.NewText(status)

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

// DeploymentHandler is a printFunc that prints a Deployments.
func DeploymentHandler(deployment *appsv1.Deployment, options Options) (component.ViewComponent, error) {
	o := NewObject(deployment)

	deployConfigGen := NewDeploymentConfiguration(deployment)
	o.RegisterConfig(func() (component.ViewComponent, error) {
		return deployConfigGen.Create()
	}, 16)

	deploySummaryGen := NewDeploymentStatus(deployment)
	o.RegisterSummary(func() (component.ViewComponent, error) {
		return deploySummaryGen.Create()
	}, 8)

	o.EnablePodTemplate(deployment.Spec.Template)

	o.EnableEvents()

	return o.ToComponent(options)
}

// DeploymentConfiguration generates deployment configuration.
type DeploymentConfiguration struct {
	deployment *appsv1.Deployment
}

// NewDeploymentConfiguration creates an instance of DeploymentConfiguration.
func NewDeploymentConfiguration(d *appsv1.Deployment) *DeploymentConfiguration {
	return &DeploymentConfiguration{
		deployment: d,
	}
}

// Create creates a deployment configuration summary.
func (dc *DeploymentConfiguration) Create() (*component.Summary, error) {
	if dc.deployment == nil {
		return nil, errors.New("deployment is nil")
	}

	sections := make([]component.SummarySection, 0)

	strategyType := dc.deployment.Spec.Strategy.Type
	sections = append(sections, component.SummarySection{
		Header:  "Deployment Strategy",
		Content: component.NewText(string(strategyType)),
	})

	switch strategyType {
	case appsv1.RollingUpdateDeploymentStrategyType:
		rollingUpdate := dc.deployment.Spec.Strategy.RollingUpdate
		if rollingUpdate == nil {
			return nil, errors.Errorf("deployment strategy type is RollingUpdate, but configuration is nil")
		}

		rollingUpdateText := fmt.Sprintf("Max Surge %s%%, Max Unavailable %s%%",
			rollingUpdate.MaxSurge.String(),
			rollingUpdate.MaxUnavailable.String(),
		)

		sections = append(sections, component.SummarySection{
			Header:  "Rolling Update Strategy",
			Content: component.NewText(rollingUpdateText),
		})

		if selector := dc.deployment.Spec.Selector; selector != nil {
			var selectors []component.Selector

			for _, lsr := range selector.MatchExpressions {
				o, err := component.MatchOperator(string(lsr.Operator))
				if err != nil {
					return nil, err
				}

				es := component.NewExpressionSelector(lsr.Key, o, lsr.Values)
				selectors = append(selectors, es)
			}

			for k, v := range selector.MatchLabels {
				ls := component.NewLabelSelector(k, v)
				selectors = append(selectors, ls)
			}

			sections = append(sections, component.SummarySection{
				Header:  "Selectors",
				Content: component.NewSelectors(selectors),
			})
		}

		minReadySeconds := fmt.Sprintf("%d", dc.deployment.Spec.MinReadySeconds)
		sections = append(sections, component.SummarySection{
			Header:  "Min Ready Seconds",
			Content: component.NewText(minReadySeconds),
		})

		if rhl := dc.deployment.Spec.RevisionHistoryLimit; rhl != nil {
			revisionHistoryLimit := fmt.Sprintf("%d", *rhl)
			sections = append(sections, component.SummarySection{
				Header:  "Revision History Limit",
				Content: component.NewText(revisionHistoryLimit),
			})
		}
	}

	var replicas int32
	if dc.deployment.Spec.Replicas != nil {
		replicas = *dc.deployment.Spec.Replicas
	}

	sections = append(sections, component.SummarySection{
		Header:  "Replicas",
		Content: component.NewText(fmt.Sprintf("%d", replicas)),
	})

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

// DeploymentStatus generates deployment status.
type DeploymentStatus struct {
	deployment *appsv1.Deployment
}

// NewDeploymentStatus creates an instance of DeploymentStatus.
func NewDeploymentStatus(d *appsv1.Deployment) *DeploymentStatus {
	return &DeploymentStatus{
		deployment: d,
	}
}

// Create generates a deployment status quadrant.
func (ds *DeploymentStatus) Create() (*component.Quadrant, error) {
	if ds.deployment == nil {
		return nil, errors.New("deployment is nil")
	}

	status := ds.deployment.Status

	quadrant := component.NewQuadrant()
	if err := quadrant.Set(component.QuadNW, "Updated", fmt.Sprintf("%d", status.UpdatedReplicas)); err != nil {
		return nil, errors.New("unable to set quadrant nw")
	}
	if err := quadrant.Set(component.QuadNE, "Total", fmt.Sprintf("%d", status.Replicas)); err != nil {
		return nil, errors.New("unable to set quadrant ne")
	}
	if err := quadrant.Set(component.QuadSW, "Unavailable", fmt.Sprintf("%d", status.UnavailableReplicas)); err != nil {
		return nil, errors.New("unable to set quadrant sw")
	}
	if err := quadrant.Set(component.QuadSE, "Available", fmt.Sprintf("%d", status.AvailableReplicas)); err != nil {
		return nil, errors.New("unable to set quadrant se")
	}

	return quadrant, nil
}

package printer

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/gridlayout"
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
		deploymentPath := gvkPath(d.TypeMeta.APIVersion, d.TypeMeta.Kind, d.Name)
		row["Name"] = component.NewLink("", d.Name, deploymentPath)
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

// DeploymentHandler is a printFunc that prints a Deployments.
// TODO: This handler is incomplete.
func DeploymentHandler(deployment *appsv1.Deployment, options Options) (component.ViewComponent, error) {
	gl := gridlayout.New()

	configSection := gl.CreateSection(8)

	deployConfigGen := NewDeploymentConfiguration(deployment)
	configView, err := deployConfigGen.Create()
	if err != nil {
		return nil, err
	}

	configSection.Add(configView, 16)

	summarySection := gl.CreateSection(8)

	deploySummaryGen := NewDeploymentStatus(deployment)
	statusView, err := deploySummaryGen.Create()
	if err != nil {
		return nil, err
	}

	summarySection.Add(statusView, 8)

	podTemplate := NewPodTemplate(deployment.Spec.Template)
	if err = podTemplate.AddToGridLayout(gl); err != nil {
		return nil, err
	}

	grid := gl.ToGrid()

	return grid, nil
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
		Content: component.NewText("", string(strategyType)),
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
			Content: component.NewText("", rollingUpdateText),
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
			Content: component.NewText("", minReadySeconds),
		})

		if rhl := dc.deployment.Spec.RevisionHistoryLimit; rhl != nil {
			revisionHistoryLimit := fmt.Sprintf("%d", *rhl)
			sections = append(sections, component.SummarySection{
				Header:  "Revision History Limit",
				Content: component.NewText("", revisionHistoryLimit),
			})
		}
	}

	creationTimestamp := dc.deployment.CreationTimestamp.Time
	sections = append(sections, component.SummarySection{
		Header:  "Age",
		Content: component.NewTimestamp(creationTimestamp),
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

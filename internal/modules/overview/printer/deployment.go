/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

var (
	deploymentConditionColumns = component.NewTableCols("Type", "Reason", "Status", "Message", "Last Update", "Last Transition")
)

// DeploymentListHandler is a printFunc that lists deployments
func DeploymentListHandler(_ context.Context, list *appsv1.DeploymentList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	tbl := component.NewTable("Deployments", "We couldn't find any deployments!", cols)

	for _, d := range list.Items {
		row := component.TableRow{}
		nameLink, err := opts.Link.ForObject(&d, d.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(d.Labels)

		status := fmt.Sprintf("%d/%d", d.Status.AvailableReplicas, d.Status.AvailableReplicas+d.Status.UnavailableReplicas)
		row["Status"] = component.NewText(status)

		ts := d.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		containers := component.NewContainers()
		for i := range d.Spec.Template.Spec.Containers {
			c := d.Spec.Template.Spec.Containers[i]
			containers.Add(c.Name, c.Image)
		}
		row["Containers"] = containers
		row["Selector"] = printSelector(d.Spec.Selector)

		tbl.Add(row)
	}
	return tbl, nil
}

// DeploymentHandler is a printFunc that prints a Deployments.
func DeploymentHandler(ctx context.Context, deployment *appsv1.Deployment, options Options) (component.Component, error) {
	o := NewObject(deployment)

	deployConfigGen := NewDeploymentConfiguration(deployment)
	configSummary, err := deployConfigGen.Create()
	if err != nil {
		return nil, err
	}

	status, err := deploymentStatus(deployment)
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(configSummary)
	o.RegisterSummary(status)
	o.RegisterItems([]ItemDescriptor{
		{
			Func: func() (component.Component, error) {
				return deploymentPods(ctx, deployment, options)
			},
			Width: component.WidthFull,
		},
		{
			Func: func() (i component.Component, e error) {
				return deploymentConditions(deployment)
			},
			Width: component.WidthFull,
		},
	}...)
	o.EnablePodTemplate(deployment.Spec.Template)
	o.EnableEvents()

	return o.ToComponent(ctx, options)
}

func deploymentStatus(deployment *appsv1.Deployment) (*component.Summary, error) {
	if deployment == nil {
		return nil, errors.New("unable to generate status from a nil deployment")
	}

	status := deployment.Status

	summary := component.NewSummary("Status", []component.SummarySection{
		{
			Header:  "Available Replicas",
			Content: component.NewText(fmt.Sprintf("%d", status.AvailableReplicas)),
		},
		{
			Header:  "Ready Replicas",
			Content: component.NewText(fmt.Sprintf("%d", status.ReadyReplicas)),
		},
		{
			Header:  "Total Replicas",
			Content: component.NewText(fmt.Sprintf("%d", status.Replicas)),
		},
		{
			Header:  "Unavailable Replicas",
			Content: component.NewText(fmt.Sprintf("%d", status.UnavailableReplicas)),
		},
		{
			Header:  "Updated Replicas",
			Content: component.NewText(fmt.Sprintf("%d", status.UpdatedReplicas)),
		},
	}...)

	return summary, nil
}

func deploymentConditions(deployment *appsv1.Deployment) (component.Component, error) {
	if deployment == nil {
		return nil, errors.New("unable to generate conditions from a nil deployment")
	}

	table := component.NewTable("Conditions", "There are no deployment conditions!", deploymentConditionColumns)

	for _, condition := range deployment.Status.Conditions {
		row := component.TableRow{
			"Type":            component.NewText(string(condition.Type)),
			"Reason":          component.NewText(condition.Reason),
			"Status":          component.NewText(string(condition.Status)),
			"Message":         component.NewText(condition.Message),
			"Last Update":     component.NewTimestamp(condition.LastUpdateTime.Time),
			"Last Transition": component.NewTimestamp(condition.LastTransitionTime.Time),
		}

		table.Add(row)
	}

	table.Sort("Type", false)

	return table, nil
}

type actionGeneratorFunction func(*appsv1.Deployment) []component.Action

// DeploymentConfiguration generates deployment configuration.
type DeploymentConfiguration struct {
	deployment       *appsv1.Deployment
	actionGenerators []actionGeneratorFunction
}

// NewDeploymentConfiguration creates an instance of DeploymentConfiguration.
func NewDeploymentConfiguration(d *appsv1.Deployment) *DeploymentConfiguration {
	return &DeploymentConfiguration{
		deployment:       d,
		actionGenerators: []actionGeneratorFunction{editDeploymentAction},
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

		rollingUpdateText := fmt.Sprintf("Max Surge %s, Max Unavailable %s",
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

	for _, generator := range dc.actionGenerators {
		actions := generator(dc.deployment)
		for _, action := range actions {
			summary.AddAction(action)
		}
	}

	return summary, nil
}

func deploymentPods(ctx context.Context, deployment *appsv1.Deployment, options Options) (component.Component, error) {
	if deployment == nil {
		return nil, errors.New("deployment is nil")
	}

	objectStore := options.DashConfig.ObjectStore()

	if objectStore == nil {
		return nil, errors.New("objectStore is nil")
	}

	selector := labels.Set(deployment.Spec.Template.ObjectMeta.Labels)

	key := store.Key{
		Namespace:  deployment.Namespace,
		APIVersion: "v1",
		Kind:       "Pod",
		Selector:   &selector,
	}

	list, _, err := objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "list all objects for key %s", key)
	}

	podList := &corev1.PodList{}
	for i := range list.Items {

		pod := &corev1.Pod{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(list.Items[i].Object, pod)
		if err != nil {
			return nil, err
		}

		if err := copyObjectMeta(pod, &list.Items[i]); err != nil {
			return nil, errors.Wrap(err, "copy object metadata")
		}

		podList.Items = append(podList.Items, *pod)
	}

	options.DisableLabels = true
	return PodListHandler(ctx, podList, options)
}

func editDeploymentAction(deployment *appsv1.Deployment) []component.Action {
	replicas := deployment.Spec.Replicas
	if replicas == nil {
		return []component.Action{}
	}

	gvk := deployment.GroupVersionKind()
	group := gvk.Group
	version := gvk.Version
	kind := gvk.Kind

	action := component.Action{
		Name:  "Edit",
		Title: "Deployment Editor",
		Form: component.Form{
			Fields: []component.FormField{
				component.NewFormFieldNumber("Replicas", "replicas", fmt.Sprintf("%d", *replicas)),
				component.NewFormFieldHidden("group", group),
				component.NewFormFieldHidden("version", version),
				component.NewFormFieldHidden("kind", kind),
				component.NewFormFieldHidden("name", deployment.Name),
				component.NewFormFieldHidden("namespace", deployment.Namespace),
				component.NewFormFieldHidden("action", "deployment/configuration"),
			},
		},
	}

	return []component.Action{action}

}

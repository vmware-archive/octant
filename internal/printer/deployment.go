/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// DeploymentListHandler is a printFunc that lists deployments
func DeploymentListHandler(ctx context.Context, list *appsv1.DeploymentList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	ot := NewObjectTable("Deployments", "We couldn't find any deployments!", cols, opts.DashConfig.ObjectStore())
	ot.EnablePluginStatus(opts.DashConfig.PluginManager())
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

		if err := ot.AddRowForObject(ctx, &d, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// DeploymentHandler is a printFunc that prints a Deployments.
func DeploymentHandler(ctx context.Context, deployment *appsv1.Deployment, options Options) (component.Component, error) {
	o := NewObject(deployment)
	o.EnableEvents()

	dh, err := newDeploymentHandler(deployment, o)
	if err != nil {
		return nil, err
	}

	if err := dh.Config(); err != nil {
		return nil, errors.Wrap(err, "print deployment configuration")
	}
	if err := dh.Status(); err != nil {
		return nil, errors.Wrap(err, "print deployment status")
	}
	if err := dh.Pods(ctx, deployment, options); err != nil {
		return nil, errors.Wrap(err, "print deployment pods")
	}

	return o.ToComponent(ctx, options)
}

func createDeploymentSummaryStatus(deployment *appsv1.Deployment) (*component.Summary, error) {
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

type actionGeneratorFunction func(*appsv1.Deployment) ([]component.Action, error)

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
		actions, err := generator(dc.deployment)
		if err != nil {
			return nil, errors.Wrap(err, "generate deployment actions")
		}
		for _, action := range actions {
			summary.AddAction(action)
		}
	}

	return summary, nil
}

func editDeploymentAction(deployment *appsv1.Deployment) ([]component.Action, error) {
	replicas := deployment.Spec.Replicas
	if replicas == nil {
		return []component.Action{}, nil
	}

	form, err := component.CreateFormForObject(octant.ActionDeploymentConfiguration, deployment,
		component.NewFormFieldNumber("Replicas", "replicas", fmt.Sprintf("%d", *replicas)),
	)
	if err != nil {
		return nil, err
	}

	action := component.Action{
		Name:  "Edit",
		Title: "Deployment Editor",
		Form:  form,
	}

	return []component.Action{action}, nil
}

type deploymentObject interface {
	Config() error
	Status() error
	Pods(ctx context.Context, object runtime.Object, options Options) error
}

type deploymentHandler struct {
	deployment  *appsv1.Deployment
	configFunc  func(*appsv1.Deployment) (*component.Summary, error)
	summaryFunc func(*appsv1.Deployment) (*component.Summary, error)
	podFunc     func(context.Context, []runtime.Object, Options) (component.Component, error)
	object      *Object
}

var _ deploymentObject = (*deploymentHandler)(nil)

func newDeploymentHandler(deployment *appsv1.Deployment, object *Object) (*deploymentHandler, error) {
	if deployment == nil {
		return nil, errors.New("can't print a nil deployment")
	}

	if object == nil {
		return nil, errors.New("can't print deployment using a nil object printer")
	}

	dh := &deploymentHandler{
		deployment:  deployment,
		configFunc:  defaultDeploymentConfig,
		summaryFunc: defaultDeploymentSummary,
		podFunc:     defaultDeploymentPods,
		object:      object,
	}

	return dh, nil
}

func (d *deploymentHandler) Config() error {
	out, err := d.configFunc(d.deployment)
	if err != nil {
		return err
	}

	d.object.RegisterConfig(out)
	return nil
}

func defaultDeploymentConfig(deployment *appsv1.Deployment) (*component.Summary, error) {
	return NewDeploymentConfiguration(deployment).Create()
}

func (d *deploymentHandler) Status() error {
	out, err := d.summaryFunc(d.deployment)
	if err != nil {
		return err
	}

	d.object.RegisterSummary(out)
	return nil
}

func defaultDeploymentSummary(deployment *appsv1.Deployment) (*component.Summary, error) {
	return createDeploymentSummaryStatus(deployment)
}

func (d *deploymentHandler) Pods(ctx context.Context, object runtime.Object, options Options) error {
	d.object.EnablePodTemplate(d.deployment.Spec.Template)

	replicaSets, err := listReplicaSetsAsObjects(ctx, d.deployment, options)
	if replicaSets == nil || err != nil {
		return err
	}

	objectList := make([]runtime.Object, len(replicaSets))
	for i := range replicaSets {
		objectList[i] = replicaSets[i]
	}

	d.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return d.podFunc(ctx, objectList, options)
		},
	})

	return nil
}

func defaultDeploymentPods(ctx context.Context, objects []runtime.Object, options Options) (component.Component, error) {
	return createRollingPodListView(ctx, objects, options)
}

func listReplicaSetsAsObjects(ctx context.Context, object runtime.Object, options Options) ([]runtime.Object, error) {
	objectStore := options.DashConfig.ObjectStore()
	var replicaSetList []*appsv1.ReplicaSet

	accessor := meta.NewAccessor()

	namespace, err := accessor.Namespace(object)
	if err != nil {
		return nil, errors.Wrap(err, "get namespace for object")
	}

	apiVersion, err := accessor.APIVersion(object)
	if err != nil {
		return nil, errors.Wrap(err, "Get apiVersion for object")
	}

	kind, err := accessor.Kind(object)
	if err != nil {
		return nil, errors.Wrap(err, "get kind for object")
	}

	name, err := accessor.Name(object)
	if err != nil {
		return nil, errors.Wrap(err, "get name for object")
	}

	key := store.Key{
		Namespace:  namespace,
		APIVersion: "apps/v1",
		Kind:       "ReplicaSet",
	}

	list, _, err := objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "list all objects for key %+v", key)
	}

	for i := range list.Items {
		replicaSet := &appsv1.ReplicaSet{}

		err := kubernetes.FromUnstructured(&list.Items[i], replicaSet)
		if err != nil {
			return nil, err
		}

		for _, ownerReference := range replicaSet.OwnerReferences {
			if ownerReference.APIVersion == apiVersion &&
				ownerReference.Kind == kind &&
				ownerReference.Name == name {
				if *(replicaSet.Spec.Replicas) == 0 {
					continue
				}
				replicaSetList = append(replicaSetList, replicaSet)
			}
		}
	}

	objectList := make([]runtime.Object, len(replicaSetList))
	for i := range replicaSetList {
		objectList[i] = replicaSetList[i]
	}

	return objectList, nil
}

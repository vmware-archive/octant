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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

// ReplicaSetListHandler is a printFunc that lists deployments
func ReplicaSetListHandler(ctx context.Context, list *appsv1.ReplicaSetList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	tbl := component.NewTable("ReplicaSets", "We couldn't find any replica sets!", cols)

	for _, rs := range list.Items {
		row := component.TableRow{}
		nameLink, err := opts.Link.ForObject(&rs, rs.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
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
func ReplicaSetHandler(ctx context.Context, replicaSet *appsv1.ReplicaSet, options Options) (component.Component, error) {
	o := NewObject(replicaSet)
	o.EnableEvents()

	rsh, err := newReplicaSetHander(replicaSet, o)
	if err != nil {
		return nil, err
	}

	if err := rsh.Config(options); err != nil {
		return nil, errors.Wrap(err, "print replicaset configuration")
	}

	if err := rsh.Status(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print replicaset status")
	}

	if err := rsh.Pods(ctx, replicaSet, options); err != nil {
		return nil, errors.Wrap(err, "print replicaset pods")
	}

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
func (rc *ReplicaSetConfiguration) Create(options Options) (*component.Summary, error) {
	if rc == nil || rc.replicaset == nil {
		return nil, errors.New("replicaset is nil")
	}

	rs := rc.replicaset

	sections := component.SummarySections{}

	if controllerRef := metav1.GetControllerOf(rs); controllerRef != nil {
		controlledBy, err := options.Link.ForOwner(rs, controllerRef)
		if err != nil {
			return nil, err
		}
		sections = append(sections, component.SummarySection{
			Header:  "Controlled By",
			Content: controlledBy,
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
	context     context.Context
	namespace   string
	selector    *metav1.LabelSelector
	uid         types.UID
	objectStore store.Store
}

// NewReplicaSetStatus creates an instance of ReplicaSetStatus
func NewReplicaSetStatus(ctx context.Context, replicaSet *appsv1.ReplicaSet, options Options) *ReplicaSetStatus {
	if err := options.DashConfig.Validate(); err != nil {
		return nil
	}
	return &ReplicaSetStatus{
		context:     ctx,
		namespace:   replicaSet.ObjectMeta.Namespace,
		selector:    replicaSet.Spec.Selector,
		uid:         replicaSet.GetUID(),
		objectStore: options.DashConfig.ObjectStore(),
	}
}

// Create generates a replicaset status quadrant
func (replicaSetStatus *ReplicaSetStatus) Create() (*component.Quadrant, error) {
	if replicaSetStatus == nil {
		return nil, errors.New("replicaset is nil")
	}

	pods, err := listPods(replicaSetStatus.context, replicaSetStatus.namespace, replicaSetStatus.selector, replicaSetStatus.uid, replicaSetStatus.objectStore)
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

type replicaSetObject interface {
	Config(options Options) error
	Status(ctx context.Context, options Options) error
	Pods(ctx context.Context, object runtime.Object, options Options) error
}

type replicaSetHandler struct {
	replicaSet *appsv1.ReplicaSet
	configFunc func(*appsv1.ReplicaSet, Options) (*component.Summary, error)
	statusFunc func(context.Context, *appsv1.ReplicaSet, Options) (*component.Quadrant, error)
	podFunc    func(context.Context, runtime.Object, Options) (component.Component, error)
	object     *Object
}

var _ replicaSetObject = (*replicaSetHandler)(nil)

func newReplicaSetHander(replicaSet *appsv1.ReplicaSet, object *Object) (*replicaSetHandler, error) {
	if replicaSet == nil {
		return nil, errors.New("can't print a nil replicaset")
	}

	if object == nil {
		return nil, errors.New("can't print a replicaset using a nil object printer")
	}

	rh := &replicaSetHandler{
		replicaSet: replicaSet,
		configFunc: defaultReplicaSetConfig,
		statusFunc: defaultReplicaSetStatus,
		podFunc:    defaultReplicaSetPods,
		object:     object,
	}

	return rh, nil
}

func (r *replicaSetHandler) Config(options Options) error {
	out, err := r.configFunc(r.replicaSet, options)
	if err != nil {
		return err
	}
	r.object.RegisterConfig(out)
	return nil
}

func defaultReplicaSetConfig(replicaSet *appsv1.ReplicaSet, options Options) (*component.Summary, error) {
	return NewReplicaSetConfiguration(replicaSet).Create(options)
}

func (r *replicaSetHandler) Status(ctx context.Context, options Options) error {
	if r.replicaSet == nil {
		return errors.New("can't display status for nil replicaset")
	}

	r.object.RegisterItems(ItemDescriptor{
		Width: component.WidthQuarter,
		Func: func() (component.Component, error) {
			return r.statusFunc(ctx, r.replicaSet, options)
		},
	})

	return nil
}

func defaultReplicaSetStatus(ctx context.Context, replicaSet *appsv1.ReplicaSet, options Options) (*component.Quadrant, error) {
	return NewReplicaSetStatus(ctx, replicaSet, options).Create()
}

func (r *replicaSetHandler) Pods(ctx context.Context, object runtime.Object, options Options) error {
	r.object.EnablePodTemplate(r.replicaSet.Spec.Template)

	r.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return r.podFunc(ctx, object, options)
		},
	})
	return nil
}

func defaultReplicaSetPods(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	return createPodListView(ctx, object, options)
}

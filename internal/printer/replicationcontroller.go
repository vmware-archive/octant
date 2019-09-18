/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

// ReplicationControllerListHandler is a printFunc that lists ReplicationControllers
func ReplicationControllerListHandler(ctx context.Context, list *corev1.ReplicationControllerList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	tbl := component.NewTable("ReplicationControllers",
		"We couldn't find any replication controllers!", cols)

	for _, rc := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&rc, rc.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink

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
	o.EnableEvents()

	rch, err := newReplicationControllerHandler(rc, o)
	if err != nil {
		return nil, err
	}

	if err := rch.Config(options); err != nil {
		return nil, errors.Wrap(err, "print replicationcontroller configuration")
	}

	if err := rch.Status(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print replicationcontroller status")
	}

	if err := rch.Pods(ctx, rc, options); err != nil {
		return nil, errors.Wrap(err, "print replicationcontroller pods")
	}

	return o.ToComponent(ctx, options)
}

// ReplicationControllerConfiguration generates a replicationcontroller configuration
type ReplicationControllerConfiguration struct {
	replicationController *corev1.ReplicationController
}

// NewReplicationControllerConfiguration creates an instance of ReplicationControllerConfiguration
func NewReplicationControllerConfiguration(rc *corev1.ReplicationController) *ReplicationControllerConfiguration {
	return &ReplicationControllerConfiguration{
		replicationController: rc,
	}
}

// Create generates a replicationcontroller configuration summary
func (rcc *ReplicationControllerConfiguration) Create(options Options) (*component.Summary, error) {
	if rcc == nil || rcc.replicationController == nil {
		return nil, errors.New("replicationcontroller is nil")
	}

	replicationController := rcc.replicationController

	sections := component.SummarySections{}

	if controllerRef := metav1.GetControllerOf(replicationController); controllerRef != nil {
		controlledBy, err := options.Link.ForOwner(replicationController, controllerRef)
		if err != nil {
			return nil, err
		}

		sections = append(sections, component.SummarySection{
			Header:  "Controlled By",
			Content: controlledBy,
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
	context     context.Context
	namespace   string
	selector    map[string]string
	uid         types.UID
	objectStore store.Store
}

// NewReplicationControllerStatus creates an instance of ReplicationControllerStatus
func NewReplicationControllerStatus(ctx context.Context, replicationController *corev1.ReplicationController, options Options) *ReplicationControllerStatus {
	return &ReplicationControllerStatus{
		context:     ctx,
		namespace:   replicationController.ObjectMeta.Namespace,
		selector:    replicationController.Spec.Selector,
		uid:         replicationController.GetUID(),
		objectStore: options.DashConfig.ObjectStore(),
	}
}

// Create generates a replicaset status quadrant
func (rcs *ReplicationControllerStatus) Create() (*component.Quadrant, error) {
	if rcs == nil {
		return nil, errors.New("replicationcontroller is nil")
	}

	selectors := metav1.LabelSelector{
		MatchLabels: rcs.selector,
	}

	pods, err := listPods(rcs.context, rcs.namespace, &selectors, rcs.uid, rcs.objectStore)
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

type replicationControllerObject interface {
	Config(options Options) error
	Status(ctx context.Context, options Options) error
	Pods(ctx context.Context, object runtime.Object, options Options) error
}

type replicationControllerHandler struct {
	replicationController *corev1.ReplicationController
	configFunc            func(*corev1.ReplicationController, Options) (*component.Summary, error)
	statusFunc            func(context.Context, *corev1.ReplicationController, Options) (*component.Quadrant, error)
	podFunc               func(context.Context, runtime.Object, Options) (component.Component, error)
	object                *Object
}

var _ replicationControllerObject = (*replicaSetHandler)(nil)

func newReplicationControllerHandler(replicationController *corev1.ReplicationController, object *Object) (*replicationControllerHandler, error) {
	if replicationController == nil {
		return nil, errors.New("can't print a nil replicationcontroller")
	}

	if object == nil {
		return nil, errors.New("can't print a replicationcontroller using a nil object printer")
	}

	rch := &replicationControllerHandler{
		replicationController: replicationController,
		configFunc:            defaultReplicationControllerConfig,
		statusFunc:            defaultReplicationControllerStatus,
		podFunc:               defaultReplicationControllerPods,
		object:                object,
	}

	return rch, nil
}

func (r *replicationControllerHandler) Config(options Options) error {
	out, err := r.configFunc(r.replicationController, options)
	if err != nil {
		return err
	}
	r.object.RegisterConfig(out)
	return nil
}

func defaultReplicationControllerConfig(replicationController *corev1.ReplicationController, options Options) (*component.Summary, error) {
	return NewReplicationControllerConfiguration(replicationController).Create(options)
}

func (r *replicationControllerHandler) Status(ctx context.Context, options Options) error {
	if r.replicationController == nil {
		return errors.New("can't display status for nil replicationcontroller")
	}

	r.object.RegisterItems(ItemDescriptor{
		Width: component.WidthQuarter,
		Func: func() (component.Component, error) {
			return r.statusFunc(ctx, r.replicationController, options)
		},
	})

	return nil
}

func defaultReplicationControllerStatus(ctx context.Context, replicationController *corev1.ReplicationController, options Options) (*component.Quadrant, error) {
	return NewReplicationControllerStatus(ctx, replicationController, options).Create()
}

func (r *replicationControllerHandler) Pods(ctx context.Context, object runtime.Object, options Options) error {
	r.object.EnablePodTemplate(*r.replicationController.Spec.Template)

	r.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return r.podFunc(ctx, object, options)
		},
	})

	return nil
}

func defaultReplicationControllerPods(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	return createPodListView(ctx, object, options)
}

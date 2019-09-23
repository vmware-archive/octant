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
	"k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

// HorizontalPodAutoscalerListHandler is a printFunc that lists horizontal pod autoscalers
func HorizontalPodAutoscalerListHandler(_ context.Context, list *v2beta2.HorizontalPodAutoscalerList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("horizontalpod handler list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Targets", "Age")
	tbl := component.NewTable("Horizontal Pod Autoscalers",
		"We couldn't find any horizontal pod autoscalers", cols)

	for _, horizontalPodAutoscaler := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&horizontalPodAutoscaler, horizontalPodAutoscaler.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(horizontalPodAutoscaler.Labels)
		row["Targets"] = component.NewText("placeholder")
		row["Age"] = component.NewTimestamp(horizontalPodAutoscaler.CreationTimestamp.Time)

		tbl.Add(row)
	}
	return tbl, nil
}

// HorizontalPodAutoscalerHandler is a printFunc that prints a HorizontalPodAutoscaler
func HorizontalPodAutoscalerHandler(ctx context.Context, horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler, options Options) (component.Component, error) {
	o := NewObject(horizontalPodAutoscaler)
	o.EnableEvents()

	hh, err := newHorizontalPodAutoscalerHandler(horizontalPodAutoscaler, o)
	if err != nil {
		return nil, err
	}

	if err := hh.Config(options); err != nil {
		return nil, errors.Wrap(err, "print horizontalpodautoscaler configuration")
	}

	if err := hh.Status(); err != nil {
		return nil, errors.Wrap(err, "print horizontalpodautoscaler status")
	}

	if err := hh.Pods(ctx, horizontalPodAutoscaler, options); err != nil {
		return nil, errors.Wrap(err, "print horizontalpodautoscaler pods")
	}

	if err := hh.Conditions(); err != nil {
		return nil, errors.Wrap(err, "print horizontalpodautoscaler conditions")
	}

	return o.ToComponent(ctx, options)
}

func createHorizontalPodAutoscalerSummaryStatus(horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler) (*component.Summary, error) {
	if horizontalPodAutoscaler == nil {
		return nil, errors.New("unable to generate status for a nil horizontalpodautoscaler")
	}

	status := horizontalPodAutoscaler.Status

	summary := component.NewSummary("Status", []component.SummarySection{
		{
			Header:  "Current Replicas",
			Content: component.NewText(fmt.Sprintf("%d", status.CurrentReplicas)),
		},
		{
			Header:  "Desired Replicas",
			Content: component.NewText(fmt.Sprintf("%d", status.DesiredReplicas)),
		},
		{
			Header:  "Targets",
			Content: component.NewText(""),
		},
	}...)

	return summary, nil
}

func createHorizontalPodAutoscalerConditionsView(horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler) (*component.Table, error) {
	if horizontalPodAutoscaler == nil {
		return nil, errors.New("unable to generate conditions from a nil horizontalpodautoscaler")
	}

	cols := component.NewTableCols("Type", "Reason", "Status", "Message", "Last Transition")
	table := component.NewTable("Conditions", "There are no horizontalpodautoscaler conditions!", cols)

	for _, condition := range horizontalPodAutoscaler.Status.Conditions {
		row := component.TableRow{
			"Type":            component.NewText(string(condition.Type)),
			"Reason":          component.NewText(condition.Reason),
			"Status":          component.NewText(string(condition.Status)),
			"Message":         component.NewText(condition.Message),
			"Last Transition": component.NewTimestamp(condition.LastTransitionTime.Time),
		}

		table.Add(row)
	}

	table.Sort("Type", false)

	return table, nil
}

// HorizontalPodAutoscalerConfiguration generates a horizontalpodautoscaler configuration
type HorizontalPodAutoscalerConfiguration struct {
	horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler
}

// NewHorizontalPodAutoscalerConfiguration creates an instance of HorizontalPodAutoscalerConfiguration
func NewHorizontalPodAutoscalerConfiguration(hpa *v2beta2.HorizontalPodAutoscaler) *HorizontalPodAutoscalerConfiguration {
	return &HorizontalPodAutoscalerConfiguration{
		horizontalPodAutoscaler: hpa,
	}
}

type horizontalPodAutoscalerObject interface {
	Config(options Options) error
	Status() error
	Pods(ctx context.Context, object runtime.Object, options Options) error
	Conditions() error
}

type horizontalPodAutoscalerHandler struct {
	horizontalPodAutoScaler *v2beta2.HorizontalPodAutoscaler
	configFunc              func(*v2beta2.HorizontalPodAutoscaler, Options) (*component.Summary, error)
	statusFunc              func(*v2beta2.HorizontalPodAutoscaler) (*component.Summary, error)
	podFunc                 func(context.Context, runtime.Object, Options) (component.Component, error)
	conditionsFunc          func(*v2beta2.HorizontalPodAutoscaler) (*component.Table, error)
	object                  *Object
}

// Create creates a horizontalpodautoscaler configuration sumamry
func (hc *HorizontalPodAutoscalerConfiguration) Create(options Options) (*component.Summary, error) {
	if hc.horizontalPodAutoscaler == nil {
		return nil, errors.New("horizontalpodautoscaler is nil")
	}

	hpa := hc.horizontalPodAutoscaler

	sections := component.SummarySections{}

	scaleTarget, err := forScaleTarget(hpa, &hpa.Spec.ScaleTargetRef, options)
	if err != nil {
		return nil, err
	}

	sections = append(sections, component.SummarySection{
		Header:  "Reference",
		Content: scaleTarget,
	})

	minReplicas := fmt.Sprintf("%d", *hpa.Spec.MinReplicas)
	maxReplicas := fmt.Sprintf("%d", hpa.Spec.MaxReplicas)
	sections.AddText("Min Replicas", minReplicas)
	sections.AddText("Max Replicas", maxReplicas)

	// for _, m := range hpa.Spec.Metrics {
	// 	TODO: implement me
	// }

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

var _ horizontalPodAutoscalerObject = (*horizontalPodAutoscalerHandler)(nil)

func newHorizontalPodAutoscalerHandler(horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler, object *Object) (*horizontalPodAutoscalerHandler, error) {
	if horizontalPodAutoscaler == nil {
		return nil, errors.New("can't print a nil horizontalpodautoscaler")
	}

	if object == nil {
		return nil, errors.New("can't print horizontalpodautoscaler using a nil object printer")
	}

	hh := &horizontalPodAutoscalerHandler{
		horizontalPodAutoScaler: horizontalPodAutoscaler,
		configFunc:              defaultHorizontalPodAutoscalerConfig,
		statusFunc:              defaultHorizontalPodAutoscalerStatus,
		podFunc:                 defaultHorizontalPodAutoscalerPods,
		conditionsFunc:          defaultHorizontalPodAutoscalerConditions,
		object:                  object,
	}

	return hh, nil
}

func (h *horizontalPodAutoscalerHandler) Config(options Options) error {
	out, err := h.configFunc(h.horizontalPodAutoScaler, options)
	if err != nil {
		return err
	}

	h.object.RegisterConfig(out)
	return nil
}

func defaultHorizontalPodAutoscalerConfig(horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler, options Options) (*component.Summary, error) {
	return NewHorizontalPodAutoscalerConfiguration(horizontalPodAutoscaler).Create(options)
}

func (h *horizontalPodAutoscalerHandler) Status() error {
	out, err := h.statusFunc(h.horizontalPodAutoScaler)
	if err != nil {
		return err
	}

	h.object.RegisterSummary(out)
	return nil
}

func defaultHorizontalPodAutoscalerStatus(horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler) (*component.Summary, error) {
	return createHorizontalPodAutoscalerSummaryStatus(horizontalPodAutoscaler)
}

func (h *horizontalPodAutoscalerHandler) Pods(ctx context.Context, object runtime.Object, options Options) error {
	o := options.DashConfig.ObjectStore()

	if o == nil {
		return errors.New("object store is nil")
	}

	if h == nil || h.horizontalPodAutoScaler == nil {
		return errors.New("horizontalpodautoscaler is nil")
	}

	targetRef := h.horizontalPodAutoScaler.Spec.ScaleTargetRef

	key := store.Key{
		Namespace:  h.horizontalPodAutoScaler.Namespace,
		APIVersion: targetRef.APIVersion,
		Kind:       targetRef.Kind,
		Name:       targetRef.Name,
	}

	target, found, err := o.Get(ctx, key)
	if err != nil {
		return errors.New("get scaled targets")
	}
	if !found {
		return nil
	}

	// TODO: what are other possible targets? Replication Controller?
	switch {
	case targetRef.Kind == "Deployment":
		deployment := &appsv1.Deployment{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(target.Object, deployment)
		if err != nil {
			return nil
		}
		h.object.EnablePodTemplate(deployment.Spec.Template)
	case targetRef.Kind == "ReplicaSet":
		replicaSet := &appsv1.ReplicaSet{}
		runtime.DefaultUnstructuredConverter.FromUnstructured(target.Object, replicaSet)
		if err != nil {
			return nil
		}
		h.object.EnablePodTemplate(replicaSet.Spec.Template)
	default:
		return nil
	}

	h.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return h.podFunc(ctx, object, options)
		},
	})

	return nil
}

func defaultHorizontalPodAutoscalerPods(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	return createPodListView(ctx, object, options)
}

func (h *horizontalPodAutoscalerHandler) Conditions() error {
	if h.horizontalPodAutoScaler == nil {
		return errors.New("can't display conditions for nil horizontalpodautoscaler")
	}

	h.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return h.conditionsFunc(h.horizontalPodAutoScaler)
		},
	})

	return nil
}

func defaultHorizontalPodAutoscalerConditions(horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler) (*component.Table, error) {
	return createHorizontalPodAutoscalerConditionsView(horizontalPodAutoscaler)
}

// forScaleTarget returns a scale target for a cross version object reference
func forScaleTarget(object runtime.Object, scaleTarget *v2beta2.CrossVersionObjectReference, options Options) (*component.Link, error) {
	if scaleTarget == nil || object == nil {
		return component.NewLink("", "none", ""), nil
	}

	accessor := meta.NewAccessor()
	ns, err := accessor.Namespace(object)
	if err != nil {
		return component.NewLink("", "none", ""), nil
	}

	return options.Link.ForGVK(
		ns,
		scaleTarget.APIVersion,
		scaleTarget.Kind,
		scaleTarget.Name,
		scaleTarget.Name,
	)
}

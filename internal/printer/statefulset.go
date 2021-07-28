/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
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

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// StatefulSetListHandler is a printFunc that list stateful sets
func StatefulSetListHandler(ctx context.Context, list *appsv1.StatefulSetList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Desired", "Current", "Age", "Selector")
	ot := NewObjectTable("StatefulSets", "We couldn't find any stateful sets!", cols, options.DashConfig.ObjectStore())
	ot.EnablePluginStatus(options.DashConfig.PluginManager())
	for _, statefulSet := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&statefulSet, statefulSet.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(statefulSet.Labels)

		desired := fmt.Sprintf("%d", *statefulSet.Spec.Replicas)
		row["Desired"] = component.NewText(desired)

		current := fmt.Sprintf("%d", statefulSet.Status.Replicas)
		row["Current"] = component.NewText(current)

		ts := statefulSet.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		row["Selector"] = printSelector(statefulSet.Spec.Selector)

		if err := ot.AddRowForObject(ctx, &statefulSet, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// StatefulSetHandler is a printFunc that prints a StatefulSet
func StatefulSetHandler(ctx context.Context, statefulSet *appsv1.StatefulSet, options Options) (component.Component, error) {
	o := NewObject(statefulSet)
	o.EnableEvents()

	sh, err := newStatefulSetHandler(statefulSet, o)
	if err != nil {
		return nil, err
	}

	if err := sh.Config(options); err != nil {
		return nil, errors.Wrap(err, "print statefulset configuration")
	}

	if err := sh.Status(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print statefulset status")
	}

	if err := sh.Pods(ctx, statefulSet, options); err != nil {
		return nil, errors.Wrap(err, "print statefulset pods")
	}

	return o.ToComponent(ctx, options)
}

// StatefulSetConfiguration generates a statefulset configuration
type StatefulSetConfiguration struct {
	statefulset *appsv1.StatefulSet
}

// NewStatefulSetConfiguration creates an instance of StatefulSetConfiguration
func NewStatefulSetConfiguration(statefulSet *appsv1.StatefulSet) *StatefulSetConfiguration {
	return &StatefulSetConfiguration{
		statefulset: statefulSet,
	}
}

// Create generates a statefulset configuration summary
func (sc *StatefulSetConfiguration) Create(options Options) (*component.Summary, error) {
	if sc == nil || sc.statefulset == nil {
		return nil, errors.New("statefulset is nil")
	}

	statefulSet := sc.statefulset

	sections := component.SummarySections{}

	sections.AddText("Update Strategy", string(statefulSet.Spec.UpdateStrategy.Type))

	if selector := statefulSet.Spec.Selector; selector != nil {
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

	total := fmt.Sprintf("%d", statefulSet.Status.Replicas)

	if desired := statefulSet.Spec.Replicas; desired != nil {
		desired := fmt.Sprintf("%d", *desired)
		status := fmt.Sprintf("%s Desired / %s Total", desired, total)
		sections.AddText("Replicas", status)
	}

	sections.AddText("Pod Management Policy", string(statefulSet.Spec.PodManagementPolicy))

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

// StatefulSetStatus generates a statefulset status
type StatefulSetStatus struct {
	context     context.Context
	namespace   string
	selector    *metav1.LabelSelector
	uid         types.UID
	objectStore store.Store
}

// NewStatefulSetStatus creates an instance of StatefulSetStatus
func NewStatefulSetStatus(ctx context.Context, statefulSet *appsv1.StatefulSet, options Options) *StatefulSetStatus {
	if err := options.DashConfig.Validate(); err != nil {
		return nil
	}

	return &StatefulSetStatus{
		context:     ctx,
		namespace:   statefulSet.ObjectMeta.Namespace,
		selector:    statefulSet.Spec.Selector,
		uid:         statefulSet.GetUID(),
		objectStore: options.DashConfig.ObjectStore(),
	}
}

// Create generates a statefulset status quadrant
func (statefulSetStatus *StatefulSetStatus) Create() (*component.Quadrant, error) {
	if statefulSetStatus == nil {
		return nil, errors.New("statefulset is nil")
	}

	pods, err := listPods(statefulSetStatus.context, statefulSetStatus.namespace, statefulSetStatus.selector, statefulSetStatus.uid, statefulSetStatus.objectStore)
	if err != nil {
		return nil, errors.Wrap(err, "list pods")
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

type statefulSetObject interface {
	Config(options Options) error
	Status(ctx context.Context, options Options) error
	Pods(ctx context.Context, object runtime.Object, options Options) error
}

type statefulSetHandler struct {
	statefulSet *appsv1.StatefulSet
	configFunc  func(*appsv1.StatefulSet, Options) (*component.Summary, error)
	statusFunc  func(context.Context, *appsv1.StatefulSet, Options) (*component.Quadrant, error)
	podFunc     func(context.Context, runtime.Object, Options) (component.Component, error)
	object      *Object
}

var _ statefulSetObject = (*statefulSetHandler)(nil)

func newStatefulSetHandler(statefulSet *appsv1.StatefulSet, object *Object) (*statefulSetHandler, error) {
	if statefulSet == nil {
		return nil, errors.New("can't print a nil statefulset")
	}

	if object == nil {
		return nil, errors.New("can't print statefulset using a nil object printer")
	}

	sh := &statefulSetHandler{
		statefulSet: statefulSet,
		configFunc:  defaultStatefulSetConfig,
		statusFunc:  defaultStatefulSetStatus,
		podFunc:     defaultStatefulSetPods,
		object:      object,
	}

	return sh, nil
}

func (s *statefulSetHandler) Config(options Options) error {
	out, err := s.configFunc(s.statefulSet, options)
	if err != nil {
		return err
	}
	s.object.RegisterConfig(out)
	return nil
}

func defaultStatefulSetConfig(statefulSet *appsv1.StatefulSet, options Options) (*component.Summary, error) {
	return NewStatefulSetConfiguration(statefulSet).Create(options)
}

func (s *statefulSetHandler) Status(ctx context.Context, options Options) error {
	if s.statefulSet == nil {
		return errors.New("can't display status for nil statefulset")
	}

	s.object.RegisterItems(ItemDescriptor{
		Width: component.WidthQuarter,
		Func: func() (component.Component, error) {
			return s.statusFunc(ctx, s.statefulSet, options)
		},
	})
	return nil
}

func defaultStatefulSetStatus(ctx context.Context, statefulSet *appsv1.StatefulSet, options Options) (*component.Quadrant, error) {
	return NewStatefulSetStatus(ctx, statefulSet, options).Create()
}

func (s *statefulSetHandler) Pods(ctx context.Context, object runtime.Object, options Options) error {
	s.object.EnablePodTemplate(s.statefulSet.Spec.Template)

	s.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return s.podFunc(ctx, object, options)
		},
	})
	return nil
}

func defaultStatefulSetPods(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	return createPodListView(ctx, object, options)
}

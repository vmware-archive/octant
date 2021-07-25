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
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// DaemonSetListHandler is a printFunc that lists daemon sets
func DaemonSetListHandler(ctx context.Context, list *appsv1.DaemonSetList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("daemon set list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Desired", "Current", "Ready",
		"Up-To-Date", "Age", "Node Selector")
	ot := NewObjectTable("Daemon Sets", "We couldn't find any daemon sets!", cols, opts.DashConfig.ObjectStore(), opts.DashConfig.TerminateThreshold())
	ot.EnablePluginStatus(opts.DashConfig.PluginManager())
	for _, daemonSet := range list.Items {
		row := component.TableRow{}
		nameLink, err := opts.Link.ForObject(&daemonSet, daemonSet.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(daemonSet.Labels)
		row["Desired"] = component.NewText(fmt.Sprintf("%d", daemonSet.Status.DesiredNumberScheduled))
		row["Current"] = component.NewText(fmt.Sprintf("%d", daemonSet.Status.CurrentNumberScheduled))
		row["Ready"] = component.NewText(fmt.Sprintf("%d", daemonSet.Status.NumberReady))
		row["Up-To-Date"] = component.NewText(fmt.Sprintf("%d", daemonSet.Status.UpdatedNumberScheduled))
		row["Age"] = component.NewTimestamp(daemonSet.ObjectMeta.CreationTimestamp.Time)
		row["Node Selector"] = printSelectorMap(daemonSet.Spec.Template.Spec.NodeSelector)

		if err := ot.AddRowForObject(ctx, &daemonSet, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// DaemonSetHandler is a printFunc that prints a daemon set
func DaemonSetHandler(ctx context.Context, daemonSet *appsv1.DaemonSet, options Options) (component.Component, error) {
	o := NewObject(daemonSet)
	o.EnableEvents()

	dsh, err := newDaemonSetHandler(daemonSet, o)
	if err != nil {
		return nil, err
	}

	if err := dsh.Config(options); err != nil {
		return nil, errors.Wrap(err, "print daemonset configuration")
	}

	if err := dsh.Status(options); err != nil {
		return nil, errors.Wrap(err, "print daemonset status")
	}

	if err := dsh.Pods(ctx, daemonSet, options); err != nil {
		return nil, errors.Wrap(err, "print daemonset pods")
	}

	return o.ToComponent(ctx, options)
}

// DaemonSetConfiguration generates a daemonset configuration
type DaemonSetConfiguration struct {
	daemonset *appsv1.DaemonSet
}

// NewDaemonSetConfiguration creates an instance of DaemonSetConfiguration
func NewDaemonSetConfiguration(ds *appsv1.DaemonSet) *DaemonSetConfiguration {
	return &DaemonSetConfiguration{
		daemonset: ds,
	}
}

// Create generates a daemonset configuration summary
func (dc *DaemonSetConfiguration) Create() (*component.Summary, error) {
	if dc == nil || dc.daemonset == nil {
		return nil, errors.New("daemon set is nil")
	}

	ds := dc.daemonset

	sections := component.SummarySections{}

	rollingUpdate := ds.Spec.UpdateStrategy.RollingUpdate
	if rollingUpdate != nil {
		rollingUpdateText := fmt.Sprintf("Max Unavailable %s",
			rollingUpdate.MaxUnavailable.String(),
		)
		sections.AddText("Update Strategy", rollingUpdateText)
	}

	if historyLimit := ds.Spec.RevisionHistoryLimit; historyLimit != nil {
		sections.AddText("Revision History Limit", fmt.Sprint(*historyLimit))
	}

	if selector := ds.Spec.Selector; selector != nil {
		sections.Add("Selectors", printSelector(selector))
	}

	if selector := ds.Spec.Template.Spec.NodeSelector; selector != nil {
		sections.Add("Node Selectors", printSelectorMap(selector))
	}

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

func createDaemonSetSummaryStatus(daemonSet *appsv1.DaemonSet) (*component.Summary, error) {
	if daemonSet == nil {
		return nil, errors.New("daemon set is nil")
	}

	sections := component.SummarySections{}

	status := daemonSet.Status
	sections.AddText("Current Number Scheduled", fmt.Sprint(status.CurrentNumberScheduled))
	sections.AddText("Desired Number Scheduled", fmt.Sprint(status.DesiredNumberScheduled))
	sections.AddText("Number Available", fmt.Sprint(status.NumberAvailable))
	sections.AddText("Number Mis-scheduled", fmt.Sprint(status.NumberMisscheduled))
	sections.AddText("Number Ready", fmt.Sprint(status.NumberReady))
	sections.AddText("Updated Number Scheduled", fmt.Sprint(status.UpdatedNumberScheduled))

	summary := component.NewSummary("Status", sections...)

	return summary, nil
}

type daemonSetObject interface {
	Config(options Options) error
	Status(options Options) error
	Pods(ctx context.Context, object runtime.Object, options Options) error
}

type daemonSetHandler struct {
	daemonSet  *appsv1.DaemonSet
	configFunc func(*appsv1.DaemonSet, Options) (*component.Summary, error)
	statusFunc func(*appsv1.DaemonSet, Options) (*component.Summary, error)
	podFunc    func(context.Context, runtime.Object, Options) (component.Component, error)
	object     *Object
}

var _ daemonSetObject = (*daemonSetHandler)(nil)

func newDaemonSetHandler(daemonSet *appsv1.DaemonSet, object *Object) (*daemonSetHandler, error) {
	if daemonSet == nil {
		return nil, errors.New("can't print a nil daemonset")
	}

	if object == nil {
		return nil, errors.New("can't print a daemonset using a nil object printer")
	}

	dh := &daemonSetHandler{
		daemonSet:  daemonSet,
		configFunc: defaultDaemonSetConfig,
		statusFunc: defaultDaemonSetSummary,
		podFunc:    defaultDaemonSetPods,
		object:     object,
	}

	return dh, nil
}

func (d *daemonSetHandler) Config(options Options) error {
	out, err := d.configFunc(d.daemonSet, options)
	if err != nil {
		return err
	}
	d.object.RegisterConfig(out)
	return nil
}

func defaultDaemonSetConfig(daemonSet *appsv1.DaemonSet, options Options) (*component.Summary, error) {
	return NewDaemonSetConfiguration(daemonSet).Create()
}

func (d *daemonSetHandler) Status(options Options) error {
	out, err := d.statusFunc(d.daemonSet, options)
	if err != nil {
		return err
	}

	d.object.RegisterSummary(out)
	return nil
}

func defaultDaemonSetSummary(daemonSet *appsv1.DaemonSet, option Options) (*component.Summary, error) {
	return createDaemonSetSummaryStatus(daemonSet)
}

func (d *daemonSetHandler) Pods(ctx context.Context, object runtime.Object, options Options) error {
	d.object.EnablePodTemplate(d.daemonSet.Spec.Template)

	d.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return d.podFunc(ctx, object, options)
		},
	})
	return nil
}

func defaultDaemonSetPods(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	return createPodListView(ctx, object, options)
}

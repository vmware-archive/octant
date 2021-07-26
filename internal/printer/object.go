/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/vmware-tanzu/octant/pkg/view/flexlayout"
)

//go:generate mockgen -destination=./fake/mock_object_interface.go -package=fake github.com/vmware-tanzu/octant/internal/printer ObjectInterface

// ObjectInterface is an interface for printing an object.
type ObjectInterface interface {
	// AddButton adds a button.
	AddButton(name string, payload action.Payload, buttonOptions ...component.ButtonOption)
	// EnableEvents enables showing events.
	EnableEvents()
	// EnableJobTemplate adds a job template.
	EnableJobTemplate(templateSpec batchv1beta1.JobTemplateSpec)
	// EnablePodTemplate adds a pod template.
	EnablePodTemplate(templateSpec corev1.PodTemplateSpec)
	// RegisterConfig registers a config summary component.
	RegisterConfig(summary *component.Summary)
	// RegisterItems registers items to be shown
	RegisterItems(items ...ItemDescriptor)
	// RegisterSummary registers a summary summary component.
	RegisterSummary(summary *component.Summary)
}

func defaultPodTemplateGen(ctx context.Context, object runtime.Object, template corev1.PodTemplateSpec, fl *flexlayout.FlexLayout, options Options) error {
	podTemplate := NewPodTemplate(object, template)
	if err := podTemplate.AddToFlexLayout(ctx, fl, options); err != nil {
		return fmt.Errorf("add pod template to layout: %w", err)
	}

	return nil
}

func defaultJobTemplateGen(ctx context.Context, object runtime.Object, template batchv1beta1.JobTemplateSpec, fl *flexlayout.FlexLayout, options Options) error {
	podTemplate := NewJobTemplate(ctx, object, template)
	if err := podTemplate.AddToFlexLayout(fl, options); err != nil {
		return fmt.Errorf("add job template to layout: %w", err)
	}

	return nil
}

func defaultEventsGen(ctx context.Context, object runtime.Object, fl *flexlayout.FlexLayout, options Options) error {
	if err := createEventsForObject(ctx, fl, object, options); err != nil {
		return fmt.Errorf("add events to layout: %w", err)
	}

	return nil
}

type conditionsGenFn func(ctx context.Context, object runtime.Object, fl *flexlayout.FlexLayout) error
type conditionsFilterFn func(*component.Table)
type mapGenFn func(runtime.Object) (map[string]interface{}, error)

func conditionsGenFactory(sortKey string, columns [][]string, mapFn mapGenFn) conditionsGenFn {
	return func(ctx context.Context, object runtime.Object, fl *flexlayout.FlexLayout) error {
		if err := createConditionsForObject(ctx, fl, object, sortKey, columns, mapFn); err != nil {
			return fmt.Errorf("add conditions to layout: %w", err)
		}

		return nil
	}
}

// ObjectPrinterFunc is a func that create a view.
type ObjectPrinterFunc func() (component.Component, error)

// ObjectPrinterLayoutFunc is a func that render a view in a flex layout.
type ObjectPrinterLayoutFunc func(*flexlayout.FlexLayout) error

// ItemDescriptor describes a func to print a view and its width.
type ItemDescriptor struct {
	Component component.Component
	Func      ObjectPrinterFunc
	Width     int
}

type podTemplateOptions struct {
	template corev1.PodTemplateSpec
}

type jobTemplateOptions struct {
	template batchv1beta1.JobTemplateSpec
}

// ObjectOpts are options for configuration Object.
type ObjectOpts func(o *Object)

// Object prints an object.
type Object struct {
	config              *component.Summary
	summary             *component.Summary
	isEventsEnabled     bool
	isConditionsEnabled bool

	itemsLists [][]ItemDescriptor

	isPodTemplateEnabled bool
	podTemplateOptions   podTemplateOptions

	isJobTemplateEnabled bool
	jobTemplateOptions   jobTemplateOptions

	object runtime.Object

	flexLayout *flexlayout.FlexLayout

	PodTemplateGen func(context.Context, runtime.Object, corev1.PodTemplateSpec, *flexlayout.FlexLayout, Options) error
	JobTemplateGen func(context.Context, runtime.Object, batchv1beta1.JobTemplateSpec, *flexlayout.FlexLayout, Options) error
	EventsGen      func(ctx context.Context, object runtime.Object, fl *flexlayout.FlexLayout, options Options) error
	ConditionsGen  func(ctx context.Context, object runtime.Object, fl *flexlayout.FlexLayout) error
}

// NewObject creates an instance of Object.
func NewObject(object runtime.Object, options ...ObjectOpts) *Object {
	o := &Object{
		object:     object,
		flexLayout: flexlayout.New(),

		PodTemplateGen: defaultPodTemplateGen,
		JobTemplateGen: defaultJobTemplateGen,
		EventsGen:      defaultEventsGen,
		ConditionsGen:  conditionsGenFactory("", nil, nil),
	}

	for _, option := range options {
		option(o)
	}

	o.isConditionsEnabled = true

	return o
}

var _ ObjectInterface = &Object{}

// RegisterConfig registers the config view for an object.
func (o *Object) RegisterConfig(summary *component.Summary) {
	o.config = summary
}

// RegisterSummary registers a summary view for an object.
func (o *Object) RegisterSummary(summary *component.Summary) {
	o.summary = summary
}

// EnablePodTemplate enables the pod template view for the object.
func (o *Object) EnablePodTemplate(templateSpec corev1.PodTemplateSpec) {
	o.isPodTemplateEnabled = true
	o.podTemplateOptions.template = templateSpec
}

// EnableJobTemplate enables the job template view for the object.
func (o *Object) EnableJobTemplate(templateSpec batchv1beta1.JobTemplateSpec) {
	o.isJobTemplateEnabled = true
	o.jobTemplateOptions.template = templateSpec
}

// EnableEvents enables the event view for the object.
func (o *Object) EnableEvents() {
	o.isEventsEnabled = true
}

// EnableConditions enables the conditions view for the object.
func (o *Object) DisableConditions() {
	o.isConditionsEnabled = false
}

// RegisterItems registers one or more items to be printed in a section.
// Each call to RegisterItems will create a new section.
func (o *Object) RegisterItems(items ...ItemDescriptor) {
	o.itemsLists = append(o.itemsLists, items)
}

func (o *Object) summaryComponent(title string, summary *component.Summary, section *flexlayout.Section, additional ...component.SummarySection) error {
	if section == nil {
		return fmt.Errorf("section is nil")
	}

	if summary == nil {
		summary = component.NewSummary(title)
	} else {
		summary.SetTitleText(title)
	}

	summary.Add(additional...)

	if len(summary.Sections()) < 1 {
		return nil
	}

	if err := section.Add(summary, component.WidthHalf); err != nil {
		return fmt.Errorf("add component to %q layout: %w", title, err)
	}

	return nil
}

// ToComponent converts Object to a view.
func (o *Object) ToComponent(ctx context.Context, options Options) (component.Component, error) {
	if o.object == nil {
		return nil, fmt.Errorf("object is nil")
	}

	summarySection := o.flexLayout.AddSection()

	pluginPrinter := options.DashConfig.PluginManager()
	if pluginPrinter == nil {
		return nil, fmt.Errorf("plugin printer is nil")
	}

	pr, err := pluginPrinter.Print(ctx, o.object)
	if err != nil {
		return nil, fmt.Errorf("plugin manager: %w", err)
	}

	if err := o.summaryComponent("Configuration", o.config, summarySection, pr.Config...); err != nil {
		return nil, fmt.Errorf("generate configuration component: %w", err)
	}

	if err := o.summaryComponent("Status", o.summary, summarySection, pr.Status...); err != nil {
		return nil, fmt.Errorf("generate summary component: %w", err)
	}

	for _, items := range o.itemsLists {
		section := o.flexLayout.AddSection()

		for _, item := range items {
			var c component.Component

			if item.Component != nil {
				c = item.Component
			} else {
				vc, err := item.Func()
				if err != nil {
					return nil, fmt.Errorf("failed to create item view: %w", err)
				}

				if vc == nil {
					// don't print nil objects
					continue
				}

				c = vc
			}

			if err := section.Add(c, item.Width); err != nil {
				return nil, fmt.Errorf("unable to add item to layout section in object printer: %w", err)
			}
		}
	}

	if len(pr.Items) > 0 {
		section := o.flexLayout.AddSection()

		for _, item := range pr.Items {
			if err := section.Add(item.View, item.Width); err != nil {
				return nil, fmt.Errorf("unable to add plugin item to layout section in object printer: %w", err)
			}
		}
	}

	if o.isPodTemplateEnabled {
		if err := o.PodTemplateGen(ctx, o.object, o.podTemplateOptions.template, o.flexLayout, options); err != nil {
			return nil, fmt.Errorf("generate pod template: %w", err)
		}
	}

	if o.isJobTemplateEnabled {
		if err := o.JobTemplateGen(ctx, o.object, o.jobTemplateOptions.template, o.flexLayout, options); err != nil {
			return nil, fmt.Errorf("generate job template: %w", err)
		}
	}

	if o.isEventsEnabled {
		if err := o.EventsGen(ctx, o.object, o.flexLayout, options); err != nil {
			return nil, fmt.Errorf("add events to layout: %w", err)
		}
	}

	if o.isConditionsEnabled {
		if err := o.ConditionsGen(ctx, o.object, o.flexLayout); err != nil {
			return nil, fmt.Errorf("add conditions to layout: %w", err)
		}
	}

	return o.flexLayout.ToComponent("Summary"), nil
}

func (o *Object) AddButton(name string, payload action.Payload, buttonOptions ...component.ButtonOption) {
	o.flexLayout.AddButton(name, payload, buttonOptions...)
}

/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"

	"github.com/pkg/errors"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/view/component"
	"github.com/vmware/octant/pkg/view/flexlayout"
)

//go:generate mockgen -source=object.go -destination=./fake/mock_object_interface.go -package=fake github.com/vmware/octant/internal/modules/overview/printer ObjectInterface

// ObjectInterface is an interface for printing an object.
type ObjectInterface interface {
	AddButton(name string, payload action.Payload, buttonOptions ...component.ButtonOption)
}

func defaultMetadataGen(object runtime.Object, fl *flexlayout.FlexLayout, options Options) error {
	metadata, err := NewMetadata(object, options.Link)
	if err != nil {
		return errors.Wrap(err, "create metadata generator")
	}

	if err := metadata.AddToFlexLayout(fl); err != nil {
		return errors.Wrap(err, "add metadata to layout")
	}

	return nil
}

func defaultPodTemplateGen(object runtime.Object, template corev1.PodTemplateSpec, fl *flexlayout.FlexLayout, options Options) error {
	podTemplate := NewPodTemplate(object, template)
	if err := podTemplate.AddToFlexLayout(fl, options); err != nil {
		return errors.Wrap(err, "add pod template to layout")
	}

	return nil
}

func defaultJobTemplateGen(object runtime.Object, template batchv1beta1.JobTemplateSpec, fl *flexlayout.FlexLayout, options Options) error {
	podTemplate := NewJobTemplate(object, template)
	if err := podTemplate.AddToFlexLayout(fl, options); err != nil {
		return errors.Wrap(err, "add job template to layout")
	}

	return nil
}

func defaultEventsGen(ctx context.Context, object runtime.Object, fl *flexlayout.FlexLayout, options Options) error {
	if err := createEventsForObject(ctx, fl, object, options); err != nil {
		return errors.Wrap(err, "add events to layout")
	}

	return nil
}

// ObjectPrinterFunc is a func that create a view.
type ObjectPrinterFunc func() (component.Component, error)

// ObjectPrinterLayoutFunc is a func that render a view in a flex layout.
type ObjectPrinterLayoutFunc func(*flexlayout.FlexLayout) error

// ItemDescriptor describes a func to print a view and its width.
type ItemDescriptor struct {
	Func  ObjectPrinterFunc
	Width int
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
	config          *component.Summary
	summary         *component.Summary
	isEventsEnabled bool

	itemsLists [][]ItemDescriptor

	isPodTemplateEnabled bool
	podTemplateOptions   podTemplateOptions

	isJobTemplateEnabled bool
	jobTemplateOptions   jobTemplateOptions

	object runtime.Object

	flexLayout *flexlayout.FlexLayout

	MetadataGen    func(runtime.Object, *flexlayout.FlexLayout, Options) error
	PodTemplateGen func(runtime.Object, corev1.PodTemplateSpec, *flexlayout.FlexLayout, Options) error
	JobTemplateGen func(runtime.Object, batchv1beta1.JobTemplateSpec, *flexlayout.FlexLayout, Options) error
	EventsGen      func(ctx context.Context, object runtime.Object, fl *flexlayout.FlexLayout, options Options) error
}

// NewObject creates an instance of Object.
func NewObject(object runtime.Object, options ...ObjectOpts) *Object {
	o := &Object{
		object:     object,
		flexLayout: flexlayout.New(),

		MetadataGen:    defaultMetadataGen,
		PodTemplateGen: defaultPodTemplateGen,
		JobTemplateGen: defaultJobTemplateGen,
		EventsGen:      defaultEventsGen,
	}

	for _, option := range options {
		option(o)
	}

	return o
}

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

// RegisterItems registers one or more items to be printed in a section.
// Each call to RegisterItems will create a new section.
func (o *Object) RegisterItems(items ...ItemDescriptor) {
	o.itemsLists = append(o.itemsLists, items)
}

func (o *Object) summaryComponent(title string, summary *component.Summary, section *flexlayout.Section, additional ...component.SummarySection) error {
	if section == nil {
		return errors.Errorf("section is nil")
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
		return errors.Wrapf(err, "add component to %q layout", title)
	}

	return nil
}

// ToComponent converts Object to a view.
func (o *Object) ToComponent(ctx context.Context, options Options) (component.Component, error) {
	if o.object == nil {
		return nil, errors.New("object is nil")
	}

	summarySection := o.flexLayout.AddSection()

	pluginPrinter := options.DashConfig.PluginManager()
	if pluginPrinter == nil {
		return nil, errors.New("plugin printer is nil")
	}

	pr, err := pluginPrinter.Print(ctx, o.object)
	if err != nil {
		return nil, errors.Wrap(err, "plugin manager")
	}

	if err := o.summaryComponent("Configuration", o.config, summarySection, pr.Config...); err != nil {
		return nil, errors.Wrap(err, "generate configuration component")
	}

	if err := o.summaryComponent("Status", o.summary, summarySection, pr.Status...); err != nil {
		return nil, errors.Wrap(err, "generate summary component")
	}

	if err := o.MetadataGen(o.object, o.flexLayout, options); err != nil {
		return nil, errors.Wrap(err, "generate metadata")
	}

	for _, items := range o.itemsLists {
		section := o.flexLayout.AddSection()

		for _, item := range items {
			vc, err := item.Func()
			if err != nil {
				return nil, errors.Wrap(err, "failed to create item view")
			}

			if err := section.Add(vc, item.Width); err != nil {
				return nil, errors.Wrap(err, "unable to add item to layout section in object printer")
			}
		}
	}

	if len(pr.Items) > 0 {
		section := o.flexLayout.AddSection()

		for _, item := range pr.Items {
			if err := section.Add(item.View, item.Width); err != nil {
				return nil, errors.Wrap(err, "unable to add plugin item to layout section in object printer")
			}
		}
	}

	if o.isPodTemplateEnabled {
		if err := o.PodTemplateGen(o.object, o.podTemplateOptions.template, o.flexLayout, options); err != nil {
			return nil, errors.Wrap(err, "generate pod template")
		}
	}

	if o.isJobTemplateEnabled {
		if err := o.JobTemplateGen(o.object, o.jobTemplateOptions.template, o.flexLayout, options); err != nil {
			return nil, errors.Wrap(err, "generate job template")
		}
	}

	if o.isEventsEnabled {
		if err := o.EventsGen(ctx, o.object, o.flexLayout, options); err != nil {
			return nil, errors.Wrap(err, "add events to layout")
		}
	}

	return o.flexLayout.ToComponent("Summary"), nil
}

func (o *Object) AddButton(name string, payload action.Payload, buttonOptions ...component.ButtonOption) {
	o.flexLayout.AddButton(name, payload, buttonOptions...)
}

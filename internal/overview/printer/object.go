package printer

import (
	"context"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/flexlayout"
	"github.com/pkg/errors"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func defaultMetadataGen(object runtime.Object, fl *flexlayout.FlexLayout) error {
	metadata, err := NewMetadata(object)
	if err != nil {
		return errors.Wrap(err, "create metadata generator")
	}

	if err := metadata.AddToFlexLayout(fl); err != nil {
		return errors.Wrap(err, "add metadata to layout")
	}

	return nil
}

func defaultPodTemplateGen(object runtime.Object, template corev1.PodTemplateSpec, fl *flexlayout.FlexLayout) error {
	podTemplate := NewPodTemplate(object, template)
	if err := podTemplate.AddToFlexLayout(fl); err != nil {
		return errors.Wrap(err, "add pod template to layout")
	}

	return nil
}

func defaultJobTemplateGen(object runtime.Object, template batchv1beta1.JobTemplateSpec, fl *flexlayout.FlexLayout) error {
	podTemplate := NewJobTemplate(object, template)
	if err := podTemplate.AddToFlexLayout(fl); err != nil {
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
type ObjectPrinterFunc func() (component.ViewComponent, error)

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
	config          ItemDescriptor
	summary         []ItemDescriptor
	isEventsEnabled bool

	itemsLists [][]ItemDescriptor

	isPodTemplateEnabled bool
	podTemplateOptions   podTemplateOptions

	isJobTemplateEnabled bool
	jobTemplateOptions   jobTemplateOptions

	object runtime.Object

	flexlayout *flexlayout.FlexLayout

	MetadataGen    func(runtime.Object, *flexlayout.FlexLayout) error
	PodTemplateGen func(runtime.Object, corev1.PodTemplateSpec, *flexlayout.FlexLayout) error
	JobTemplateGen func(runtime.Object, batchv1beta1.JobTemplateSpec, *flexlayout.FlexLayout) error
	EventsGen      func(ctx context.Context, object runtime.Object, fl *flexlayout.FlexLayout, options Options) error
}

// NewObject creates an instance of Object.
func NewObject(object runtime.Object, opts ...ObjectOpts) *Object {
	o := &Object{
		object:     object,
		flexlayout: flexlayout.New(),

		MetadataGen:    defaultMetadataGen,
		PodTemplateGen: defaultPodTemplateGen,
		JobTemplateGen: defaultJobTemplateGen,
		EventsGen:      defaultEventsGen,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// RegisterConfig registers the config view for an object.
func (o *Object) RegisterConfig(fn ObjectPrinterFunc, width int) {
	o.config = ItemDescriptor{Func: fn, Width: width}
}

// RegisterSummary registers a summary view for an object. You can
// call this multiple times. Summaries are printed in the same section
// as the config view.
func (o *Object) RegisterSummary(fn ObjectPrinterFunc, width int) {
	o.summary = append(o.summary, ItemDescriptor{Func: fn, Width: width})
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

// ToComponent converts Object to a view.
func (o *Object) ToComponent(ctx context.Context, options Options) (component.ViewComponent, error) {
	if o.object == nil {
		return nil, errors.New("object is nil")
	}

	if o.config.Func != nil {
		configView, err := o.config.Func()
		if err != nil {
			return nil, errors.Wrap(err, "generate config view")
		}

		configSection := o.flexlayout.AddSection()
		if configView != nil {
			if err := configSection.Add(configView, o.config.Width); err != nil {
				return nil, errors.Wrap(err, "add config view to layout")
			}
		}

		for _, summaryItem := range o.summary {
			view, err := summaryItem.Func()
			if err != nil {
				return nil, errors.Wrap(err, "generate summary item view")
			}

			if view != nil {
				if err := configSection.Add(view, summaryItem.Width); err != nil {
					return nil, errors.Wrap(err, "add summary view to layout")
				}
			}
		}
	}

	if err := o.MetadataGen(o.object, o.flexlayout); err != nil {
		return nil, errors.Wrap(err, "generate metadata")
	}

	for _, items := range o.itemsLists {
		section := o.flexlayout.AddSection()

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

	if o.isPodTemplateEnabled {
		if err := o.PodTemplateGen(o.object, o.podTemplateOptions.template, o.flexlayout); err != nil {
			return nil, errors.Wrap(err, "generate pod template")
		}
	}

	if o.isJobTemplateEnabled {
		if err := o.JobTemplateGen(o.object, o.jobTemplateOptions.template, o.flexlayout); err != nil {
			return nil, errors.Wrap(err, "generate job template")
		}
	}

	if o.isEventsEnabled {
		if err := o.EventsGen(ctx, o.object, o.flexlayout, options); err != nil {
			return nil, errors.Wrap(err, "add events to layout")
		}
	}

	return o.flexlayout.ToComponent("Summary"), nil
}

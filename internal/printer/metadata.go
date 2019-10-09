/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	link2 "github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/pkg/view/component"
	"github.com/vmware/octant/pkg/view/flexlayout"
)

type Metadata struct {
	object runtime.Object
	link   link2.Interface
}

func NewMetadata(object runtime.Object, l link2.Interface) (*Metadata, error) {
	if object == nil {
		return nil, errors.New("object is nil")
	}

	if l == nil {
		return nil, errors.New("link generator is nil")
	}

	return &Metadata{
		object: object,
		link:   l,
	}, nil
}

func (m *Metadata) AddToFlexLayout(fl *flexlayout.FlexLayout) error {
	if fl == nil {
		return errors.New("flex layout is nil")
	}

	section := fl.AddSection()

	summary, err := m.createSummary()
	if err != nil {
		return errors.Wrap(err, "create summary")
	}

	if err := section.Add(summary, component.WidthFull); err != nil {
		return errors.Wrap(err, "add summary to layout")
	}

	return nil
}

func (m *Metadata) createSummary() (*component.Summary, error) {
	sections := component.SummarySections{}

	object, ok := m.object.(metav1.Object)
	if !ok {
		return nil, errors.New("object is a meta v1 object")
	}

	sections.Add("Age", component.NewTimestamp(object.GetCreationTimestamp().Time))

	if labels := object.GetLabels(); len(labels) > 0 {
		sections.Add("Labels", component.NewLabels(labels))
	}

	if annotations := object.GetAnnotations(); len(annotations) > 0 {
		sections.Add("Annotations", component.NewAnnotations(annotations))
	}

	ownerReference := metav1.GetControllerOf(object)
	if ownerReference != nil {
		controlledBy, err := m.link.ForGVK(
			object.GetNamespace(),
			ownerReference.APIVersion,
			ownerReference.Kind,
			ownerReference.Name,
			ownerReference.Name,
		)
		if err != nil {
			return nil, err
		}
		sections.Add("Controlled By", controlledBy)
	}

	summary := component.NewSummary("Metadata", sections...)
	return summary, nil
}

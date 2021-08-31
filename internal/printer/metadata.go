/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/vmware-tanzu/octant/pkg/view/flexlayout"
)

//  MetadataHandler converts object metadata to a flex layout containing object metadata.
func MetadataHandler(object runtime.Object, linkGenerator link.Interface) (*component.FlexLayout, error) {
	if object == nil {
		return nil, fmt.Errorf("can't create metadata view for nil object")
	}

	if linkGenerator == nil {
		return nil, fmt.Errorf("link generator is required")
	}

	layout := flexlayout.New()

	metadata, err := NewMetadata(object, linkGenerator)
	if err != nil {
		return nil, fmt.Errorf("create metadata generator: %v", err)
	}

	if err := metadata.AddToFlexLayout(layout); err != nil {
		return nil, fmt.Errorf("add metadata to layout: %w", err)
	}

	return layout.ToComponent("Metadata"), nil
}

// Metadata represents object metadata.
type Metadata struct {
	object runtime.Object
	link   link.Interface
}

// NewMetadata creates an instance of Metadata.
func NewMetadata(object runtime.Object, l link.Interface) (*Metadata, error) {
	if object == nil {
		return nil, fmt.Errorf("object is nil")
	}

	if l == nil {
		return nil, fmt.Errorf("link generator is nil")
	}

	return &Metadata{
		object: object,
		link:   l,
	}, nil
}

// AddToFlexLayout adds metadata to a flex layout.
func (m *Metadata) AddToFlexLayout(fl *flexlayout.FlexLayout) error {
	if fl == nil {
		return fmt.Errorf("flex layout is nil")
	}

	section := fl.AddSection()

	summary, err := m.createSummary()
	if err != nil {
		return fmt.Errorf("create summary: %w", err)
	}

	if err := section.Add(summary, component.WidthFull); err != nil {
		return fmt.Errorf("add summary to layout: %w", err)
	}

	fieldSummaryList, err := m.managedFields()
	if err != nil {
		return err
	}

	for i, _ := range fieldSummaryList {
		if err := section.Add(&fieldSummaryList[i], component.WidthFull); err != nil {
			return fmt.Errorf("add managedFields to layout: %w", err)
		}
	}
	return nil
}

func (m *Metadata) createSummary() (*component.Summary, error) {
	sections := component.SummarySections{}

	object, ok := m.object.(metav1.Object)
	if !ok {
		return nil, fmt.Errorf("object is a meta v1 object")
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

func (m *Metadata) managedFields() ([]component.Summary, error) {
	a, err := meta.Accessor(m.object)
	if err != nil {
		return nil, err
	}

	var summaryList []component.Summary
	for _, field := range a.GetManagedFields() {
		fields, err := convertFieldsToFormattedString(field.FieldsV1)
		if err != nil {
			return nil, err
		}

		var timestamp *component.Timestamp
		if field.Time != nil {
			timestamp = component.NewTimestamp(field.Time.Rfc3339Copy().UTC())
		}

		summary := component.NewSummary(field.Manager, []component.SummarySection{
			{
				Header:  "Operation",
				Content: component.NewText(string(field.Operation)),
			},
			{
				Header:  "Updated",
				Content: timestamp,
			},
			{
				Header:  "Fields",
				Content: component.NewJSONEditor(fields, false),
			},
		}...)
		summaryList = append(summaryList, *summary)
	}
	return summaryList, nil
}

// convertFieldsToFormattedString formats managed fields
func convertFieldsToFormattedString(field *metav1.FieldsV1) (string, error) {
	if field == nil {
		return "", fmt.Errorf("cannot convert nil field")
	}
	var out bytes.Buffer
	if err := json.Indent(&out, field.Raw, "", "\t"); err != nil {
		return "", err
	}
	return string(out.Bytes()), nil
}

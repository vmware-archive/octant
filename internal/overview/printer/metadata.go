package printer

import (
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/heptio/developer-dash/pkg/view/flexlayout"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Metadata struct {
	object runtime.Object
}

func NewMetadata(object runtime.Object) (*Metadata, error) {
	if object == nil {
		return nil, errors.New("object is nil")
	}

	return &Metadata{
		object: object,
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
	var sections component.SummarySections

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
		sections.Add("Controlled By",
			link.ForGVK(
				object.GetNamespace(),
				ownerReference.APIVersion,
				ownerReference.Kind,
				ownerReference.Name,
				ownerReference.Name,
			),
		)
	}

	summary := component.NewSummary("Metadata", sections...)
	return summary, nil
}

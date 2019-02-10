package printer

import (
	"github.com/heptio/developer-dash/internal/view/flexlayout"
	corev1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/gridlayout"
)

type PodTemplate struct {
	podTemplateSpec corev1.PodTemplateSpec
}

func NewPodTemplate(podTemplateSpec corev1.PodTemplateSpec) *PodTemplate {
	return &PodTemplate{
		podTemplateSpec: podTemplateSpec,
	}
}

func (pt *PodTemplate) AddToFlexLayout(fl *flexlayout.FlexLayout) error {
	if fl == nil {
		return errors.New("flex layout is nil")
	}

	headerSection := fl.AddSection()
	podTemplateHeader := NewPodTemplateHeader(pt.podTemplateSpec.ObjectMeta.Labels)
	headerLabels, err := podTemplateHeader.Create()
	if err != nil {
		return err
	}

	if err := headerSection.Add(headerLabels, 23); err != nil {
		return errors.Wrap(err, "add pod template header")
	}

	containerSection := fl.AddSection()

	for _, container := range pt.podTemplateSpec.Spec.Containers {
		containerConfig := NewContainerConfiguration(&container)
		summary, err := containerConfig.Create()
		if err != nil {
			return err
		}

		if err := containerSection.Add(summary, 16); err != nil {
			return errors.Wrap(err, "add container")
		}
	}

	return nil
}

func (pt *PodTemplate) AddToGridLayout(gl *gridlayout.GridLayout) error {
	if gl == nil {
		return errors.New("grid layout is nil")
	}

	headerSection := gl.CreateSection(2)

	podTemplateHeader := NewPodTemplateHeader(pt.podTemplateSpec.ObjectMeta.Labels)
	headerLabels, err := podTemplateHeader.Create()
	if err != nil {
		return err
	}

	headerSection.Add(headerLabels, 23)

	containerSection := gl.CreateSection(16)

	for _, container := range pt.podTemplateSpec.Spec.Containers {
		containerConfig := NewContainerConfiguration(&container)
		summary, err := containerConfig.Create()
		if err != nil {
			return err
		}

		containerSection.Add(summary, 16)
	}

	return nil
}

// PodTemplateHeader creates a pod template header. It consists of a
// selectors component with title `Pod Template` and the associated
// match selectors.
type PodTemplateHeader struct {
	labels map[string]string
}

// NewPodTemplateHeader creates an instance of PodTemplateHeader.
func NewPodTemplateHeader(labels map[string]string) *PodTemplateHeader {
	return &PodTemplateHeader{
		labels: labels,
	}
}

// Create creates a labels component.
func (pth *PodTemplateHeader) Create() (*component.Labels, error) {
	view := component.NewLabels(pth.labels)
	view.Metadata.SetTitleText("Pod Template")

	return view, nil
}

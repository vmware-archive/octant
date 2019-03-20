package printer

import (
	"github.com/heptio/developer-dash/pkg/view/flexlayout"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/pkg/view/component"
)

type PodTemplate struct {
	parent          runtime.Object
	podTemplateSpec corev1.PodTemplateSpec
}

func NewPodTemplate(parent runtime.Object, podTemplateSpec corev1.PodTemplateSpec) *PodTemplate {
	return &PodTemplate{
		parent:          parent,
		podTemplateSpec: podTemplateSpec,
	}
}

func (pt *PodTemplate) AddToFlexLayout(fl *flexlayout.FlexLayout, options Options) error {
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
		containerConfig := NewContainerConfiguration(pt.parent, &container, options.PortForward, false)
		summary, err := containerConfig.Create()
		if err != nil {
			return err
		}

		if err := containerSection.Add(summary, 12); err != nil {
			return errors.Wrap(err, "add container")
		}
	}

	podSection := fl.AddSection()

	volumeTable, err := printVolumes(pt.podTemplateSpec.Spec.Volumes)
	if err != nil {
		return errors.Wrap(err, "print volumes")
	}
	if err := podSection.Add(volumeTable, 12); err != nil {
		return err
	}

	tolerationList, err := printTolerations(pt.podTemplateSpec.Spec)
	if err != nil {
		return errors.Wrap(err, "print tolerations")
	}
	if err := podSection.Add(tolerationList, 12); err != nil {
		return err
	}

	affinityList, err := printAffinity(pt.podTemplateSpec.Spec)
	if err != nil {
		return errors.Wrap(err, "print affinities")
	}
	if err := podSection.Add(affinityList, 12); err != nil {
		return err
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

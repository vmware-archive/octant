package printer

import (
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/flexlayout"
	"github.com/pkg/errors"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type JobTemplate struct {
	parent          runtime.Object
	jobTemplateSpec batchv1beta1.JobTemplateSpec
}

func NewJobTemplate(parent runtime.Object, jobTemplateSpec batchv1beta1.JobTemplateSpec) *JobTemplate {
	return &JobTemplate{
		parent:          parent,
		jobTemplateSpec: jobTemplateSpec,
	}
}

func (jt *JobTemplate) AddToFlexLayout(fl *flexlayout.FlexLayout, options Options) error {
	if fl == nil {
		return errors.New("flex layout is nil")
	}

	headerSection := fl.AddSection()
	jobTemplateHeader := NewJobTemplateHeader(jt.jobTemplateSpec.ObjectMeta.Labels)
	headerLabels, err := jobTemplateHeader.Create()
	if err != nil {
		return err
	}

	if err := headerSection.Add(headerLabels, 23); err != nil {
		return errors.Wrap(err, "add job template header")
	}

	containerSection := fl.AddSection()

	for _, container := range jt.jobTemplateSpec.Spec.Template.Spec.Containers {
		containerConfig := NewContainerConfiguration(jt.parent, &container, options.PortForward, false)
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

// JobTemplateHeader creates a job template header
type JobTemplateHeader struct {
	labels map[string]string
}

// NewJobTemplateHeader creates an instance of JobTemplateHeader
func NewJobTemplateHeader(labels map[string]string) *JobTemplateHeader {
	return &JobTemplateHeader{
		labels: labels,
	}
}

// Create creates a label component
func (jth *JobTemplateHeader) Create() (*component.Labels, error) {
	view := component.NewLabels(jth.labels)
	view.Metadata.SetTitleText("Job Template")

	return view, nil
}

package printer

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/view/component"
)

// PodTemplateHeader creates a pod template header. It consists of a
// selectors component with title `Pod Template` and the associated
// match selectors.
type PodTemplateHeader struct {
	labelSelector *metav1.LabelSelector
}

// NewPodTemplateHeader creates an instance of PodTemplateHeader.
func NewPodTemplateHeader(ls *metav1.LabelSelector) *PodTemplateHeader {
	return &PodTemplateHeader{
		labelSelector: ls,
	}
}

// Create creates a selectors component.
func (pth *PodTemplateHeader) Create() (*component.Selectors, error) {
	selectors, err := buildSelectors(pth.labelSelector)
	if err != nil {
		return nil, err
	}

	selectors.Metadata.Title = "Pod Template"

	return selectors, nil
}

package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type ContainerSummary struct{}

var _ View = (*ContainerSummary)(nil)

func NewContainerSummary(prefix, namespace string, c clock.Clock) View {
	return &ContainerSummary{}
}

func (js *ContainerSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	podTemplate, err := podTemplateSpec(object)
	if err != nil {
		return nil, err
	}

	contents, err := printPodTemplate(podTemplate, nil)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func podTemplateSpec(object runtime.Object) (*core.PodTemplateSpec, error) {
	switch o := object.(type) {
	case *batch.CronJob:
		return &o.Spec.JobTemplate.Spec.Template, nil
	case *extensions.DaemonSet:
		return &o.Spec.Template, nil
	case *extensions.Deployment:
		return &o.Spec.Template, nil
	case *batch.Job:
		return &o.Spec.Template, nil
	case *extensions.ReplicaSet:
		return &o.Spec.Template, nil
	case *core.ReplicationController:
		return o.Spec.Template, nil
	case *apps.StatefulSet:
		return &o.Spec.Template, nil
	default:
		return nil, errors.Errorf("can't extract pod template spec from %T", object)
	}
}

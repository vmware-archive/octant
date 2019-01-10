package overview

import (
	"context"
	"github.com/heptio/developer-dash/internal/cache"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type ContainerSummary struct{}

var _ View = (*ContainerSummary)(nil)

func NewContainerSummary(prefix, namespace string, c clock.Clock) View {
	return &ContainerSummary{}
}

func (js *ContainerSummary) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
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

func podTemplateSpec(object runtime.Object) (*corev1.PodTemplateSpec, error) {
	switch o := object.(type) {
	case *batchv1beta1.CronJob:
		return &o.Spec.JobTemplate.Spec.Template, nil
	case *appsv1.DaemonSet:
		return &o.Spec.Template, nil
	case *appsv1.Deployment:
		return &o.Spec.Template, nil
	case *batchv1.Job:
		return &o.Spec.Template, nil
	case *appsv1.ReplicaSet:
		return &o.Spec.Template, nil
	case *corev1.ReplicationController:
		return o.Spec.Template, nil
	case *appsv1.StatefulSet:
		return &o.Spec.Template, nil
	default:
		return nil, errors.Errorf("can't extract pod template spec from %T", object)
	}
}

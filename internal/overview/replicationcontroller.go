package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/content"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type ReplicationControllerSummary struct{}

var _ View = (*ReplicationControllerSummary)(nil)

func NewReplicationControllerSummary(prefix, namespace string, c clock.Clock) View {
	return &ReplicationControllerSummary{}
}

func (js *ReplicationControllerSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	rc, err := retrieveReplicationController(object)
	if err != nil {
		return nil, err
	}

	s := &metav1.LabelSelector{
		MatchLabels: rc.Spec.Selector,
	}

	pods, err := listPods(rc.GetNamespace(), s, rc.GetUID(), c)
	if err != nil {
		return nil, err
	}

	detail, err := printReplicationControllerSummary(rc, pods)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	contents := []content.Content{
		&summary,
	}

	return contents, nil
}

func retrieveReplicationController(object runtime.Object) (*corev1.ReplicationController, error) {
	rc, ok := object.(*corev1.ReplicationController)
	if !ok {
		return nil, errors.Errorf("expected object to be a ReplicationController, it was %T", object)
	}

	return rc, nil
}

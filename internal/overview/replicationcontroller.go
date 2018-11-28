package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/content"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/apis/core"
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

func retrieveReplicationController(object runtime.Object) (*core.ReplicationController, error) {
	rc, ok := object.(*core.ReplicationController)
	if !ok {
		return nil, errors.Errorf("expected object to be a ReplicationController, it was %T", object)
	}

	return rc, nil
}

package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type ReplicaSetSummary struct{}

var _ view.View = (*ReplicaSetSummary)(nil)

func NewReplicaSetSummary(prefix, namespace string, c clock.Clock) view.View {
	return &ReplicaSetSummary{}
}

func (rss *ReplicaSetSummary) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	replicaSet, err := retrieveReplicaSet(object)
	if err != nil {
		return nil, err
	}

	return rss.summary(replicaSet, c)
}

func (rss *ReplicaSetSummary) summary(replicaSet *appsv1.ReplicaSet, c cache.Cache) ([]content.Content, error) {
	pods, err := listPods(replicaSet.GetNamespace(), replicaSet.Spec.Selector, replicaSet.GetUID(), c)
	if err != nil {
		return nil, err
	}

	section, err := printReplicaSetSummary(replicaSet, pods)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{section})
	contents := []content.Content{
		&summary,
	}

	return contents, nil
}

func retrieveReplicaSet(object runtime.Object) (*appsv1.ReplicaSet, error) {
	replicaSet, ok := object.(*appsv1.ReplicaSet)
	if !ok {
		return nil, errors.Errorf("expected object to be a ReplicaSet, it was %T", object)
	}

	return replicaSet, nil
}

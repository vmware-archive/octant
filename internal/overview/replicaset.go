package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type ReplicaSetSummary struct{}

var _ View = (*ReplicaSetSummary)(nil)

func NewReplicaSetSummary() *ReplicaSetSummary {
	return &ReplicaSetSummary{}
}

func (rss *ReplicaSetSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	replicaSet, err := retrieveReplicaSet(object)
	if err != nil {
		return nil, err
	}

	return rss.summary(replicaSet, c)
}

func (rss *ReplicaSetSummary) summary(replicaSet *extensions.ReplicaSet, c Cache) ([]content.Content, error) {
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

func retrieveReplicaSet(object runtime.Object) (*extensions.ReplicaSet, error) {
	replicaSet, ok := object.(*extensions.ReplicaSet)
	if !ok {
		return nil, errors.Errorf("expected object to be a ReplicaSet, it was %T", object)
	}

	return replicaSet, nil
}

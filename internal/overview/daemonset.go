package overview

import (
	"context"
	"github.com/heptio/developer-dash/internal/cache"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type DaemonSetSummary struct{}

var _ View = (*DaemonSetSummary)(nil)

func NewDaemonSetSummary(prefix, namespace string, c clock.Clock) View {
	return &DaemonSetSummary{}
}

func (rss *DaemonSetSummary) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	replicaSet, err := retrieveDaemonSet(object)
	if err != nil {
		return nil, err
	}

	return rss.summary(replicaSet, c)
}

func (rss *DaemonSetSummary) summary(replicaSet *appsv1.DaemonSet, c cache.Cache) ([]content.Content, error) {
	pods, err := listPods(replicaSet.GetNamespace(), replicaSet.Spec.Selector, replicaSet.GetUID(), c)
	if err != nil {
		return nil, err
	}

	section, err := printDaemonSetSummary(replicaSet, pods)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{section})
	contents := []content.Content{
		&summary,
	}

	return contents, nil
}

func retrieveDaemonSet(object runtime.Object) (*appsv1.DaemonSet, error) {
	replicaSet, ok := object.(*appsv1.DaemonSet)
	if !ok {
		return nil, errors.Errorf("expected object to be a DaemonSet, it was %T", object)
	}

	return replicaSet, nil
}

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

type StatefulSetSummary struct{}

var _ view.View = (*StatefulSetSummary)(nil)

func NewStatefulSetSummary(prefix, namespace string, c clock.Clock) view.View {
	return &StatefulSetSummary{}
}

func (js *StatefulSetSummary) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	ss, err := retrieveStatefulSet(object)
	if err != nil {
		return nil, err
	}

	pods, err := listPods(ss.GetNamespace(), ss.Spec.Selector, ss.GetUID(), c)
	if err != nil {
		return nil, err
	}

	detail, err := printStatefulSetSummary(ss, pods)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	contents := []content.Content{
		&summary,
	}

	return contents, nil
}

func retrieveStatefulSet(object runtime.Object) (*appsv1.StatefulSet, error) {
	rc, ok := object.(*appsv1.StatefulSet)
	if !ok {
		return nil, errors.Errorf("expected object to be a StatefulSet, it was %T", object)
	}

	return rc, nil
}

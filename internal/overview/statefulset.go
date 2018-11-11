package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/content"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/apps"
)

type StatefulSetSummary struct{}

var _ View = (*StatefulSetSummary)(nil)

func NewStatefulSetSummary() *StatefulSetSummary {
	return &StatefulSetSummary{}
}

func (js *StatefulSetSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
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
	return []content.Content{
		&summary,
	}, nil
}

func retrieveStatefulSet(object runtime.Object) (*apps.StatefulSet, error) {
	rc, ok := object.(*apps.StatefulSet)
	if !ok {
		return nil, errors.Errorf("expected object to be a StatefulSet, it was %T", object)
	}

	return rc, nil
}

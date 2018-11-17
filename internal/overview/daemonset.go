package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type DaemonSetSummary struct{}

var _ View = (*DaemonSetSummary)(nil)

func NewDaemonSetSummary() *DaemonSetSummary {
	return &DaemonSetSummary{}
}

func (rss *DaemonSetSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	replicaSet, err := retrieveDaemonSet(object)
	if err != nil {
		return nil, err
	}

	return rss.summary(replicaSet, c)
}

func (rss *DaemonSetSummary) summary(replicaSet *extensions.DaemonSet, c Cache) ([]content.Content, error) {
	pods, err := listPods(replicaSet.GetNamespace(), replicaSet.Spec.Selector, replicaSet.GetUID(), c)
	if err != nil {
		return nil, err
	}

	section, err := printDaemonSetSummary(replicaSet, pods)
	if err != nil {
		return nil, err
	}

	podTemplate, err := printPodTemplate(&replicaSet.Spec.Template, nil)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{section})
	contents := []content.Content{
		&summary,
	}

	contents = append(contents, podTemplate...)
	return contents, nil
}

func retrieveDaemonSet(object runtime.Object) (*extensions.DaemonSet, error) {
	replicaSet, ok := object.(*extensions.DaemonSet)
	if !ok {
		return nil, errors.Errorf("expected object to be a DaemonSet, it was %T", object)
	}

	return replicaSet, nil
}

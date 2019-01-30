package printer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
)

type TolerationDescriber struct {
	podSpec corev1.PodSpec
}

func NewTolerationDescriber(podSpec corev1.PodSpec) *TolerationDescriber {
	return &TolerationDescriber{
		podSpec: podSpec,
	}
}

func (td *TolerationDescriber) Create() (*component.List, error) {
	var items []component.ViewComponent

	for _, toleration := range td.podSpec.Tolerations {
		msg, err := td.describe(toleration)
		if err != nil {
			return nil, err
		}

		if evictSecs := toleration.TolerationSeconds; evictSecs != nil {
			msg = fmt.Sprintf("%s Evict after %d seconds.", msg, *evictSecs)
		}

		items = append(items, component.NewText("", msg))
	}

	list := component.NewList("", items)

	return list, nil
}

const (
	taintFmt    = "Schedule on nodes with %s taint."
	wildcardFmt = "Schedule on all nodes."
)

func (td *TolerationDescriber) describe(toleration corev1.Toleration) (string, error) {
	msgFmt := taintFmt
	var a []interface{}

	if toleration.Effect != "" && toleration.Key == "" && toleration.Value == "" {
		a = append(a, string(toleration.Effect))
	} else if toleration.Key != "" && toleration.Value != "" && toleration.Effect == "" {
		a = append(a, fmt.Sprintf("%s:%s", toleration.Key, toleration.Value))
	} else if toleration.Key != "" && toleration.Value != "" && toleration.Effect != "" {
		a = append(a, fmt.Sprintf("%s:%s:%s", toleration.Key, toleration.Value, toleration.Effect))
	} else if toleration.Key != "" && toleration.Operator == corev1.TolerationOpExists {
		a = append(a, fmt.Sprintf("%s", toleration.Key))
	} else if toleration.Key == "" && toleration.Operator == corev1.TolerationOpExists {
		msgFmt = wildcardFmt
	} else {
		return "", errors.Errorf("unable to describe toleration")
	}

	return fmt.Sprintf(msgFmt, a...), nil
}

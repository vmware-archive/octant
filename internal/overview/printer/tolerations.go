package printer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
)

func printTolerations(podSpec corev1.PodSpec) (component.Component, error) {
	td := &tolerationDescriber{
		podSpec: podSpec,
	}

	return td.Create()
}

type tolerationDescriber struct {
	podSpec corev1.PodSpec
}

func (td *tolerationDescriber) Create() (*component.Table, error) {
	cols := component.NewTableCols("Description")
	table := component.NewTable("Taints and Tolerations", cols)

	for _, toleration := range td.podSpec.Tolerations {
		msg, err := td.describe(toleration)
		if err != nil {
			return nil, err
		}

		if evictSecs := toleration.TolerationSeconds; evictSecs != nil {
			msg = fmt.Sprintf("%s Evict after %d seconds.", msg, *evictSecs)
		}

		table.Add(component.TableRow{"Description": component.NewText(msg)})
	}

	return table, nil
}

const (
	taintFmt    = "Schedule on nodes with %s taint."
	wildcardFmt = "Schedule on all nodes."
)

func (td *tolerationDescriber) describe(toleration corev1.Toleration) (string, error) {
	msgFmt := taintFmt
	var a []interface{}

	if toleration.Effect != "" && toleration.Key == "" && toleration.Value == "" {
		a = append(a, string(toleration.Effect))
	} else if toleration.Key != "" && toleration.Value != "" && toleration.Effect == "" {
		a = append(a, fmt.Sprintf("%s:%s", toleration.Key, toleration.Value))
	} else if toleration.Key != "" && toleration.Value != "" && toleration.Effect != "" {
		a = append(a, fmt.Sprintf("%s:%s:%s", toleration.Key, toleration.Value, toleration.Effect))
	} else if toleration.Key != "" && toleration.Operator == corev1.TolerationOpExists {
		a = append(a, toleration.Key)
	} else if toleration.Key == "" && toleration.Operator == corev1.TolerationOpExists {
		msgFmt = wildcardFmt
	} else {
		return "", errors.Errorf("unable to describe toleration")
	}

	return fmt.Sprintf(msgFmt, a...), nil
}

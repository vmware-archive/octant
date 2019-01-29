package printer

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
)

func buildSelectors(selector *metav1.LabelSelector) (*component.Selectors, error) {
	if selector == nil {
		return nil, errors.Errorf("selectors was nil")
	}

	var selectors []component.Selector

	for _, lsr := range selector.MatchExpressions {
		o, err := component.MatchOperator(string(lsr.Operator))
		if err != nil {
			return nil, err
		}

		es := component.NewExpressionSelector(lsr.Key, o, lsr.Values)
		selectors = append(selectors, es)
	}

	for k, v := range selector.MatchLabels {
		ls := component.NewLabelSelector(k, v)
		selectors = append(selectors, ls)
	}

	selectorsComponent := component.NewSelectors(selectors)

	return selectorsComponent, nil
}

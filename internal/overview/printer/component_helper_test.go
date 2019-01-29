package printer

import (
	"testing"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_buildSelectors(t *testing.T) {
	cases := []struct {
		name          string
		labelSelector *metav1.LabelSelector
		expected      *component.Selectors
		isErr         bool
	}{
		{
			name: "in general",
			labelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "key",
						Operator: "In",
						Values:   []string{"value"},
					},
				},
				MatchLabels: map[string]string{
					"key": "value",
				},
			},
			expected: component.NewSelectors([]component.Selector{
				component.NewExpressionSelector("key", component.OperatorIn, []string{"value"}),
				component.NewLabelSelector("key", "value"),
			}),
		},
		{
			name: "invalid expression operator",
			labelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "key",
						Operator: "invalid",
						Values:   []string{"value"},
					},
				},
			},
			isErr: true,
		},
		{
			name:          "nil label selector",
			labelSelector: nil,
			isErr:         true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildSelectors(tc.labelSelector)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

package component_test

import (
	"testing"

	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchOperator(t *testing.T) {
	cases := []struct {
		name     string
		s        string
		isErr    bool
		expected component.Operator
	}{
		{
			name:     "existing operator",
			s:        "In",
			expected: component.OperatorIn,
		},
		{
			name:  "invalid operator",
			s:     "Invalid",
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o, err := component.MatchOperator(tc.s)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, o)
		})
	}
}

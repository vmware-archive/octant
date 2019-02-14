package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SelectorFromFilters(t *testing.T) {
	tests := []struct {
		name      string
		filters   []string
		expected  string
		expectErr bool
	}{
		{
			name:     "simple",
			filters:  []string{"app:nginx"},
			expected: "app=nginx",
		},
		{
			name:     "empty",
			filters:  []string{},
			expected: "",
		},
		{
			name:     "multiple",
			filters:  []string{"app:nginx", "env:production"},
			expected: "app=nginx,env=production",
		},
		{
			name:      "invalid",
			filters:   []string{"app=nginx"},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := selectorFromFilters(tc.filters)
			hadErr := (err != nil)
			if !assert.Equalf(t, tc.expectErr, hadErr, "unexpected error: %v", err) || hadErr {
				return
			}
			assert.Equal(t, tc.expected, s.String())
		})
	}
}

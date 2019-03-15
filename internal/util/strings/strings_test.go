package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	cases := []struct {
		name     string
		s        string
		sl       []string
		expected bool
	}{
		{
			name:     "does contain",
			s:        "1",
			sl:       []string{"1", "2", "3"},
			expected: true,
		},
		{
			name:     "does not contain",
			s:        "4",
			sl:       []string{"1", "2", "3"},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Contains(tc.s, tc.sl)
			assert.Equal(t, tc.expected, got)
		})
	}
}

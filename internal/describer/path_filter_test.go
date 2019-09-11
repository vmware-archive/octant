package describer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathFilter_Match(t *testing.T) {
	tests := []struct {
		name        string
		filterPath  string
		contentPath string
		isMatch     bool
	}{
		{
			name:        "case 1",
			filterPath:  "/foo/bar",
			contentPath: "/foo/bar",
			isMatch:     true,
		},
		{
			name:        "case 2",
			filterPath:  "/foo/bar",
			contentPath: "/namespace/default/foo/bar",
			isMatch:     true,
		},
		{
			name:        "case 3",
			filterPath:  "/",
			contentPath: "/namespace/default",
			isMatch:     true,
		},
		{
			name:        "case 4",
			filterPath:  "/workloads/cron-jobs/(?P<name>^[^.]+)",
			contentPath: "/workloads/cron-jobs",
			isMatch:     false,
		},
		{
			name:        "case 5",
			filterPath:  "/",
			contentPath: "/namespace/kube-system/workloads/pods/coredns-5c98db65d4-7hmsh",
			isMatch:     false,
		},
		{
			name:        "case 6",
			filterPath:  `/workloads/pods/(?P<name>.*?)`,
			contentPath: "/namespace/kube-system/workloads/pods/coredns-5c98db65d4-7hmsh",
			isMatch:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pf := NewPathFilter(test.filterPath, nil)
			assert.Equal(t, test.isMatch, pf.Match(test.contentPath),
				"did not match for filterPath:[%s], contentPath[%s]",
				test.filterPath, test.contentPath)
		})
	}
}

func TestPathFilter_Fields(t *testing.T) {
	tests := []struct {
		name        string
		filterPath  string
		contentPath string
		expected    map[string]string
	}{
		{
			name:        "path",
			filterPath:  `/workloads/cron-jobs/(?P<name>.*?)`,
			contentPath: "/workloads/cron-jobs/name",
			expected: map[string]string{
				"name":      "name",
				"namespace": "",
			},
		},
		{
			name:        "path with namespace",
			filterPath:  `/workloads/cron-jobs/(?P<name>.*?)`,
			contentPath: "/namespace/default/workloads/cron-jobs/name",
			expected: map[string]string{
				"name":      "name",
				"namespace": "default",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pf := NewPathFilter(test.filterPath, nil)
			got := pf.Fields(test.contentPath)

			assert.Equal(t, test.expected, got)
		})
	}
}

package plugin

import (
	"testing"

	"github.com/heptio/developer-dash/internal/gvk"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestCapabilities_SupportsPrinter(t *testing.T) {
	cases := []struct {
		name         string
		in           schema.GroupVersionKind
		capabilities Capabilities
		hasSupport   bool
	}{
		{
			name: "with printer support",
			in:   gvk.PodGVK,
			capabilities: Capabilities{
				SupportsPrinterConfig: []schema.GroupVersionKind{gvk.PodGVK},
			},
			hasSupport: true,
		},
		{
			name: "with out printer support",
			in:   gvk.DeploymentGVK,
			capabilities: Capabilities{
				SupportsPrinterConfig: []schema.GroupVersionKind{gvk.PodGVK},
			},
			hasSupport: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.hasSupport, tc.capabilities.SupportsPrinter(tc.in))
		})
	}
}

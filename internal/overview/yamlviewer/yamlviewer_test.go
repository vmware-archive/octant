package yamlviewer

import (
	"testing"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
)

func Test_ToComponent(t *testing.T) {
	object := &corev1.Pod{}

	got, err := ToComponent(object)
	require.NoError(t, err)

	expected := &component.YAML{
		Metadata: component.Metadata{
			Title: []component.TitleViewComponent{
				component.NewText("YAML"),
			},
			Type: "yaml",
		},
		Config: component.YAMLConfig{
			Data: "---\nmetadata:\n  creationTimestamp: null\nspec:\n  containers: null\nstatus: {}\n",
		},
	}

	assert.Equal(t, expected, got)
}

package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/heptio/developer-dash/pkg/view/component"
)

func TestPodTemplateHeader(t *testing.T) {
	labels := map[string]string{
		"key": "value",
	}

	pth := NewPodTemplateHeader(labels)
	got, err := pth.Create()

	require.NoError(t, err)

	assert.Len(t, got.Config.Labels, 1)

	expected := component.Title(component.NewText("Pod Template"))

	assert.Equal(t, expected, got.Metadata.Title)
}

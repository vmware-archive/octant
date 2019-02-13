package printer_test

import (
	"testing"

	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobTemplateHeader(t *testing.T) {
	labels := map[string]string{
		"app": "myapp",
	}

	jth := printer.NewJobTemplateHeader(labels)
	got, err := jth.Create()

	require.NoError(t, err)

	assert.Len(t, got.Config.Labels, 1)

	expected := []component.TitleViewComponent{
		component.NewText("Job Template"),
	}

	assert.Equal(t, expected, got.Metadata.Title)
}

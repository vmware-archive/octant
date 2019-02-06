package printer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/gridlayout"
)

func TestPodTemplate(t *testing.T) {
	gl := gridlayout.New()
	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"key": "value",
			},
		},
	}

	pt := printer.NewPodTemplate(podTemplateSpec)
	err := pt.AddToGridLayout(gl)
	require.NoError(t, err)

	got := gl.ToGrid()

	headerLabels := component.NewLabels(map[string]string{"key": "value"})
	headerLabels.Metadata.Title = "Pod Template"
	headerLabelsPanel := component.NewPanel("", headerLabels)
	headerLabelsPanel.Position(0, 0, 23, 2)

	panels := []component.Panel{
		*headerLabelsPanel,
	}

	expected := component.NewGrid("Summary", panels...)

	assert.Equal(t, expected, got)
}

func TestPodTemplateHeader(t *testing.T) {
	labels := map[string]string{
		"key": "value",
	}

	pth := printer.NewPodTemplateHeader(labels)
	got, err := pth.Create()

	require.NoError(t, err)

	assert.Len(t, got.Config.Labels, 1)
	assert.Equal(t, "Pod Template", got.Metadata.Title)
}

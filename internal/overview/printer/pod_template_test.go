package printer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/overview/printer"
)

func TestPodTemplateHeader(t *testing.T) {
	labelSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"key": "value",
		},
	}

	pth := printer.NewPodTemplateHeader(labelSelector)
	got, err := pth.Create()

	require.NoError(t, err)

	assert.Len(t, got.Config.Selectors, 1)
	assert.Equal(t, "Pod Template", got.Metadata.Title)
}

func TestPodTemplateHeader_nil_label_selector(t *testing.T) {
	pth := printer.NewPodTemplateHeader(nil)
	_, err := pth.Create()
	assert.Error(t, err)
}

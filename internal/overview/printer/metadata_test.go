package printer_test

import (
	"testing"

	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/flexlayout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Metadata(t *testing.T) {
	fl := flexlayout.New()

	deployment := testutil.CreateDeployment("deployment")
	metadata, err := printer.NewMetadata(deployment)
	require.NoError(t, err)

	require.NoError(t, metadata.AddToFlexLayout(fl))

	got := fl.ToComponent("Summary")

	expected := &component.FlexLayout{
		Metadata: component.Metadata{
			Title: []component.TitleViewComponent{component.NewText("Summary")},
			Type:  "flexlayout",
		},
		Config: component.FlexLayoutConfig{
			Sections: []component.FlexLayoutSection{
				{
					{
						Width: 16,
						View: component.NewSummary("Metadata", component.SummarySections{
							{
								Header:  "Age",
								Content: component.NewTimestamp(deployment.CreationTimestamp.Time),
							},
						}...),
					},
				},
			},
		},
	}

	assert.Equal(t, expected, got)
}

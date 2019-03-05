package printer

import (
	"testing"

	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/flexlayout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Metadata(t *testing.T) {
	fl := flexlayout.New()

	deployment := testutil.CreateDeployment("deployment")
	metadata, err := NewMetadata(deployment)
	require.NoError(t, err)

	require.NoError(t, metadata.AddToFlexLayout(fl))

	got := fl.ToComponent("Summary")

	expected := component.NewFlexLayout("Summary")
	expected.AddSections([]component.FlexLayoutSection{
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
	}...)

	assert.Equal(t, expected, got)
}

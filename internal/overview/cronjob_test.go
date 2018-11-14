package overview

import (
	"context"
	"testing"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCronJobSummary_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewCronJobSummary())
}

func TestCronJobSummary(t *testing.T) {
	s := NewCronJobSummary()

	ctx := context.Background()
	cache := NewMemoryCache()

	cronJob := loadFromFile(t, "cronjob-1.yaml")
	cronJob = convertToInternal(t, cronJob)

	storeFromFile(t, "job-1.yaml", cache)

	contents, err := s.Content(ctx, cronJob, cache)
	require.NoError(t, err)

	sections := []content.Section{
		{
			Items: []content.Item{
				content.TextItem("Name", "hello"),
				content.TextItem("Namespace", "default"),
				content.LabelsItem("Labels", map[string]string{"overview": "default"}),
				content.LabelsItem("Annotations", map[string]string{}),
				content.TimeItem("Create Time", "2018-09-18T12:30:09Z"),
				content.TextItem("Active", "0"),
				content.TextItem("Schedule", "*/1 * * * *"),
				content.TextItem("Suspend", "false"),
				content.TimeItem("Last Schedule", "2018-11-02T09:45:00Z"),
				content.TextItem("Concurrency Policy", "Allow"),
				content.TextItem("Starting Deadline Seconds", "<unset>"),
			},
		},
	}
	summary := content.NewSummary("Details", sections)

	expected := []content.Content{
		&summary,
	}

	assert.Equal(t, expected, contents)
}

func TestCronJobJobs(t *testing.T) {
	cjj := NewCronJobJobs()

	ctx := context.Background()
	cache := NewMemoryCache()

	cronJob := loadFromFile(t, "cronjob-1.yaml")
	cronJob = convertToInternal(t, cronJob)

	storeFromFile(t, "job-1.yaml", cache)

	contents, err := cjj.Content(ctx, cronJob, cache)
	require.NoError(t, err)

	jobColumns := tableCols("Name", "Desired", "Successful", "Age", "Containers",
		"Images", "Selector", "Labels")

	activeTable := content.NewTable("Active Jobs")
	activeTable.Columns = jobColumns

	inactiveTable := content.NewTable("Inactive Jobs")
	inactiveTable.Columns = jobColumns
	inactiveTable.AddRow(content.TableRow{
		"Age":        content.NewStringText("1d"),
		"Containers": content.NewStringText("hello"),
		"Desired":    content.NewStringText("1"),
		"Images":     content.NewStringText("busybox"),
		"Labels":     content.NewStringText("controller-uid=f20be17b-de8b-11e8-889a-025000000001,job-name=hello-1541155320"),
		"Name":       content.NewLinkText("hello-1541155320", "/content/overview/workloads/jobs/hello-1541155320"),
		"Selector":   content.NewStringText("controller-uid=f20be17b-de8b-11e8-889a-025000000001"),
		"Successful": content.NewStringText("1"),
	})

	expected := []content.Content{
		&activeTable,
		&inactiveTable,
	}

	assert.Equal(t, expected, contents)
}

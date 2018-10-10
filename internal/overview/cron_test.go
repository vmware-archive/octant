package overview

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestCronJobsDescriber(t *testing.T) {
	namespace := "default"

	cache := NewMemoryCache()
	loadUnstructured(t, cache, namespace, "cronjob.yaml")

	d := NewCronJobsDescriber()
	got, err := d.Describe("/prefix", namespace, cache, nil)
	require.NoError(t, err)

	require.Len(t, got, 1)
	tbl, ok := got[0].(table)
	require.True(t, ok)

	assert.Equal(t, tbl.Title, "Cron Jobs")
	assert.Len(t, tbl.Rows, 1)
}

func TestCronJobDescriber(t *testing.T) {
	namespace := "default"

	cache := NewMemoryCache()
	loadUnstructured(t, cache, namespace, "cronjob.yaml")
	loadUnstructured(t, cache, namespace, "event-1.yaml")

	fields := map[string]string{
		"name": "hello",
	}

	d := NewCronJobDescriber()
	got, err := d.Describe("/prefix", namespace, cache, fields)
	require.NoError(t, err)

	require.Len(t, got, 2)
	cjTable, ok := got[0].(table)
	require.True(t, ok)

	assert.Equal(t, cjTable.Title, "Cron Job")
	assert.Len(t, cjTable.Rows, 1)

	eventsTable, ok := got[1].(table)
	require.True(t, ok)

	assert.Equal(t, eventsTable.Title, "Events")
	assert.Len(t, eventsTable.Rows, 1)
}

func Test_printCronJob(t *testing.T) {
	ti := time.Unix(1538828130, 0)
	c := clock.NewFakeClock(ti)

	cases := []struct {
		name     string
		path     string
		expected tableRow
	}{
		{
			name: "not scheduled",
			path: "cronjob.yaml",
			expected: tableRow{
				"active":        newStringText("0"),
				"age":           newStringText("<unknown>"),
				"labels":        newLabelsText(nil),
				"last_schedule": newStringText("<none>"),
				"name":          newLinkText("hello", "/api/workloads/cron-jobs/hello?namespace=default"),
				"schedule":      newStringText("*/1 * * * *"),
			},
		},
		{
			name: "scheduled",
			path: "cronjob-scheduled.yaml",
			expected: tableRow{
				"active":        newStringText("0"),
				"age":           newStringText("<unknown>"),
				"labels":        newLabelsText(nil),
				"last_schedule": newStringText("30s"),
				"name":          newLinkText("hello", "/api/workloads/cron-jobs/hello?namespace=default"),
				"schedule":      newStringText("*/1 * * * *"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cj, ok := loadType(t, tc.path).(*batchv1beta1.CronJob)
			require.True(t, ok)

			got := printCronJob(cj, "/api", "default", c)
			assert.Equal(t, tc.expected, got)
		})
	}
}

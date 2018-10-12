package overview

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestEventsDescriber(t *testing.T) {
	namespace := "default"

	cache := NewMemoryCache()
	loadUnstructured(t, cache, namespace, "event-1.yaml")
	loadUnstructured(t, cache, namespace, "event-2.yaml")

	d := NewEventsDescriber("/events")
	got, err := d.Describe("/prefix", namespace, cache, nil)
	require.NoError(t, err)

	require.Len(t, got, 1)
	tbl, ok := got[0].(table)
	require.True(t, ok)

	assert.Equal(t, tbl.Title, "Events")
	assert.Len(t, tbl.Rows, 2)
}

func Test_printEvent(t *testing.T) {
	ti := time.Unix(1538828130, 0)
	c := clock.NewFakeClock(ti)

	cases := []struct {
		name     string
		path     string
		expected tableRow
	}{
		{
			name: "event",
			path: "event-1.yaml",
			expected: tableRow{
				"message":    newStringText("(combined from similar events): Saw completed job: hello-1538868300"),
				"source":     newStringText("cronjob-controller"),
				"sub_object": newStringText(""),
				"count":      newStringText("24973"),
				"first_seen": newStringText("2018-09-18T12:40:18Z"),
				"last_seen":  newStringText("2018-10-06T23:25:55Z"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			event, ok := loadType(t, tc.path).(*corev1.Event)
			require.True(t, ok)

			got := printEvent(event, "/api", "default", c)
			assert.Equal(t, tc.expected, got)
		})
	}
}

package overview

import (
	"context"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/content"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestEventsDescriber(t *testing.T) {
	namespace := "default"

	c := cache.NewMemoryCache()
	loadUnstructured(t, c, namespace, "event-1.yaml")
	loadUnstructured(t, c, namespace, "event-2.yaml")

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache: c,
	}

	ctx := context.Background()
	d := NewEventsDescriber("/events")
	cResponse, err := d.Describe(ctx, "/prefix", namespace, clusterClient, options)
	require.NoError(t, err)

	table := content.NewTable("Events", "Namespace Events does not contain any events")
	table.Columns = []content.TableColumn{
		{Name: "Message", Accessor: "message"},
		{Name: "Source", Accessor: "source"},
		{Name: "Sub-Object", Accessor: "sub_object"},
		{Name: "Count", Accessor: "count"},
		{Name: "First Seen", Accessor: "first_seen"},
		{Name: "Last Seen", Accessor: "last_seen"},
	}
	table.AddRow(content.TableRow{
		"count":      content.NewStringText("24973"),
		"first_seen": content.NewStringText("2018-09-18T12:40:18Z"),
		"last_seen":  content.NewStringText("2018-10-06T23:25:55Z"),
		"message":    content.NewStringText("(combined from similar events): Saw completed job: hello-1538868300"),
		"source":     content.NewStringText("cronjob-controller"),
		"sub_object": content.NewStringText(""),
	})
	table.AddRow(content.TableRow{
		"count":      content.NewStringText("24973"),
		"first_seen": content.NewStringText("2018-09-18T12:40:18Z"),
		"last_seen":  content.NewStringText("2018-10-06T23:25:55Z"),
		"message":    content.NewStringText("(combined from similar events): Saw completed job: hello-1538868300"),
		"source":     content.NewStringText("cronjob-controller"),
		"sub_object": content.NewStringText(""),
	})

	expected := component.ContentResponse{}
	assert.Equal(t, expected, cResponse)
}

func Test_printEvent(t *testing.T) {
	ti := time.Unix(1538828130, 0)
	c := clock.NewFakeClock(ti)

	cases := []struct {
		name     string
		path     string
		expected content.TableRow
	}{
		{
			name: "event",
			path: "event-1.yaml",
			expected: content.TableRow{
				"message":    content.NewStringText("(combined from similar events): Saw completed job: hello-1538868300"),
				"source":     content.NewStringText("cronjob-controller"),
				"sub_object": content.NewStringText(""),
				"count":      content.NewStringText("24973"),
				"first_seen": content.NewStringText("2018-09-18T12:40:18Z"),
				"last_seen":  content.NewStringText("2018-10-06T23:25:55Z"),
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

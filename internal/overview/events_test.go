package overview

import (
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestEventsDescriber(t *testing.T) {
	namespace := "default"

	cache := NewMemoryCache()
	loadUnstructured(t, cache, namespace, "event-1.yaml")
	loadUnstructured(t, cache, namespace, "event-2.yaml")

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache: cache,
	}

	d := NewEventsDescriber("/events")
	got, err := d.Describe("/prefix", namespace, clusterClient, options)
	require.NoError(t, err)

	require.Len(t, got, 1)
	tbl, ok := got[0].(*content.Table)
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

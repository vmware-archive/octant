package overview

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/content"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

// EventsDescriber creates content for a list of events.
type EventsDescriber struct {
	*baseDescriber

	path      string
	title     string
	cacheKeys []CacheKey
}

// NewEventsDescriber creates an instance of EventsDescriber.
func NewEventsDescriber(p string) *EventsDescriber {
	return &EventsDescriber{
		baseDescriber: newBaseDescriber(),
		path:          p,
		title:         "Events",
		cacheKeys: []CacheKey{
			{
				APIVersion: "v1",
				Kind:       "Event",
			},
		},
	}
}

// Describe creates content.
func (d *EventsDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (ContentResponse, error) {
	objects, err := loadObjects(ctx, options.Cache, namespace, options.Fields, d.cacheKeys)
	if err != nil {
		return emptyContentResponse, err
	}

	var contents []content.Content

	t := newEventTable(d.title)

	sort.Slice(objects, func(i, j int) bool {
		tsI := objects[i].GetCreationTimestamp()
		tsJ := objects[j].GetCreationTimestamp()

		return tsI.Before(&tsJ)
	})

	for _, object := range objects {
		event := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, event)
		if err != nil {
			return emptyContentResponse, err
		}

		t.Rows = append(t.Rows, printEvent(event, prefix, namespace, d.clock()))
	}

	contents = append(contents, &t)

	return ContentResponse{
		Contents: contents,
		Title:    d.title,
	}, nil
}

func (d *EventsDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

func newEventTable(name string) content.Table {
	t := content.NewTable(name)

	t.Columns = []content.TableColumn{
		{Name: "Message", Accessor: "message"},
		{Name: "Source", Accessor: "source"},
		{Name: "Sub-Object", Accessor: "sub_object"},
		{Name: "Count", Accessor: "count"},
		{Name: "First Seen", Accessor: "first_seen"},
		{Name: "Last Seen", Accessor: "last_seen"},
	}

	return t
}

func printEvent(event *corev1.Event, prefix, namespace string, c clock.Clock) content.TableRow {
	firstSeen := event.FirstTimestamp.UTC().Format(time.RFC3339)
	lastSeen := event.LastTimestamp.UTC().Format(time.RFC3339)

	return content.TableRow{
		"message":    content.NewStringText(event.Message),
		"source":     content.NewStringText(event.Source.Component),
		"sub_object": content.NewStringText(""), // TODO: where does this come from?
		"count":      content.NewStringText(fmt.Sprint(event.Count)),
		"first_seen": content.NewStringText(firstSeen),
		"last_seen":  content.NewStringText(lastSeen),
	}
}

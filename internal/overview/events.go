package overview

import (
	"fmt"
	"sort"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

// EventsDescriber creates content for a list of events.
type EventsDescriber struct {
	*baseDescriber

	cacheKeys []CacheKey
}

// NewEventsDescriber creates an instance of EventsDescriber.
func NewEventsDescriber() *EventsDescriber {
	return &EventsDescriber{
		baseDescriber: newBaseDescriber(),
		cacheKeys: []CacheKey{
			{
				APIVersion: "v1",
				Kind:       "Event",
			},
		},
	}
}

// Describe creates content.
func (d *EventsDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error) {
	objects, err := loadObjects(cache, namespace, fields, d.cacheKeys)
	if err != nil {
		return nil, err
	}

	var contents []Content

	t := newEventTable("Events")

	sort.Slice(objects, func(i, j int) bool {
		tsI := objects[i].GetCreationTimestamp()
		tsJ := objects[j].GetCreationTimestamp()

		return tsI.Before(&tsJ)
	})

	for _, object := range objects {
		event := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, event)
		if err != nil {
			return nil, err
		}

		t.Rows = append(t.Rows, printEvent(event, prefix, namespace, d.clock()))
	}

	contents = append(contents, t)

	return contents, nil
}

func newEventTable(name string) table {
	t := newTable(name)

	t.Columns = []tableColumn{
		{Name: "Message", Accessor: "message"},
		{Name: "Source", Accessor: "source"},
		{Name: "Sub-Object", Accessor: "sub_object"},
		{Name: "Count", Accessor: "count"},
		{Name: "First Seen", Accessor: "first_seen"},
		{Name: "Last Seen", Accessor: "last_seen"},
	}

	return t
}

func printEvent(event *corev1.Event, prefix, namespace string, c clock.Clock) tableRow {
	firstSeen := event.FirstTimestamp.UTC().Format(time.RFC3339)
	lastSeen := event.LastTimestamp.UTC().Format(time.RFC3339)

	return tableRow{
		"message":    newStringText(event.Message),
		"source":     newStringText(event.Source.Component),
		"sub_object": newStringText(""), // TODO: where does this come from?
		"count":      newStringText(fmt.Sprint(event.Count)),
		"first_seen": newStringText(firstSeen),
		"last_seen":  newStringText(lastSeen),
	}
}

package printer

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/heptio/developer-dash/pkg/store"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/heptio/developer-dash/pkg/view/flexlayout"
)

var (
	objectEventCols = component.NewTableCols("Message", "Reason", "Type", "First Seen", "Last Seen", "From", "Count")
)

// EventListHandler is a printFunc that lists events.
func EventListHandler(ctx context.Context, list *corev1.EventList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Kind", "Message", "Reason", "Type",
		"First Seen", "Last Seen")
	table := component.NewTable("Events", cols)

	for _, event := range list.Items {
		row := component.TableRow{}

		objectPath, err := ObjectReferencePath(event.InvolvedObject)
		if err != nil {
			return nil, err
		}

		var kind component.Component = component.NewLink("",
			fmt.Sprintf("%s (%d)", event.InvolvedObject.Name, event.Count), objectPath)
		if objectPath == "" {
			kind = component.NewText(
				fmt.Sprintf("%s (%d)", event.InvolvedObject.Name, event.Count))
		}

		row["Kind"] = kind

		messageLink, err := opts.Link.ForObject(&event, event.Message)
		if err != nil {
			return nil, err
		}

		row["Message"] = messageLink
		row["Reason"] = component.NewText(event.Reason)
		row["Type"] = component.NewText(event.Type)
		row["First Seen"] = component.NewTimestamp(event.FirstTimestamp.Time)
		row["Last Seen"] = component.NewTimestamp(event.LastTimestamp.Time)

		table.Add(row)
	}

	table.Sort("Last Seen", true)

	return table, nil
}

func EventHandler(ctx context.Context, event *corev1.Event, opts Options) (component.Component, error) {
	if event == nil {
		return nil, errors.New("event can not be nil")
	}

	var detailSections []component.SummarySection

	detailSections = append(detailSections, component.SummarySection{
		Header:  "Last Seen",
		Content: component.NewTimestamp(event.LastTimestamp.Time),
	})

	detailSections = append(detailSections, component.SummarySection{
		Header:  "First Seen",
		Content: component.NewTimestamp(event.FirstTimestamp.Time),
	})

	detailSections = append(detailSections, component.SummarySection{
		Header:  "Count",
		Content: component.NewText(fmt.Sprintf("%d", event.Count)),
	})

	detailSections = append(detailSections, component.SummarySection{
		Header:  "Message",
		Content: component.NewText(event.Message),
	})

	detailSections = append(detailSections, component.SummarySection{
		Header:  "Kind",
		Content: component.NewText(event.InvolvedObject.Kind),
	})

	// NOTE: object reference can contain a field path to the object,
	// and that is not reported.
	refPath, err := ObjectReferencePath(event.InvolvedObject)
	if err != nil {
		return nil, err
	}

	if refPath != "" {
		detailSections = append(detailSections, component.SummarySection{
			Header:  "Involved Object",
			Content: component.NewLink("", event.InvolvedObject.Name, refPath),
		})
	}

	detailSections = append(detailSections, component.SummarySection{
		Header:  "Type",
		Content: component.NewText(event.Type),
	})

	detailSections = append(detailSections, component.SummarySection{
		Header:  "Reason",
		Content: component.NewText(event.Reason),
	})

	sourceMsg := event.Source.Component
	if event.Source.Host != "" {
		sourceMsg = fmt.Sprintf("%s on %s",
			event.Source.Component, event.Source.Host)
	}
	detailSections = append(detailSections, component.SummarySection{
		Header:  "Source",
		Content: component.NewText(sourceMsg),
	})

	summary := component.NewSummary("Event Detail", detailSections...)

	fl := flexlayout.New()
	summarySection := fl.AddSection()
	if err := summarySection.Add(summary, component.WidthFull); err != nil {
		return nil, errors.Wrap(err, "add event to layout")
	}

	return fl.ToComponent("Event"), nil
}

// PrintEvents collects events for a resource
func PrintEvents(list *corev1.EventList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	table := component.NewTable("Events", objectEventCols)

	for _, event := range list.Items {
		row := component.TableRow{}

		row["Message"] = component.NewText(event.Message)
		row["Reason"] = component.NewText(event.Reason)
		row["Type"] = component.NewText(event.Type)

		row["First Seen"] = component.NewTimestamp(event.FirstTimestamp.Time)
		row["Last Seen"] = component.NewTimestamp(event.LastTimestamp.Time)

		row["From"] = component.NewText(formatEventSource(event.Source))

		count := fmt.Sprintf("%d", event.Count)
		row["Count"] = component.NewText(count)

		table.Add(row)
	}

	return table, nil
}

// formatEventSource formats EventSource as a comma separated string excluding Host when empty
func formatEventSource(es corev1.EventSource) string {
	EventSourceString := []string{es.Component}
	if len(es.Host) > 0 {
		EventSourceString = append(EventSourceString, es.Host)
	}
	return strings.Join(EventSourceString, ", ")
}

func createEventsForObject(ctx context.Context, fl *flexlayout.FlexLayout, object runtime.Object, opts Options) error {
	objectStore := opts.DashConfig.ObjectStore()
	eventList, err := eventsForObject(ctx, object, objectStore)
	if err != nil {
		return errors.Wrap(err, "list events for object")
	}

	if len(eventList.Items) > 0 {
		eventTable, err := PrintEvents(eventList, opts)
		if err != nil {
			return errors.Wrap(err, "create event table for object")
		}

		eventsSection := fl.AddSection()
		if err := eventsSection.Add(eventTable, 24); err != nil {
			return errors.Wrap(err, "add event table to layout")
		}
	}

	return nil
}

func eventsForObject(ctx context.Context, object runtime.Object, o store.Store) (*corev1.EventList, error) {
	accessor := meta.NewAccessor()

	namespace, err := accessor.Namespace(object)
	if err != nil {
		return nil, errors.Wrap(err, "get namespace for object")
	}

	apiVersion, err := accessor.APIVersion(object)
	if err != nil {
		return nil, errors.Wrap(err, "Get apiVersion for object")
	}

	kind, err := accessor.Kind(object)
	if err != nil {
		return nil, errors.Wrap(err, "get kind for object")
	}

	name, err := accessor.Name(object)
	if err != nil {
		return nil, errors.Wrap(err, "get name for object")
	}

	key := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Event",
	}

	list, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "list events for object")
	}

	eventList := &corev1.EventList{}

	for _, unstructuredEvent := range list {
		event := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredEvent.Object, event)
		if err != nil {
			return nil, err
		}

		involvedObject := event.InvolvedObject
		if involvedObject.Namespace == namespace &&
			involvedObject.APIVersion == apiVersion &&
			involvedObject.Kind == kind &&
			involvedObject.Name == name {
			eventList.Items = append(eventList.Items, *event)
		}
	}

	sort.SliceStable(eventList.Items, func(i, j int) bool {
		a := eventList.Items[i]
		b := eventList.Items[j]

		return a.LastTimestamp.After(b.LastTimestamp.Time)
	})

	return eventList, nil
}

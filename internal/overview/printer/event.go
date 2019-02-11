package printer

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/gridlayout"
)

// EventListHandler is a printFunc that lists events.
func EventListHandler(list *corev1.EventList, opts Options) (component.ViewComponent, error) {
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

		infoItems := []component.ViewComponent{
			component.NewLink("", event.InvolvedObject.Name, objectPath),
			component.NewText(fmt.Sprintf("%d", event.Count)),
		}
		info := component.NewList("", infoItems)

		row["Kind"] = info
		eventPath := gvkPath(event.Namespace, event.TypeMeta.APIVersion, event.TypeMeta.Kind, event.Name)
		row["Message"] = component.NewLink("", event.Message, eventPath)
		row["Reason"] = component.NewText(event.Reason)
		row["Type"] = component.NewText(event.Type)
		row["First Seen"] = component.NewTimestamp(event.FirstTimestamp.Time)
		row["Last Seen"] = component.NewTimestamp(event.LastTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}

func EventHandler(event *corev1.Event, opts Options) (component.ViewComponent, error) {
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
	detailSections = append(detailSections, component.SummarySection{
		Header:  "Involved Object",
		Content: component.NewLink("", event.InvolvedObject.Name, refPath),
	})

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

	gl := gridlayout.New()

	summarySection := gl.CreateSection(10)
	summarySection.Add(summary, 24)

	grid := gl.ToGrid()
	return grid, nil
}

// PrintEvents collects events for a resource
func PrintEvents(list *corev1.EventList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Type", "Reason", "Age", "From", "Message")
	table := component.NewTable("Events", cols)

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

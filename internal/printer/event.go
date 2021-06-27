/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	oerrors "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/vmware-tanzu/octant/pkg/view/flexlayout"
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
	table := component.NewTable("Events", "We couldn't find any events!", cols)

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

		if !event.FirstTimestamp.Time.IsZero() {
			row["First Seen"] = component.NewTimestamp(event.FirstTimestamp.Time)
		} else {
			row["First Seen"] = component.NewText("<unknown>")
		}

		if !event.LastTimestamp.Time.IsZero() {
			row["Last Seen"] = component.NewTimestamp(event.LastTimestamp.Time)
		} else {
			row["Last Seen"] = component.NewText("<unknown>")
		}

		table.Add(row)
	}

	table.Sort("Last Seen")
	table.Reverse()

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

// PrintError prints the events table with the error
func PrintError(err error) (component.Component, error) {
	errStr := fmt.Sprintf("%s", err)

	var ae *oerrors.AccessError
	if errors.As(err, &ae) {
		errStr = fmt.Sprintf("Access Error, failed to %s: %s", ae.Verb(), ae.Key())
	}

	c := component.NewCard(component.TitleFromString("Events"))
	c.SetBody(component.NewText(errStr))

	return c, nil
}

// PrintEvents collects events for a resource
func PrintEvents(list *corev1.EventList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	table := component.NewTable("Events", "There are no events!", objectEventCols)

	sortFailed := false
	sort.SliceStable(list.Items, func(i, j int) bool {
		a, err := strconv.Atoi(list.Items[i].ResourceVersion)
		if err != nil {
			sortFailed = true
		}
		b, err := strconv.Atoi(list.Items[j].ResourceVersion)
		if err != nil {
			sortFailed = true
		}

		return b < a
	})

	if sortFailed {
		return nil, fmt.Errorf("detected invalid event resource version")
	}

	for _, event := range list.Items {
		row := component.TableRow{}

		row["Message"] = component.NewText(event.Message)
		row["Reason"] = component.NewText(event.Reason)
		row["Type"] = component.NewText(event.Type)

		if !event.FirstTimestamp.Time.IsZero() {
			row["First Seen"] = component.NewTimestamp(event.FirstTimestamp.Time)
		} else {
			row["First Seen"] = component.NewText("<unknown>")
		}

		if !event.LastTimestamp.Time.IsZero() {
			row["Last Seen"] = component.NewTimestamp(event.LastTimestamp.Time)
		} else {
			row["Last Seen"] = component.NewText("<unknown>")
		}

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
	eventList, err := store.EventsForObject(ctx, object, objectStore)
	if err != nil {
		eventError, err := PrintError(err)
		if err != nil {
			return fmt.Errorf("create event list error: %w", err)
		}
		eventSection := fl.AddSection()
		if err := eventSection.Add(eventError, component.WidthFull); err != nil {
			return fmt.Errorf("add event error to layout: %w", err)
		}
		return nil
	}

	if len(eventList.Items) > 0 {
		eventTable, err := PrintEvents(eventList, opts)
		if err != nil {
			return errors.Wrap(err, "create event table for object")
		}

		eventsSection := fl.AddSection()
		if err := eventsSection.Add(eventTable, component.WidthFull); err != nil {
			return errors.Wrap(err, "add event table to layout")
		}
	}

	return nil
}

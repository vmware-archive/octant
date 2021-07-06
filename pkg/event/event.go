/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

//go:generate mockgen -destination=./fake/mock_event.go -package=fake github.com/vmware-tanzu/octant/pkg/event WSClientGetter,WSEventSender

import (
	"fmt"
)

const (
	// EventTypeBuildInfo is a buildInfo event
	EventTypeBuildInfo EventType = "event.octant.dev/buildInfo"

	// EventTypeKubeConfigPath carries location of kube config with it
	EventTypeKubeConfigPath EventType = "event.octant.dev/kubeConfigPath"

	// EventTypeContent is a content event.
	EventTypeContent EventType = "event.octant.dev/content"

	// EventTypeNamespaces is a namespaces event.
	EventTypeNamespaces EventType = "event.octant.dev/namespaces"

	// EventTypeNavigation is a navigation event.
	EventTypeNavigation EventType = "event.octant.dev/navigation"

	// EventTypeObjectNotFound is an object not found event.
	EventTypeObjectNotFound EventType = "event.octant.dev/objectNotFound"

	// EventTypeCurrentNamespace is a current namespace event.
	EventTypeCurrentNamespace EventType = "event.octant.dev/currentNamespace"

	// EventTypeUnknown is an unknown event.
	EventTypeUnknown EventType = "event.octant.dev/unknown"

	// EventTypeNamespace is a namespace event.
	EventTypeNamespace EventType = "event.octant.dev/namespace"

	// EventTypeContext is a context event.
	EventTypeContext EventType = "event.octant.dev/context"

	// EventTypeKubeConfig is an event for updating kube contexts on the front end.
	EventTypeKubeConfig EventType = "event.octant.dev/kubeConfig"

	// EventTypeContentPath is a content path event.
	EventTypeContentPath EventType = "event.octant.dev/contentPath"

	// EventTypeFilters is a filters event.
	EventTypeFilters EventType = "event.octant.dev/filters"

	// EventTypeAlert is an alert event.
	EventTypeAlert EventType = "event.octant.dev/alert"

	// EventTypeRefresh is a refresh event.
	EventTypeRefresh EventType = "event.octant.dev/refresh"

	// EventTypeLoading is a loading event.
	EventTypeLoading EventType = "event.octant.dev/loading"

	// EventTypeAppLogs is an app logs event.
	EventTypeAppLogs EventType = "event.octant.dev/app-logs"

	// EventTypeTerminalFormat is a string with format specifiers to assist in generating
	// a terminal event type.
	EventTypeTerminalFormat string = "event.octant.dev/terminals/namespace/%s/pod/%s/container/%s"

	// EventTypeLoggingFormat is a string with format specifiers to assist in generating
	// a logging event type.
	EventTypeLoggingFormat string = "event.octant.dev/logging/namespace/%s/pod/%s"

	// EventTypeNotification sends information saved on error store
	EventTypeNotification = "event.octant.dev/notification"
)

// NewTerminalEventType returns an event type for a specific terminal instance.
// This is the Event.Type that an Octant client will watch for to read the terminal stream.
func NewTerminalEventType(namespace, pod, container string) EventType {
	return EventType(fmt.Sprintf(EventTypeTerminalFormat, namespace, pod, container))
}

// NewLoggingEventType returns an event type for pod logs.
// This is the Event.Type that an Octant client will watch for to read the logging stream.
func NewLoggingEventType(namespace, pod string) EventType {
	return EventType(fmt.Sprintf(EventTypeLoggingFormat, namespace, pod))
}

type EventType string

// Event is an event for the dash frontend.
type Event struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
	Err  error
}

func CreateEvent(eventType EventType, fields map[string]interface{}) Event {
	return Event{
		Type: eventType,
		Data: fields,
	}
}

func FindEvent(events []Event, evType EventType) (Event, error) {
	var result Event

	for _, evt := range events {
		if evt.Type == evType {
			result = evt
		}
	}

	if (result == Event{}) {
		return result, fmt.Errorf("Could not find event of type: %s", evType)
	} else {
		return result, nil
	}
}

type WSClientGetter interface {
	Get(id string) WSEventSender
}

type WSEventSender interface {
	Send(event Event)
}

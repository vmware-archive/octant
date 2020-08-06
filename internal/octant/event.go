/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package octant

import (
	"fmt"
	"github.com/vmware-tanzu/octant/pkg/event"
)

const (
	// EventTypeBuildInfo is a buildInfo event
	EventTypeBuildInfo event.EventType = "event.octant.dev/buildInfo"

	// EventTypeContent is a content event.
	EventTypeContent event.EventType = "event.octant.dev/content"

	// EventTypeNamespaces is a namespaces event.
	EventTypeNamespaces event.EventType = "event.octant.dev/namespaces"

	// EventTypeNavigation is a navigation event.
	EventTypeNavigation event.EventType = "event.octant.dev/navigation"

	// EventTypeObjectNotFound is an object not found event.
	EventTypeObjectNotFound event.EventType = "event.octant.dev/objectNotFound"

	// EventTypeCurrentNamespace is a current namespace event.
	EventTypeCurrentNamespace event.EventType = "event.octant.dev/currentNamespace"

	// EventTypeUnknown is an unknown event.
	EventTypeUnknown event.EventType = "event.octant.dev/unknown"

	// EventTypeNamespace is a namespace event.
	EventTypeNamespace event.EventType = "event.octant.dev/namespace"

	// EventTypeContext is a context event.
	EventTypeContext event.EventType = "event.octant.dev/context"

	// EventTypeKubeConfig is an event for updating kube contexts on the front end.
	EventTypeKubeConfig event.EventType = "event.octant.dev/kubeConfig"

	// EventTypeContentPath is a content path event.
	EventTypeContentPath event.EventType = "event.octant.dev/contentPath"

	// EventTypeFilters is a filters event.
	EventTypeFilters event.EventType = "event.octant.dev/filters"

	// EventTypeAlert is an alert event.
	EventTypeAlert event.EventType = "event.octant.dev/alert"

	// EventTypeRefresh is a refresh event.
	EventTypeRefresh event.EventType = "event.octant.dev/refresh"

	// EventTypeLoading is a loading event.
	EventTypeLoading event.EventType = "event.octant.dev/loading"

	// EventTypeTerminalFormat is a string with format specifiers to assist in generating
	// a terminal event type.
	EventTypeTerminalFormat string = "event.octant.dev/terminals/namespace/%s/pod/%s/container/%s"

	// EventTypeLoggingFormat is a string with format specifiers to assist in generating
	// a logging event type.
	EventTypeLoggingFormat string = "event.octant.dev/logging/namespace/%s/pod/%s"
)

// NewTerminalEventType returns an event type for a specific terminal instance.
// This is the Event.Type that an Octant client will watch for to read the terminal stream.
func NewTerminalEventType(namespace, pod, container string) event.EventType {
	return event.EventType(fmt.Sprintf(EventTypeTerminalFormat, namespace, pod, container))
}

// NewLoggingEventType returns an event type for pod logs.
// This is the Event.Type that an Octant client will watch for to read the logging stream.
func NewLoggingEventType(namespace, pod string) event.EventType {
	return event.EventType(fmt.Sprintf(EventTypeLoggingFormat, namespace, pod))
}

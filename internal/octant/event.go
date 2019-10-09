/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package octant

type EventType string

const (
	// EventTypeContent is a content event.
	EventTypeContent EventType = "content"

	// EventTypeNamespaces is a namespaces event.
	EventTypeNamespaces EventType = "namespaces"

	// EventTypeNavigation is a navigation event.
	EventTypeNavigation EventType = "navigation"

	// EventTypeObjectNotFound is an object not found event.
	EventTypeObjectNotFound EventType = "objectNotFound"

	// EventTypeCurrentNamespace is a current namespace event.
	EventTypeCurrentNamespace EventType = "currentNamespace"

	// EventTypeUnknown is an unknown event.
	EventTypeUnknown EventType = "unknown"

	// EventTypeNamespace is a namespace event.
	EventTypeNamespace EventType = "namespace"

	// EventTypeContext is a context event.
	EventTypeContext EventType = "context"

	// EventTypeKubeConfig is an event for updating kube contexts on the front end.
	EventTypeKubeConfig EventType = "kubeConfig"

	// EventTypeContentPath is a content path event.
	EventTypeContentPath EventType = "contentPath"

	// EventTypeFilters is a filters event.
	EventTypeFilters EventType = "filters"

	// EventTypeAlert is an alert event.
	EventTypeAlert EventType = "alert"
)

// Event is an event for the dash frontend.
type Event struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
	Err  error
}

const ()

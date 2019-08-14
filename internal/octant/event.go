/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
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
)

// Event is an event for the dash frontend.
type Event struct {
	Type EventType
	Data []byte
	Err  error
}

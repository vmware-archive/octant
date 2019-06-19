package octant

type EventType string

const (
	// EventTypeContent is a content event.
	EventTypeContent EventType = "content"
	// EventTypeNamespaces is a namespaces events.
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

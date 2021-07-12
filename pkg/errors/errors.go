package errors

import "time"

// InternalError represents an internal Octant error.
type InternalError interface {
	ID() string
	Error() string
	Timestamp() time.Time
	Name() string
}

type ErrorStore interface {
	List() []InternalError
	Get(string) (InternalError, bool)
	Add(InternalError) (found bool)
	AddError(error) InternalError
}

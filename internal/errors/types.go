package errors

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vmware/octant/pkg/action"
)

type InternalError struct {
	id          uuid.UUID
	timestamp   time.Time
	payload     action.Payload
	requestType string
	err         error
}

func NewInternalError(requestType string, payload action.Payload, err error) *InternalError {
	id, _ := uuid.NewUUID()

	return &InternalError{
		requestType: requestType,
		payload:     payload,
		err:         err,
		timestamp:   time.Now(),
		id:          id,
	}
}

// Error returns an error string.
func (o *InternalError) Error() string {
	return fmt.Sprintf("%s: %s", o.requestType, o.err)
}

// Client returns a client if one is available.
func (o *InternalError) RequestType() string {
	return o.requestType
}

// Request returns the payload that generated the error, if available.
func (o *InternalError) Payload() action.Payload {
	return o.payload
}

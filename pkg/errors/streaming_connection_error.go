package errors

import "github.com/vmware-tanzu/octant/internal/errors"

type StreamError struct {
	*errors.GenericError
	Fatal bool
}

func NewStreamError(err error) *StreamError {
	return &StreamError{
		errors.NewGenericError(err),
		false,
	}
}

func FatalStreamError(err error) *StreamError {
	return &StreamError{
		errors.NewGenericError(err),
		true,
	}
}

func IsFatalStreamError(err error) bool {
	switch sErr := err.(type) {
	case StreamError:
		return sErr.Fatal
	case *StreamError:
		return sErr.Fatal
	default:
		return false
	}
}

const StreamingConnectionError = "StreamingConnectionError"

func (s *StreamError) Name() string {
	return StreamingConnectionError
}

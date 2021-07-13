package errors

type StreamError struct {
	*GenericError
	Fatal bool
}

func NewStreamError(err error) *StreamError {
	return &StreamError{
		NewGenericError(err),
		false,
	}
}

func FatalStreamError(err error) *StreamError {
	return &StreamError{
		NewGenericError(err),
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

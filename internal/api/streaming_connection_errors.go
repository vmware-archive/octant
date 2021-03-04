package api

type StreamError struct {
	error
	Fatal bool
}

func FatalStreamError(err error) error {
	return StreamError{
		err,
		true,
	}
}

func IsFatalStreamError(err error) bool {
	sErr, ok := err.(StreamError)
	if !ok {
		return false
	}

	return sErr.Fatal
}

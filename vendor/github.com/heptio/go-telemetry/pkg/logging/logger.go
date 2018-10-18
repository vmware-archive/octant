package logging

import (
	"fmt"
)

type (
	// Logger has debug and info leveled logging
	Logger interface {
		infoLogger
		debugLogger
	}

	infoLogger  interface{ Infof(string, ...interface{}) }
	debugLogger interface{ Debugf(string, ...interface{}) }
	plainLogger interface{ Printf(string, ...interface{}) }

	plainAdapter struct{ plainLogger }
	infoAdapter  struct{ infoLogger }
	debugAdapter struct{ debugLogger }

	nilLogger struct{}
)

func (a *nilLogger) Debugf(string, ...interface{}) {}

func (a *nilLogger) Infof(string, ...interface{}) {}

func (a *plainAdapter) Debugf(s string, args ...interface{}) {
	a.plainLogger.Printf(s, args...)
}

func (a *plainAdapter) Infof(s string, args ...interface{}) {
	a.plainLogger.Printf(s, args...)
}

func (a *infoAdapter) Debugf(s string, args ...interface{}) {
	a.infoLogger.Infof(s, args...)
}

func (a *debugAdapter) Infof(s string, args ...interface{}) {
	a.debugLogger.Debugf(s, args...)
}

// Adapt takes a logger and adapts it to this package
func Adapt(logger interface{}) (Logger, error) {
	switch logger.(type) {
	case Logger:
		return logger.(Logger), nil
	case infoLogger:
		return &infoAdapter{infoLogger: logger.(infoLogger)}, nil
	case debugLogger:
		return &debugAdapter{debugLogger: logger.(debugLogger)}, nil
	case plainLogger:
		return &plainAdapter{plainLogger: logger.(plainLogger)}, nil
	case nil:
		return &nilLogger{}, nil
	default:
		return nil, fmt.Errorf("logger %+v doesn't provide Printf, Infof, or Debugf", logger)
	}
}

package astilog

import (
	"github.com/sirupsen/logrus"
)

// Logrus represents a logrus logger
type Logrus struct {
	*logrus.Logger
}

// NewLogrus creates a new logrus logger
func NewLogrus() *Logrus {
	return &Logrus{Logger: logrus.New()}
}

// WithField implements the LoggerWithFields interface
func (l *Logrus) WithField(k, v string) {
	l.AddHook(newWithFieldHook(k, v))
}

// WithFields implements the LoggerWithFields interface
func (l *Logrus) WithFields(fs Fields) {
	for k, v := range fs {
		l.WithField(k, v)
	}
}

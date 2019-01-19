package astilog

import (
	"io"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

// Logger represents a logger
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

// LoggerSetter represents a logger setter
type LoggerSetter interface {
	SetLogger(l Logger)
}

// Fields represents logger fields
type Fields map[string]string

// LoggerWithField represents a logger that can have fields
type LoggerWithFields interface {
	WithField(k, v string)
	WithFields(fs Fields)
}

// Outs
const (
	OutFile   = "file"
	OutStdOut = "stdout"
	OutSyslog = "syslog"
)

// New creates a new Logger
func New(c Configuration) Logger {
	// Init
	var l = NewLogrus()

	// Hooks
	l.AddHook(newWithFieldHook("app_name", c.AppName))

	// Out
	var out string
	l.Out, out = Out(c)

	// Formatter
	l.Formatter = Formatter(c, out)

	// Level
	l.Level = Level(c)
	return l
}

// Out returns the out based on the configuration
func Out(c Configuration) (w io.Writer, out string) {
	switch c.Out {
	case OutStdOut:
		return stdOut(), c.Out
	case OutSyslog:
		return syslogOut(c), c.Out
	default:
		if isTerminal(os.Stdout) {
			w = stdOut()
			out = OutStdOut
		} else {
			w = syslogOut(c)
			out = OutSyslog
		}
		if len(c.Filename) > 0 {
			f, err := os.OpenFile(c.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Println(errors.Wrapf(err, "astilog: creating %s failed", c.Filename))
			} else {
				w = f
				out = OutFile
			}
		}
		return
	}
}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

// Formats
const (
	FormatJSON = "json"
	FormatText = "text"
)

// Formatter returns the formatter based on the configuration
func Formatter(c Configuration, out string) logrus.Formatter {
	switch c.Format {
	case FormatJSON:
		return jsonFormatter(c)
	case FormatText:
		return textFormatter(c, out)
	default:
		switch out {
		case OutFile, OutStdOut:
			return textFormatter(c, out)
		default:
			return jsonFormatter(c)
		}
	}
}

func jsonFormatter(c Configuration) logrus.Formatter {
	f := &logrus.JSONFormatter{
		FieldMap:        make(logrus.FieldMap),
		TimestampFormat: c.TimestampFormat,
	}
	if len(c.MessageKey) > 0 {
		f.FieldMap[logrus.FieldKeyMsg] = c.MessageKey
	}
	return f
}

func textFormatter(c Configuration, out string) logrus.Formatter {
	return &logrus.TextFormatter{
		DisableColors:   c.DisableColors || out == OutFile,
		ForceColors:     !c.DisableColors && out != OutFile,
		FullTimestamp:   c.FullTimestamp,
		TimestampFormat: c.TimestampFormat,
	}
}

func Level(c Configuration) logrus.Level {
	if c.Verbose {
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}

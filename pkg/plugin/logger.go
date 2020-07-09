/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	golog "log"

	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/vmware-tanzu/octant/pkg/log"
)

type zapAdapter struct {
	dashLogger log.Logger
}

var _ hclog.Logger = (*zapAdapter)(nil)

// Args are alternating key, val pairs
// keys must be strings
// vals can be any type, but display is implementation specific
// Emit a message and key/value pairs at a provided log level
func (za *zapAdapter) Log(level hclog.Level, msg string, args ...interface{}) {
	// no-op: zap doesn't handle log
}

// Emit a message and key/value pairs at the TRACE level
func (za *zapAdapter) Trace(msg string, args ...interface{}) {
	// no-op: zap doesn't handle trace
}

// Emit a message and key/value pairs at the DEBUG level
func (za *zapAdapter) Debug(msg string, args ...interface{}) {
	za.dashLogger.With(args...).Debugf(msg)
}

// Emit a message and key/value pairs at the INFO level
func (za *zapAdapter) Info(msg string, args ...interface{}) {
	za.dashLogger.With(args...).Infof(msg)
}

// Emit a message and key/value pairs at the WARN level
func (za *zapAdapter) Warn(msg string, args ...interface{}) {
	za.dashLogger.With(args...).Warnf(msg)
}

// Emit a message and key/value pairs at the ERROR level
func (za *zapAdapter) Error(msg string, args ...interface{}) {
	za.dashLogger.With(args).Errorf(msg)
}

// Indicate if TRACE logs would be emitted. This and the other Is* guards
// are used to elide expensive logging code based on the current level.
func (za *zapAdapter) IsTrace() bool {
	return false
}

// Indicate if DEBUG logs would be emitted. This and the other Is* guards
func (za *zapAdapter) IsDebug() bool {
	return true
}

// Indicate if INFO logs would be emitted. This and the other Is* guards
func (za *zapAdapter) IsInfo() bool {
	return true
}

// Indicate if WARN logs would be emitted. This and the other Is* guards
func (za *zapAdapter) IsWarn() bool {
	return true
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards
func (za *zapAdapter) IsError() bool {
	return true
}

// ImpliedArgs returns With key/value pairs
func (za *zapAdapter) ImpliedArgs() []interface{} {
	// no-op
	return nil
}

// Creates a sublogger that will always have the given key/value pairs
func (za *zapAdapter) With(args ...interface{}) hclog.Logger {
	return &zapAdapter{
		dashLogger: za.dashLogger.With(args...),
	}
}

// Returns the Name of the logger
func (za *zapAdapter) Name() string {
	// no-op
	return ""
}

// Create a logger that will prepend the name string on the front of all messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (za *zapAdapter) Named(name string) hclog.Logger {
	return &zapAdapter{
		dashLogger: za.dashLogger.Named(name),
	}
}

// Create a logger that will prepend the name string on the front of all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (za *zapAdapter) ResetNamed(name string) hclog.Logger {
	return &zapAdapter{
		dashLogger: za.dashLogger.Named(name),
	}
}

// Updates the level. This should affect all sub-loggers as well. If an
// implementation cannot update the level on the fly, it should no-op.
func (za *zapAdapter) SetLevel(level hclog.Level) {
	// no-op
}

// Return a value that conforms to the stdlib log.Logger interface
func (za *zapAdapter) StandardLogger(opts *hclog.StandardLoggerOptions) *golog.Logger {
	if opts == nil {
		opts = &hclog.StandardLoggerOptions{}
	}

	return golog.New(za.StandardWriter(opts), "", 0)
}

// Return a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (za *zapAdapter) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return nil
}

/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package log

import "github.com/vmware-tanzu/octant/pkg/event"

// Logger is an interface for logging
type Logger interface {
	// Debugf uses fmt.Sprintf to log a templated message.
	Debugf(template string, args ...interface{})

	// Infof uses fmt.Sprintf to log a templated message.
	Infof(template string, args ...interface{})

	// Warnf uses fmt.Sprintf to log a templated message.
	Warnf(template string, args ...interface{})

	// Errorf uses fmt.Sprintf to log a templated message.
	Errorf(template string, args ...interface{})

	With(args ...interface{}) Logger

	WithErr(err error) Logger

	Named(name string) Logger

	Stream(readyCh <-chan struct{}) (<-chan event.Event, func())
}

// LoggerCloser is an interface that wraps a Logger and a close function.
type LoggerCloser interface {
	Logger

	// Close closes the logger.
	Close()
}

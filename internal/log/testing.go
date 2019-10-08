/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package log

import "testing"

type testingLogger struct {
	t *testing.T
}

// TestLogger returns a logger for tests
func TestLogger(t *testing.T) Logger {
	return &testingLogger{t: t}
}

func (t *testingLogger) Debugf(format string, args ...interface{}) {
	t.t.Logf(format, args...)
}
func (t *testingLogger) Infof(format string, args ...interface{}) {
	t.t.Logf(format, args...)
}
func (t *testingLogger) Warnf(format string, args ...interface{}) {
	t.t.Logf(format, args...)
}
func (t *testingLogger) Errorf(format string, args ...interface{}) {
	t.t.Errorf(format, args...)
}
func (t *testingLogger) With(args ...interface{}) Logger {
	return t
}
func(t *testingLogger) WithErr(err error) Logger {
	return t.With("err", err)
}
func (t *testingLogger) Named(string) Logger {
	return t
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package log

import (
	"testing"

	"github.com/vmware-tanzu/octant/pkg/log"
)

type testingLogger struct {
	t *testing.T
}

// TestLogger returns a logger for tests
func TestLogger(t *testing.T) log.Logger {
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
func (t *testingLogger) With(args ...interface{}) log.Logger {
	return t
}
func (t *testingLogger) WithErr(err error) log.Logger {
	return t.With("err", err)
}
func (t *testingLogger) Name() string {
	return t.Name()
}
func (t *testingLogger) Named(string) log.Logger {
	return t
}

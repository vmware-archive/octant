package astilog

import "os"

// NopLogger returns a nop logger
func NopLogger() Logger {
	return &nop{}
}

// nop is a nop logger
type nop struct{}

func (n nop) Debug(v ...interface{})                 {}
func (n nop) Debugf(format string, v ...interface{}) {}
func (n nop) Info(v ...interface{})                  {}
func (n nop) Infof(format string, v ...interface{})  {}
func (n nop) Warn(v ...interface{})                  {}
func (n nop) Warnf(format string, v ...interface{})  {}
func (n nop) Error(v ...interface{})                 {}
func (n nop) Errorf(format string, v ...interface{}) {}
func (n nop) Fatal(v ...interface{})                 { os.Exit(1) }
func (n nop) Fatalf(format string, v ...interface{}) { os.Exit(1) }

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vmware-tanzu/octant/pkg/log"
)

type key string

var contextKey = key("com.heptio.logger")

// sugaredLogWrapper adapts a zap.SugaredLogger to the Logger interface
type sugaredLogWrapper struct {
	*zap.SugaredLogger
}

func (s *sugaredLogWrapper) WithErr(err error) log.Logger {
	return &sugaredLogWrapper{s.SugaredLogger.With("err", err.Error())}
}

func (s *sugaredLogWrapper) With(args ...interface{}) log.Logger {
	return &sugaredLogWrapper{s.SugaredLogger.With(args...)}
}

func (s *sugaredLogWrapper) Named(name string) log.Logger {
	return &sugaredLogWrapper{s.SugaredLogger.Named(name)}
}

var _ log.Logger = (*sugaredLogWrapper)(nil)

// Wrap zap.SugaredLogger as Logger interface
func Wrap(z *zap.SugaredLogger) log.Logger {
	return &sugaredLogWrapper{z}
}

// NopLogger constructs a nop logger
func NopLogger() log.Logger {
	return Wrap(zap.NewNop().Sugar())
}

// WithLoggerContext returns a new context with a set logger
func WithLoggerContext(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

// From extracts a logger from the supplied context, or returns a NopLogger if none is found.
func From(ctx context.Context) log.Logger {
	if ctx == nil {
		return NopLogger()
	}
	v := ctx.Value(contextKey)
	l, ok := v.(log.Logger)
	if !ok || l == nil {
		return NopLogger()
	}
	return l
}

// InitOption is a functional option for configuring a logger.
type InitOption func(config zap.Config) zap.Config

// Init initializes a logger with options.
func Init(logLevel int, options ...InitOption) (*zap.Logger, error) {
	z, err := newZapLogger(logLevel, options...)
	if err != nil {
		return nil, fmt.Errorf("create zap logger: %w", err)
	}

	return z, nil

}

// Returns a new zap logger, setting level according to the provided
// verbosity level as an offset of the base level, Info.
// i.e. verboseLevel==0, level==Info
//      verboseLevel==1, level==Debug
func newZapLogger(verboseLevel int, options ...InitOption) (*zap.Logger, error) {
	level := zapcore.InfoLevel - zapcore.Level(verboseLevel)
	if level < zapcore.DebugLevel || level > zapcore.FatalLevel {
		level = zapcore.DebugLevel
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	for _, option := range options {
		cfg = option(cfg)
	}

	return cfg.Build()
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package log

import (
	"context"
	"fmt"
	"net/url"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vmware-tanzu/octant/pkg/event"
	"github.com/vmware-tanzu/octant/pkg/log"
)

type key string

var contextKey = key("com.heptio.logger")

type sugaredLogWrapperOption func(wrapper *sugaredLogWrapper)

func sugaredLogStreamer(s EventStreamer) sugaredLogWrapperOption {
	return func(wrapper *sugaredLogWrapper) {
		wrapper.streamer = s
	}
}

// sugaredLogWrapper adapts a zap.SugaredLogger to the Logger interface
type sugaredLogWrapper struct {
	*zap.SugaredLogger

	streamer EventStreamer
}

var _ log.LoggerCloser = &sugaredLogWrapper{}

func newSugaredLogWrapper(sl *zap.SugaredLogger, options ...sugaredLogWrapperOption) *sugaredLogWrapper {
	slw := &sugaredLogWrapper{
		SugaredLogger: sl,
	}

	for _, option := range options {
		option(slw)
	}

	return slw
}

func (s *sugaredLogWrapper) fork(sl *zap.SugaredLogger, options ...sugaredLogWrapperOption) *sugaredLogWrapper {
	return newSugaredLogWrapper(sl, append([]sugaredLogWrapperOption{sugaredLogStreamer(s.streamer)}, options...)...)
}

func (s *sugaredLogWrapper) WithErr(err error) log.Logger {
	return s.fork(s.SugaredLogger.With("err", err.Error()))
}

func (s *sugaredLogWrapper) With(args ...interface{}) log.Logger {
	return s.fork(s.SugaredLogger.With(args...))
}

func (s *sugaredLogWrapper) Named(name string) log.Logger {
	return s.fork(s.SugaredLogger.Named(name))
}

func (s *sugaredLogWrapper) Close() {
	// this fails, but it should be safe to ignore according
	// to https://github.com/uber-go/zap/issues/328
	_ = s.SugaredLogger.Sync()
}

func (s *sugaredLogWrapper) Stream(readyCh <-chan struct{}) (<-chan event.Event, func()) {
	return s.streamer.Stream(readyCh)
}

// Wrap zap.SugaredLogger as Logger interface
func Wrap(z *zap.SugaredLogger) log.LoggerCloser {
	return newSugaredLogWrapper(z)
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
func Init(logLevel int, options ...InitOption) (log.LoggerCloser, error) {
	octantSink := NewOctantSink()

	err := zap.RegisterSink("octant-stream", func(url *url.URL) (zap.Sink, error) {
		return octantSink, nil
	})
	if err != nil {
		return nil, fmt.Errorf("register octant log sink: %w", err)
	}

	z, err := newZapLogger(logLevel, options...)
	if err != nil {
		return nil, fmt.Errorf("create zap logger: %w", err)
	}

	logger := newSugaredLogWrapper(z.Sugar(), func(wrapper *sugaredLogWrapper) {
		wrapper.streamer = NewStreamer(octantSink)
	})

	return logger, nil

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
		OutputPaths:      []string{"stderr", "octant-stream://output"},
		ErrorOutputPaths: []string{"stderr", "octant-stream://output"},
	}

	for _, option := range options {
		cfg = option(cfg)
	}

	return cfg.Build()
}

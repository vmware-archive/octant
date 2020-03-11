/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package log

import (
	"context"

	"go.uber.org/zap"

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

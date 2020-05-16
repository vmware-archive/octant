/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package electron

import (
	"context"

	"github.com/asticode/go-astikit"

	log2 "github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/log"
)

type loggerAdapter struct {
	logger log.Logger
}

var _ astikit.StdLogger = &loggerAdapter{}
var _ astikit.SeverityLogger = &loggerAdapter{}

func loggerFromContext(ctx context.Context) *loggerAdapter {
	logger := log2.From(ctx)

	l := loggerAdapter{
		logger: logger,
	}

	return &l
}

func (l *loggerAdapter) Print(v ...interface{}) {
	l.logger.Infof("", v...)
}
func (l *loggerAdapter) Printf(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

func (l *loggerAdapter) Debug(v ...interface{}) {
	l.logger.Debugf("", v...)
}

func (l *loggerAdapter) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

func (l *loggerAdapter) Error(v ...interface{}) {
	l.logger.Errorf("", v...)
}

func (l *loggerAdapter) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

func (l *loggerAdapter) Info(v ...interface{}) {
	l.logger.Infof("", v...)
}

func (l *loggerAdapter) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

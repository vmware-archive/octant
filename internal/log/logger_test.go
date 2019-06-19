/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerFromContext(t *testing.T) {
	expected := NopLogger()
	ctx := WithLoggerContext(context.Background(), expected)
	actual := From(ctx)
	assert.True(t, actual == expected, "unexpected logger instance from context")
}

func TestMissingLoggerFromContext(t *testing.T) {
	notExpected := NopLogger()
	ctx := context.Background()
	actual := From(ctx)
	assert.True(t, actual != notExpected, "unexpected logger instance from context")
	assert.NotNil(t, actual, "expected non-nil logger")
}

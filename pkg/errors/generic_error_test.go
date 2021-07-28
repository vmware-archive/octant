/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGenericError(t *testing.T) {
	err := fmt.Errorf("access denied")

	intErr := NewGenericError(err)
	assert.Equal(t, "access denied", intErr.Error())
	assert.EqualError(t, intErr.err, "access denied")
	assert.NotEmpty(t, intErr.Timestamp())
	assert.NotZero(t, intErr.ID())
}

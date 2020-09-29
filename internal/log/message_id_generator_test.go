/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package log

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUUIDMessageIDGenerator_Generate(t *testing.T) {
	gen := &UUIDMessageIDGenerator{}
	got := gen.Generate()
	require.True(t, isValidUUID(got))
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

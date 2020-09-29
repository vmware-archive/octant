/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package log

import "github.com/google/uuid"

// MessageIDGenerator is an interface for generating message IDs.
type MessageIDGenerator interface {
	Generate() string
}

// UUIDMessageIDGenerator generates messages IDs as UUIDs.
type UUIDMessageIDGenerator struct{}

var _ MessageIDGenerator = &UUIDMessageIDGenerator{}

// Generate generates a UUID as a string.
func (d UUIDMessageIDGenerator) Generate() string {
	return uuid.New().String()
}

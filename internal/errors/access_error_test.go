/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/octant/pkg/store"
)

func TestNewAccessError(t *testing.T) {
	key := store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Pod",
	}
	verb := "watch"
	err := fmt.Errorf("access denied")

	intErr := NewAccessError(key, verb, err)
	assert.Equal(t, key, intErr.Key())
	assert.Equal(t, verb, intErr.Verb())
	assert.Equal(t, fmt.Sprintf("%s: %s: %s", verb, key, err), intErr.Error())
	assert.EqualError(t, err, "access denied")
	assert.NotEmpty(t, intErr.Timestamp())
	assert.NotZero(t, intErr.ID())
}

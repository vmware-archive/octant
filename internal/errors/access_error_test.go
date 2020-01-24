/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/octant/pkg/store"
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
	assert.Equal(t, fmt.Sprintf("%s: %s (error: %s)", verb, key, err.Error()), intErr.Error())
	assert.EqualError(t, err, "access denied")
	assert.NotEmpty(t, intErr.Timestamp())
	assert.NotZero(t, intErr.ID())
}

func TestFormattedAccessError(t *testing.T) {
	key := store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Pod",
	}
	verb := "watch"
	err := fmt.Errorf("access denied")

	intErr := NewAccessError(key, verb, err)
	newErr := fmt.Errorf("%w", intErr)

	var e *AccessError
	assert.True(t, goerrors.As(newErr, &e))
}

func TestNilErrAccessError(t *testing.T) {
	key := store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Pod",
	}
	verb := "watch"

	intErr := NewAccessError(key, verb, nil)
	assert.Equal(t, fmt.Sprintf("%s: %s", verb, key), intErr.Error())
}

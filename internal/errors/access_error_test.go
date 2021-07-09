/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/pkg/errors"

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

func TestDefaultKeyAccessError(t *testing.T) {
	key := store.Key{}
	verb := "watch"

	intErr := NewAccessError(key, verb, errors.New("an error message"))
	assert.Equal(t, fmt.Sprintf("%s: %s (error: %s)", verb, key, errors.New("an error message")), intErr.Error())
}

func TestIdNilErrAccessError(t *testing.T) {
	key := store.Key{}
	verb := "watch"

	intErr := NewAccessError(key, verb, nil)
	assert.Equal(t, intErr.ID(), "4163636573734572726f723a2077617463683a2043616368654b65795b41504956657273696f6e3d27272c204b696e643d27275dcbf29ce484222325")
}

func TestIdDefaultKeyAccessError(t *testing.T) {
	key := store.Key{}
	verb := "watch"

	intErr := NewAccessError(key, verb, errors.New("an error message"))
	assert.Equal(t, intErr.ID(), "4163636573734572726f723a20616e206572726f72206d657373616765cbf29ce484222325")
}

func TestIdAccessError(t *testing.T) {
	key := store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Pod",
	}
	verb := "watch"

	intErr := NewAccessError(key, verb, nil)
	assert.Equal(t, intErr.ID(), "4163636573734572726f723a2077617463683a2043616368654b65795b4e616d6573706163653d2764656661756c74272c2041504956657273696f6e3d277631272c204b696e643d27506f64275dcbf29ce484222325")
}

/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/pkg/action"
	oerrors "github.com/vmware-tanzu/octant/pkg/errors"
)

func TestNewErrorStore(t *testing.T) {
	_, err := NewErrorStore()
	require.NoError(t, err)
}

func TestNewError(t *testing.T) {
	errStore, err := NewErrorStore()
	require.NoError(t, err)

	requestType := "setNavigation"
	payload := action.Payload{}
	err = fmt.Errorf("setNavigation error")

	intErr := NewActionError(requestType, payload, err)
	errStore.Add(intErr)

	_, found := errStore.Get(intErr.ID())
	assert.True(t, found)
}

func TestErrorStore_Accessors(t *testing.T) {
	errStore, err := NewErrorStore()
	require.NoError(t, err)

	requestType := "setNamespace"
	payload := action.Payload{}
	err = fmt.Errorf("setNamespace error")

	i := NewActionError(requestType, payload, err)
	_, found := errStore.Get(i.ID())
	assert.False(t, found)

	errStore.Add(i)

	e, found := errStore.Get(i.ID())
	assert.True(t, found)
	assert.Equal(t, i.ID(), e.ID())

	l := errStore.List()
	assert.Len(t, l, 1)
	assert.Equal(t, i.ID(), l[0].ID())
}

func TestErrorStore_ListOrder(t *testing.T) {
	errStore, err := NewErrorStore()
	require.NoError(t, err)

	requestType := "setContext"
	payload := action.Payload{}
	err = fmt.Errorf("setContext error")

	older := NewActionError(requestType, payload, err)
	newer := oerrors.NewGenericError(err)

	errStore.Add(newer)
	errStore.Add(older)

	l := errStore.List()
	assert.Len(t, l, 2)
	assert.Equal(t, older.ID(), l[1].ID())
	assert.Equal(t, newer.ID(), l[0].ID())
}

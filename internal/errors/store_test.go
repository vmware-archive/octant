package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/octant/pkg/action"
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

	intErr := errStore.NewError(requestType, payload, err)
	_, found := errStore.Get(intErr.id)
	assert.True(t, found)
}

func TestErrorStore_Accessors(t *testing.T) {
	errStore, err := NewErrorStore()
	require.NoError(t, err)

	requestType := "setNamespace"
	payload := action.Payload{}
	err = fmt.Errorf("setNamespace error")

	i := NewInternalError(requestType, payload, err)
	_, found := errStore.Get(i.id)
	assert.False(t, found)

	errStore.Add(i)

	e, found := errStore.Get(i.id)
	assert.True(t, found)
	assert.Equal(t, i.id, e.id)

	l := errStore.List()
	assert.Len(t, l, 1)
	assert.Equal(t, i.id, l[0].id)
}

func TestErrorStore_ListOrder(t *testing.T) {
	errStore, err := NewErrorStore()
	require.NoError(t, err)

	requestType := "setContext"
	payload := action.Payload{}
	err = fmt.Errorf("setContext error")

	older := errStore.NewError(requestType, payload, err)
	newer := errStore.NewError(requestType, payload, err)

	l := errStore.List()
	assert.Len(t, l, 2)
	assert.Equal(t, older.id, l[1].id)
	assert.Equal(t, newer.id, l[0].id)
}

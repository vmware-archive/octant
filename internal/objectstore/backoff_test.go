package objectstore

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/vmware-tanzu/octant/pkg/store"
)

func Test_waiting(t *testing.T) {
	key := store.Key{APIVersion: "apiVersion"}
	entry := newBackoffEntry(key, defaultBackoff)

	assert.False(t, entry.isWaiting())
	entry.setWaiting(true)
	assert.True(t, entry.isWaiting())
}

func Test_wait(t *testing.T) {
	key := store.Key{APIVersion: "apiVersion"}
	entry := newBackoffEntry(key, defaultBackoff)
	d := entry.wait()
	assert.True(t, d > time.Second)
	assert.True(t, d < (d + (time.Millisecond * 500)))
}

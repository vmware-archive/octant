package objectstore

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/pkg/store"
)

func TestDynamicCache_backoff(t *testing.T) {
	d := &DynamicCache{
		factories: initFactoriesCache(),
	}

	ctx := context.TODO()
	key := store.Key{APIVersion: gvk.Pod.Version, Kind: gvk.Pod.Kind}

	tD := d.backoff(ctx, key)
	require.True(t, d.isBackingOff(ctx, key))
	<-time.After(tD + (time.Millisecond * 250))
	assert.False(t, d.isBackingOff(ctx, key))
}

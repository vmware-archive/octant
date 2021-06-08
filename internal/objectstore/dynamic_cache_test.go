package objectstore

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
)

func TestNewDynamicCache(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	defer controller.Finish()

	clusterClient := clusterFake.NewMockClientInterface(controller)
	clusterClient.EXPECT().DynamicClient()

	dc, err := NewDynamicCache(ctx, clusterClient)
	require.NoError(t, err)

	dc.knownInformers.Store(clusterFake.NewMockSharedIndexInformer(controller), make(chan struct{}))
	require.Equal(t, 1, mapLength(&dc.knownInformers))

	dc.stopAllInformers()

	require.Equal(t, 0, mapLength(&dc.knownInformers))
}

func TestDynamicCache_UpdateClusterClient(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	defer controller.Finish()

	clusterClient := clusterFake.NewMockClientInterface(controller)
	clusterClient.EXPECT().DynamicClient()

	dc, err := NewDynamicCache(ctx, clusterClient)
	require.NoError(t, err)

	dc.informerFactories.Store("default", clusterFake.NewMockDynamicSharedInformerFactory(controller))

	err = dc.UpdateClusterClient(ctx, clusterClient)
	require.NoError(t, err)

	require.Equal(t, 0, mapLength(&dc.informerFactories))
}

func mapLength(sMap *sync.Map) int {
	var l int
	sMap.Range(func(k, v interface{}) bool {
		l += 1
		return true
	})
	return l
}

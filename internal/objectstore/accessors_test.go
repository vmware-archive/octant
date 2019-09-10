package objectstore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/cluster/fake"
	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/pkg/store"
)

func Test_informerSynced(t *testing.T) {
	c := initInformerSynced()
	key := store.Key{APIVersion: "apiVersion"}
	otherKey := store.Key{APIVersion: "apiVersion2"}
	require.True(t, c.hasSynced(key))

	c.setSynced(key, true)
	require.True(t, c.hasSynced(key))
	require.True(t, c.hasSynced(otherKey))

	c.setSynced(key, false)
	require.False(t, c.hasSynced(key))
}

func Test_factoriesCache(t *testing.T) {
	const namespaceName = "test"

	controller := gomock.NewController(t)
	defer controller.Finish()

	dynamicClient := fake.NewMockDynamicInterface(controller)

	client := fake.NewMockClientInterface(controller)
	client.EXPECT().
		DynamicClient().
		Return(dynamicClient, nil)

	c := initFactoriesCache()

	ctx := context.Background()
	factory, err := initInformerFactory(ctx, client, namespaceName)
	require.NoError(t, err)

	c.set(namespaceName, factory)

	got, isFound := c.get(namespaceName)
	require.True(t, isFound)
	require.Equal(t, factory, got)

	c.delete(namespaceName)
	_, isFound = c.get(namespaceName)
	require.False(t, isFound)
}

func Test_seenGVKsCache(t *testing.T) {
	c := initSeenGVKsCache()
	c.setSeen("test", gvk.Pod, true)

	tests := []struct {
		name      string
		namespace string
		gvk       schema.GroupVersionKind
		expected  bool
	}{
		{
			name:      "gvk that has been seen",
			namespace: "test",
			gvk:       gvk.Pod,
			expected:  true,
		},
		{
			name:      "namespace that has not been seen",
			namespace: "other",
			gvk:       gvk.Pod,
			expected:  false,
		},
		{
			name:      "gvk that has not been seen",
			namespace: "test",
			gvk:       gvk.Deployment,
			expected:  false,
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			got := c.hasSeen(test.namespace, test.gvk)
			require.Equal(t, test.expected, got)
		})
	}
}

func Test_informerContextCache(t *testing.T) {
	c := initInformerContextCache()

	gvr1 := schema.GroupVersionResource{
		Group:    "group",
		Version:  "version",
		Resource: "resource1",
	}
	gvr2 := schema.GroupVersionResource{
		Group:    "group",
		Version:  "version",
		Resource: "resource2",
	}
	_ = c.addChild(gvr1)
	assert.Len(t, c.cache, 1)
	_ = c.addChild(gvr1)
	assert.Len(t, c.cache, 1)
	_ = c.addChild(gvr2)
	assert.Len(t, c.cache, 2)
	c.delete(gvr1)
	assert.Len(t, c.cache, 1)
	c.reset()
	assert.Len(t, c.cache, 0)
}

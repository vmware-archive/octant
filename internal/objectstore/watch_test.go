package objectstore

import (
	"context"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/cluster"
	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

type watchMocks struct {
	controller *gomock.Controller

	informerFactory     *clusterfake.MockDynamicSharedInformerFactory
	informer            *clusterfake.MockGenericInformer
	client              *clusterfake.MockClientInterface
	sharedIndexInformer *clusterfake.MockSharedIndexInformer
	backendObjectStore  *storefake.MockObjectStore
}

func newWatchMocks(t *testing.T) *watchMocks {
	controller := gomock.NewController(t)
	m := &watchMocks{
		controller:          controller,
		informerFactory:     clusterfake.NewMockDynamicSharedInformerFactory(controller),
		informer:            clusterfake.NewMockGenericInformer(controller),
		client:              clusterfake.NewMockClientInterface(controller),
		backendObjectStore:  storefake.NewMockObjectStore(controller),
		sharedIndexInformer: clusterfake.NewMockSharedIndexInformer(controller),
	}

	return m
}

func Test_WatchList_not_cached(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mocks := newWatchMocks(t)
	defer mocks.controller.Finish()

	mocks.informer.EXPECT().Informer().Return(mocks.sharedIndexInformer)

	mocks.sharedIndexInformer.EXPECT().AddEventHandler(gomock.Any())

	podGVR := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}

	mocks.informerFactory.EXPECT().ForResource(gomock.Eq(podGVR)).Return(mocks.informer)

	podGK := schema.GroupKind{
		Kind: "Pod",
	}
	mocks.client.EXPECT().Resource(gomock.Eq(podGK)).Return(podGVR, nil)

	mocks.informerFactory.EXPECT().Start(gomock.Eq(ctx.Done()))

	pod1 := testutil.CreatePod("pod1")
	pod1.Namespace = "test"
	pod2 := testutil.CreatePod("pod2")
	pod2.Namespace = "test"

	listKey := objectstoreutil.Key{Namespace: "test", APIVersion: "v1", Kind: "Pod"}
	objects := []*unstructured.Unstructured{
		testutil.ToUnstructured(t, pod1),
		testutil.ToUnstructured(t, pod2),
	}

	mocks.backendObjectStore.EXPECT().List(gomock.Any(), gomock.Eq(listKey)).Return(objects, nil)

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}
	}

	setBackendFunc := func(w *Watch) {
		w.backendObjectStore = mocks.backendObjectStore
	}

	watch, err := NewWatch(mocks.client, ctx.Done(), factoryFunc, setBackendFunc)
	require.NoError(t, err)

	got, err := watch.List(ctx, listKey)
	require.NoError(t, err)

	expected := []*unstructured.Unstructured{
		testutil.ToUnstructured(t, pod1),
		testutil.ToUnstructured(t, pod2),
	}
	assert.Equal(t, expected, got)
}

func Test_WatchList_stored(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mocks := newWatchMocks(t)
	defer mocks.controller.Finish()

	pod1 := testutil.CreatePod("pod1")
	pod1.Namespace = "test"
	pod2 := testutil.CreatePod("pod2")
	pod2.Namespace = "test"

	listKey := objectstoreutil.Key{Namespace: "test", APIVersion: "v1", Kind: "Pod"}

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}
	}

	setBackendFunc := func(w *Watch) {
		w.backendObjectStore = mocks.backendObjectStore
	}

	cacheKeyFunc := func(w *Watch) {
		gvk := listKey.GroupVersionKind()
		w.watchedGVKs[gvk] = true
		w.cachedObjects[gvk] = map[types.UID]*unstructured.Unstructured{
			pod1.UID: testutil.ToUnstructured(t, pod1),
			pod2.UID: testutil.ToUnstructured(t, pod2),
		}
	}

	watch, err := NewWatch(mocks.client, ctx.Done(), factoryFunc, setBackendFunc, cacheKeyFunc)
	require.NoError(t, err)

	got, err := watch.List(ctx, listKey)
	require.NoError(t, err)

	expected := []*unstructured.Unstructured{
		testutil.ToUnstructured(t, pod1),
		testutil.ToUnstructured(t, pod2),
	}

	sort.Slice(got, func(i, j int) bool {
		return got[i].GetName() < got[j].GetName()
	})

	assert.Equal(t, expected, got)
}

func Test_WatchList_stored_with_selector(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mocks := newWatchMocks(t)
	defer mocks.controller.Finish()

	pod1 := testutil.CreatePod("pod1")
	pod1.Namespace = "test"
	pod1.Labels = map[string]string{
		"app": "app1",
	}
	pod2 := testutil.CreatePod("pod2")
	pod2.Namespace = "test"

	ls := &labels.Set{
		"app": "app1",
	}

	listKey := objectstoreutil.Key{
		Namespace:  "test",
		APIVersion: "v1",
		Kind:       "Pod",
		Selector:   ls,
	}

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}
	}

	setBackendFunc := func(w *Watch) {
		w.backendObjectStore = mocks.backendObjectStore
	}

	cacheKeyFunc := func(w *Watch) {
		gvk := listKey.GroupVersionKind()
		w.watchedGVKs[gvk] = true
		w.cachedObjects[gvk] = map[types.UID]*unstructured.Unstructured{
			pod1.UID: testutil.ToUnstructured(t, pod1),
			pod2.UID: testutil.ToUnstructured(t, pod2),
		}
	}

	watch, err := NewWatch(mocks.client, ctx.Done(), factoryFunc, setBackendFunc, cacheKeyFunc)
	require.NoError(t, err)

	got, err := watch.List(ctx, listKey)
	require.NoError(t, err)

	expected := []*unstructured.Unstructured{
		testutil.ToUnstructured(t, pod1),
	}

	sort.Slice(got, func(i, j int) bool {
		return got[i].GetName() < got[j].GetName()
	})

	assert.Equal(t, expected, got)
}

func Test_WatchGet_not_stored(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mocks := newWatchMocks(t)
	defer mocks.controller.Finish()

	mocks.informer.EXPECT().Informer().Return(mocks.sharedIndexInformer)

	mocks.sharedIndexInformer.EXPECT().AddEventHandler(gomock.Any())

	podGVR := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}

	mocks.informerFactory.EXPECT().ForResource(gomock.Eq(podGVR)).Return(mocks.informer)

	podGK := schema.GroupKind{
		Kind: "Pod",
	}
	mocks.client.EXPECT().Resource(gomock.Eq(podGK)).Return(podGVR, nil)

	mocks.informerFactory.EXPECT().Start(gomock.Eq(ctx.Done()))

	pod1 := testutil.CreatePod("pod1")
	pod1.Namespace = "test"

	getKey := objectstoreutil.Key{Namespace: "test", APIVersion: "v1", Kind: "Pod", Name: pod1.Name}
	mocks.backendObjectStore.EXPECT().Get(gomock.Any(), gomock.Eq(getKey)).Return(testutil.ToUnstructured(t, pod1), nil)

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}
	}

	setBackendFunc := func(w *Watch) {
		w.backendObjectStore = mocks.backendObjectStore
	}

	watch, err := NewWatch(mocks.client, ctx.Done(), factoryFunc, setBackendFunc)
	require.NoError(t, err)

	got, err := watch.Get(ctx, getKey)
	require.NoError(t, err)

	expected := testutil.ToUnstructured(t, pod1)
	assert.Equal(t, expected, got)
}

func Test_WatchGet_stored(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mocks := newWatchMocks(t)
	defer mocks.controller.Finish()

	pod1 := testutil.CreatePod("pod1")
	pod1.Namespace = "test"

	getKey := objectstoreutil.Key{Namespace: "test", APIVersion: "v1", Kind: "Pod", Name: pod1.Name}

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}
	}

	setBackendFunc := func(w *Watch) {
		w.backendObjectStore = mocks.backendObjectStore
	}

	cacheKeyFunc := func(w *Watch) {
		gvk := getKey.GroupVersionKind()
		w.watchedGVKs[gvk] = true
		w.cachedObjects[gvk] = map[types.UID]*unstructured.Unstructured{
			pod1.UID: testutil.ToUnstructured(t, pod1),
		}
	}

	watch, err := NewWatch(mocks.client, ctx.Done(), factoryFunc, setBackendFunc, cacheKeyFunc)
	require.NoError(t, err)

	got, err := watch.Get(ctx, getKey)
	require.NoError(t, err)

	expected := testutil.ToUnstructured(t, pod1)

	assert.Equal(t, expected, got)
}

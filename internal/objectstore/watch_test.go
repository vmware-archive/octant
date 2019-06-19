/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstore

import (
	"context"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"

	"github.com/vmware/octant/internal/cluster"
	clusterfake "github.com/vmware/octant/internal/cluster/fake"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/store"
	objectStoreFake "github.com/vmware/octant/pkg/store/fake"
	"github.com/vmware/octant/third_party/k8s.io/client-go/dynamic/dynamicinformer"
)

type watchMocks struct {
	controller *gomock.Controller

	informerFactory     *clusterfake.MockDynamicSharedInformerFactory
	informer            *clusterfake.MockGenericInformer
	client              *clusterfake.MockClientInterface
	sharedIndexInformer *clusterfake.MockSharedIndexInformer
	namespaceClient     *clusterfake.MockNamespaceInterface
	backendObjectStore  *objectStoreFake.MockStore
}

func newWatchMocks(t *testing.T) *watchMocks {
	controller := gomock.NewController(t)
	m := &watchMocks{
		controller:          controller,
		informerFactory:     clusterfake.NewMockDynamicSharedInformerFactory(controller),
		informer:            clusterfake.NewMockGenericInformer(controller),
		client:              clusterfake.NewMockClientInterface(controller),
		backendObjectStore:  objectStoreFake.NewMockStore(controller),
		sharedIndexInformer: clusterfake.NewMockSharedIndexInformer(controller),
		namespaceClient:     clusterfake.NewMockNamespaceInterface(controller),
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

	mocks.informerFactory.EXPECT().Start(gomock.Any())

	pod1 := testutil.CreatePod("pod1")
	pod1.Namespace = "test"
	pod2 := testutil.CreatePod("pod2")
	pod2.Namespace = "test"

	listKey := store.Key{Namespace: "test", APIVersion: "v1", Kind: "Pod"}
	objects := []*unstructured.Unstructured{
		testutil.ToUnstructured(t, pod1),
		testutil.ToUnstructured(t, pod2),
	}

	mocks.backendObjectStore.EXPECT().HasAccess(gomock.Any(), "list").Return(nil)
	mocks.backendObjectStore.EXPECT().List(gomock.Any(), gomock.Eq(listKey)).Return(objects, nil)

	nsKey := store.Key{APIVersion: "v1", Kind: "Namespace"}
	mocks.backendObjectStore.EXPECT().Watch(gomock.Any(), nsKey, gomock.Any()).Return(nil)

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}
	}

	setBackendFunc := func(w *Watch) {
		w.backendObjectStore = mocks.backendObjectStore
	}

	mocks.client.EXPECT().NamespaceClient().Return(mocks.namespaceClient, nil)
	namespaces := []string{"test"}
	mocks.namespaceClient.EXPECT().Names().Return(namespaces, nil)

	watch, err := NewWatch(ctx, mocks.client, factoryFunc, setBackendFunc)
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

	listKey := store.Key{Namespace: "test", APIVersion: "v1", Kind: "Pod"}

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}
	}

	setBackendFunc := func(w *Watch) {
		w.backendObjectStore = mocks.backendObjectStore
	}

	cacheKeyFunc := func(w *Watch) {
		w.watchedGVKs[listKey.Namespace] = make(map[schema.GroupVersionKind]bool)
		w.cachedObjects[listKey.Namespace] = make(map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured)

		gvk := listKey.GroupVersionKind()
		w.watchedGVKs[listKey.Namespace][gvk] = true
		w.cachedObjects[listKey.Namespace][gvk] = map[types.UID]*unstructured.Unstructured{
			pod1.UID: testutil.ToUnstructured(t, pod1),
			pod2.UID: testutil.ToUnstructured(t, pod2),
		}
	}

	nsKey := store.Key{APIVersion: "v1", Kind: "Namespace"}
	mocks.backendObjectStore.EXPECT().Watch(gomock.Any(), nsKey, gomock.Any()).Return(nil)

	mocks.client.EXPECT().NamespaceClient().Return(mocks.namespaceClient, nil)
	namespaces := []string{"test"}
	mocks.namespaceClient.EXPECT().Names().Return(namespaces, nil)

	watch, err := NewWatch(ctx, mocks.client, factoryFunc, setBackendFunc, cacheKeyFunc)
	require.NoError(t, err)

	mocks.backendObjectStore.EXPECT().HasAccess(gomock.Any(), "list").Return(nil)
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

	listKey := store.Key{
		Namespace:  "test",
		APIVersion: "v1",
		Kind:       "Pod",
		Selector:   ls,
	}

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}
	}

	setBackendFunc := func(w *Watch) {
		w.backendObjectStore = mocks.backendObjectStore
	}

	cacheKeyFunc := func(w *Watch) {
		w.watchedGVKs[listKey.Namespace] = make(map[schema.GroupVersionKind]bool)
		w.cachedObjects[listKey.Namespace] = make(map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured)

		gvk := listKey.GroupVersionKind()
		w.watchedGVKs[listKey.Namespace][gvk] = true
		w.cachedObjects[listKey.Namespace][gvk] = map[types.UID]*unstructured.Unstructured{
			pod1.UID: testutil.ToUnstructured(t, pod1),
			pod2.UID: testutil.ToUnstructured(t, pod2),
		}
	}

	mocks.client.EXPECT().NamespaceClient().Return(mocks.namespaceClient, nil)
	namespaces := []string{"test"}
	mocks.namespaceClient.EXPECT().Names().Return(namespaces, nil)

	nsKey := store.Key{APIVersion: "v1", Kind: "Namespace"}
	mocks.backendObjectStore.EXPECT().Watch(gomock.Any(), nsKey, gomock.Any()).Return(nil)

	watch, err := NewWatch(ctx, mocks.client, factoryFunc, setBackendFunc, cacheKeyFunc)
	require.NoError(t, err)

	mocks.backendObjectStore.EXPECT().HasAccess(gomock.Any(), "list").Return(nil)
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

	mocks.informerFactory.EXPECT().Start(gomock.Any())

	pod1 := testutil.CreatePod("pod1")
	pod1.Namespace = "test"

	getKey := store.Key{Namespace: "test", APIVersion: "v1", Kind: "Pod", Name: pod1.Name}
	mocks.backendObjectStore.EXPECT().Get(gomock.Any(), gomock.Eq(getKey)).Return(testutil.ToUnstructured(t, pod1), nil)

	nsKey := store.Key{APIVersion: "v1", Kind: "Namespace"}
	mocks.backendObjectStore.EXPECT().Watch(gomock.Any(), nsKey, gomock.Any()).Return(nil)

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}
	}

	setBackendFunc := func(w *Watch) {
		w.backendObjectStore = mocks.backendObjectStore
	}

	mocks.client.EXPECT().NamespaceClient().Return(mocks.namespaceClient, nil)
	namespaces := []string{"test"}
	mocks.namespaceClient.EXPECT().Names().Return(namespaces, nil)

	watch, err := NewWatch(ctx, mocks.client, factoryFunc, setBackendFunc)
	require.NoError(t, err)

	mocks.backendObjectStore.EXPECT().HasAccess(gomock.Any(), "get").Return(nil)
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

	getKey := store.Key{Namespace: "test", APIVersion: "v1", Kind: "Pod", Name: pod1.Name}

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}
	}

	setBackendFunc := func(w *Watch) {
		w.backendObjectStore = mocks.backendObjectStore
		w.watchedGVKs[getKey.Namespace] = make(map[schema.GroupVersionKind]bool)
		w.cachedObjects[getKey.Namespace] = make(map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured)
	}

	cacheKeyFunc := func(w *Watch) {
		w.watchedGVKs[getKey.Namespace] = make(map[schema.GroupVersionKind]bool)
		w.cachedObjects[getKey.Namespace] = make(map[schema.GroupVersionKind]map[types.UID]*unstructured.Unstructured)

		gvk := getKey.GroupVersionKind()
		w.watchedGVKs[getKey.Namespace][gvk] = true
		w.cachedObjects[getKey.Namespace][gvk] = map[types.UID]*unstructured.Unstructured{
			pod1.UID: testutil.ToUnstructured(t, pod1),
		}
	}

	mocks.client.EXPECT().NamespaceClient().Return(mocks.namespaceClient, nil)
	namespaces := []string{"test"}
	mocks.namespaceClient.EXPECT().Names().Return(namespaces, nil)

	nsKey := store.Key{APIVersion: "v1", Kind: "Namespace"}
	mocks.backendObjectStore.EXPECT().Watch(gomock.Any(), nsKey, gomock.Any()).Return(nil)

	watch, err := NewWatch(ctx, mocks.client, factoryFunc, setBackendFunc, cacheKeyFunc)
	require.NoError(t, err)

	mocks.backendObjectStore.EXPECT().HasAccess(gomock.Any(), "get").Return(nil)
	got, err := watch.Get(ctx, getKey)
	require.NoError(t, err)

	expected := testutil.ToUnstructured(t, pod1)

	assert.Equal(t, expected, got)
}

func TestWatch_UpdateClusterClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mocks := newWatchMocks(t)
	defer mocks.controller.Finish()

	factoryFunc := func(c *Watch) {
		c.initFactoryFunc = func(cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return mocks.informerFactory, nil
		}

		c.initBackendFunc = func(w *Watch) (store.Store, error) {
			return mocks.backendObjectStore, nil
		}
	}

	namespaces := []string{"test"}
	mocks.namespaceClient.EXPECT().Names().Return(namespaces, nil).AnyTimes()

	mocks.client.EXPECT().NamespaceClient().Return(mocks.namespaceClient, nil)

	nsKey := store.Key{APIVersion: "v1", Kind: "Namespace"}
	mocks.backendObjectStore.EXPECT().Watch(gomock.Any(), nsKey, gomock.Any()).Return(nil)

	watch, err := NewWatch(ctx, mocks.client, factoryFunc)
	require.NoError(t, err)

	newClient := clusterfake.NewMockClientInterface(mocks.controller)
	newClient.EXPECT().NamespaceClient().Return(mocks.namespaceClient, nil)

	newBackendStore := objectStoreFake.NewMockStore(mocks.controller)
	newBackendStore.EXPECT().Watch(gomock.Any(), nsKey, gomock.Any()).Return(nil)

	watch.initBackendFunc = func(*Watch) (store.Store, error) {
		return newBackendStore, nil
	}

	err = watch.UpdateClusterClient(ctx, newClient)
	require.NoError(t, err)

	assert.Equal(t, newBackendStore, watch.backendObjectStore)
}

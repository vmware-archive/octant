/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicFake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/informers"
	clientGoTesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"

	"github.com/vmware-tanzu/octant/internal/cluster"
	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	"github.com/vmware-tanzu/octant/internal/objectstore/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
)

var (
	podGVR = schema.GroupVersionResource{Version: "v1", Resource: "pods"}
)

type fakeLister struct {
	listObjects []runtime.Object
	listErr     error

	getObject runtime.Object
	getErr    error
}

var _ cache.GenericLister = (*fakeLister)(nil)

func (l fakeLister) List(selector kLabels.Selector) ([]runtime.Object, error) {
	return l.listObjects, l.listErr
}

func (l fakeLister) Get(name string) (runtime.Object, error) {
	return l.getObject, l.getErr
}

func (l fakeLister) ByNamespace(namespace string) cache.GenericNamespaceLister {
	return l
}

func expectNamespaceAccess(
	accessClient *clusterFake.MockSelfSubjectAccessReviewInterface,
	authClient *clusterFake.MockAuthorizationV1Interface,
	namespaceCount int,
) {
	authClient.EXPECT().SelfSubjectAccessReviews().Return(accessClient).MaxTimes(namespaceCount)
	accessResp := &authorizationv1.SelfSubjectAccessReview{
		Status: authorizationv1.SubjectAccessReviewStatus{
			Allowed: true,
		},
	}
	accessClient.EXPECT().Create(gomock.Any()).Return(accessResp, nil).MaxTimes(namespaceCount)
}

func Test_DynamicCache_List(t *testing.T) {
	h := initDynamicCacheTestHarness(t)
	defer h.finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pod := testutil.CreatePod("pod")

	objects := []runtime.Object{testutil.ToUnstructured(t, pod)}

	l := &fakeLister{listObjects: objects}
	h.setupLister(podGVR, l)

	h.mapResources(pod.GroupVersionKind(), podGVR)

	c, err := h.factory(ctx)
	require.NoError(t, err)

	h.setSynced(t, c, pod)

	key := h.keyFromObject(t, pod)

	got, isLoading, err := c.List(ctx, key)
	require.NoError(t, err)
	require.False(t, isLoading)

	expected := testutil.ToUnstructuredList(t, pod)
	assert.Equal(t, expected, got)
}

func Test_DynamicCache_Get(t *testing.T) {
	h := initDynamicCacheTestHarness(t)
	defer h.finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pod := testutil.CreatePod("pod")

	l := &fakeLister{getObject: testutil.ToUnstructured(t, pod)}
	h.setupLister(podGVR, l)
	h.mapResources(pod.GroupVersionKind(), podGVR)

	c, err := h.factory(ctx)
	require.NoError(t, err)

	h.setSynced(t, c, pod)
	key := h.keyFromObject(t, pod)

	got, found, err := c.Get(ctx, key)
	require.NoError(t, err)
	require.True(t, found)

	expected := testutil.ToUnstructured(t, pod)

	assert.Equal(t, expected, got)
}

func TestDynamicCache_Update(t *testing.T) {
	h := initDynamicCacheTestHarness(t)
	defer h.finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))

	podInformer := h.informerFor(podGVR)
	h.mapResources(pod.GroupVersionKind(), podGVR)

	l := &fakeLister{getObject: pod}
	podInformer.EXPECT().Lister().Return(l)

	scheme := runtime.NewScheme()

	dc := dynamicFake.NewSimpleDynamicClient(scheme, pod)

	h.client.EXPECT().DynamicClient().Return(dc, nil)

	c, err := h.factory(ctx)
	require.NoError(t, err)

	key := h.keyFromObject(t, pod)

	err = c.Update(ctx, key, func(*unstructured.Unstructured) error {
		return nil
	})
	require.NoError(t, err)

	assert.Len(t, dc.Actions(), 1)

	action := dc.Actions()[0]
	assert.Equal(t, "update", action.GetVerb())
}

func TestDynamicCache_Delete(t *testing.T) {
	h := initDynamicCacheTestHarness(t)
	defer h.finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))
	h.mapResources(pod.GroupVersionKind(), podGVR)

	scheme := runtime.NewScheme()

	dc := dynamicFake.NewSimpleDynamicClient(scheme, pod)
	h.client.EXPECT().DynamicClient().Return(dc, nil)

	c, err := h.factory(ctx)
	require.NoError(t, err)

	key := h.keyFromObject(t, pod)

	err = c.Delete(ctx, key)
	require.NoError(t, err)

	assert.Len(t, dc.Actions(), 1)

	expected := clientGoTesting.DeleteActionImpl{
		ActionImpl: clientGoTesting.ActionImpl{
			Namespace: pod.GetNamespace(),
			Verb:      "delete",
			Resource:  podGVR,
		},
		Name: pod.GetName(),
	}

	got := dc.Actions()[0]
	assert.Equal(t, expected, got)
}

func TestDynamicCache_Unwatch(t *testing.T) {
	h := initDynamicCacheTestHarness(t)
	defer h.finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pod := testutil.CreatePod("pod")
	h.mapResources(pod.GroupVersionKind(), podGVR)

	c, err := h.factory(ctx)
	require.NoError(t, err)

	h.informerFactory.EXPECT().Delete(podGVR)

	err = c.Unwatch(ctx, pod.GroupVersionKind())
	require.NoError(t, err)
}

type dynamicCacheTestHarness struct {
	controller       *gomock.Controller
	client           *clusterFake.MockClientInterface
	informerFactory  *fake.MockInformerFactory
	kubernetesClient *clusterFake.MockKubernetesInterface
	authClient       *clusterFake.MockAuthorizationV1Interface
	namespaceClient  *clusterFake.MockNamespaceInterface
	accessClient     *clusterFake.MockSelfSubjectAccessReviewInterface
}

func initDynamicCacheTestHarness(t *testing.T) *dynamicCacheTestHarness {
	controller := gomock.NewController(t)
	client := clusterFake.NewMockClientInterface(controller)
	informerFactory := fake.NewMockInformerFactory(controller)
	kubernetesClient := clusterFake.NewMockKubernetesInterface(controller)
	authClient := clusterFake.NewMockAuthorizationV1Interface(controller)
	accessClient := clusterFake.NewMockSelfSubjectAccessReviewInterface(controller)
	namespaceClient := clusterFake.NewMockNamespaceInterface(controller)

	namespaces := []string{"test", ""}
	client.EXPECT().NamespaceClient().Return(namespaceClient, nil).AnyTimes()
	namespaceClient.EXPECT().Names().Return(namespaces, nil).MaxTimes(2)
	expectNamespaceAccess(accessClient, authClient, len(namespaces))

	client.EXPECT().KubernetesClient().Return(kubernetesClient, nil).AnyTimes()
	kubernetesClient.EXPECT().AuthorizationV1().Return(authClient).AnyTimes()

	return &dynamicCacheTestHarness{
		controller:       controller,
		client:           client,
		informerFactory:  informerFactory,
		kubernetesClient: kubernetesClient,
		authClient:       authClient,
		accessClient:     accessClient,
		namespaceClient:  namespaceClient,
	}
}

func (h *dynamicCacheTestHarness) finish() {
	h.controller.Finish()
}

func (h *dynamicCacheTestHarness) factory(ctx context.Context) (*DynamicCache, error) {
	factoryFunc := func(c *DynamicCache) {
		c.initFactoryFunc = func(i context.Context, clientInterface cluster.ClientInterface, s string) (factory InformerFactory, e error) {
			return h.informerFactory, nil

		}
		c.waitForSyncFunc = func(context.Context, store.Key, *DynamicCache, informers.GenericInformer, chan bool) {
			return
		}
	}

	resourceAccess := NewResourceAccess(h.client)
	return NewDynamicCache(ctx, h.client, factoryFunc, Access(resourceAccess))
}

func (h *dynamicCacheTestHarness) informerFor(gvr schema.GroupVersionResource) *clusterFake.MockGenericInformer {
	informer := clusterFake.NewMockGenericInformer(h.controller)
	h.informerFactory.EXPECT().
		ForResource(gvr).
		Return(informer)

	return informer
}

func (h *dynamicCacheTestHarness) setupLister(gvr schema.GroupVersionResource, l cache.GenericLister) {
	informer := h.informerFor(gvr)
	informer.EXPECT().Lister().Return(l)
}

func (h *dynamicCacheTestHarness) mapResources(groupVersionKind schema.GroupVersionKind, resource schema.GroupVersionResource) {
	h.client.EXPECT().
		Resource(groupVersionKind.GroupKind()).
		Return(resource, nil).
		AnyTimes()
}

func (h *dynamicCacheTestHarness) setSynced(t *testing.T, c *DynamicCache, object runtime.Object) {
	accessor, err := meta.Accessor(object)
	require.NoError(t, err)

	apiVersion, kind := object.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()

	key := store.Key{
		Namespace:  accessor.GetNamespace(),
		APIVersion: apiVersion,
		Kind:       kind,
	}
	c.informerSynced.setSynced(key, true)
}

func (h *dynamicCacheTestHarness) keyFromObject(t *testing.T, object runtime.Object) store.Key {
	key, err := store.KeyFromObject(object)
	require.NoError(t, err)
	return key
}

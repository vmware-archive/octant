/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/informers"
	testing2 "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"

	"github.com/vmware/octant/internal/cluster"
	clusterfake "github.com/vmware/octant/internal/cluster/fake"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/store"
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
	accessClient *clusterfake.MockSelfSubjectAccessReviewInterface,
	authClient *clusterfake.MockAuthorizationV1Interface,
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	controller := gomock.NewController(t)
	defer controller.Finish()

	client := clusterfake.NewMockClientInterface(controller)
	informerFactory := clusterfake.NewMockDynamicSharedInformerFactory(controller)
	genericInformer := clusterfake.NewMockGenericInformer(controller)
	kubernetesClient := clusterfake.NewMockKubernetesInterface(controller)
	authClient := clusterfake.NewMockAuthorizationV1Interface(controller)
	accessClient := clusterfake.NewMockSelfSubjectAccessReviewInterface(controller)
	namespaceClient := clusterfake.NewMockNamespaceInterface(controller)

	pod := testutil.CreatePod("pod")
	objects := []runtime.Object{testutil.ToUnstructured(t, pod)}

	podGVR := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}
	informerFactory.EXPECT().ForResource(gomock.Eq(podGVR)).Return(genericInformer)

	podGK := schema.GroupKind{
		Kind: "Pod",
	}
	// CheckAccess and currentInformer
	client.EXPECT().Resource(gomock.Eq(podGK)).Return(podGVR, nil).MaxTimes(2)

	client.EXPECT().KubernetesClient().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().AuthorizationV1().Return(authClient)

	client.EXPECT().NamespaceClient().Return(namespaceClient, nil).MaxTimes(2)
	namespaces := []string{"test", ""}
	namespaceClient.EXPECT().Names().Return(namespaces, nil).MaxTimes(2)

	expectNamespaceAccess(accessClient, authClient, len(namespaces))

	informerFactory.EXPECT().Start(gomock.Eq(ctx.Done()))

	l := &fakeLister{listObjects: objects}
	genericInformer.EXPECT().Lister().Return(l)

	factoryFunc := func(c *DynamicCache) {
		c.initFactoryFunc = func(context.Context, cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return informerFactory, nil
		}
		c.waitForSyncFunc = func(context.Context, store.Key, *DynamicCache, informers.GenericInformer, <-chan struct{}, chan bool) {
			return
		}
	}

	resourceAccess := NewResourceAccess(client)

	c, err := NewDynamicCache(ctx, client, factoryFunc, Access(resourceAccess))
	require.NoError(t, err)

	key := store.Key{
		Namespace:  "test",
		APIVersion: "v1",
		Kind:       "Pod",
	}

	got, isLoading, err := c.List(ctx, key)
	require.NoError(t, err)
	require.False(t, isLoading)

	expected := testutil.ToUnstructuredList(t, pod)
	assert.Equal(t, expected, got)
}

func Test_DynamicCache_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	controller := gomock.NewController(t)
	defer controller.Finish()

	client := clusterfake.NewMockClientInterface(controller)
	informerFactory := clusterfake.NewMockDynamicSharedInformerFactory(controller)
	genericInformer := clusterfake.NewMockGenericInformer(controller)
	kubernetesClient := clusterfake.NewMockKubernetesInterface(controller)
	authClient := clusterfake.NewMockAuthorizationV1Interface(controller)
	accessClient := clusterfake.NewMockSelfSubjectAccessReviewInterface(controller)
	namespaceClient := clusterfake.NewMockNamespaceInterface(controller)

	pod := testutil.CreatePod("pod")

	podGVR := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}
	informerFactory.EXPECT().ForResource(gomock.Eq(podGVR)).Return(genericInformer)

	podGK := schema.GroupKind{
		Kind: "Pod",
	}
	// CheckAccess and currentInformer
	client.EXPECT().Resource(gomock.Eq(podGK)).Return(podGVR, nil).MaxTimes(2)

	client.EXPECT().KubernetesClient().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().AuthorizationV1().Return(authClient)

	client.EXPECT().NamespaceClient().Return(namespaceClient, nil).MaxTimes(2)
	namespaces := []string{"test", ""}
	namespaceClient.EXPECT().Names().Return(namespaces, nil).MaxTimes(2)

	expectNamespaceAccess(accessClient, authClient, len(namespaces))

	informerFactory.EXPECT().Start(gomock.Eq(ctx.Done()))

	l := &fakeLister{getObject: testutil.ToUnstructured(t, pod)}
	genericInformer.EXPECT().Lister().Return(l)

	factoryFunc := func(c *DynamicCache) {
		c.initFactoryFunc = func(context.Context, cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return informerFactory, nil
		}
		c.waitForSyncFunc = func(context.Context, store.Key, *DynamicCache, informers.GenericInformer, <-chan struct{}, chan bool) {
			return
		}
	}

	resourceAccess := NewResourceAccess(client)
	c, err := NewDynamicCache(ctx, client, factoryFunc, Access(resourceAccess))
	require.NoError(t, err)

	key := store.Key{
		Namespace:  "test",
		APIVersion: "v1",
		Kind:       "Pod",
		Name:       "pod",
	}

	got, found, err := c.Get(ctx, key)
	require.NoError(t, err)
	require.True(t, found)

	expected := testutil.ToUnstructured(t, pod)

	assert.Equal(t, expected, got)
}

func TestDynamicCache_Update(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))

	controller := gomock.NewController(t)
	defer controller.Finish()

	client := clusterfake.NewMockClientInterface(controller)
	informerFactory := clusterfake.NewMockDynamicSharedInformerFactory(controller)
	genericInformer := clusterfake.NewMockGenericInformer(controller)
	kubernetesClient := clusterfake.NewMockKubernetesInterface(controller)
	authClient := clusterfake.NewMockAuthorizationV1Interface(controller)
	accessClient := clusterfake.NewMockSelfSubjectAccessReviewInterface(controller)
	namespaceClient := clusterfake.NewMockNamespaceInterface(controller)

	podGVR := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}
	informerFactory.EXPECT().ForResource(gomock.Eq(podGVR)).Return(genericInformer)

	podGK := schema.GroupKind{
		Kind: "Pod",
	}
	// CheckAccess and currentInformer
	client.EXPECT().Resource(gomock.Eq(podGK)).Return(podGVR, nil).AnyTimes()

	client.EXPECT().KubernetesClient().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().AuthorizationV1().Return(authClient)

	client.EXPECT().NamespaceClient().Return(namespaceClient, nil).AnyTimes()
	namespaces := []string{"test", ""}
	namespaceClient.EXPECT().Names().Return(namespaces, nil).AnyTimes()

	expectNamespaceAccess(accessClient, authClient, len(namespaces))

	informerFactory.EXPECT().Start(gomock.Eq(ctx.Done()))

	l := &fakeLister{getObject: pod}
	genericInformer.EXPECT().Lister().Return(l)

	factoryFunc := func(c *DynamicCache) {
		c.initFactoryFunc = func(context.Context, cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return informerFactory, nil
		}
		c.waitForSyncFunc = func(context.Context, store.Key, *DynamicCache, informers.GenericInformer, <-chan struct{}, chan bool) {
			return
		}
	}

	scheme := runtime.NewScheme()

	dc := dynamicfake.NewSimpleDynamicClient(scheme, pod)

	client.EXPECT().DynamicClient().Return(dc, nil)

	resourceAccess := NewResourceAccess(client)
	c, err := NewDynamicCache(ctx, client, factoryFunc, Access(resourceAccess))
	require.NoError(t, err)

	key, err := store.KeyFromObject(pod)
	require.NoError(t, err)

	err = c.Update(ctx, key, func(*unstructured.Unstructured) error {
		return nil
	})
	require.NoError(t, err)

	assert.Len(t, dc.Actions(), 1)

	action := dc.Actions()[0]
	assert.Equal(t, "update", action.GetVerb())
}

func TestDynamicCache_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))

	controller := gomock.NewController(t)
	defer controller.Finish()

	client := clusterfake.NewMockClientInterface(controller)
	informerFactory := clusterfake.NewMockDynamicSharedInformerFactory(controller)
	kubernetesClient := clusterfake.NewMockKubernetesInterface(controller)
	authClient := clusterfake.NewMockAuthorizationV1Interface(controller)
	accessClient := clusterfake.NewMockSelfSubjectAccessReviewInterface(controller)
	namespaceClient := clusterfake.NewMockNamespaceInterface(controller)

	podGVR := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}

	podGK := schema.GroupKind{
		Kind: "Pod",
	}
	// CheckAccess and currentInformer
	client.EXPECT().Resource(gomock.Eq(podGK)).Return(podGVR, nil).AnyTimes()

	client.EXPECT().KubernetesClient().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().AuthorizationV1().Return(authClient)

	client.EXPECT().NamespaceClient().Return(namespaceClient, nil).AnyTimes()
	namespaces := []string{"test", ""}
	namespaceClient.EXPECT().Names().Return(namespaces, nil).AnyTimes()

	expectNamespaceAccess(accessClient, authClient, len(namespaces))

	factoryFunc := func(c *DynamicCache) {
		c.initFactoryFunc = func(context.Context, cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return informerFactory, nil
		}
		c.waitForSyncFunc = func(context.Context, store.Key, *DynamicCache, informers.GenericInformer, <-chan struct{}, chan bool) {
			return
		}
	}

	scheme := runtime.NewScheme()

	dc := dynamicfake.NewSimpleDynamicClient(scheme, pod)

	client.EXPECT().DynamicClient().Return(dc, nil)

	resourceAccess := NewResourceAccess(client)
	c, err := NewDynamicCache(ctx, client, factoryFunc, Access(resourceAccess))
	require.NoError(t, err)

	key, err := store.KeyFromObject(pod)
	require.NoError(t, err)

	err = c.Delete(ctx, key)
	require.NoError(t, err)

	assert.Len(t, dc.Actions(), 1)

	expected := testing2.DeleteActionImpl{
		ActionImpl: testing2.ActionImpl{
			Namespace: pod.GetNamespace(),
			Verb:      "delete",
			Resource:  podGVR,
		},
		Name: pod.GetName(),
	}

	got := dc.Actions()[0]
	assert.Equal(t, expected, got)
}

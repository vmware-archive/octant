package objectstore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/cluster"
	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/cacheutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
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

func Test_DynamicCache_List(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	controller := gomock.NewController(t)
	defer controller.Finish()

	client := clusterfake.NewMockClientInterface(controller)
	informerFactory := clusterfake.NewMockDynamicSharedInformerFactory(controller)
	sharedIndexInformer := clusterfake.NewMockSharedIndexInformer(controller)
	informer := clusterfake.NewMockGenericInformer(controller)

	informer.EXPECT().Informer().Return(sharedIndexInformer)
	sharedIndexInformer.EXPECT().HasSynced().Return(true)

	pod := testutil.CreatePod("pod")
	objects := []runtime.Object{pod}

	podGVR := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}
	informerFactory.EXPECT().ForResource(gomock.Eq(podGVR)).Return(informer)

	podGK := schema.GroupKind{
		Kind: "Pod",
	}
	client.EXPECT().Resource(gomock.Eq(podGK)).Return(podGVR, nil)

	informerFactory.EXPECT().Start(gomock.Eq(ctx.Done()))

	l := &fakeLister{listObjects: objects}
	informer.EXPECT().Lister().Return(l)

	factoryFunc := func(c *DynamicCache) {
		c.initFactoryFunc = func(cluster.ClientInterface) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return informerFactory, nil
		}
	}

	c, err := NewDynamicCache(client, ctx.Done(), factoryFunc)
	require.NoError(t, err)

	key := cacheutil.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Pod",
	}

	got, err := c.List(ctx, key)
	require.NoError(t, err)

	expected := []*unstructured.Unstructured{
		testutil.ToUnstructured(t, pod),
	}

	assert.Equal(t, expected, got)
}

func Test_DynamicCache_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	controller := gomock.NewController(t)
	defer controller.Finish()

	client := clusterfake.NewMockClientInterface(controller)
	informerFactory := clusterfake.NewMockDynamicSharedInformerFactory(controller)
	informer := clusterfake.NewMockGenericInformer(controller)
	sharedIndexInformer := clusterfake.NewMockSharedIndexInformer(controller)
	informer.EXPECT().Informer().Return(sharedIndexInformer)
	sharedIndexInformer.EXPECT().HasSynced().Return(true)

	pod := testutil.CreatePod("pod")

	podGVR := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}
	informerFactory.EXPECT().ForResource(gomock.Eq(podGVR)).Return(informer)

	podGK := schema.GroupKind{
		Kind: "Pod",
	}
	client.EXPECT().Resource(gomock.Eq(podGK)).Return(podGVR, nil)

	informerFactory.EXPECT().Start(gomock.Eq(ctx.Done()))

	l := &fakeLister{getObject: pod}
	informer.EXPECT().Lister().Return(l)

	factoryFunc := func(c *DynamicCache) {
		c.initFactoryFunc = func(cluster.ClientInterface) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return informerFactory, nil
		}
	}

	c, err := NewDynamicCache(client, ctx.Done(), factoryFunc)
	require.NoError(t, err)

	key := cacheutil.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Pod",
		Name:       "pod",
	}

	got, err := c.Get(ctx, key)
	require.NoError(t, err)

	expected := testutil.ToUnstructured(t, pod)

	assert.Equal(t, expected, got)
}

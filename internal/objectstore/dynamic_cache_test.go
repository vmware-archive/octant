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
	"k8s.io/client-go/tools/cache"

	"github.com/heptio/developer-dash/internal/cluster"
	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/store"
	"github.com/heptio/developer-dash/third_party/k8s.io/client-go/dynamic/dynamicinformer"
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
	sharedIndexInformer := clusterfake.NewMockSharedIndexInformer(controller)
	informer := clusterfake.NewMockGenericInformer(controller)
	kubernetesClient := clusterfake.NewMockKubernetesInterface(controller)
	authClient := clusterfake.NewMockAuthorizationV1Interface(controller)
	accessClient := clusterfake.NewMockSelfSubjectAccessReviewInterface(controller)
	namespaceClient := clusterfake.NewMockNamespaceInterface(controller)

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
	informer.EXPECT().Lister().Return(l)

	factoryFunc := func(c *DynamicCache) {
		c.initFactoryFunc = func(cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return informerFactory, nil
		}
	}

	c, err := NewDynamicCache(client, ctx.Done(), factoryFunc)
	require.NoError(t, err)

	key := store.Key{
		Namespace:  "test",
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
	kubernetesClient := clusterfake.NewMockKubernetesInterface(controller)
	authClient := clusterfake.NewMockAuthorizationV1Interface(controller)
	accessClient := clusterfake.NewMockSelfSubjectAccessReviewInterface(controller)
	namespaceClient := clusterfake.NewMockNamespaceInterface(controller)

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
	// CheckAccess and currentInformer
	client.EXPECT().Resource(gomock.Eq(podGK)).Return(podGVR, nil).MaxTimes(2)

	client.EXPECT().KubernetesClient().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().AuthorizationV1().Return(authClient)

	client.EXPECT().NamespaceClient().Return(namespaceClient, nil).MaxTimes(2)
	namespaces := []string{"test", ""}
	namespaceClient.EXPECT().Names().Return(namespaces, nil).MaxTimes(2)

	expectNamespaceAccess(accessClient, authClient, len(namespaces))

	informerFactory.EXPECT().Start(gomock.Eq(ctx.Done()))

	l := &fakeLister{getObject: pod}
	informer.EXPECT().Lister().Return(l)

	factoryFunc := func(c *DynamicCache) {
		c.initFactoryFunc = func(cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return informerFactory, nil
		}
	}

	c, err := NewDynamicCache(client, ctx.Done(), factoryFunc)
	require.NoError(t, err)

	key := store.Key{
		Namespace:  "test",
		APIVersion: "v1",
		Kind:       "Pod",
		Name:       "pod",
	}

	got, err := c.Get(ctx, key)
	require.NoError(t, err)

	expected := testutil.ToUnstructured(t, pod)

	assert.Equal(t, expected, got)
}

func Test_DynamicCache_HasAccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	controller := gomock.NewController(t)
	defer controller.Finish()

	client := clusterfake.NewMockClientInterface(controller)
	informerFactory := clusterfake.NewMockDynamicSharedInformerFactory(controller)
	namespaceClient := clusterfake.NewMockNamespaceInterface(controller)

	client.EXPECT().NamespaceClient().Return(namespaceClient, nil).MaxTimes(3)
	namespaces := []string{"test", ""}
	namespaceClient.EXPECT().Names().Return(namespaces, nil).MaxTimes(3)

	factoryFunc := func(c *DynamicCache) {
		c.initFactoryFunc = func(cluster.ClientInterface, string) (dynamicinformer.DynamicSharedInformerFactory, error) {
			return informerFactory, nil
		}
	}

	scenarios := []struct {
		name       string
		resource   string
		key        store.Key
		accessFunc func(c *DynamicCache)
		expectErr  bool
	}{
		{
			name:     "pods",
			resource: "pods",
			key: store.Key{
				APIVersion: "apps/v1",
				Kind:       "Pod",
			},
			accessFunc: func(c *DynamicCache) {
				access := make(accessMap)
				aKey := accessKey{
					Namespace: "",
					Group:     "apps",
					Resource:  "pods",
					Verb:      "get",
				}
				access[aKey] = true
				c.access = access
			},
			expectErr: false,
		},
		{
			name:     "crds",
			resource: "customresourcedefinitions",
			key: store.Key{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			accessFunc: func(c *DynamicCache) {
				access := make(accessMap)
				aKey := accessKey{
					Namespace: "",
					Group:     "apiextensions.k8s.io",
					Resource:  "customresourcedefinitions",
					Verb:      "get",
				}
				access[aKey] = true
				c.access = access
			},
			expectErr: false,
		},
		{
			name:     "no access crds",
			resource: "customresourcedefinitions",
			key: store.Key{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			accessFunc: func(c *DynamicCache) {
				access := make(accessMap)
				aKey := accessKey{
					Namespace: "",
					Group:     "apiextensions.k8s.io",
					Resource:  "customresourcedefinitions",
					Verb:      "get",
				}
				access[aKey] = false
				c.access = access
			},
			expectErr: true,
		},
	}

	for _, ts := range scenarios {
		t.Run(ts.name, func(t *testing.T) {
			c, err := NewDynamicCache(client, ctx.Done(), factoryFunc, ts.accessFunc)
			require.NoError(t, err)

			gvk := ts.key.GroupVersionKind()
			podGVR := schema.GroupVersionResource{
				Group:    gvk.Group,
				Version:  gvk.Version,
				Resource: ts.resource,
			}
			client.EXPECT().Resource(gomock.Eq(gvk.GroupKind())).Return(podGVR, nil)

			if ts.expectErr {
				require.Error(t, c.HasAccess(ts.key, "get"))
			} else {
				require.NoError(t, c.HasAccess(ts.key, "get"))
			}
		})
	}
}

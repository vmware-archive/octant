package objectstore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/cluster"
	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authorizationapi "k8s.io/api/authorization/v1"
	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	kubernetesClient := clusterfake.NewMockKubernetesInterface(controller)
	authClient := clusterfake.NewMockAuthorizationV1Interface(controller)
	accessClient := clusterfake.NewMockSelfSubjectAccessReviewInterface(controller)

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

	client.EXPECT().KubernetesClient().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().AuthorizationV1().Return(authClient)

	authResp := &authorizationapi.SelfSubjectAccessReview{
		Status: authorizationapi.SubjectAccessReviewStatus{
			Allowed: true,
		},
	}

	verbs := []string{"get", "list", "watch"}
	for _, verb := range verbs {
		resourceAttributes := &authorizationv1.ResourceAttributes{
			Verb:      verb,
			Version:   "v1",
			Resource:  "pods",
			Namespace: metav1.NamespaceAll,
		}

		sar := &authorizationv1.SelfSubjectAccessReview{
			Spec: authorizationv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: resourceAttributes,
			},
		}
		authClient.EXPECT().SelfSubjectAccessReviews().Return(accessClient)
		accessClient.EXPECT().Create(sar).Return(authResp, nil)
	}

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

	key := objectstoreutil.Key{
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
	kubernetesClient := clusterfake.NewMockKubernetesInterface(controller)
	authClient := clusterfake.NewMockAuthorizationV1Interface(controller)
	accessClient := clusterfake.NewMockSelfSubjectAccessReviewInterface(controller)

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

	client.EXPECT().KubernetesClient().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().AuthorizationV1().Return(authClient)

	authResp := &authorizationapi.SelfSubjectAccessReview{
		Status: authorizationapi.SubjectAccessReviewStatus{
			Allowed: true,
		},
	}

	verbs := []string{"get", "list", "watch"}
	for _, verb := range verbs {
		resourceAttributes := &authorizationv1.ResourceAttributes{
			Verb:      verb,
			Version:   "v1",
			Resource:  "pods",
			Namespace: metav1.NamespaceAll,
		}

		sar := &authorizationv1.SelfSubjectAccessReview{
			Spec: authorizationv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: resourceAttributes,
			},
		}
		authClient.EXPECT().SelfSubjectAccessReviews().Return(accessClient)
		accessClient.EXPECT().Create(sar).Return(authResp, nil)
	}

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

	key := objectstoreutil.Key{
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

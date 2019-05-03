package objectstore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/cluster"
	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/third_party/k8s.io/client-go/dynamic/dynamicinformer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

func verbs() []string {
	return []string{"get", "list", "watch"}
}

func expectNamespaceRules(
	namespaces []string,
	authClient *clusterfake.MockAuthorizationV1Interface,
	rulesClient *clusterfake.MockSelfSubjectRulesReviewInterface,
) {
	resourceRules := []authorizationv1.ResourceRule{
		authorizationv1.ResourceRule{
			Verbs:     verbs(),
			APIGroups: []string{""},
			Resources: []string{"pods"},
		},
	}

	rulesResp := &authorizationv1.SelfSubjectRulesReview{
		Status: authorizationv1.SubjectRulesReviewStatus{
			ResourceRules: resourceRules,
		},
	}

	for _, namespace := range namespaces {
		srr := &authorizationv1.SelfSubjectRulesReview{
			Spec: authorizationv1.SelfSubjectRulesReviewSpec{
				Namespace: namespace,
			},
		}
		authClient.EXPECT().SelfSubjectRulesReviews().Return(rulesClient)
		rulesClient.EXPECT().Create(srr).Return(rulesResp, nil)
	}
}

func expectClusterScopedRoles(
	accessClient *clusterfake.MockSelfSubjectAccessReviewInterface,
	authClient *clusterfake.MockAuthorizationV1Interface,
) {
	// Cluster scoped access checking
	for _, gvr := range []struct {
		group    string
		version  string
		resource string
		key      string
	}{
		{"apiextensions.k8s.io", "v1beta1", "CustomResourceDefinition", "customresourcedefinitions"},
		{"rbac.authorization.k8s.io", "v1", "ClusterRole", "clusterroles"},
		{"rbac.authorization.k8s.io", "v1", "ClusterRoleBinding", "clusterrolebindings"},
	} {
		for _, verb := range verbs() {
			resourceAttributes := &authorizationv1.ResourceAttributes{
				Verb:     verb,
				Group:    gvr.group,
				Version:  gvr.version,
				Resource: gvr.resource,
			}

			sar := &authorizationv1.SelfSubjectAccessReview{
				Spec: authorizationv1.SelfSubjectAccessReviewSpec{
					ResourceAttributes: resourceAttributes,
				},
			}

			accessResp := &authorizationv1.SelfSubjectAccessReview{
				Status: authorizationv1.SubjectAccessReviewStatus{
					Allowed: true,
				},
			}
			accessClient.EXPECT().Create(sar).Return(accessResp, nil)
			authClient.EXPECT().SelfSubjectAccessReviews().Return(accessClient)
		}
	}
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
	rulesClient := clusterfake.NewMockSelfSubjectRulesReviewInterface(controller)
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

	expectNamespaceRules(namespaces, authClient, rulesClient)
	expectClusterScopedRoles(accessClient, authClient)

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

	key := objectstoreutil.Key{
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
	rulesClient := clusterfake.NewMockSelfSubjectRulesReviewInterface(controller)
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

	expectNamespaceRules(namespaces, authClient, rulesClient)
	expectClusterScopedRoles(accessClient, authClient)

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

	key := objectstoreutil.Key{
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

func Test_DynamicCache_hasGetListWatch(t *testing.T) {
	scenarios := []struct {
		name     string
		verbs    []string
		expected bool
	}{
		// scenario 1
		{
			name:     "only get, list, watch",
			verbs:    verbs(),
			expected: true,
		},
		// scenario 2
		{
			name:     "get, list, but missing watch",
			verbs:    []string{"get", "list", "create"},
			expected: false,
		},
		// scenario 3
		{
			name:     "* (all verbs)",
			verbs:    []string{"*"},
			expected: true,
		},
		// scenario 4
		{
			name:     "create, watch, list, get",
			verbs:    []string{"create", "watch", "list", "get"},
			expected: true,
		},
	}

	for _, ts := range scenarios {
		t.Run(ts.name, func(t *testing.T) {
			hasAccess := hasGetListWatch(ts.verbs)
			if ts.expected != hasAccess {
				t.Errorf("expected %t got %t", ts.expected, hasAccess)
			}
		})
	}
}

func Test_DynamicCache_CheckAccess(t *testing.T) {
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
		key        objectstoreutil.Key
		accessFunc func(c *DynamicCache)
		expectErr  bool
	}{
		{
			name:     "pods",
			resource: "pods",
			key: objectstoreutil.Key{
				APIVersion: "apps/v1",
				Kind:       "Pod",
			},
			accessFunc: func(c *DynamicCache) {
				access := make(accessMap)
				access[""] = make(map[string]map[string]bool)
				access[""]["apps"] = map[string]bool{"pods": true}
				c.access = access
			},
			expectErr: false,
		},
		{
			name:     "crds",
			resource: "customresourcedefinitions",
			key: objectstoreutil.Key{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			accessFunc: func(c *DynamicCache) {
				access := make(accessMap)
				access[""] = make(map[string]map[string]bool)
				access[""]["apiextensions.k8s.io"] = map[string]bool{"customresourcedefinitions": true}
				c.access = access
			},
			expectErr: false,
		},
		{
			name:     "no access crds",
			resource: "customresourcedefinitions",
			key: objectstoreutil.Key{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			accessFunc: func(c *DynamicCache) {
				access := make(accessMap)
				access[""] = make(map[string]map[string]bool)
				access[""]["apiextensions.k8s.io"] = map[string]bool{"customresourcedefinitions": false}
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
				require.Error(t, c.CheckAccess(ts.key))
			} else {
				require.NoError(t, c.CheckAccess(ts.key))
			}
		})
	}
}

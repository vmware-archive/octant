package cache

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var resources = []*metav1.APIResourceList{
	{
		GroupVersion: "apps/v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "deployments",
				SingularName: "deployment",
				Group:        "apps",
				Version:      "v1",
				Kind:         "Deployment",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
	{
		GroupVersion: "extensions/v1beta1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "ingresses",
				SingularName: "ingress",
				Group:        "extensions",
				Version:      "v1beta1",
				Kind:         "Ingress",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
	{
		GroupVersion: "v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "services",
				SingularName: "service",
				Group:        "",
				Version:      "v1",
				Kind:         "Service",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
			metav1.APIResource{
				Name:         "secrets",
				SingularName: "secret",
				Group:        "",
				Version:      "v1",
				Kind:         "Secret",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
	{
		GroupVersion: "bar/v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "bars",
				SingularName: "bar",
				Group:        "bar",
				Version:      "v1",
				Kind:         "Bar",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
	{
		GroupVersion: "foo/v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "kinds",
				SingularName: "kind",
				Group:        "foo",
				Version:      "v1",
				Kind:         "Kind",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
			metav1.APIResource{
				Name:         "foos",
				SingularName: "foo",
				Group:        "foo",
				Version:      "v1",
				Kind:         "Foo",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
			metav1.APIResource{
				Name:         "others",
				SingularName: "other",
				Group:        "foo",
				Version:      "v1",
				Kind:         "Other",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
}

type cancelFunc func()

func cancelNop() {}

func newCache(t *testing.T, objects []runtime.Object) (*InformerCache, cancelFunc, error) {
	scheme := newScheme()

	client, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)
	if err != nil {
		return nil, cancelNop, err
	}
	stopCh := make(chan struct{})

	restMapper, err := client.RESTMapper()
	require.NoError(t, err, "fetching RESTMapper")
	return NewInformerCache(stopCh, client.FakeDynamic, restMapper), func() { close(stopCh) }, nil
}

func TestInformerCache_List(t *testing.T) {
	objects := []runtime.Object{}
	for _, u := range genObjectsSeed() {
		objects = append(objects, u)
	}

	cases := []struct {
		name        string
		key         Key
		expectedLen int
		expectErr   bool
	}{
		{
			name: "ns, apiVersion, kind, name",
			key: Key{
				Namespace:  "default",
				APIVersion: "foo/v1",
				Kind:       "Kind",
				Name:       "foo1",
			},
			expectedLen: 1,
		},
		{
			name: "ns, apiVersion, kind",
			key: Key{
				Namespace:  "default",
				APIVersion: "foo/v1",
				Kind:       "Kind",
			},
			expectedLen: 2,
		},
		{
			name: "ns, apiVersion: error because we require kind",
			key: Key{
				Namespace:  "default",
				APIVersion: "foo/v1",
			},
			expectErr: true,
		},
		{
			name: "not found",
			key: Key{
				Namespace:  "default",
				APIVersion: "foo/v1",
				Kind:       "Kind",
				Name:       "does-not-exist",
			},
			expectedLen: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, cancel, err := newCache(t, objects)
			require.NoError(t, err)

			objs, err := c.List(tc.key)
			hadErr := (err != nil)
			assert.Equalf(t, tc.expectErr, hadErr, "error mismatch: %v", err)
			assert.Len(t, objs, tc.expectedLen)
			cancel()
		})
	}
}

func TestInformerCache_Watch(t *testing.T) {
	scheme := newScheme()

	objects := []runtime.Object{
		newUnstructured("apps/v1", "Deployment", "default", "deploy3"),
	}

	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	discoveryClient := clusterClient.FakeDiscovery
	discoveryClient.Resources = resources

	dynamicClient := clusterClient.FakeDynamic

	notifyCh := make(chan Notification)
	notifyDone := make(chan struct{})

	restMapper, err := clusterClient.RESTMapper()
	require.NoError(t, err, "fetching RESTMapper")

	cache := NewInformerCache(notifyDone, dynamicClient, restMapper,
		InformerCacheNotificationOpt(notifyCh, notifyDone),
		InformerCacheLoggerOpt(log.TestLogger(t)),
	)

	defer func() {
		close(notifyDone)
	}()

	// verify predefined objects are present
	cacheKey := Key{Namespace: "default", APIVersion: "apps/v1", Kind: "Deployment"}
	found, err := cache.List(cacheKey)
	require.NoError(t, err)

	require.Len(t, found, 1)

	// drain initial object notifications (we expect an ADD followed by an UPDATE)
	for i := 0; i < 2; i++ {
		select {
		case <-time.After(500 * time.Millisecond):
			t.Fatal("timed out wating for create object to notify")
		case <-notifyCh:
		}
	}

	// define new object
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("apps/v1")
	obj.SetKind("Deployment")
	obj.SetName("deploy2")
	obj.SetNamespace("default")

	res := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	resClient := dynamicClient.Resource(res).Namespace("default")

	// create object
	_, err = resClient.Create(obj, metav1.CreateOptions{})
	require.NoError(t, err)

	// wait for cache to store an item before proceeding.
	select {
	case <-time.After(10 * time.Second):
		t.Fatal("timed out wating for create object to notify")
	case <-notifyCh:
	}

	found, err = cache.List(cacheKey)
	require.NoError(t, err)

	// 2 == initial + the new object
	require.Len(t, found, 2)

	annotations := map[string]string{"update": "update"}
	obj.SetAnnotations(annotations)

	// update object
	_, err = resClient.Update(obj, metav1.UpdateOptions{})
	require.NoError(t, err)

	// wait for cache to store an item before proceeding.
	select {
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for update object to notify")
	case <-notifyCh:
	}

	found, err = cache.List(cacheKey)
	require.NoError(t, err)

	require.Len(t, found, 2)

	// Find the object we updated
	var match bool
	for _, u := range found {
		if u.GetName() == obj.GetName() && u.GroupVersionKind() == obj.GroupVersionKind() {
			match = true
			require.Equal(t, annotations, u.GetAnnotations())
		}
	}
	require.True(t, match, "unable to find object from fetched results")
}

func TestInformerCache_Watch_Stop(t *testing.T) {
	scheme := newScheme()

	objects := []runtime.Object{}

	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	discoveryClient := clusterClient.FakeDiscovery
	discoveryClient.Resources = resources

	dynamicClient := clusterClient.FakeDynamic

	notifyCh := make(chan Notification)
	notifyDone := make(chan struct{})

	restMapper, err := clusterClient.RESTMapper()
	require.NoError(t, err, "fetching RESTMapper")

	cache := NewInformerCache(notifyDone, dynamicClient, restMapper,
		InformerCacheNotificationOpt(notifyCh, notifyDone),
		InformerCacheLoggerOpt(log.TestLogger(t)),
	)

	// verify predefined objects are present
	cacheKey := Key{Namespace: "default", APIVersion: "apps/v1", Kind: "Deployment"}
	found, err := cache.List(cacheKey)
	require.NoError(t, err)

	require.Len(t, found, 0)

	// define new object
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("apps/v1")
	obj.SetKind("Deployment")
	obj.SetName("deploy2")
	obj.SetNamespace("default")

	res := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	resClient := dynamicClient.Resource(res).Namespace("default")

	// Stop notifications
	close(notifyDone)

	// Drain notifications
	closeDone := make(chan struct{})
	go func() {
		for range notifyCh {
		}
		close(closeDone)
	}()

	// Wait for informers to shutdown
	select {
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for notification channel to close")
	case <-closeDone:
	}

	// create object
	_, err = resClient.Create(obj, metav1.CreateOptions{})
	require.NoError(t, err)

	found, err = cache.List(cacheKey)
	require.NoError(t, err)

	// The second object is not seen because we shutdown the informer
	require.Len(t, found, 0)
}

func TestChannelContext(t *testing.T) {
	parentCh := make(chan struct{})
	done1 := make(chan struct{})
	done2 := make(chan struct{})
	var count int32 = 2
	ctx1, cancel1 := channelContext(parentCh)
	defer cancel1()
	ctx2, cancel2 := channelContext(parentCh)
	defer cancel2()

	go func() {
		<-ctx1.Done()
		atomic.AddInt32(&count, -1)
		close(done1)
	}()
	go func() {
		<-ctx2.Done()
		atomic.AddInt32(&count, -1)
		close(done2)
	}()

	// Initial state
	assert.Equal(t, int32(2), atomic.LoadInt32(&count))

	// Canceling ctx1 (a child) should only cancel that context, not its siblings
	cancel1()
	<-done1
	assert.Equal(t, int32(1), atomic.LoadInt32(&count))

	// Canceling parentCh should cancel all remaining child contexts
	close(parentCh)
	<-done2
	assert.Equal(t, int32(0), atomic.LoadInt32(&count))
}

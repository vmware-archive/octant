package fake

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/third_party/dynamicfake"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/dynamic"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/testing"
)

var (
	scheme         = runtime.NewScheme()
	codecs         = serializer.NewCodecFactory(scheme)
	parameterCodec = runtime.NewParameterCodec(scheme)
)

// Client implements cluster.Interface.
type Client struct {
	FakeDynamic   *dynamicfake.FakeDynamicClient
	FakeDiscovery *fakediscovery.FakeDiscovery
}

// NewClient creates an instance of Client.
func NewClient(scheme *runtime.Scheme, resources []*metav1.APIResourceList, objects []runtime.Object) (*Client, error) {
	client := fakeclientset.NewSimpleClientset()
	fakeDiscovery, ok := client.Discovery().(*fakediscovery.FakeDiscovery)
	if !ok {
		return nil, errors.New("couldn't convert Discovery() to *FakeDiscovery")
	}
	fakeDiscovery.Resources = resources

	restMapper, err := restMapper(fakeDiscovery)
	if err != nil {
		return nil, errors.Wrap(err, "constructing RESTMapper - check resources")
	}
	dynamicClient := NewSimpleDynamicClient(scheme, restMapper, objects...)

	return &Client{
		FakeDynamic:   dynamicClient,
		FakeDiscovery: fakeDiscovery,
	}, nil
}

// DynamicClient returns a dynamic client or an error.
func (c *Client) DynamicClient() (dynamic.Interface, error) {
	return c.FakeDynamic, nil
}

// DiscoveryClient returns a discovery client or an error.
func (c *Client) DiscoveryClient() (discovery.DiscoveryInterface, error) {
	return c.FakeDiscovery, nil
}

// NamespaceClient returns a namspace client or an error.
func (c *Client) NamespaceClient() (cluster.NamespaceInterface, error) {
	return &NamespaceClient{}, nil
}

// RESTMapper returns a RESTMapper using the client's discovery interface.
// The mappings depend on the resources supplied in NewClient.
func (c *Client) RESTMapper() (meta.RESTMapper, error) {
	return restMapper(c.FakeDiscovery)
}

func restMapper(discoveryClient discovery.DiscoveryInterface) (meta.RESTMapper, error) {
	resources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return nil, err
	}

	return restmapper.NewDiscoveryRESTMapper(resources), nil
}

// registerListKind registers a List kind in the provided scheme as
// a list container for the provided non-list kind.
// This container list must be present in the scheme for a list operation to succeed.
func registerListKind(scheme *runtime.Scheme, gvk schema.GroupVersionKind) {
	// Heuristic for list kind: original kind + List suffix. Might
	// not always be true but this tracker has a pretty limited
	// understanding of the actual API model.
	listGVK := gvk
	listGVK.Kind = listGVK.Kind + "List"
	// GVK does have the concept of "internal version". The scheme recognizes
	// the runtime.APIVersionInternal, but not the empty string.
	if listGVK.Version == "" {
		listGVK.Version = runtime.APIVersionInternal
	}

	if scheme.Recognizes(listGVK) {
		return
	}

	scheme.AddKnownTypeWithName(listGVK, &unstructured.UnstructuredList{})
}

// NewSimpleDynamicClient creates a FakeDynamicClient which fixes behavior from dynamicfake.NewSimpleDynamicClient - we properly forward
// ADDED events for preexisting objects when adding watches.
func NewSimpleDynamicClient(scheme *runtime.Scheme, restMapper meta.RESTMapper, objects ...runtime.Object) *dynamicfake.FakeDynamicClient {
	// In order to use List with this client, you have to have the v1.List registered in your scheme. Neat thing though
	// it does NOT have to be the *same* list
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{Group: "fake-dynamic-client-group", Version: "v1", Kind: "List"}, &unstructured.UnstructuredList{})

	codecs := serializer.NewCodecFactory(scheme)
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	cs := &dynamicfake.FakeDynamicClient{}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		w, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}

		gvk, err := restMapper.KindFor(gvr)
		if err != nil {
			fmt.Printf("OH NO THE SKY IS FALLING %#v -> %#v\n", gvr, err)
			return false, nil, fmt.Errorf("no registered kind for resource: %v", gvr.String())
		}

		// JIT register *List kinds in the scheme to support List operations
		registerListKind(scheme, gvk)

		l, err := o.List(gvr, gvk, ns)
		if err != nil {
			return false, nil, errors.Wrap(err, "listing existing objects")
		}

		// Replay existing objects
		rfw, ok := w.(*watch.RaceFreeFakeWatcher)
		if !ok {
			return false, nil, fmt.Errorf("unexpected watch type: %T", w)
		}

		ul, ok := l.(*unstructured.UnstructuredList)
		if !ok {
			return false, nil, errors.Errorf("wrong type for list: %T\n", l)
		}
		err = ul.EachListItem(func(obj runtime.Object) error {
			rfw.Add(obj)
			return nil
		})

		return true, w, nil
	})

	return cs
}

package fake

import (
	"errors"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/third_party/dynamicfake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/dynamic"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

// Client implements cluster.Interface.
type Client struct {
	FakeDynamic   *dynamicfake.FakeDynamicClient
	FakeDiscovery *fakediscovery.FakeDiscovery
}

// NewClient creates an instance of Client.
func NewClient(scheme *runtime.Scheme, objects []runtime.Object) (*Client, error) {
	dynamicClient := dynamicfake.NewSimpleDynamicClient(scheme, objects...)

	client := fakeclientset.NewSimpleClientset()
	fakeDiscovery, ok := client.Discovery().(*fakediscovery.FakeDiscovery)
	if !ok {
		return nil, errors.New("couldn't convert Discovery() to *FakeDiscovery")
	}

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

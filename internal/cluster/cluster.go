/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/install"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	internalLog "github.com/vmware-tanzu/octant/internal/log"
	clusterTypes "github.com/vmware-tanzu/octant/pkg/cluster"
	"github.com/vmware-tanzu/octant/pkg/log"

	// auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

//go:generate mockgen -source=../../pkg/cluster/client.go -destination=./fake/mock_client_interface.go -package=fake github.com/vmware-tanzu/octant/pkg/cluster ClientInterface
//go:generate mockgen -source=../../pkg/cluster/namespace.go -destination=./fake/mock_namespace_interface.go -package=fake github.com/vmware-tanzu/octant/pkg/cluster NamespaceInterface
//go:generate mockgen -source=../../pkg/cluster/info.go -destination=./fake/mock_info_interface.go -package=fake github.com/vmware-tanzu/octant/pkg/cluster InfoInterface
//go:generate mockgen -source=../../vendor/k8s.io/client-go/informers/generic.go -destination=./fake/mock_genericinformer.go -package=fake k8s.io/client-go/informers GenericInformer
//go:generate mockgen -source=../../vendor/k8s.io/client-go/dynamic/dynamicinformer/interface.go -destination=./fake/mock_dynamicsharedinformerfactory.go -package=fake k8s.io/client-go/dynamic/dynamicinformer DynamicSharedInformerFactory
//go:generate mockgen -source=../../vendor/k8s.io/client-go/discovery/discovery_client.go -imports=openapi_v2=github.com/googleapis/gnostic/openapiv2 -destination=./fake/mock_discoveryinterface.go -package=fake k8s.io/client-go/discovery DiscoveryInterface
//go:generate mockgen -source=../../vendor/k8s.io/client-go/kubernetes/clientset.go -destination=./fake/mock_kubernetes_client.go -package=fake -mock_names=Interface=MockKubernetesInterface k8s.io/client-go/kubernetes Interface
//go:generate mockgen -destination=./fake/mock_sharedindexinformer.go -package=fake k8s.io/client-go/tools/cache SharedIndexInformer
//go:generate mockgen -destination=./fake/mock_authorization.go -package=fake k8s.io/client-go/kubernetes/typed/authorization/v1 AuthorizationV1Interface,SelfSubjectAccessReviewInterface,SelfSubjectAccessReviewsGetter,SelfSubjectRulesReviewInterface,SelfSubjectRulesReviewsGetter
//go:generate mockgen -source=../../vendor/k8s.io/client-go/dynamic/interface.go -destination=./fake/mock_dynamic_client.go -package=fake -imports=github.com/vmware-tanzu/octant/vendor/k8s.io/client-go/dynamic=k8s.io/client-go/dynamic -mock_names=Interface=MockDynamicInterface k8s.io/client-go/dynamic Interface

// Cluster is a client for cluster operations
type Cluster struct {
	clientConfig clientcmd.ClientConfig
	restConfig   *rest.Config
	logger       log.Logger

	kubernetesClient kubernetes.Interface
	dynamicClient    dynamic.Interface
	discoveryClient  discovery.DiscoveryInterface

	restMapper *restmapper.DeferredDiscoveryRESTMapper

	closeFn context.CancelFunc

	defaultNamespace   string
	providedNamespaces []string
}

var _ clusterTypes.ClientInterface = (*Cluster)(nil)

func newCluster(ctx context.Context, clientConfig clientcmd.ClientConfig, restClient *rest.Config, defaultNamespace string, providedNamespaces []string) (*Cluster, error) {
	logger := internalLog.From(ctx).With("component", "cluster client")

	install.Install(scheme.Scheme)
	_ = admissionregistrationv1.AddToScheme(scheme.Scheme)
	_ = apiregistrationv1.AddToScheme(scheme.Scheme)

	kubernetesClient, err := kubernetes.NewForConfig(restClient)
	if err != nil {
		return nil, errors.Wrap(err, "create kubernetes client")
	}

	dynamicClient, err := dynamic.NewForConfig(restClient)
	if err != nil {
		return nil, errors.Wrap(err, "create dynamic client")
	}

	dir, err := ioutil.TempDir("", "octant")
	if err != nil {
		return nil, errors.Wrap(err, "create temp directory")
	}

	logger.With("dir", dir).Debugf("created temp directory")

	cachedDiscoveryClient, err := disk.NewCachedDiscoveryClientForConfig(
		restClient,
		dir,
		dir,
		180*time.Second,
	)
	if err != nil {
		return nil, errors.Wrap(err, "create cached discovery client")
	}

	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)

	c := &Cluster{
		clientConfig:       clientConfig,
		restConfig:         restClient,
		kubernetesClient:   kubernetesClient,
		dynamicClient:      dynamicClient,
		discoveryClient:    cachedDiscoveryClient,
		restMapper:         restMapper,
		logger:             internalLog.From(ctx),
		defaultNamespace:   defaultNamespace,
		providedNamespaces: providedNamespaces,
	}

	ctx, cancel := context.WithCancel(ctx)
	c.closeFn = cancel

	go func() {
		<-ctx.Done()
		logger.Infof("removing cluster client temporary directory")

		if err := os.RemoveAll(dir); err != nil {
			logger.WithErr(err).Errorf("closing temporary directory")
		}
	}()

	return c, nil
}

func (c *Cluster) Close() {
	if c.closeFn != nil {
		c.closeFn()
	}
}

func (c *Cluster) DefaultNamespace() string {
	return c.defaultNamespace
}

func (c *Cluster) ResourceExists(gvr schema.GroupVersionResource) bool {
	restMapper := c.restMapper
	_, err := restMapper.KindFor(gvr)
	return err == nil
}

func (c *Cluster) Resource(gk schema.GroupKind) (schema.GroupVersionResource, bool, error) {
	restMapping, err := c.restMapper.RESTMapping(gk)
	if err != nil {
		return schema.GroupVersionResource{}, false, err
	}
	return restMapping.Resource, restMapping.Scope.Name() == meta.RESTScopeNameNamespace, nil
}

func (c *Cluster) ResetMapper() {
	c.restMapper.Reset()
}

// KubernetesClient returns a Kubernetes client.
func (c *Cluster) KubernetesClient() (kubernetes.Interface, error) {
	return c.kubernetesClient, nil
}

// NamespaceClient returns a namespace client.
func (c *Cluster) NamespaceClient() (clusterTypes.NamespaceInterface, error) {
	rc, err := c.RESTClient()
	if err != nil {
		return nil, err
	}

	dc, err := c.DynamicClient()
	if err != nil {
		return nil, err
	}

	ns, _, err := c.clientConfig.Namespace()
	if err != nil {
		return nil, errors.Wrap(err, "resolving initial namespace")
	}
	return newNamespaceClient(dc, rc, ns, c.providedNamespaces), nil
}

// DynamicClient returns a dynamic client.
func (c *Cluster) DynamicClient() (dynamic.Interface, error) {
	return c.dynamicClient, nil
}

// DiscoveryClient returns a DiscoveryClient for the cluster.
func (c *Cluster) DiscoveryClient() (discovery.DiscoveryInterface, error) {
	return c.discoveryClient, nil
}

// InfoClient returns an InfoClient for the cluster.
func (c *Cluster) InfoClient() (clusterTypes.InfoInterface, error) {
	return newClusterInfo(c.clientConfig), nil
}

// RESTClient returns a RESTClient for the cluster.
func (c *Cluster) RESTClient() (rest.Interface, error) {
	return rest.RESTClientFor(c.restConfig)
}

// RESTConfig returns configuration for communicating with the cluster.
func (c *Cluster) RESTConfig() *rest.Config {
	return c.restConfig
}

// Version returns a ServerVersion for the cluster.
func (c *Cluster) Version() (string, error) {
	dc, err := c.DiscoveryClient()
	if err != nil {
		return "", err
	}
	serverVersion, err := dc.ServerVersion()
	if err != nil {
		return "", err
	}
	return fmt.Sprint(serverVersion), nil
}

type clusterOptions struct {
	InitialNamespace   string
	ProvidedNamespaces []string
	RESTConfigOptions  clusterTypes.RESTConfigOptions
}

type ClusterOption func(*clusterOptions)

func WithInitialNamespace(initialNamespace string) ClusterOption {
	return func(clusterOptions *clusterOptions) {
		clusterOptions.InitialNamespace = initialNamespace
	}
}

func WithProvidedNamespaces(providedNamespaces []string) ClusterOption {
	return func(clusterOptions *clusterOptions) {
		clusterOptions.ProvidedNamespaces = providedNamespaces
	}
}

func WithClientQPS(qps float32) ClusterOption {
	return func(clusterOptions *clusterOptions) {
		clusterOptions.RESTConfigOptions.QPS = qps
	}
}

func WithClientBurst(burst int) ClusterOption {
	return func(clusterOptions *clusterOptions) {
		clusterOptions.RESTConfigOptions.Burst = burst
	}
}

func WithClientUserAgent(userAgent string) ClusterOption {
	return func(clusterOptions *clusterOptions) {
		clusterOptions.RESTConfigOptions.UserAgent = userAgent
	}
}

func WithRESTConfigOptions(restConfigOptions clusterTypes.RESTConfigOptions) ClusterOption {
	return func(clusterOptions *clusterOptions) {
		clusterOptions.RESTConfigOptions = restConfigOptions
	}
}

// FromClientConfig creates a Cluster from a k8s.io/client-go ClientConfig
func FromClientConfig(
	ctx context.Context,
	clientConfig clientcmd.ClientConfig,
	opts ...ClusterOption,
) (*Cluster, error) {
	options := clusterOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	var defaultNamespace string

	if options.InitialNamespace == "" {
		defaultNamespace, _, err = clientConfig.Namespace()
		if err != nil {
			return nil, err
		}
	} else {
		defaultNamespace = options.InitialNamespace
	}

	logger := internalLog.From(ctx)
	logger.With("client-qps", options.RESTConfigOptions.QPS, "client-burst", options.RESTConfigOptions.Burst).
		Debugf("initializing REST client configuration")

	restConfig = withConfigDefaults(restConfig, options.RESTConfigOptions)

	return newCluster(ctx, clientConfig, restConfig, defaultNamespace, options.ProvidedNamespaces)
}

// withConfigDefaults returns an extended rest.Config object with additional defaults applied
// See core_client.go#setConfigDefaults
func withConfigDefaults(inConfig *rest.Config, options clusterTypes.RESTConfigOptions) *rest.Config {
	config := rest.CopyConfig(inConfig)
	config.QPS = options.QPS
	config.Burst = options.Burst
	config.APIPath = "/api"
	if config.GroupVersion == nil || config.GroupVersion.Group != scheme.Scheme.PrioritizedVersionsForGroup("")[0].Group {
		gv := scheme.Scheme.PrioritizedVersionsForGroup("")[0]
		config.GroupVersion = &gv
	}
	codec := runtime.NoopEncoder{Decoder: scheme.Codecs.UniversalDecoder()}
	config.NegotiatedSerializer = serializer.NegotiatedSerializerWrapper(runtime.SerializerInfo{Serializer: codec})

	if options.UserAgent != "" {
		config.UserAgent = options.UserAgent
	}

	return config
}

/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package dash

import (
	"net"

	"k8s.io/client-go/dynamic/dynamicinformer"

	internalCluster "github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/kubeconfig"
	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/cluster"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/store"
)

type Options struct {
	BrowserPath            string
	BuildInfo              config.BuildInfo
	ClientBurst            int
	ClientQPS              float32
	Context                string
	DisableClusterOverview bool
	EnableMemStats         bool
	EnableOpenCensus       bool
	FrontendURL            string
	KubeConfig             string
	Listener               net.Listener
	Namespace              string
	Namespaces             []string
	UserAgent              string

	clusterClient          cluster.ClientInterface
	factory                dynamicinformer.DynamicSharedInformerFactory
	objectStore            store.Store
	streamingClientFactory api.StreamingClientFactory
}

type RunnerOption struct {
	kubeConfigOption kubeconfig.KubeConfigOption
	nonClusterOption func(*Options)
}

func WithBrowserPath(browserPath string) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.Noop(),
		nonClusterOption: func(o *Options) {
			o.BrowserPath = browserPath
		},
	}
}

func WithBuildInfo(buildInfo config.BuildInfo) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.Noop(),
		nonClusterOption: func(o *Options) {
			o.BuildInfo = buildInfo
		},
	}
}

func WithClientBurst(burst int) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.FromClusterOption(internalCluster.WithClientBurst(burst)),
		nonClusterOption: func(o *Options) {},
	}
}

func WithClientQPS(qps float32) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.FromClusterOption(internalCluster.WithClientQPS(qps)),
		nonClusterOption: func(o *Options) {},
	}
}

func WithClientUserAgent(userAgent string) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.FromClusterOption(internalCluster.WithClientUserAgent(userAgent)),
		nonClusterOption: func(o *Options) {},
	}
}

func WithClusterClient(client cluster.ClientInterface) RunnerOption {
	return RunnerOption{
		nonClusterOption: func(o *Options) {
			o.clusterClient = client
		},
	}
}

func WithoutClusterOverview() RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.Noop(),
		nonClusterOption: func(o *Options) {
			o.DisableClusterOverview = true
		},
	}
}

func WithContext(context string) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.WithContextName(context),
		nonClusterOption: func(o *Options) {
			o.Context = context
		},
	}
}

// WithDynamicSharedInformerFactory allows overriding the DynamicSharedInformerFactory that is used
// by the default ObjectStore implementation. This is used for embedded uses of Octant and testing. If the
// WithObjectStore option is provided, this option is ignored.
func WithDynamicSharedInformerFactory(factory dynamicinformer.DynamicSharedInformerFactory) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.Noop(),
		nonClusterOption: func(o *Options) {
			o.factory = factory
		},
	}
}

func WithFrontendURL(frontendURL string) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.Noop(),
		nonClusterOption: func(o *Options) {
			o.FrontendURL = frontendURL
		},
	}
}

func WithKubeConfig(kubeConfig string) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.WithKubeConfigList(kubeConfig),
		nonClusterOption: func(o *Options) {
			o.KubeConfig = kubeConfig
		},
	}
}

func WithListener(listener net.Listener) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.Noop(),
		nonClusterOption: func(o *Options) {
			o.Listener = listener
		},
	}
}

func WithMemStats() RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.Noop(),
		nonClusterOption: func(o *Options) {
			o.EnableMemStats = true
		},
	}
}

func WithNamespace(namespace string) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.FromClusterOption(internalCluster.WithInitialNamespace(namespace)),
		nonClusterOption: func(o *Options) {
			o.Namespace = namespace
		},
	}
}

func WithNamespaces(namespaces []string) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.FromClusterOption(internalCluster.WithProvidedNamespaces(namespaces)),
		nonClusterOption: func(o *Options) {
			o.Namespaces = namespaces
		},
	}
}

// WithObjectStore allows overriding the default ObjectStore implementation with a custom object store. This is for
// embedded uses of Octant and testing purposes. When using this option the WithDynamicSharedInformerFactory option
// is ignored.
func WithObjectStore(objectStore store.Store) RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.Noop(),
		nonClusterOption: func(o *Options) {
			o.objectStore = objectStore
		},
	}
}

func WithOpenCensus() RunnerOption {
	return RunnerOption{
		kubeConfigOption: kubeconfig.Noop(),
		nonClusterOption: func(o *Options) {
			o.EnableOpenCensus = true
		},
	}
}

func WithStreamingClientFactory(factory api.StreamingClientFactory) RunnerOption {
	return RunnerOption{
		nonClusterOption: func(o *Options) {
			o.streamingClientFactory = factory
		},
	}
}

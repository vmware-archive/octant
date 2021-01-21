/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package kubeconfig

import (
	"context"
	"path/filepath"
	"sort"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/util/strings"
)

type KubeConfigOption struct {
	kubeConfigOption func(*kubeConfigOptions)
	clusterOption    cluster.ClusterOption
}

type kubeConfigOptions struct {
	KubeConfigList string
	ContextName    string
}

func Noop() KubeConfigOption {
	return KubeConfigOption{
		kubeConfigOption: func(*kubeConfigOptions) {},
	}
}

func WithKubeConfigList(kubeConfigList string) KubeConfigOption {
	return KubeConfigOption{
		kubeConfigOption: func(kubeConfigOptions *kubeConfigOptions) {
			kubeConfigOptions.KubeConfigList = kubeConfigList
		},
	}
}

func WithContextName(contextName string) KubeConfigOption {
	return KubeConfigOption{
		kubeConfigOption: func(kubeConfigOptions *kubeConfigOptions) {
			kubeConfigOptions.ContextName = contextName
		},
	}
}

func FromClusterOption(clusterOption cluster.ClusterOption) KubeConfigOption {
	return KubeConfigOption{
		kubeConfigOption: func(*kubeConfigOptions) {},
		clusterOption:    clusterOption,
	}
}

func NewKubeConfigContextManager(ctx context.Context, opts ...KubeConfigOption) (*KubeConfigContextManager, error) {
	options := kubeConfigOptions{}
	clusterOptions := []cluster.ClusterOption{}
	for _, opt := range opts {
		opt.kubeConfigOption(&options)
		clusterOptions = append(clusterOptions, opt.clusterOption)
	}
	chain := strings.Deduplicate(filepath.SplitList(options.KubeConfigList))
	rules := &clientcmd.ClientConfigLoadingRules{
		Precedence: chain,
	}

	overrides := &clientcmd.ConfigOverrides{}
	if options.ContextName != "" {
		overrides.CurrentContext = options.ContextName
	}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
	clusterClient, err := cluster.FromClientConfig(ctx, clientConfig, clusterOptions...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create cluster client")
	}

	config, err := clientConfig.RawConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to load kube config")
	}
	var contextList []Context

	for name := range config.Contexts {
		contextList = append(contextList, Context{Name: name})
	}

	sort.Slice(contextList, func(i, j int) bool {
		return contextList[i].Name < contextList[j].Name
	})

	contextName := options.ContextName
	if contextName == "" {
		contextName = config.CurrentContext
	}

	return &KubeConfigContextManager{
		configLoadingRules: &clientcmd.ClientConfigLoadingRules{
			Precedence: chain,
		},
		currentContext: contextName,
		contexts:       contextList,
		clusterClient:  clusterClient,
		clusterOptions: clusterOptions,
	}, nil
}

type KubeConfigContextManager struct {
	configLoadingRules *clientcmd.ClientConfigLoadingRules
	currentContext     string
	contexts           []Context
	clusterClient      cluster.ClientInterface
	clusterOptions     []cluster.ClusterOption
}

// Context describes a kube config context.
type Context struct {
	Name string `json:"name"`
}

// UseFSContext is used to indicate a context switch to the file system Kubeconfig context
const UseFSContext = ""

func (k *KubeConfigContextManager) CurrentContext() string {
	return k.currentContext
}

func (k *KubeConfigContextManager) Contexts() []Context {
	return k.contexts
}

func (k *KubeConfigContextManager) SwitchContext(ctx context.Context, contextName string) error {
	if k.clusterClient != nil {
		k.clusterClient.Close()
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		k.configLoadingRules,
		&clientcmd.ConfigOverrides{CurrentContext: contextName},
	)

	var err error
	k.clusterClient, err = cluster.FromClientConfig(ctx, clientConfig, k.clusterOptions...)
	if err != nil {
		return errors.Wrap(err, "unable to create cluster client")
	}

	if contextName == UseFSContext {
		rawConfig, err := clientConfig.RawConfig()
		if err != nil {
			return errors.Wrap(err, "unable to infer context name from kube config")
		}
		contextName = rawConfig.CurrentContext
	}

	k.currentContext = contextName
	return nil
}

func (k *KubeConfigContextManager) ClusterClient() cluster.ClientInterface {
	return k.clusterClient
}

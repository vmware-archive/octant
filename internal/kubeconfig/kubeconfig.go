/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package kubeconfig

import (
	"context"
	"path/filepath"
	"sort"
	"sync/atomic"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/pkg/errors"

	internalCluster "github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/util/strings"
	"github.com/vmware-tanzu/octant/pkg/cluster"
)

type KubeConfigOption struct {
	kubeConfigOption func(*kubeConfigOptions)
	clusterOption    internalCluster.ClusterOption
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

func FromClusterOption(clusterOption internalCluster.ClusterOption) KubeConfigOption {
	return KubeConfigOption{
		kubeConfigOption: func(*kubeConfigOptions) {},
		clusterOption:    clusterOption,
	}
}

func NewKubeConfigContextManager(ctx context.Context, opts ...KubeConfigOption) (*KubeConfigContextManager, error) {
	options := kubeConfigOptions{}
	clusterOptions := []internalCluster.ClusterOption{}
	for _, opt := range opts {
		if opt.kubeConfigOption != nil {
			opt.kubeConfigOption(&options)
			clusterOptions = append(clusterOptions, opt.clusterOption)
		}
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
	clusterClient, err := internalCluster.FromClientConfig(ctx, clientConfig, clusterOptions...)
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

	kubeConfigCtxMgr := &KubeConfigContextManager{
		configLoadingRules: &clientcmd.ClientConfigLoadingRules{
			Precedence: chain,
		},
		currentContext: contextName,
		contexts:       contextList,
		clusterOptions: clusterOptions,
	}
	kubeConfigCtxMgr.clusterClient.Store(clusterClient)
	return kubeConfigCtxMgr, nil
}

type KubeConfigContextManager struct {
	configLoadingRules *clientcmd.ClientConfigLoadingRules
	currentContext     string
	contexts           []Context
	clusterClient      atomic.Value // cluster.ClientInterface
	clusterOptions     []internalCluster.ClusterOption
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
	v := k.clusterClient.Load()
	if v != nil {
		clusterClient := v.(cluster.ClientInterface)
		clusterClient.Close()
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		k.configLoadingRules,
		&clientcmd.ConfigOverrides{CurrentContext: contextName},
	)

	var err error
	clusterClient, err := internalCluster.FromClientConfig(ctx, clientConfig, k.clusterOptions...)
	if err != nil {
		return errors.Wrap(err, "unable to create cluster client")
	}
	k.clusterClient.Store(clusterClient)

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
	v := k.clusterClient.Load()
	if v == nil {
		return nil
	}
	return v.(cluster.ClientInterface)
}

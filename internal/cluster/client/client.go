/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/util/strings"
)

// ClusterClientManager is a manager for cluster clients
type ClusterClientManager interface {
	//ClusterClients() []cluster.ClientInterface
	//Contexts() map{string}string
	Get(ctx context.Context, contextName string) (cluster.ClientInterface, error)
}

type clusterClientManager struct {
	clusterClients map[string]cluster.ClientInterface
}

var _ ClusterClientManager = (*clusterClientManager)(nil)

// NewClusterClientManager creates an instance of clusterClientManager
func NewClusterClientManager(ctx context.Context, kubeConfig string, options cluster.RESTConfigOptions) (ClusterClientManager, error) {
	clusterClients := map[string]cluster.ClientInterface{}

	// TODO: refactor
	chain := strings.Deduplicate(filepath.SplitList(kubeConfig))
	rules := &clientcmd.ClientConfigLoadingRules{
		Precedence: chain,
	}
	overrides := &clientcmd.ConfigOverrides{}
	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
	rawConfig, err := cc.RawConfig()
	if err != nil {
		return nil, errors.Wrap(err, "clusterClientManager rawConfig")
	}
	// TODO: refactor

	for contextName := range rawConfig.Contexts {
		client, err := cluster.FromKubeConfig(ctx, kubeConfig, contextName, options)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create client from kubeconfig")
		}
		clusterClients[contextName] = client
	}
	ccm := &clusterClientManager{
		clusterClients: clusterClients,
	}
	return ccm, nil
}

// Get gets a cluster client
func (c *clusterClientManager) Get(ctx context.Context, contextName string) (cluster.ClientInterface, error) {
	clusterClient, ok := c.clusterClients[contextName]
	if !ok {
		return nil, errors.New("can't get cluster client")
	}
	return clusterClient, nil
}

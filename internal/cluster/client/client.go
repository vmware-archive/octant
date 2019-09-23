/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/log"
)

//go:generate mockgen -source=client.go -destination=./fake/mock_cluster_client_manager.go -package=fake github.com/vmware/octant/internal/cluster/client ClusterClientManager

// ClusterClientManager is a manager for cluster clients
type ClusterClientManager interface {
	// ClusterClients() []cluster.ClientInterface
	// Contexts() map{string}string
	SetDefault(ctx context.Context, contextName string)
	Get(ctx context.Context, contextName string) (cluster.ClientInterface, error)
}

type clusterClientManager struct {
	clusterClients map[string]cluster.ClientInterface
	defaultContext string
	contexts       []string
}

var _ ClusterClientManager = (*clusterClientManager)(nil)

// NewClusterClientManager creates an instance of clusterClientManager
func NewClusterClientManager(ctx context.Context, kubeConfig string, options cluster.RESTConfigOptions) (ClusterClientManager, error) {
	clusterClients := map[string]cluster.ClientInterface{}

	cc := cluster.ClientConfigFromKubeConfig(kubeConfig, nil)
	rawConfig, err := cc.RawConfig()
	if err != nil {
		return nil, errors.Wrap(err, "clusterClientManager rawConfig")
	}

	logger := log.From(ctx)
	contexts := []string{}
	for contextName := range rawConfig.Contexts {
		contexts = append(contexts, contextName)
		client, err := cluster.FromKubeConfig(ctx, kubeConfig, contextName, options)
		if err != nil {
			logger.WithErr(err).Warnf(fmt.Sprintf("unablet to create client for %s from kubeconfig", contextName))
			continue
		}
		clusterClients[contextName] = client
	}
	ccm := &clusterClientManager{
		clusterClients: clusterClients,
		defaultContext: rawConfig.CurrentContext,
		contexts:       contexts,
	}
	return ccm, nil
}

// SetDefault
func (c *clusterClientManager) SetDefault(ctx context.Context, contextName string) {
	c.defaultContext = contextName
}

// Get gets a cluster client
func (c *clusterClientManager) Get(ctx context.Context, contextName string) (cluster.ClientInterface, error) {
	if contextName == "" {
		contextName = c.defaultContext
	}
	clusterClient, ok := c.clusterClients[contextName]
	if !ok {
		return nil, errors.New("can't get cluster client")
	}
	return clusterClient, nil
}

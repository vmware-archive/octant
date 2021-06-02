/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/cluster"
)

//go:generate mockgen -destination=./fake/mock_pod_metrics_loader.go -package=fake github.com/vmware-tanzu/octant/internal/octant PodMetricsLoader
//go:generate mockgen -destination=./fake/mock_pod_metrics_crud.go -package=fake github.com/vmware-tanzu/octant/internal/octant PodMetricsCRUD

// PodMetricsCRUD contains CRUD methods for accessing pod metrics.
type PodMetricsCRUD interface {
	// Get returns pod metrics for a pod. If pod is not found, isFound will be false.
	Get(ctx context.Context, namespace, name string) (pod *unstructured.Unstructured, isFound bool, err error)
}

type clusterPodMetricsCRUD struct {
	clusterClient cluster.ClientInterface
}

var _ PodMetricsCRUD = (*clusterPodMetricsCRUD)(nil)

func newClusterPodMetricsCRUD(clusterClient cluster.ClientInterface) (*clusterPodMetricsCRUD, error) {
	if clusterClient == nil {
		return nil, fmt.Errorf("cluster client is nil")
	}

	return &clusterPodMetricsCRUD{clusterClient: clusterClient}, nil
}

func (c *clusterPodMetricsCRUD) Get(ctx context.Context, namespace, name string) (*unstructured.Unstructured, bool, error) {
	client, err := c.clusterClient.DynamicClient()
	if err != nil {
		return nil, false, fmt.Errorf("get dynamic client: %w", err)
	}

	options := metav1.GetOptions{}
	object, err := client.Resource(PodMetricsResource).Namespace(namespace).Get(context.TODO(), name, options)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, false, nil
		}
		if errors.IsServiceUnavailable(err) {
			logger := log.From(ctx)
			logger.Warnf("service unavailable : %w", err)
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("get pod metrics: %w", err)
	}

	return object, true, nil
}

type noPodMetricsSupport interface {
	NoPodMetricsSupport() bool
}

// NoPodMetricsErr is an error signifying a cluster does not have pod metrics support.
type NoPodMetricsErr struct{}

var _ noPodMetricsSupport = (*NoPodMetricsErr)(nil)
var _ error = (*NoPodMetricsErr)(nil)

func (e *NoPodMetricsErr) NoPodMetricsSupport() bool {
	return true
}

func (e *NoPodMetricsErr) Error() string {
	return "pod metrics are not supported by this cluster"
}

// IsPodMetricsNotSupported returns true if error is pod metrics not supported.
func IsPodMetricsNotSupported(err error) bool {
	e, ok := err.(noPodMetricsSupport)
	return ok && e.NoPodMetricsSupport()
}

// PodMetricsLoader loads metrics for a pod.
type PodMetricsLoader interface {
	// Load loads metrics for a pod given namespace and a name. It returns false if the
	// object is not found.
	Load(ctx context.Context, namespace, name string) (object *unstructured.Unstructured, isFound bool, err error)
	// SupportsMetrics returns true if the cluster has metrics support.
	SupportsMetrics(ctx context.Context) (bool, error)
}

// ClusterPodMetricsLoaderOption is an option for configuring ClusterPodMetricsLoader.
type ClusterPodMetricsLoaderOption func(loader *ClusterPodMetricsLoader)

// ClusterPodMetricsLoader loads metrics for a pod using a cluster client.
type ClusterPodMetricsLoader struct {
	PodMetricsCRUD PodMetricsCRUD

	clusterClient cluster.ClientInterface
	supportsOnce  sync.Once

	hasPodMetricsSupport bool
}

var _ PodMetricsLoader = (*ClusterPodMetricsLoader)(nil)

// NewClusterPodMetricsLoader creates an instance of ClusterPodMetricsLoader.
func NewClusterPodMetricsLoader(clusterClient cluster.ClientInterface, options ...ClusterPodMetricsLoaderOption) (*ClusterPodMetricsLoader, error) {
	if clusterClient == nil {
		return nil, fmt.Errorf("cluster client is nil")
	}

	pml := &ClusterPodMetricsLoader{
		clusterClient: clusterClient,
		supportsOnce:  sync.Once{},
	}

	for _, option := range options {
		option(pml)
	}

	if pml.PodMetricsCRUD == nil {
		pmc, err := newClusterPodMetricsCRUD(clusterClient)
		if err != nil {
			return nil, fmt.Errorf("create pod metrics CRUD client: %w", err)
		}
		pml.PodMetricsCRUD = pmc
	}

	return pml, nil
}

var (
	// PodMetricsResource is resource for pod metrics.
	PodMetricsResource = schema.GroupVersionResource{Group: "metrics.k8s.io", Version: "v1beta1", Resource: "pods"}
)

// Load loads metrics for a pod given namespace and a name.
func (ml *ClusterPodMetricsLoader) Load(ctx context.Context, namespace, name string) (*unstructured.Unstructured, bool, error) {
	return ml.PodMetricsCRUD.Get(ctx, namespace, name)
}

func (ml *ClusterPodMetricsLoader) SupportsMetrics(ctx context.Context) (bool, error) {
	var sErr error
	ml.supportsOnce.Do(func() {
		discoveryClient, err := ml.clusterClient.DiscoveryClient()
		if err != nil {
			sErr = fmt.Errorf("get discovery cluster: %w", err)
			return
		}

		lists, err := discoveryClient.ServerPreferredNamespacedResources()
		if err != nil {
			if discovery.IsGroupDiscoveryFailedError(err) && strings.Contains(err.Error(), "metrics") {
				logger := log.From(ctx)
				logger.Debugf("metrics discovery failed: %w", err)
				ml.hasPodMetricsSupport = false
				return
			}
			sErr = fmt.Errorf("get preferred namespaced resources: %w", err)
			return
		}

		podMetricString := fmt.Sprintf("%s %s", gvk.PodMetrics.GroupVersion().String(), gvk.PodMetrics.Kind)

		for _, list := range lists {
			for i := range list.APIResources {
				s := fmt.Sprintf("%s %s", list.GroupVersion, list.APIResources[i].Kind)

				if s == podMetricString {
					ml.hasPodMetricsSupport = true
				}
			}
		}
	})

	if sErr != nil {
		return false, sErr
	}

	return ml.hasPodMetricsSupport, nil
}

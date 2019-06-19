/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package resourceviewer

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/modules/overview/objectstatus"
	"github.com/vmware/octant/internal/modules/overview/objectvisitor"
	"github.com/vmware/octant/pkg/store"
	dashStrings "github.com/vmware/octant/internal/util/strings"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/view/component"
)

const defaultPodGroupPrefix = "pods-"

// CollectorOption is an option for configuring Collector.
type CollectorOption func(c *Collector)

// Collector collects objects to construct a resource viewer.
type Collector struct {
	edges  map[string][]string
	nodes  map[string]component.Node
	logger log.Logger

	podGroupPrefix string

	// groupPods sets the pod grouping. If it is true, group pods in one
	// graph node. If not, show them separately.
	groupPods bool

	// podGroupIDs maps a pod to a pod group
	podGroupIDs map[string]string

	// podStats counts pods in a replica set.
	podStats map[string]int

	podNodes map[string]component.PodStatus

	objectStore   store.Store
	link          link.Interface
	pluginPrinter plugin.ManagerInterface

	mu sync.Mutex
}

var _ objectvisitor.ObjectHandler = (*Collector)(nil)

// NewCollector creates an instance of Collector.
func NewCollector(dashConfig config.Dash, options ...CollectorOption) (*Collector, error) {
	l, err := link.NewFromDashConfig(dashConfig)
	if err != nil {
		return nil, err
	}

	collector := &Collector{
		groupPods:      true,
		podGroupPrefix: defaultPodGroupPrefix,
		objectStore:    dashConfig.ObjectStore(),
		link:           l,
		pluginPrinter:  dashConfig.PluginManager(),
	}

	for _, option := range options {
		option(collector)
	}

	collector.Reset()

	return collector, nil
}

// Reset resets a Collector's nodes and edges.
func (c *Collector) Reset() {
	c.edges = make(map[string][]string)
	c.nodes = make(map[string]component.Node)
	c.podNodes = make(map[string]component.PodStatus)
	c.podGroupIDs = make(map[string]string)
	c.podStats = make(map[string]int)
}

// Process process an object by saving the object to a map.
func (c *Collector) Process(ctx context.Context, object objectvisitor.ClusterObject) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var nodeID string
	var node component.Node
	var err error

	if c.isObjectPod(object) && c.groupPods {
		pod := &corev1.Pod{}
		if err := scheme.Scheme.Convert(object, pod, 0); err != nil {
			return errors.Wrap(err, "unable to convert object to pod")
		}

		if ownerReference := metav1.GetControllerOf(pod); ownerReference != nil {
			id := fmt.Sprintf("%s-%s-%s", ownerReference.APIVersion, ownerReference.Kind, ownerReference.Name)
			c.podStats[id]++

		}

		nodeID, node, err = c.createPodGroupNode(ctx, object)
	} else {
		nodeID, node, err = c.createObjectNode(ctx, object)
	}

	if err != nil {
		if isSkippedNode(err) {
			return nil
		}

		gvk := object.GetObjectKind().GroupVersionKind()
		accessor := meta.NewAccessor()
		name, err := accessor.Name(object)
		if err == nil {
			return errors.Wrapf(err, "processing unknown %s", gvk.String())
		}

		return errors.Wrapf(err, "processing %s %s", gvk.String(), name)
	}

	if _, ok := c.nodes[nodeID]; !ok {
		c.nodes[nodeID] = node
	}

	return nil
}

func (c *Collector) createPodGroupNode(ctx context.Context, object objectvisitor.ClusterObject) (string, component.Node, error) {
	pgd, err := c.podGroupDetails(object)
	if err != nil {
		return "", component.Node{}, errors.Wrap(err, "getting pod group id for pod")
	}

	accessor := meta.NewAccessor()

	name, err := accessor.Name(object)
	if err != nil {
		return "", component.Node{}, errors.Wrap(err, "getting name for pod")
	}

	status, err := objectstatus.Status(ctx, object, c.objectStore)
	if err != nil {
		return "", component.Node{}, errors.Wrap(err, "getting status for pod")
	}

	objectKind := object.GetObjectKind()
	apiVersion, kind := objectKind.GroupVersionKind().ToAPIVersionAndKind()

	podStatus, ok := c.podNodes[pgd.id]
	if !ok {
		podStatus = *component.NewPodStatus()
		c.podNodes[pgd.id] = podStatus
	}

	podStatus.AddSummary(name, status.Details, status.Status())

	node := component.Node{
		Name:       pgd.name,
		APIVersion: apiVersion,
		Kind:       kind,
		Status:     podStatus.Status(),
		Details: []component.Component{
			&podStatus,
		},
	}

	node, err = pluginStatus(object, node, c.pluginPrinter)
	if err != nil {
		return "", component.Node{}, err
	}

	//c.podGroupIDs[string(uid)] = pgd.id
	c.podGroupIDs[name] = pgd.id

	return pgd.id, node, nil
}

func pluginStatus(object objectvisitor.ClusterObject, node component.Node, pluginPrinter plugin.ManagerInterface) (component.Node, error) {
	osr, err := pluginPrinter.ObjectStatus(object.(runtime.Object))
	if err != nil {
		return component.Node{}, errors.Wrap(err, "plugin object status error")
	}
	if osr.ObjectStatus.Status != "" {
		node.Status = osr.ObjectStatus.Status
	}
	for _, detail := range osr.ObjectStatus.Details {
		node.Details = append(node.Details, detail)
	}
	return node, nil
}

type isSkipped interface {
	IsSkipped() bool
}

func isSkippedNode(err error) bool {
	sn, ok := err.(isSkipped)
	return ok && sn.IsSkipped()
}

type skipNode struct{}

func (e skipNode) IsSkipped() bool {
	return true
}

func (e skipNode) Error() string {
	return "skip node"
}

func (c *Collector) createObjectNode(ctx context.Context, object objectvisitor.ClusterObject) (string, component.Node, error) {
	objectKind := object.GetObjectKind()
	gvk := objectKind.GroupVersionKind()
	apiVersion, kind := gvk.ToAPIVersionAndKind()

	accessor := meta.NewAccessor()

	if (gvk.Group == "apps" || gvk.Group == "extensions") &&
		gvk.Kind == "ReplicaSet" {
		apiVersion = "extensions/v1beta1"
		replicaSet := &appsv1.ReplicaSet{}
		if err := scheme.Scheme.Convert(object, replicaSet, nil); err != nil {
			return "", component.Node{}, errors.Wrap(err, "convert object to Replica Set")
		}

		replicas := replicaSet.Spec.Replicas
		if replicas == nil || *replicas < 1 {
			return "", component.Node{}, &skipNode{}
		}
	}

	name, err := accessor.Name(object)
	if err != nil {
		return "", component.Node{}, errors.New("unable to get object name")
	}

	var nodeStatus component.NodeStatus

	status, err := objectstatus.Status(ctx, object, c.objectStore)
	if err != nil {
		c.log().Errorf("error retrieving object status: %v", err)
		nodeStatus = component.NodeStatusOK
	} else {
		nodeStatus = status.Status()
	}

	q := url.Values{}
	objectPath, err := c.link.ForObjectWithQuery(object, name, q)
	if err != nil {
		return "", component.Node{}, err
	}

	node := component.Node{
		Name:       name,
		APIVersion: apiVersion,
		Kind:       kind,
		Status:     nodeStatus,
		Details:    status.Details,
		Path:       objectPath,
	}

	nodeID, err := genNodeID(object)
	if err != nil {
		return "", component.Node{}, errors.New("unable to get object name")
	}

	node, err = pluginStatus(object, node, c.pluginPrinter)
	if err != nil {
		return "", component.Node{}, err
	}

	return string(nodeID), node, nil
}

// AddChild adds children for an object to create edges. Pods are collated to a single object.
func (c *Collector) AddChild(parent objectvisitor.ClusterObject, children ...objectvisitor.ClusterObject) error {
	if c.isObjectPod(parent) {
		// reverse the relationship, so the pod group details don't need to be accounted for.
		for _, child := range children {
			if err := c.AddChild(child, parent); err != nil {
				return err
			}
		}

		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	pid, err := genNodeID(parent)
	if err != nil {
		return err
	}

	for _, child := range children {
		var cid string

		if c.isObjectPod(child) && c.groupPods {
			pgd, err := c.podGroupDetails(child)
			if err != nil {
				return errors.Wrap(err, "find pod group id for pod")
			}

			cid = pgd.id
		} else {
			cid, err = genNodeID(child)
			if err != nil {
				return err
			}
		}

		if !dashStrings.Contains(cid, c.edges[pid]) {
			c.edges[pid] = append(c.edges[pid], cid)
		}
	}

	return nil
}

func (c *Collector) isObjectPod(object objectvisitor.ClusterObject) bool {
	objectKind := object.GetObjectKind()
	objectGVK := objectKind.GroupVersionKind()

	return objectGVK.String() == gvk.PodGVK.String()
}

type podGroupDetails struct {
	id   string
	name string
}

func (c *Collector) podGroupDetails(object objectvisitor.ClusterObject) (podGroupDetails, error) {
	if !c.isObjectPod(object) {
		return podGroupDetails{}, errors.Errorf("can't create pod group details for a %T", object)
	}
	obj, err := meta.Accessor(object)
	if err != nil {
		return podGroupDetails{}, err
	}

	reference := metav1.GetControllerOf(obj)
	if reference == nil {
		fmt.Println("creating pod group details for pod without parent")
		return podGroupDetails{
			id:   string(obj.GetName()),
			name: obj.GetName(),
		}, nil
	}

	id := fmt.Sprintf("%s%s-%s-%s", c.podGroupPrefix, reference.APIVersion, reference.Kind, reference.Name)

	pgd := podGroupDetails{
		id:   id,
		name: fmt.Sprintf("%s pods", reference.Name),
	}

	return pgd, nil
}

func (c *Collector) Component(selected string) (component.Component, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	nodes := make(map[string]component.Node)
	for k, v := range c.nodes {
		nodes[k] = v
	}

	rv := component.NewResourceViewer("Resource Viewer")

	var nodeIDs []string
	for nodeID, node := range nodes {
		if strings.HasPrefix(nodeID, c.podGroupPrefix) {
			ownerID := strings.TrimPrefix(nodeID, c.podGroupPrefix)
			node.Details = append(node.Details,
				component.NewText(fmt.Sprintf("Pod count: %d", c.podStats[ownerID])))
			nodes[nodeID] = node
		}

		rv.AddNode(nodeID, node)
		nodeIDs = append(nodeIDs, nodeID)
	}

	for nodeID, edges := range c.edges {
		sort.Strings(edges)
		for _, edgeID := range edges {
			if dashStrings.Contains(edgeID, nodeIDs) {
				if err := rv.AddEdge(nodeID, edgeID, component.EdgeTypeExplicit); err != nil {
					c.log().WithErr(err).Errorf("unable to add edge to object graph")
				}
			}
		}
	}

	podGroupID, ok := c.podGroupIDs[selected]
	if ok {
		selected = podGroupID
	}

	rv.Select(selected)

	return rv, nil
}

func (c *Collector) log() log.Logger {
	if c.logger != nil {
		return c.logger
	}

	return log.NopLogger()
}

func genNodeID(object runtime.Object) (string, error) {
	if object == nil {
		return "", errors.New("can't generate node id for nil object")
	}

	accessor := meta.NewAccessor()

	name, err := accessor.Name(object)
	if err != nil {
		return "", err
	}

	apiVersion, kind := object.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()

	return fmt.Sprintf("%s-%s-%s", apiVersion, kind, name), nil
}

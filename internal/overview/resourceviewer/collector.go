package resourceviewer

import (
	"fmt"
	"strings"

	"github.com/heptio/developer-dash/internal/overview/link"

	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/overview/objectstatus"
	"github.com/heptio/developer-dash/internal/overview/objectvisitor"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

// Collector collects objects to construct a resource viewer.
type Collector struct {
	edges  map[string][]string
	nodes  map[string]component.Node
	logger log.Logger

	// podStats counts pods in a replica set.
	podStats map[string]int
}

var _ objectvisitor.ObjectHandler = (*Collector)(nil)

// NewCollector creates an instance of Collector.
func NewCollector() *Collector {
	c := &Collector{
		podStats: make(map[string]int),
	}
	c.Reset()

	return c
}

// Reset resets a Collector's nodes and edges.
func (c *Collector) Reset() {
	c.edges = make(map[string][]string)
	c.nodes = make(map[string]component.Node)
}

// Process process an object by saving the object to a map.
func (c *Collector) Process(object objectvisitor.ClusterObject) error {
	var uid string
	var node component.Node
	var err error

	if c.isPod(object) {
		pod := &corev1.Pod{}
		if err := scheme.Scheme.Convert(object, pod, 0); err != nil {
			return errors.Wrap(err, "unable to convert object to pod")
		}

		if ownerReference := metav1.GetControllerOf(pod); ownerReference != nil {
			c.podStats[string(ownerReference.UID)]++
		}

		uid, node, err = c.createPodGroupNode(object)
	} else {
		uid, node, err = c.createObjectNode(object)
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

	if _, ok := c.nodes[uid]; !ok {
		c.nodes[uid] = node
	}

	return nil
}

func (c *Collector) createPodGroupNode(object objectvisitor.ClusterObject) (string, component.Node, error) {
	pgd, err := c.podGroupDetails(object)
	if err != nil {
		return "", component.Node{}, errors.Wrap(err, "getting pod group id for pod")
	}

	objectKind := object.GetObjectKind()
	apiVersion, kind := objectKind.GroupVersionKind().ToAPIVersionAndKind()

	node := component.Node{
		Name:       pgd.name,
		APIVersion: apiVersion,
		Kind:       kind,
		// set all statuses to ok until we start checking for the real status.
		Status: component.NodeStatusOK,
	}

	return pgd.id, node, nil
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

func (c *Collector) createObjectNode(object objectvisitor.ClusterObject) (string, component.Node, error) {
	objectKind := object.GetObjectKind()
	gvk := objectKind.GroupVersionKind()
	apiVersion, kind := gvk.ToAPIVersionAndKind()

	accessor := meta.NewAccessor()

	if (gvk.Group == "apps" || gvk.Group == "extensions") &&
		gvk.Kind == "ReplicaSet" {
		replicaSet := &appsv1.ReplicaSet{}
		if err := scheme.Scheme.Convert(object, replicaSet, nil); err != nil {
			return "", component.Node{}, errors.Wrap(err, "convert object to Replica Set")
		}

		replicas := replicaSet.Spec.Replicas
		if replicas == nil || *replicas < 1 {
			return "", component.Node{}, &skipNode{}
		}
	}

	uid, err := accessor.UID(object)
	if err != nil {
		return "", component.Node{}, err
	}

	name, err := accessor.Name(object)
	if err != nil {
		return "", component.Node{}, errors.New("unable to get object name")
	}

	objectStatus, err := objectstatus.Status(object)
	if err != nil {
		c.log().Errorf("error retrieving object status: %v", err)
		objectStatus.NodeStatus = component.NodeStatusOK
	}

	node := component.Node{
		Name:       name,
		APIVersion: apiVersion,
		Kind:       kind,
		Status:     objectStatus.NodeStatus,
		Details:    objectStatus.Details,
		Path:       link.ForObject(object, name),
	}

	return string(uid), node, nil
}

// AddChild adds children for an object to create edges. Pods are collated to a single object.
func (c *Collector) AddChild(parent objectvisitor.ClusterObject, children ...objectvisitor.ClusterObject) error {
	accessor := meta.NewAccessor()
	pid, err := accessor.UID(parent)
	if err != nil {
		return err
	}

	for _, child := range children {
		var cid string

		if c.isPod(child) {
			pgd, err := c.podGroupDetails(child)
			if err != nil {
				return errors.Wrap(err, "find pod group id for pod")
			}

			cid = pgd.id
		} else {
			id, err := accessor.UID(child)
			if err != nil {
				return err
			}

			cid = string(id)
		}

		if !listContains(c.edges[string(pid)], cid) {
			c.edges[string(pid)] = append(c.edges[string(pid)], cid)
		}
	}

	return nil
}

func (c *Collector) isPod(object objectvisitor.ClusterObject) bool {
	objectKind := object.GetObjectKind()
	gvk := objectKind.GroupVersionKind()

	return gvk.Group == "" &&
		gvk.Version == "v1" &&
		gvk.Kind == "Pod"
}

type podGroupDetails struct {
	id   string
	name string
}

func (c *Collector) podGroupDetails(object objectvisitor.ClusterObject) (podGroupDetails, error) {
	obj, err := meta.Accessor(object)
	if err != nil {
		return podGroupDetails{}, err
	}

	reference := metav1.GetControllerOf(obj)
	if reference == nil {

		return podGroupDetails{
			id:   string(obj.GetUID()),
			name: obj.GetName(),
		}, nil
	}

	id := fmt.Sprintf("pods-%s", reference.UID)

	pgd := podGroupDetails{
		id:   id,
		name: fmt.Sprintf("%s pods", reference.Name),
	}

	return pgd, nil
}

func (c *Collector) ViewComponent() (component.ViewComponent, error) {
	rv := component.NewResourceViewer("Resource Viewer")

	var nodeIDs []string
	for nodeID, node := range c.nodes {
		if strings.HasPrefix(nodeID, "pods-") {
			ownerID := strings.TrimPrefix(nodeID, "pods-")
			node.Details = append(node.Details,
				component.NewText(fmt.Sprintf("Pod count: %d", c.podStats[ownerID])))
			c.nodes[nodeID] = node
		}

		rv.AddNode(nodeID, node)
		nodeIDs = append(nodeIDs, nodeID)
	}

	for nodeID, edges := range c.edges {
		for _, edgeID := range edges {
			if listContains(nodeIDs, edgeID) {
				rv.AddEdge(nodeID, edgeID, component.EdgeTypeExplicit)
			}
		}
	}

	return rv, nil
}

func (c *Collector) log() log.Logger {
	if c.logger != nil {
		return c.logger
	}

	return log.NopLogger()
}

func listContains(lst []string, s string) bool {
	for _, item := range lst {
		if item == s {
			return true
		}
	}
	return false
}

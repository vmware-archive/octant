package resourceviewer

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/internal/modules/overview/objectstatus"
	"github.com/vmware/octant/internal/modules/overview/objectvisitor"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

//go:generate mockgen -destination=./fake/mock_object_status.go -package=fake github.com/vmware/octant/internal/modules/overview/resourceviewer ObjectStatus

// HandlerOption is an option for configuring Handler.
type HandlerOption func(h *Handler)

// SetHandlerObjectStatus configures handler to use a custom object status generator.
func SetHandlerObjectStatus(objectStatus ObjectStatus) HandlerOption {
	return func(h *Handler) {
		h.objectStatus = objectStatus
	}
}

// Handler is a visitor handler.
type Handler struct {
	objectStore   store.Store
	link          link.Interface
	pluginPrinter plugin.ManagerInterface

	adjList map[types.UID][]runtime.Object
	nodes   map[types.UID]runtime.Object

	mu           sync.Mutex
	objectStatus ObjectStatus
}

var _ objectvisitor.ObjectHandler = (*Handler)(nil)

// NewHandler creates an instance of Handler.
func NewHandler(dashConfig config.Dash, options ...HandlerOption) (*Handler, error) {
	l, err := link.NewFromDashConfig(dashConfig)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		objectStore:   dashConfig.ObjectStore(),
		link:          l,
		pluginPrinter: dashConfig.PluginManager(),
		adjList:       make(map[types.UID][]runtime.Object),
		nodes:         make(map[types.UID]runtime.Object),
		objectStatus:  NewHandlerObjectStatus(dashConfig.ObjectStore(), dashConfig.PluginManager()),
	}

	for _, option := range options {
		option(h)
	}

	return h, nil
}

// AddEdge adds edges to the graph.
func (h *Handler) AddEdge(v1, v2 runtime.Object) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	v1Accessor, err := meta.Accessor(v1)
	if err != nil {
		return err
	}
	v1UID := v1Accessor.GetUID()

	h.nodes[v1UID] = v1

	cur := h.adjList[v1UID]
	cur = append(cur, v2)
	h.adjList[v1UID] = cur

	return nil
}

// Process adds nodes to the dependency graph.
func (h *Handler) Process(ctx context.Context, object runtime.Object) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	accessor, err := meta.Accessor(object)
	if err != nil {
		return err
	}

	uid := accessor.GetUID()
	h.nodes[uid] = object

	return nil
}

// AdjacencyList creates an adjacency list from the handler.
func (h *Handler) AdjacencyList() (*component.AdjList, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	adjList := component.AdjList{}

	podGroupEdges := make(map[string]map[string]component.Edge)

	// iterate over our internal list to build an exportable version.
	for parentUID, connections := range h.adjList {
		parent := h.nodes[parentUID]

		inGroup, err := isPodInGroup(parent)
		if err != nil {
			return nil, err
		}

		// if this node is a pod group, track the edges, so they can be added at the
		// the end.
		if inGroup {
			name, err := podGroupName(parent)
			if err != nil {
				return nil, err
			}

			for i := range connections {
				child := connections[i]

				edges, ok := podGroupEdges[name]
				if !ok {
					edges = make(map[string]component.Edge)
				}

				childAccessor, err := meta.Accessor(child)
				if err != nil {
					return nil, err
				}

				isParent, err := isObjectParent(parent, child)
				if err != nil {
					return nil, err
				}

				if isParent {
					continue
				}

				id := string(childAccessor.GetUID())
				edges[id] = component.Edge{
					Node: id,
					Type: component.EdgeTypeExplicit,
				}

				podGroupEdges[name] = edges
			}

			continue
		}

		edgeMap := make(map[string]component.Edge)

		for i := range connections {
			child := connections[i]

			inGroup, err := isPodInGroup(parent)
			if err != nil {
				return nil, err
			}

			if inGroup {
				name, err := podGroupName(parent)
				if err != nil {
					return nil, err
				}

				edge := component.Edge{
					Node: name,
					Type: component.EdgeTypeExplicit,
				}

				edgeMap[name] = edge
				continue
			}

			isParent, err := isObjectParent(parent, child)
			if err != nil {
				return nil, err
			}

			if isParent {
				continue
			}

			name, err := edgeName(child)
			if err != nil {
				if isSkippedNode(err) {
					continue
				}
				return nil, err
			}

			edge := component.Edge{
				Node: name,
				Type: component.EdgeTypeExplicit,
			}

			key := string(parentUID)
			adjList[key] = append(adjList[key], edge)
		}
	}

	for key, m := range podGroupEdges {
		for _, edge := range m {
			adjList[key] = append(adjList[key], edge)
		}
	}

	adjList = deDupEdges(adjList)

	return &adjList, nil
}

// Nodes generates nodes from the handler.
func (h *Handler) Nodes(ctx context.Context) (component.Nodes, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	nodes := component.Nodes{}

	var podsInAGroup []runtime.Object

	for uid, node := range h.nodes {

		ok, err := isPodInGroup(node)
		if err != nil {
			return nil, err
		}

		if ok {
			podsInAGroup = append(podsInAGroup, node)
			continue
		}

		onc := objectNode{
			link:          h.link,
			pluginPrinter: h.pluginPrinter,
			objectStatus:  h.objectStatus,
		}

		componentNode, err := onc.Create(ctx, node)
		if err != nil {
			if isSkippedNode(err) {
				continue
			}
			return nil, err
		}

		nodes[string(uid)] = *componentNode
	}

	nameMap, e := h.buildPodGroups(podsInAGroup)
	if e != nil {
		return nil, e
	}

	for podGroupName, objects := range nameMap {
		pgn := podGroupNode{
			objectStatus: h.objectStatus,
		}
		group, err := pgn.Create(ctx, podGroupName, objects)
		if err != nil {
			return nil, err
		}
		nodes[podGroupName] = *group
	}

	return nodes, nil
}

func (h *Handler) buildPodGroups(podsInAGroup []runtime.Object) (map[string][]runtime.Object, error) {
	nameMap := make(map[string][]runtime.Object)
	for _, object := range podsInAGroup {
		name, err := podGroupName(object)
		if err != nil {
			return nil, err
		}

		nameMap[name] = append(nameMap[name], object)
	}
	return nameMap, nil
}

func edgeName(object runtime.Object) (string, error) {
	if object == nil {
		return "", errors.New("can't build edge name for nil object")
	}

	ok, err := isPodInGroup(object)
	if err != nil {
		return "", err
	}

	accessor, err := meta.Accessor(object)
	if err != nil {
		return "", err
	}

	if ok {
		// If pod has owner references, associate this pod with a grouping. The name will be
		// constructed from the pod's labels.
		return podGroupName(object)
	}

	isReplicaSet, err := isObjectReplicaSet(object)
	if err != nil {
		return "", err
	}
	if isReplicaSet {
		if err := checkReplicaCount(object); err != nil {
			return "", err
		}
	}

	return string(accessor.GetUID()), nil
}

func isPodInGroup(object runtime.Object) (bool, error) {
	if !isObjectPod(object) {
		return false, nil
	}

	pod, err := convertObjectToPod(object)
	if err != nil {
		return false, err
	}

	return len(pod.OwnerReferences) > 0, nil
}

func convertObjectToPod(object runtime.Object) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	if err := scheme.Scheme.Convert(object, pod, 0); err != nil {
		return nil, err
	}
	return pod, nil
}

func podGroupName(object runtime.Object) (string, error) {
	pod, err := convertObjectToPod(object)
	if err != nil {
		return "", err
	}

	if len(pod.OwnerReferences) < 1 {
		return "", errors.Errorf("pod %s has no owner references", pod.Name)
	}

	ownerReference := pod.OwnerReferences[0]
	return fmt.Sprintf("%s pods", ownerReference.Name), nil
}

func isObjectPod(object runtime.Object) bool {
	if object == nil {
		return false
	}

	objectGVK := object.GetObjectKind().GroupVersionKind()
	podGVK := gvk.PodGVK

	return podGVK.String() == objectGVK.String()
}

func isObjectParent(child, parent runtime.Object) (bool, error) {
	childAccessor, err := meta.Accessor(child)
	if err != nil {
		return false, err
	}

	parentAccessor, err := meta.Accessor(parent)
	if err != nil {
		return false, err
	}

	for _, ownerReference := range childAccessor.GetOwnerReferences() {
		if parentAccessor.GetUID() == ownerReference.UID {
			return true, nil
		}
	}

	return false, nil
}

type ObjectStatus interface {
	Status(ctx context.Context, object runtime.Object) (*objectstatus.ObjectStatus, error)
}

type HandlerObjectStatus struct {
	objectStore   store.Store
	pluginManager plugin.ManagerInterface
}

var _ ObjectStatus = (*HandlerObjectStatus)(nil)

func NewHandlerObjectStatus(objectStore store.Store, pluginManager plugin.ManagerInterface) *HandlerObjectStatus {
	return &HandlerObjectStatus{
		objectStore:   objectStore,
		pluginManager: pluginManager,
	}
}

func (h *HandlerObjectStatus) Status(ctx context.Context, object runtime.Object) (*objectstatus.ObjectStatus, error) {
	status, err := objectstatus.Status(ctx, object, h.objectStore)
	if err != nil {
		return nil, err
	}

	pluginStatus, err := h.pluginManager.ObjectStatus(object)
	if err != nil {
		return nil, err
	}

	status.Details = append(status.Details, pluginStatus.ObjectStatus.Details...)

	return &status, nil
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

func checkReplicaCount(object runtime.Object) error {
	u, ok := object.(*unstructured.Unstructured)
	if !ok {
		return errors.Errorf("expected unstructured object; got %T", object)
	}

	i, ok, err := unstructured.NestedInt64(u.Object, "spec", "replicas")
	if err != nil {
		return err
	}

	if !ok || i < 1 {
		return &skipNode{}
	}

	return nil
}

func isObjectReplicaSet(object runtime.Object) (bool, error) {
	if object == nil {
		return false, errors.New("can't check if nil object is a replica set")
	}

	groupVersionKind := object.GetObjectKind().GroupVersionKind()

	return (groupVersionKind.Group == "apps" || groupVersionKind.Group == "extensions") &&
		groupVersionKind.Kind == "ReplicaSet", nil
}

func deDupEdges(adjList component.AdjList) component.AdjList {
	list := component.AdjList{}

	lookup := make(map[string]string)
	edgeTypeLookup := make(map[string]component.EdgeType)

	var keys []string
	for k := range adjList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := adjList[k]
		for _, edge := range v {
			edgeTypeLookup[edge.Node] = edge.Type

			if lookup[k] == edge.Node {
				continue
			}

			lookup[edge.Node] = k
		}
	}

	for k, v := range lookup {
		list[v] = append(list[v], component.Edge{Node: k, Type: edgeTypeLookup[k]})
	}

	for k := range list {
		cur := list[k]
		sort.Slice(cur, func(i, j int) bool {
			return cur[i].Node < cur[j].Node
		})
		list[k] = cur
	}

	return list
}

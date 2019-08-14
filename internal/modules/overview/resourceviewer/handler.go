package resourceviewer

import (
	"context"
	"fmt"
	"sort"
	"strings"
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
	"github.com/vmware/octant/internal/util/kubernetes"
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

type nodesStorage map[types.UID]*unstructured.Unstructured

type adjListStorage map[string]map[string]*unstructured.Unstructured

func (als adjListStorage) hasKey(uid string) bool {
	for k := range als {
		if k == uid {
			return true
		}
	}

	return false
}

func (als adjListStorage) hasEdgeForKey(keyUID, edgeUID string) bool {
	edges, ok := als[keyUID]
	if !ok {
		return false
	}

	_, ok = edges[edgeUID]
	return ok
}

func (als adjListStorage) isEdge(uid string) bool {
	for k := range als {
		for edgeUID := range als[k] {
			if uid == edgeUID {
				return true
			}
		}
	}

	return false
}

func (als adjListStorage) addEdgeForKey(uid, edgeUID string, object *unstructured.Unstructured) {
	if _, ok := als[uid]; !ok {
		als[uid] = make(map[string]*unstructured.Unstructured)
	}

	als[uid][edgeUID] = object
}

// Handler is a visitor handler.
type Handler struct {
	objectStore   store.Store
	link          link.Interface
	pluginPrinter plugin.ManagerInterface

	nodes   nodesStorage
	adjList adjListStorage

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
		adjList:       adjListStorage{},
		nodes:         nodesStorage{},
		objectStatus:  NewHandlerObjectStatus(dashConfig.ObjectStore(), dashConfig.PluginManager()),
	}

	for _, option := range options {
		option(h)
	}

	return h, nil
}

// AddEdge adds edges to the graph.
func (h *Handler) AddEdge(ctx context.Context, from, to *unstructured.Unstructured) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	fromName, err := edgeName(from)
	if err != nil {
		if isSkippedNode(err) {
			return nil
		}
		return errors.Wrap(err, "could not generate from edge")
	}

	toName, err := edgeName(to)
	if err != nil {
		if isSkippedNode(err) {
			return nil
		}
		return errors.Wrap(err, "could not generate to edge")
	}

	// is from a key in the adjacency list?
	if h.adjList.hasKey(fromName) {
		if !h.adjList.hasEdgeForKey(toName, fromName) {
			// add to to from
			h.adjList.addEdgeForKey(fromName, toName, to)
		}
	} else {
		// is to a key in the adjacency list?
		if h.adjList.hasKey(toName) {
			// add from to to
			h.adjList.addEdgeForKey(toName, fromName, from)
		} else {
			// add to to from
			h.adjList.addEdgeForKey(fromName, toName, to)
		}
	}

	h.addNode(fromName, from)
	h.addNode(toName, to)

	return nil
}

func (h *Handler) addNode(name string, object *unstructured.Unstructured) {
	h.nodes[types.UID(name)] = object
}

// Process adds nodes to the dependency graph.
func (h *Handler) Process(ctx context.Context, object *unstructured.Unstructured) error {
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

func (h *Handler) AdjacencyList() (*component.AdjList, error) {
	list := component.AdjList{}

	for k, v := range h.adjList {
		for edgeName := range v {

			list[k] = append(list[k], component.Edge{
				Node: edgeName,
				Type: component.EdgeTypeExplicit,
			})
		}

		// sort the edges by node to make them easier to compare
		sort.Slice(list[k], func(i, j int) bool {
			return list[k][i].Node < list[k][j].Node
		})
	}

	return &list, nil
}

// Nodes generates nodes from the handler.
func (h *Handler) Nodes(ctx context.Context) (component.Nodes, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	nodes := component.Nodes{}

	var podsInAGroup []unstructured.Unstructured

	for uid, node := range h.nodes {

		ok, err := isPodInGroup(node)
		if err != nil {
			return nil, err
		}

		if ok {
			podsInAGroup = append(podsInAGroup, *node)
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

func (h *Handler) buildPodGroups(podsInAGroup []unstructured.Unstructured) (map[string][]unstructured.Unstructured, error) {
	nameMap := make(map[string][]unstructured.Unstructured)
	for _, object := range podsInAGroup {
		name, err := podGroupName(&object)
		if err != nil {
			return nil, err
		}

		nameMap[name] = append(nameMap[name], object)
	}
	return nameMap, nil
}

func edgeName(object *unstructured.Unstructured) (string, error) {
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

func isPodInGroup(object *unstructured.Unstructured) (bool, error) {
	if !isObjectPod(object) {
		return false, nil
	}

	pod, err := convertObjectToPod(object)
	if err != nil {
		return false, err
	}

	return len(pod.OwnerReferences) > 0, nil
}

func convertObjectToPod(object *unstructured.Unstructured) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	if err := scheme.Scheme.Convert(object, pod, 0); err != nil {
		return nil, err
	}
	return pod, nil
}

func podGroupName(object *unstructured.Unstructured) (string, error) {
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

func isObjectPod(object *unstructured.Unstructured) bool {
	if object == nil {
		return false
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	objectGVK := object.GroupVersionKind()
	podGVK := gvk.Pod

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

	pluginStatus, err := h.pluginManager.ObjectStatus(ctx, object)
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

func checkReplicaCount(object *unstructured.Unstructured) error {
	if object == nil {
		return errors.Errorf("unable to check for replica count in nil object")
	}

	i, ok, err := unstructured.NestedInt64(object.Object, "spec", "replicas")
	if err != nil {
		return err
	}

	if !ok || i < 1 {
		return &skipNode{}
	}

	return nil
}

func isObjectReplicaSet(object *unstructured.Unstructured) (bool, error) {
	if object == nil {
		return false, errors.New("can't check if nil object is a replica set")
	}

	groupVersionKind := object.GroupVersionKind()

	return (groupVersionKind.Group == "apps" || groupVersionKind.Group == "extensions") &&
		groupVersionKind.Kind == "ReplicaSet", nil
}

func (h *Handler) summarizeNodeList() string {
	var sb strings.Builder

	header := "nodes"
	fmt.Fprintf(&sb, "%s\n%s\n", header, strings.Repeat("=", len(header)))

	var uids []string

	for uid := range h.nodes {
		uids = append(uids, string(uid))
	}

	sort.Strings(uids)

	for _, uid := range uids {
		fmt.Fprintf(&sb, "%s: %s\n", uid, kubernetes.PrintObject(h.nodes[types.UID(uid)]))
	}

	sb.WriteString("===== end ====\n")

	return sb.String()
}

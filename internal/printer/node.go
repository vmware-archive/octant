/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	// labelNodeRolePrefix is a label prefix for node roles
	// It's copied over to here until it's merged in core: https://github.com/kubernetes/kubernetes/pull/39112
	labelNodeRolePrefix = "node-role.kubernetes.io/"

	// nodeLabelRole specifies the role of a node
	nodeLabelRole = "kubernetes.io/role"
)

var (
	nodeListColumns = component.NewTableCols("Name", "Labels", "Status", "Roles", "Age", "Version")
)

// NodeListHandler is a printFunc that prints nodes
func NodeListHandler(ctx context.Context, list *corev1.NodeList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("node list is nil")
	}

	table := component.NewTable("Nodes", "We couldn't find any nodes!", nodeListColumns)

	for _, node := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&node, node.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(node.Labels)
		row["Status"] = component.NewText(nodeStatusMessage(node))
		row["Roles"] = component.NewText(nodeRoles(node))
		row["Age"] = component.NewTimestamp(node.CreationTimestamp.Time)
		row["Version"] = component.NewText(node.Status.NodeInfo.KubeletVersion)

		table.Add(row)
	}

	return table, nil
}

// NodeHandler is a printFunc that prints nodes
func NodeHandler(ctx context.Context, node *corev1.Node, options Options) (component.Component, error) {
	o := NewObject(node)

	nh, err := newNodeHandler(node, o)
	if err != nil {
		return nil, err
	}

	if err := nh.Config(options); err != nil {
		return nil, errors.Wrap(err, "print node configuration")
	}
	if err := nh.Addresses(options); err != nil {
		return nil, errors.Wrap(err, "print node addresses")
	}
	if err := nh.Resources(options); err != nil {
		return nil, errors.Wrap(err, "print node resources")
	}
	if err := nh.Conditions(options); err != nil {
		return nil, errors.Wrap(err, "print node conditions")
	}
	if err := nh.Images(options); err != nil {
		return nil, errors.Wrap(err, "print node images")
	}
	return o.ToComponent(ctx, options)
}

type nodeResource struct {
	CPU              string
	Memory           string
	EphemeralStorage string
	Pods             string
}

func parseResourceList(resourceList corev1.ResourceList) nodeResource {
	nr := nodeResource{}

	if cpu := resourceList.Cpu(); cpu != nil {
		nr.CPU = cpu.String()
	}

	if memory := resourceList.Memory(); memory != nil {
		nr.Memory = memory.String()
	}

	if ephemeralStorage := resourceList.StorageEphemeral(); ephemeralStorage != nil {
		nr.EphemeralStorage = ephemeralStorage.String()
	}

	if pods := resourceList.Pods(); pods != nil {
		nr.Pods = pods.String()
	}

	return nr
}

var (
	nodeResourcesColumns = component.NewTableCols("Key", "Capacity", "Allocatable")
)

func createNodeResourcesView(node *corev1.Node) (*component.Table, error) {
	if node == nil {
		return nil, errors.New("nil nodes don't have resources")
	}

	table := component.NewTable("Resources", "There are no resources!", nodeResourcesColumns)

	allocatable := parseResourceList(node.Status.Allocatable)
	capacity := parseResourceList(node.Status.Capacity)

	table.Add([]component.TableRow{
		{
			"Key":         component.NewText("CPU"),
			"Capacity":    component.NewText(capacity.CPU),
			"Allocatable": component.NewText(allocatable.CPU),
		},
		{
			"Key":         component.NewText("Memory"),
			"Capacity":    component.NewText(capacity.Memory),
			"Allocatable": component.NewText(allocatable.Memory),
		},
		{
			"Key":         component.NewText("Ephemeral Storage"),
			"Capacity":    component.NewText(capacity.EphemeralStorage),
			"Allocatable": component.NewText(allocatable.EphemeralStorage),
		},
		{
			"Key":         component.NewText("Pods"),
			"Capacity":    component.NewText(capacity.Pods),
			"Allocatable": component.NewText(allocatable.Pods),
		},
	}...)

	return table, nil
}

var (
	nodeAddressesColumns = component.NewTableCols("Type", "Address")
)

func createNodeAddressesView(node *corev1.Node) (*component.Table, error) {
	table := component.NewTable("Addresses", "There are no addresses!", nodeAddressesColumns)

	for _, address := range node.Status.Addresses {
		row := component.TableRow{}
		row["Type"] = component.NewText(string(address.Type))
		row["Address"] = component.NewText(address.Address)

		table.Add(row)
	}

	return table, nil
}

func nodeStatusMessage(node corev1.Node) string {
	conditionMap := make(map[corev1.NodeConditionType]*corev1.NodeCondition)
	NodeAllConditions := []corev1.NodeConditionType{corev1.NodeReady}
	for i := range node.Status.Conditions {
		cond := node.Status.Conditions[i]
		conditionMap[cond.Type] = &cond
	}
	var status []string
	for _, validCondition := range NodeAllConditions {
		if condition, ok := conditionMap[validCondition]; ok {
			if condition.Status == corev1.ConditionTrue {
				status = append(status, string(condition.Type))
			} else {
				status = append(status, "Not"+string(condition.Type))
			}
		}
	}
	if len(status) == 0 {
		status = append(status, "Unknown")
	}
	if node.Spec.Unschedulable {
		status = append(status, "SchedulingDisabled")
	}

	return strings.Join(status, ",")
}
func nodeRoles(node corev1.Node) string {
	roles := strings.Join(findNodeRoles(node), ",")
	if roles == "" {
		return "<none>"
	}

	return roles
}

// findNodeRoles returns the roles of a given node.
// The roles are determined by looking for:
// * a node-role.kubernetes.io/<role>="" label
// * a kubernetes.io/role="<role>" label
func findNodeRoles(node corev1.Node) []string {
	roles := sets.NewString()
	for k, v := range node.Labels {
		switch {
		case strings.HasPrefix(k, labelNodeRolePrefix):
			if role := strings.TrimPrefix(k, labelNodeRolePrefix); len(role) > 0 {
				roles.Insert(role)
			}

		case k == nodeLabelRole && v != "":
			roles.Insert(v)
		}
	}
	return roles.List()
}

// NodeConfiguration generates a node configuration
type NodeConfiguration struct {
	node *corev1.Node
}

// NewNodeConfiguration creates an instance of NodeConfiguration
func NewNodeConfiguration(node *corev1.Node) *NodeConfiguration {
	return &NodeConfiguration{
		node: node,
	}
}

// Create creates a node configuration summary
func (n *NodeConfiguration) Create(options Options) (*component.Summary, error) {
	if n == nil || n.node == nil {
		return nil, errors.New("cannot generate status for nil node")
	}
	nodeInfo := n.node.Status.NodeInfo

	summary := component.NewSummary("Status", []component.SummarySection{
		{
			Header:  "Architecture",
			Content: component.NewText(nodeInfo.Architecture),
		},
		{
			Header:  "Boot ID",
			Content: component.NewText(nodeInfo.BootID),
		},
		{
			Header:  "Container Runtime Version",
			Content: component.NewText(nodeInfo.ContainerRuntimeVersion),
		},
		{
			Header:  "Kernel Version",
			Content: component.NewText(nodeInfo.KernelVersion),
		},
		{
			Header:  "KubeProxy Version",
			Content: component.NewText(nodeInfo.KubeProxyVersion),
		},
		{
			Header:  "Kubelet Version",
			Content: component.NewText(nodeInfo.KubeletVersion),
		},
		{
			Header:  "Machine ID",
			Content: component.NewText(nodeInfo.MachineID),
		},
		{
			Header:  "Operating System",
			Content: component.NewText(nodeInfo.OperatingSystem),
		},
		{
			Header:  "OS Image",
			Content: component.NewText(nodeInfo.OSImage),
		},
		{
			Header:  "Pod CIDR",
			Content: component.NewText(n.node.Spec.PodCIDR),
		},
		{
			Header:  "System UUID",
			Content: component.NewText(nodeInfo.SystemUUID),
		},
	}...)

	return summary, nil
}

var (
	nodeConditionsColumns = component.NewTableCols("Type", "Reason", "Status", "Message", "Last Heartbeat", "Last Transition")
)

func createNodeConditionsView(node *corev1.Node) (*component.Table, error) {
	if node == nil {
		return nil, errors.New("cannot generate conditions for nil node")
	}

	table := component.NewTable("Conditions", "There are no node conditions!", nodeConditionsColumns)

	for _, condition := range node.Status.Conditions {
		row := component.TableRow{
			"Type":            component.NewText(string(condition.Type)),
			"Reason":          component.NewText(condition.Reason),
			"Status":          component.NewText(string(condition.Status)),
			"Message":         component.NewText(condition.Message),
			"Last Heartbeat":  component.NewTimestamp(condition.LastHeartbeatTime.Time),
			"Last Transition": component.NewTimestamp(condition.LastTransitionTime.Time),
		}

		table.Add(row)
	}

	table.Sort("Type")

	return table, nil
}

var (
	nodeImagesColumns = component.NewTableCols("Names", "Size")
)

func createNodeImagesView(node *corev1.Node) (*component.Table, error) {
	if node == nil {
		return nil, errors.New("cannot generate images for nil node")
	}

	table := component.NewTable("Images", "There are no images!", nodeImagesColumns)

	for _, containerImage := range node.Status.Images {
		row := component.TableRow{
			"Names": component.NewMarkdownText(strings.Join(containerImage.Names, "\n")),
			"Size":  component.NewText(fmt.Sprintf("%d", containerImage.SizeBytes)),
		}

		table.Add(row)
	}

	table.Sort("Names")

	return table, nil
}

type nodeObject interface {
	Config(options Options) error
	Addresses(options Options) error
	Resources(options Options) error
	Conditions(options Options) error
	Images(options Options) error
}

type nodeHandler struct {
	node           *corev1.Node
	configFunc     func(*corev1.Node, Options) (*component.Summary, error)
	addressesFunc  func(*corev1.Node, Options) (*component.Table, error)
	resourcesFunc  func(*corev1.Node, Options) (*component.Table, error)
	conditionsFunc func(*corev1.Node, Options) (*component.Table, error)
	imagesFunc     func(*corev1.Node, Options) (*component.Table, error)
	object         *Object
}

var _ nodeObject = (*nodeHandler)(nil)

func newNodeHandler(node *corev1.Node, object *Object) (*nodeHandler, error) {
	if node == nil {
		return nil, errors.New("can't print a nil node")
	}

	if object == nil {
		return nil, errors.New("can't print node using a nil object printer")
	}

	nh := &nodeHandler{
		node:           node,
		configFunc:     defaultNodeConfig,
		addressesFunc:  defaultNodeAddresses,
		resourcesFunc:  defaultNodeResources,
		conditionsFunc: defaultNodeConditions,
		imagesFunc:     defaultNodeImages,
		object:         object,
	}
	return nh, nil
}

func (n *nodeHandler) Config(options Options) error {
	out, err := n.configFunc(n.node, options)
	if err != nil {
		return err
	}
	n.object.RegisterConfig(out)
	return nil
}

func defaultNodeConfig(node *corev1.Node, options Options) (*component.Summary, error) {
	return NewNodeConfiguration(node).Create(options)
}

func (n *nodeHandler) Addresses(options Options) error {
	if n.node == nil {
		return errors.New("can't display addresses for nil node")
	}

	n.object.RegisterItems(ItemDescriptor{
		Width: component.WidthHalf,
		Func: func() (component.Component, error) {
			return n.addressesFunc(n.node, options)
		},
	})
	return nil
}

func defaultNodeAddresses(node *corev1.Node, options Options) (*component.Table, error) {
	return createNodeAddressesView(node)
}

func (n *nodeHandler) Resources(options Options) error {
	if n.node == nil {
		return errors.New("can't display resources for nil node")
	}

	n.object.RegisterItems(ItemDescriptor{
		Width: component.WidthHalf,
		Func: func() (component.Component, error) {
			return n.resourcesFunc(n.node, options)
		},
	})
	return nil
}

func defaultNodeResources(node *corev1.Node, options Options) (*component.Table, error) {
	return createNodeResourcesView(node)
}

func (n *nodeHandler) Conditions(options Options) error {
	if n.node == nil {
		return errors.New("can't display resources for nil node")
	}

	n.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return n.conditionsFunc(n.node, options)
		},
	})
	return nil
}

func defaultNodeConditions(node *corev1.Node, options Options) (*component.Table, error) {
	return createNodeConditionsView(node)
}

func (n *nodeHandler) Images(options Options) error {
	if n.node == nil {
		return errors.New("can't display resources for nil node")
	}

	n.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return n.imagesFunc(n.node, options)
		},
	})
	return nil
}

func defaultNodeImages(node *corev1.Node, options Options) (*component.Table, error) {
	return createNodeImagesView(node)
}

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestNodeListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	node := testutil.CreateNode("node-1")
	node.Status.NodeInfo.KubeletVersion = "1.15.1"
	node.CreationTimestamp = *testutil.CreateTimestamp()

	tpo.PathForObject(node, node.Name, "/node")

	list := &corev1.NodeList{
		Items: []corev1.Node{
			*node,
		},
	}

	ctx := context.Background()
	got, err := NodeListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	expected := component.NewTableWithRows("Nodes", "We couldn't find any nodes!", nodeListColumns, []component.TableRow{
		{
			"Age":     component.NewTimestamp(node.CreationTimestamp.Time),
			"Name":    component.NewLink("", "node-1", "/node"),
			"Labels":  component.NewLabels(make(map[string]string)),
			"Version": component.NewText("1.15.1"),
			"Status":  component.NewText("Unknown"),
			"Roles":   component.NewText("<none>"),
		},
	})

	component.AssertEqual(t, expected, got)
}

func Test_NodeConfiguration(t *testing.T) {
	node := testutil.CreateNode("node")
	node.Status.NodeInfo.Architecture = "amd64"
	node.Status.NodeInfo.BootID = "7eee89e0-b78a-4c30-a1bc-d43ad479b35a"
	node.Status.NodeInfo.ContainerRuntimeVersion = "containerd://1.2.6-0ubuntu1"
	node.Status.NodeInfo.KernelVersion = "4.15.0-58-generic"
	node.Status.NodeInfo.KubeProxyVersion = "v1.15.3"
	node.Status.NodeInfo.KubeletVersion = "v1.15.3"
	node.Status.NodeInfo.MachineID = "87050f150cca41c0ab58b7672b5dbc11"
	node.Status.NodeInfo.OperatingSystem = "linux"
	node.Status.NodeInfo.OSImage = "Ubuntu Disco Dingo (development branch)"
	node.Spec.PodCIDR = "10.244.0.0/24"
	node.Status.NodeInfo.SystemUUID = "5AF9704C-2606-11B2-A85C-C7F92F2B85CA"

	cases := []struct {
		name     string
		node     *corev1.Node
		isErr    bool
		expected *component.Summary
	}{
		{
			name: "general",
			node: node,
			expected: component.NewSummary("Status", []component.SummarySection{
				{
					Header:  "Architecture",
					Content: component.NewText("amd64"),
				},
				{
					Header:  "Boot ID",
					Content: component.NewText("7eee89e0-b78a-4c30-a1bc-d43ad479b35a"),
				},
				{
					Header:  "Container Runtime Version",
					Content: component.NewText("containerd://1.2.6-0ubuntu1"),
				},
				{
					Header:  "Kernel Version",
					Content: component.NewText("4.15.0-58-generic"),
				},
				{
					Header:  "KubeProxy Version",
					Content: component.NewText("v1.15.3"),
				},
				{
					Header:  "Kubelet Version",
					Content: component.NewText("v1.15.3"),
				},
				{
					Header:  "Machine ID",
					Content: component.NewText("87050f150cca41c0ab58b7672b5dbc11"),
				},
				{
					Header:  "Operating System",
					Content: component.NewText("linux"),
				},
				{
					Header:  "OS Image",
					Content: component.NewText("Ubuntu Disco Dingo (development branch)"),
				},
				{
					Header:  "Pod CIDR",
					Content: component.NewText("10.244.0.0/24"),
				},
				{
					Header:  "System UUID",
					Content: component.NewText("5AF9704C-2606-11B2-A85C-C7F92F2B85CA"),
				},
			}...),
		},
		{
			name:  "pod is nil",
			node:  nil,
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			nc := NewNodeConfiguration(tc.node)

			summary, err := nc.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createAddressesView(t *testing.T) {
	node := testutil.CreateNode("node-1")
	node.Status.Addresses = []corev1.NodeAddress{
		{
			Type:    corev1.NodeHostName,
			Address: "host.local",
		},
		{
			Type:    corev1.NodeInternalIP,
			Address: "192.168.1.1",
		},
	}

	got, err := createNodeAddressesView(node)
	require.NoError(t, err)

	expected := component.NewTableWithRows("Addresses", "There are no addresses!", nodeAddressesColumns, []component.TableRow{
		{
			"Type":    component.NewText("Hostname"),
			"Address": component.NewText("host.local"),
		},
		{
			"Type":    component.NewText("InternalIP"),
			"Address": component.NewText("192.168.1.1"),
		},
	})

	component.AssertEqual(t, expected, got)
}

func Test_createNodeResourcesView(t *testing.T) {
	node := testutil.CreateNode("node-1")

	capacityQuantity, err := resource.ParseQuantity("1")
	require.NoError(t, err)
	allocatedQuantity, err := resource.ParseQuantity("2")
	require.NoError(t, err)

	resourceNames := []corev1.ResourceName{corev1.ResourceCPU, corev1.ResourceMemory, corev1.ResourcePods,
		corev1.ResourceEphemeralStorage}

	node.Status.Allocatable = map[corev1.ResourceName]resource.Quantity{}
	node.Status.Capacity = map[corev1.ResourceName]resource.Quantity{}

	for _, resourceName := range resourceNames {
		node.Status.Allocatable[resourceName] = allocatedQuantity
		node.Status.Capacity[resourceName] = capacityQuantity
	}

	got, err := createNodeResourcesView(node)
	require.NoError(t, err)

	expected := component.NewTableWithRows("Resources", "There are no resources!", nodeResourcesColumns, []component.TableRow{
		{
			"Key":         component.NewText("CPU"),
			"Capacity":    component.NewText("1"),
			"Allocatable": component.NewText("2"),
		},
		{
			"Key":         component.NewText("Memory"),
			"Capacity":    component.NewText("1"),
			"Allocatable": component.NewText("2"),
		},
		{
			"Key":         component.NewText("Ephemeral Storage"),
			"Capacity":    component.NewText("1"),
			"Allocatable": component.NewText("2"),
		},
		{
			"Key":         component.NewText("Pods"),
			"Capacity":    component.NewText("1"),
			"Allocatable": component.NewText("2"),
		},
	})

	component.AssertEqual(t, expected, got)
}

func Test_createNodeConditionsView(t *testing.T) {

	node := testutil.CreateNode("node-1")
	node.Status.Conditions = []corev1.NodeCondition{
		{
			Type:               "type",
			Status:             "status",
			LastHeartbeatTime:  *testutil.CreateTimestamp(),
			LastTransitionTime: *testutil.CreateTimestamp(),
			Reason:             "reason",
			Message:            "message",
		},
	}

	got, err := createNodeConditionsView(node)
	require.NoError(t, err)

	expected := component.NewTableWithRows("Conditions", "There are no node conditions!", nodeConditionsColumns, []component.TableRow{
		{
			"Type":            component.NewText("type"),
			"Reason":          component.NewText("reason"),
			"Status":          component.NewText("status"),
			"Message":         component.NewText("message"),
			"Last Heartbeat":  component.NewTimestamp(node.Status.Conditions[0].LastHeartbeatTime.Time),
			"Last Transition": component.NewTimestamp(node.Status.Conditions[0].LastTransitionTime.Time),
		},
	})

	component.AssertEqual(t, expected, got)
}

func Test_createNodeImagesView(t *testing.T) {

	node := testutil.CreateNode("node-1")
	node.Status.Images = []corev1.ContainerImage{
		{
			Names:     []string{"a"},
			SizeBytes: 1000000,
		},
		{
			Names:     []string{"b-1", "b-2"},
			SizeBytes: 1000000,
		},
	}

	got, err := createNodeImagesView(node)
	require.NoError(t, err)

	expected := component.NewTableWithRows("Images", "There are no images!", nodeImagesColumns, []component.TableRow{
		{
			"Names":     component.NewMarkdownText("a"),
			"Size (GB)": component.NewText("0.001"),
		},
		{
			"Names":     component.NewMarkdownText("b-1\nb-2"),
			"Size (GB)": component.NewText("0.001"),
		},
	})

	component.AssertEqual(t, expected, got)
}

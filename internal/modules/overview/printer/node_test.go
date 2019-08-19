package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
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

func Test_nodeAddresses(t *testing.T) {
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

	got, err := nodeAddresses(node)
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

func Test_nodeResources(t *testing.T) {
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

	got, err := nodeResources(node)
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

func Test_nodeConditions(t *testing.T) {

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

	got, err := nodeConditions(node)
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

func Test_nodeImages(t *testing.T) {

	node := testutil.CreateNode("node-1")
	node.Status.Images = []corev1.ContainerImage{
		{
			Names:     []string{"a"},
			SizeBytes: 10,
		},
		{
			Names:     []string{"b-1", "b-2"},
			SizeBytes: 10,
		},
	}

	got, err := nodeImages(node)
	require.NoError(t, err)

	expected := component.NewTableWithRows("Images", "There are no images!", nodeImagesColumns, []component.TableRow{
		{
			"Names": component.NewMarkdownText("a"),
			"Size":  component.NewText("10"),
		},
		{
			"Names": component.NewMarkdownText("b-1\nb-2"),
			"Size":  component.NewText("10"),
		},
	})

	component.AssertEqual(t, expected, got)
}

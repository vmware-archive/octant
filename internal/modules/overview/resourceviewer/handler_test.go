package resourceviewer

import (
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"

	configFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/modules/overview/objectstatus"
	"github.com/vmware/octant/internal/modules/overview/resourceviewer/fake"
	"github.com/vmware/octant/internal/testutil"
	pluginFake "github.com/vmware/octant/pkg/plugin/fake"
	storeFake "github.com/vmware/octant/pkg/store/fake"
	"github.com/vmware/octant/pkg/view/component"
)

func TestHandler(t *testing.T) {
	cr := testutil.CreateClusterRole("cluster-role")

	deployment := testutil.CreateDeployment("deployment")
	deployment.SetOwnerReferences(testutil.ToOwnerReferences(t, cr))

	replicaSet1 := testutil.CreateAppReplicaSet("replica-set-1")
	replicaSet1.SetOwnerReferences(testutil.ToOwnerReferences(t, deployment))
	replicaSet1.Spec.Replicas = pointer.Int32Ptr(1)

	replicaSet2 := testutil.CreateAppReplicaSet("replica-set-2")
	replicaSet2.SetOwnerReferences(testutil.ToOwnerReferences(t, deployment))

	replicaSet3 := testutil.CreateExtReplicaSet("replica-set-3")
	replicaSet3.SetOwnerReferences(testutil.ToOwnerReferences(t, deployment))
	replicaSet3.Spec.Replicas = pointer.Int32Ptr(1)

	serviceAccount := testutil.CreateServiceAccount("service-account")

	pod1 := testutil.CreatePod("pod-1")
	pod1.Spec.ServiceAccountName = serviceAccount.Name
	pod1.SetOwnerReferences(testutil.ToOwnerReferences(t, replicaSet1))
	pod2 := testutil.CreatePod("pod-2")
	pod2.Spec.ServiceAccountName = serviceAccount.Name
	pod2.SetOwnerReferences(testutil.ToOwnerReferences(t, replicaSet1))
	pod3 := testutil.CreatePod("pod-3")
	pod3.Spec.ServiceAccountName = serviceAccount.Name
	pod4 := testutil.CreatePod("pod-4")
	pod4.Spec.ServiceAccountName = serviceAccount.Name
	pod4.SetOwnerReferences(testutil.ToOwnerReferences(t, replicaSet3))

	service1 := testutil.CreateService("service1")

	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	objectStore := storeFake.NewMockStore(controller)
	dashConfig.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

	pluginManager := pluginFake.NewMockManagerInterface(controller)
	dashConfig.EXPECT().PluginManager().Return(pluginManager).AnyTimes()

	objectStatus := fake.NewMockObjectStatus(controller)
	objectStatus.EXPECT().
		Status(gomock.Any(), gomock.Any()).
		Return(&objectstatus.ObjectStatus{}, nil).
		AnyTimes()

	handler, err := NewHandler(dashConfig, SetHandlerObjectStatus(objectStatus))
	require.NoError(t, err)

	ctx := context.Background()
	mockRelations := func(a runtime.Object, objects ...runtime.Object) {
		for _, b := range objects {
			require.NoError(t, handler.AddEdge(ctx, a, b))
			require.NoError(t, handler.AddEdge(ctx, b, a))
			require.NoError(t, handler.Process(ctx, b))
		}
		require.NoError(t, handler.Process(ctx, a))
	}

	mockRelations(cr, deployment)
	mockRelations(deployment, replicaSet1, replicaSet2, replicaSet3)
	mockRelations(replicaSet1, pod1, pod2)
	mockRelations(replicaSet3, pod4)
	mockRelations(service1, pod1, pod2)
	mockRelations(service1, pod4)

	require.NoError(t, handler.Process(ctx, pod3))
	require.NoError(t, handler.AddEdge(ctx, pod1, serviceAccount))
	require.NoError(t, handler.AddEdge(ctx, pod2, serviceAccount))
	require.NoError(t, handler.AddEdge(ctx, pod3, serviceAccount))
	require.NoError(t, handler.AddEdge(ctx, pod4, serviceAccount))
	require.NoError(t, handler.Process(ctx, serviceAccount))

	mockLinkPath(t, dashConfig, cr)
	mockLinkPath(t, dashConfig, deployment)
	mockLinkPath(t, dashConfig, replicaSet1)
	mockLinkPath(t, dashConfig, replicaSet3)
	mockLinkPath(t, dashConfig, pod1)
	mockLinkPath(t, dashConfig, pod2)
	mockLinkPath(t, dashConfig, pod3)
	mockLinkPath(t, dashConfig, pod4)
	mockLinkPath(t, dashConfig, serviceAccount)
	mockLinkPath(t, dashConfig, service1)

	expectedAdjList := &component.AdjList{
		string(cr.UID): {
			{Node: string(deployment.UID), Type: component.EdgeTypeExplicit},
		},
		string(deployment.UID): {
			{Node: string(replicaSet1.UID), Type: component.EdgeTypeExplicit},
			{Node: string(replicaSet3.UID), Type: component.EdgeTypeExplicit},
		},
		fmt.Sprintf("%s pods", replicaSet1.Name): {
			{Node: string(serviceAccount.UID), Type: component.EdgeTypeExplicit},
		},
		fmt.Sprintf("%s pods", replicaSet3.Name): {
			{Node: string(serviceAccount.UID), Type: component.EdgeTypeExplicit},
		},
		string(pod3.UID): {
			{Node: string(serviceAccount.UID), Type: component.EdgeTypeExplicit},
		},
		string(replicaSet1.UID): {
			{Node: fmt.Sprintf("%s pods", replicaSet1.Name), Type: component.EdgeTypeExplicit},
		},
		string(replicaSet3.UID): {
			{Node: fmt.Sprintf("%s pods", replicaSet3.Name), Type: component.EdgeTypeExplicit},
		},
		string(service1.UID): {
			{Node: fmt.Sprintf("%s pods", replicaSet1.Name), Type: component.EdgeTypeExplicit},
			{Node: fmt.Sprintf("%s pods", replicaSet3.Name), Type: component.EdgeTypeExplicit},
		},
	}

	list, err := handler.AdjacencyList()
	require.NoError(t, err)
	require.Equal(t, expectedAdjList, list, "adjacency lists don't match")

	objectPath := func(t *testing.T, object runtime.Object) *component.Link {
		accessor, err := meta.Accessor(object)
		require.NoError(t, err)
		name := accessor.GetName()
		return component.NewLink("", name, path.Join("/", name))
	}

	podStatus1 := component.NewPodStatus()
	podStatus1.AddSummary(pod1.Name, nil, component.NodeStatusOK)
	podStatus1.AddSummary(pod2.Name, nil, component.NodeStatusOK)

	podStatus2 := component.NewPodStatus()
	podStatus2.AddSummary(pod4.Name, nil, component.NodeStatusOK)

	expectedNodes := component.Nodes{
		string(cr.UID): {
			Name:       cr.Name,
			APIVersion: cr.APIVersion,
			Kind:       cr.Kind,
			Status:     component.NodeStatusOK,
			Path:       objectPath(t, cr),
		},
		string(deployment.UID): {
			Name:       deployment.Name,
			APIVersion: deployment.APIVersion,
			Kind:       deployment.Kind,
			Status:     component.NodeStatusOK,
			Path:       objectPath(t, deployment),
		},
		string(replicaSet1.UID): {
			Name:       replicaSet1.Name,
			APIVersion: "apps/v1",
			Kind:       replicaSet1.Kind,
			Status:     component.NodeStatusOK,
			Path:       objectPath(t, replicaSet1),
		},
		string(replicaSet3.UID): {
			Name:       replicaSet3.Name,
			APIVersion: "extensions/v1beta1",
			Kind:       replicaSet3.Kind,
			Status:     component.NodeStatusOK,
			Path:       objectPath(t, replicaSet3),
		},
		string(pod3.UID): {
			Name:       pod3.Name,
			APIVersion: pod3.APIVersion,
			Kind:       pod3.Kind,
			Status:     component.NodeStatusOK,
			Path:       objectPath(t, pod3),
		},
		fmt.Sprintf("%s pods", replicaSet1.Name): {
			Name:       fmt.Sprintf("%s pods", replicaSet1.Name),
			APIVersion: "v1",
			Kind:       "Pod",
			Status:     component.NodeStatusOK,
			Details:    []component.Component{podStatus1},
		},
		fmt.Sprintf("%s pods", replicaSet3.Name): {
			Name:       fmt.Sprintf("%s pods", replicaSet3.Name),
			APIVersion: "v1",
			Kind:       "Pod",
			Status:     component.NodeStatusOK,
			Details:    []component.Component{podStatus2},
		},
		string(serviceAccount.UID): {
			Name:       serviceAccount.Name,
			APIVersion: serviceAccount.APIVersion,
			Kind:       serviceAccount.Kind,
			Status:     component.NodeStatusOK,
			Path:       objectPath(t, serviceAccount),
		},
		string(service1.UID): {
			Name:       service1.Name,
			APIVersion: service1.APIVersion,
			Kind:       service1.Kind,
			Status:     component.NodeStatusOK,
			Path:       objectPath(t, service1),
		},
	}

	nodes, err := handler.Nodes(ctx)
	require.NoError(t, err)

	testutil.AssertJSONEqual(t, expectedNodes, nodes)
}

func Test_edgeName(t *testing.T) {
	replicaSet := testutil.CreateAppReplicaSet("replica-set")
	replicaSetPod := testutil.CreatePod("pod")
	replicaSetPod.SetOwnerReferences(testutil.ToOwnerReferences(t, replicaSet))

	tests := []struct {
		name     string
		object   runtime.Object
		expected string
		isErr    bool
	}{
		{
			name:     "in general",
			object:   testutil.CreateDeployment("deployment"),
			expected: "deployment",
		},
		{
			name:     "pod",
			object:   testutil.CreatePod("pod"),
			expected: "pod",
		},
		{
			name:     "pod in replica set",
			object:   replicaSetPod,
			expected: "replica-set pods",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			got, err := edgeName(test.object)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.expected, got)
		})
	}
}

func Test_isObjectParent(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment")

	replicaSet1 := testutil.CreateAppReplicaSet("replica-set-1")
	replicaSet1.SetOwnerReferences(testutil.ToOwnerReferences(t, deployment))

	replicaSet2 := testutil.CreateAppReplicaSet("replica-set-2")

	tests := []struct {
		name     string
		parent   runtime.Object
		child    runtime.Object
		expected bool
		wantErr  bool
	}{
		{
			name:     "is parent",
			parent:   deployment,
			child:    replicaSet1,
			expected: true,
		},
		{
			name:     "is not parent",
			parent:   deployment,
			child:    replicaSet2,
			expected: false,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			got, err := isObjectParent(test.child, test.parent)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, got)
		})
	}
}

func Test_adjListStorage(t *testing.T) {
	als := &adjListStorage{}

	assert.False(t, als.hasEdgeForKey("1", "2"))
	als.addEdgeForKey("1", "2", nil)
	assert.True(t, als.hasEdgeForKey("1", "2"))
	assert.True(t, als.hasKey("1"))
}

func mockLinkPath(t *testing.T, dashConfig *configFake.MockDash, object runtime.Object) {
	accessor, err := meta.Accessor(object)
	require.NoError(t, err)

	apiVersion, kind := object.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()

	label := path.Join("/", accessor.GetName())

	dashConfig.EXPECT().
		ObjectPath(accessor.GetNamespace(), apiVersion, kind, accessor.GetName()).
		Return(label, nil).
		AnyTimes()
}

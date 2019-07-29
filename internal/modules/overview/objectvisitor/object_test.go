package objectvisitor_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	clusterFake "github.com/vmware/octant/internal/cluster/fake"
	configFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/modules/overview/objectvisitor"
	"github.com/vmware/octant/internal/modules/overview/objectvisitor/fake"
	queryerFake "github.com/vmware/octant/internal/queryer/fake"
	"github.com/vmware/octant/internal/testutil"
)

func TestObject_Visit(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	resourceList := &metav1.APIResourceList{}
	resourceList.APIResources = []metav1.APIResource{
		{
			Kind:       "Deployment",
			Namespaced: true,
		},
	}

	discoveryClient := clusterFake.NewMockDiscoveryInterface(controller)
	discoveryClient.EXPECT().
		ServerResourcesForGroupVersion("apps/v1").Return(resourceList, nil)

	clusterClient := clusterFake.NewMockClientInterface(controller)
	clusterClient.EXPECT().DiscoveryClient().Return(discoveryClient, nil)

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().ClusterClient().Return(clusterClient)

	deployment := testutil.CreateDeployment("deployment")

	replicaSet := testutil.CreateAppReplicaSet("replica-set")
	replicaSet.SetOwnerReferences(testutil.ToOwnerReferences(t, deployment))

	pod := testutil.CreatePod("pod")
	pod.SetOwnerReferences(testutil.ToOwnerReferences(t, replicaSet))

	u := testutil.ToUnstructured(t, replicaSet)

	q := queryerFake.NewMockQueryer(controller)
	q.EXPECT().
		OwnerReference(gomock.Any(), deployment.Namespace, replicaSet.OwnerReferences[0]).
		Return(deployment, nil)
	q.EXPECT().
		Children(gomock.Any(), u).
		Return([]runtime.Object{pod}, nil)

	handler := fake.NewMockObjectHandler(controller)
	handler.EXPECT().AddEdge(gomock.Any(), u, pod).Return(nil)
	handler.EXPECT().AddEdge(gomock.Any(), u, deployment).Return(nil)
	handler.EXPECT().Process(gomock.Any(), u).Return(nil)

	var visited []runtime.Object
	visitor := fake.NewMockVisitor(controller)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler).
		DoAndReturn(func(ctx context.Context, object runtime.Object, handler objectvisitor.ObjectHandler) error {
			visited = append(visited, object)
			return nil
		}).
		AnyTimes()

	object := objectvisitor.NewObject(dashConfig, q)

	ctx := context.Background()
	err := object.Visit(ctx, u, handler, visitor)

	sortObjectsByName(t, visited)

	expected := []runtime.Object{deployment, pod}
	assert.Equal(t, expected, visited)
	assert.NoError(t, err)
}

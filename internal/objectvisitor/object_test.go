package objectvisitor_test

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	configFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/objectvisitor"
	"github.com/vmware/octant/internal/objectvisitor/fake"
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

	dashConfig := configFake.NewMockDash(controller)

	deployment := testutil.ToUnstructured(t, testutil.CreateDeployment("deployment"))
	replicaSet := testutil.ToUnstructured(t, testutil.CreateAppReplicaSet("replica-set"))
	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))

	q := queryerFake.NewMockQueryer(controller)
	q.EXPECT().
		OwnerReference(gomock.Any(), replicaSet).
		Return(true, deployment, nil)
	q.EXPECT().
		Children(gomock.Any(), replicaSet).Return(testutil.ToUnstructuredList(t, pod), nil)

	handler := fake.NewMockObjectHandler(controller)
	handler.EXPECT().AddEdge(gomock.Any(), replicaSet, pod).Return(nil)
	handler.EXPECT().AddEdge(gomock.Any(), replicaSet, deployment).Return(nil)
	handler.EXPECT().Process(gomock.Any(), replicaSet).Return(nil)

	var visited []unstructured.Unstructured
	var mu sync.Mutex
	visitor := fake.NewMockVisitor(controller)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler, gomock.Any()).
		DoAndReturn(func(ctx context.Context, object *unstructured.Unstructured, handler objectvisitor.ObjectHandler, _ bool) error {
			mu.Lock()
			defer mu.Unlock()
			visited = append(visited, *object)
			return nil
		}).
		Times(2)

	object := objectvisitor.NewObject(dashConfig, q)

	ctx := context.Background()
	err := object.Visit(ctx, replicaSet, handler, visitor, true)
	require.NoError(t, err)

	sortObjectsByName(t, visited)
	expected := testutil.ToUnstructuredList(t, deployment, pod)
	assert.Equal(t, expected.Items, visited)
}

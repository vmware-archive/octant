package objectvisitor_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware/octant/internal/objectvisitor"
	"github.com/vmware/octant/internal/objectvisitor/fake"
	queryerFake "github.com/vmware/octant/internal/queryer/fake"
	"github.com/vmware/octant/internal/testutil"
)

func TestService_Visit(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	object := testutil.CreateService("service")
	u := testutil.ToUnstructured(t, object)

	q := queryerFake.NewMockQueryer(controller)
	ingress := testutil.CreateIngress("ingress")
	q.EXPECT().
		IngressesForService(gomock.Any(), object).
		Return([]*extv1beta1.Ingress{ingress}, nil)
	pod := testutil.CreatePod("pod")
	q.EXPECT().
		PodsForService(gomock.Any(), object).
		Return([]*corev1.Pod{pod}, nil)

	handler := fake.NewMockObjectHandler(controller)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, ingress)).
		Return(nil)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, pod)).
		Return(nil)

	var visited []unstructured.Unstructured
	visitor := fake.NewMockVisitor(controller)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler, true).
		DoAndReturn(func(ctx context.Context, object *unstructured.Unstructured, handler objectvisitor.ObjectHandler, _ bool) error {
			visited = append(visited, *object)
			return nil
		}).
		AnyTimes()

	service := objectvisitor.NewService(q)

	ctx := context.Background()

	err := service.Visit(ctx, u, handler, visitor, true)

	sortObjectsByName(t, visited)
	expected := testutil.ToUnstructuredList(t, ingress, pod)
	assert.Equal(t, expected.Items, visited)
	assert.NoError(t, err)
}

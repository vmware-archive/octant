package objectvisitor_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/internal/modules/overview/objectvisitor"
	"github.com/vmware/octant/internal/modules/overview/objectvisitor/fake"
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
		AddEdge(u, ingress).
		Return(nil)
	handler.EXPECT().
		AddEdge(u, pod).
		Return(nil)
	handler.EXPECT().
		Process(gomock.Any(), u).Return(nil)

	var visited []runtime.Object
	visitor := fake.NewMockVisitor(controller)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler).
		DoAndReturn(func(ctx context.Context, object runtime.Object, handler objectvisitor.ObjectHandler) error {
			visited = append(visited, object)
			return nil
		}).
		AnyTimes()

	service := objectvisitor.NewService(q)

	ctx := context.Background()

	err := service.Visit(ctx, u, handler, visitor)

	sortObjectsByName(t, visited)
	expected := []runtime.Object{ingress, pod}
	assert.Equal(t, expected, visited)
	assert.NoError(t, err)
}

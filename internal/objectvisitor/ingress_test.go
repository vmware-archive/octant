package objectvisitor_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/objectvisitor"
	"github.com/vmware-tanzu/octant/internal/objectvisitor/fake"
	queryerFake "github.com/vmware-tanzu/octant/internal/queryer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestIngress_Visit(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	object := testutil.CreateIngress("ingress")
	u := testutil.ToUnstructured(t, object)

	q := queryerFake.NewMockQueryer(controller)
	service := testutil.CreateService("service")
	q.EXPECT().
		ServicesForIngress(gomock.Any(), object).
		Return(testutil.ToUnstructuredList(t, service), nil)

	handler := fake.NewMockObjectHandler(controller)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, service), gomock.Any()).
		Return(nil)

	var visited []unstructured.Unstructured
	visitor := fake.NewMockVisitor(controller)
	handler.EXPECT().SetLevel(gomock.Any(), 1).Return(2)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler, true, gomock.Any()).
		DoAndReturn(func(ctx context.Context, object *unstructured.Unstructured, handler objectvisitor.ObjectHandler, _ bool, _ int) error {
			visited = append(visited, *object)
			return nil
		})

	ingress := objectvisitor.NewIngress(q)

	ctx := context.Background()
	err := ingress.Visit(ctx, u, handler, visitor, true, 1)

	sortObjectsByName(t, visited)

	expected := testutil.ToUnstructuredList(t, service)
	assert.Equal(t, expected.Items, visited)
	assert.NoError(t, err)
}

func TestIngress_Visit_invalid_service_name(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	object := testutil.CreateIngress("ingress")
	u := testutil.ToUnstructured(t, object)

	q := queryerFake.NewMockQueryer(controller)
	q.EXPECT().
		ServicesForIngress(gomock.Any(), object).
		Return(testutil.ToUnstructuredList(t), nil)

	handler := fake.NewMockObjectHandler(controller)
	handler.EXPECT().SetLevel(gomock.Any(), 1).Return(2)

	visitor := fake.NewMockVisitor(controller)

	ingress := objectvisitor.NewIngress(q)

	ctx := context.Background()
	err := ingress.Visit(ctx, u, handler, visitor, true, 1)

	assert.NoError(t, err)

}

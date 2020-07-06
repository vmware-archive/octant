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

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/objectvisitor"
	"github.com/vmware-tanzu/octant/internal/objectvisitor/fake"
	queryerFake "github.com/vmware-tanzu/octant/internal/queryer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestObject_Visit(t *testing.T) {
	deployment := testutil.ToUnstructured(t, testutil.CreateDeployment("deployment"))
	replicaSet := testutil.ToUnstructured(t, testutil.CreateAppReplicaSet("replica-set"))
	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))

	type ctorArgs struct {
		dashConfig func(ctrl *gomock.Controller) *configFake.MockDash
		queryer    func(ctrl *gomock.Controller) *queryerFake.MockQueryer
	}
	tests := []struct {
		name        string
		ctorArgs    ctorArgs
		handler     func(ctrl *gomock.Controller) *fake.MockObjectHandler
		expected    *unstructured.UnstructuredList
		visitObject *unstructured.Unstructured
	}{
		{
			name: "single owner reference",
			ctorArgs: ctorArgs{
				dashConfig: func(ctrl *gomock.Controller) *configFake.MockDash {
					dashConfig := configFake.NewMockDash(ctrl)
					return dashConfig
				},
				queryer: func(ctrl *gomock.Controller) *queryerFake.MockQueryer {
					q := queryerFake.NewMockQueryer(ctrl)
					q.EXPECT().
						OwnerReference(gomock.Any(), replicaSet).
						Return(true, []*unstructured.Unstructured{deployment}, nil)
					q.EXPECT().
						Children(gomock.Any(), replicaSet).Return(testutil.ToUnstructuredList(t, pod), nil)
					return q
				},
			},
			handler: func(ctrl *gomock.Controller) *fake.MockObjectHandler {
				handler := fake.NewMockObjectHandler(ctrl)
				handler.EXPECT().AddEdge(gomock.Any(), replicaSet, pod).Return(nil)
				handler.EXPECT().AddEdge(gomock.Any(), replicaSet, deployment).Return(nil)
				handler.EXPECT().Process(gomock.Any(), replicaSet).Return(nil)
				return handler
			},
			visitObject: replicaSet,
			expected:    testutil.ToUnstructuredList(t, deployment, pod),
		},
		{
			name: "multiple owner reference",
			ctorArgs: ctorArgs{
				dashConfig: func(ctrl *gomock.Controller) *configFake.MockDash {
					dashConfig := configFake.NewMockDash(ctrl)
					return dashConfig
				},
				queryer: func(ctrl *gomock.Controller) *queryerFake.MockQueryer {
					q := queryerFake.NewMockQueryer(ctrl)
					q.EXPECT().
						OwnerReference(gomock.Any(), pod).
						Return(true, []*unstructured.Unstructured{deployment, replicaSet}, nil)
					q.EXPECT().
						Children(gomock.Any(), pod).Return(&unstructured.UnstructuredList{}, nil)
					return q
				},
			},
			visitObject: pod,
			handler: func(ctrl *gomock.Controller) *fake.MockObjectHandler {
				handler := fake.NewMockObjectHandler(ctrl)
				handler.EXPECT().AddEdge(gomock.Any(), pod, replicaSet).Return(nil)
				handler.EXPECT().AddEdge(gomock.Any(), pod, deployment).Return(nil)
				handler.EXPECT().Process(gomock.Any(), pod).Return(nil)
				return handler
			},
			expected: testutil.ToUnstructuredList(t, deployment, replicaSet),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			resourceList := &metav1.APIResourceList{}
			resourceList.APIResources = []metav1.APIResource{
				{
					Kind:       "Deployment",
					Namespaced: true,
				},
			}

			handler := tt.handler(controller)

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

			object := objectvisitor.NewObject(tt.ctorArgs.dashConfig(controller), tt.ctorArgs.queryer(controller))

			ctx := context.Background()
			err := object.Visit(ctx, tt.visitObject, handler, visitor, true)
			require.NoError(t, err)

			sortObjectsByName(t, visited)
			assert.Equal(t, tt.expected.Items, visited)

		})
	}
}

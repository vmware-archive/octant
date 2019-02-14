package objectstatus

import (
	"testing"

	"github.com/heptio/developer-dash/internal/overview/objectvisitor"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func Test_status(t *testing.T) {
	deployObjectStatus := ObjectStatus{
		NodeStatus: component.NodeStatusOK,
		Details:    component.Title(component.NewText("apps/v1 Deployment is OK")),
	}

	lookup := statusLookup{
		{apiVersion: "v1", kind: "Object"}: func(runtime.Object) (ObjectStatus, error) {
			return deployObjectStatus, nil
		},
	}

	cases := []struct {
		name     string
		object   runtime.Object
		lookup   statusLookup
		expected ObjectStatus
		isErr    bool
	}{
		{
			name:     "in general",
			object:   createDeployment("deployment"),
			lookup:   lookup,
			expected: deployObjectStatus,
		},
		{
			name:   "nil object",
			object: nil,
			lookup: lookup,
			isErr:  true,
		},
		{
			name:   "nil lookup",
			object: createDeployment("deployment"),
			lookup: nil,
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := status(tc.object, tc.lookup)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}

func createDeployment(name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta:   genTypeMeta(objectvisitor.DeploymentGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

func genTypeMeta(gvk schema.GroupVersionKind) metav1.TypeMeta {
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	return metav1.TypeMeta{
		APIVersion: apiVersion,
		Kind:       kind,
	}
}

func genObjectMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: "namespace",
		UID:       types.UID(name),
	}
}

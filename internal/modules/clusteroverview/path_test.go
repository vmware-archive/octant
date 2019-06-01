package clusteroverview

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/internal/gvk"
	"github.com/heptio/developer-dash/internal/testutil"
)

func Test_objectPath_SupportedGroupVersionKind(t *testing.T) {
	tests := []struct {
		name string
		gvk  schema.GroupVersionKind
	}{
		{
			name: "cluster role",
			gvk:  gvk.ClusterRoleGVK,
		},
		{
			name: "cluster role binding",
			gvk:  gvk.ClusterRoleBindingGVK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			op := objectPath{}
			supported := op.SupportedGroupVersionKind()
			requireGVKPresent(t, test.gvk, supported)
		})
	}
}

func Test_objectPath_GroupVersionKindPath(t *testing.T) {
	tests := []struct {
		name     string
		object   runtime.Object
		isErr    bool
		expected string
	}{
		{
			name:     "cluster role",
			object:   testutil.CreateClusterRole("object"),
			expected: buildObjectPath("/rbac/cluster-roles/object"),
		},
		{
			name:     "cluster role binding",
			object:   testutil.CreateClusterRoleBinding("object", "roleName", nil),
			expected: buildObjectPath("/rbac/cluster-role-bindings/object"),
		},
		{
			name: "unknown",
			object: testutil.CreateEvent("object"),
			isErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			op := objectPath{}

			apiVersion, kind, name := objectDetails(t, test.object)

			got, err := op.GroupVersionKindPath(testutil.DefaultNamespace, apiVersion, kind, name)
			if test.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, test.expected, got)
		})
	}
}

func requireGVKPresent(t *testing.T, gvk schema.GroupVersionKind, list []schema.GroupVersionKind) {
	for _, current := range list {
		if current.Group == gvk.Group &&
			current.Version == gvk.Version &&
			current.Kind == gvk.Kind {
			return
		}
	}

	t.Fatalf("%s was not present", gvk.String())
}

func buildObjectPath(rest string) string {
	return path.Join("/content/cluster-overview", rest)
}

func objectDetails(t *testing.T, object runtime.Object) (string, string, string) {
	accessor := meta.NewAccessor()

	apiVersion, err := accessor.APIVersion(object)
	require.NoError(t, err)

	kind, err := accessor.Kind(object)
	require.NoError(t, err)

	name, err := accessor.Name(object)
	require.NoError(t, err)

	return apiVersion, kind, name
}

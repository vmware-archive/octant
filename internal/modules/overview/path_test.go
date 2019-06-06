package overview

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
			name: "cron job",
			gvk:  gvk.CronJobGVK,
		},
		{
			name: "daemon set",
			gvk:  gvk.DaemonSetGVK,
		},
		{
			name: "deployment",
			gvk:  gvk.DeploymentGVK,
		},
		{
			name: "job",
			gvk:  gvk.JobGVK,
		},
		{
			name: "pod",
			gvk:  gvk.PodGVK,
		},
		{
			name: "replica set",
			gvk:  gvk.ReplicaSetGVK,
		},
		{
			name: "replication controller",
			gvk:  gvk.ReplicationControllerGVK,
		},
		{
			name: "stateful set",
			gvk:  gvk.StatefulSetGVK,
		},
		{
			name: "ingress",
			gvk:  gvk.IngressGVK,
		},
		{
			name: "service",
			gvk:  gvk.ServiceGVK,
		},
		{
			name: "config map",
			gvk:  gvk.ConfigMapGVK,
		},
		{
			name: "secret",
			gvk:  gvk.SecretGVK,
		},
		{
			name: "persistent volume claim",
			gvk:  gvk.PersistentVolumeClaimGVK,
		},
		{
			name: "service account",
			gvk:  gvk.ServiceAccountGVK,
		},
		{
			name: "role",
			gvk:  gvk.RoleGVK,
		},
		{
			name: "role binding",
			gvk:  gvk.RoleBindingGVK,
		},
		{
			name: "event",
			gvk:  gvk.Event,
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
		namespace string
		object   runtime.Object
		isErr    bool
		expected string
	}{
		{
			name:     "cron job",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateCronJob("object"),
			expected: buildObjectPath("/workloads/cron-jobs/object"),
		},
		{
			name:     "daemon set",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateDaemonSet("object"),
			expected: buildObjectPath("/workloads/daemon-sets/object"),
		},
		{
			name:     "deployment",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateDeployment("object"),
			expected: buildObjectPath("/workloads/deployments/object"),
		},
		{
			name:     "job",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateJob("object"),
			expected: buildObjectPath("/workloads/jobs/object"),
		},
		{
			name:     "pod",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreatePod("object"),
			expected: buildObjectPath("/workloads/pods/object"),
		},
		{
			name:     "replica set",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateReplicaSet("object"),
			expected: buildObjectPath("/workloads/replica-sets/object"),
		},
		{
			name:     "replication controller",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateReplicationController("object"),
			expected: buildObjectPath("/workloads/replication-controllers/object"),
		},
		{
			name:     "stateful set",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateStatefulSet("object"),
			expected: buildObjectPath("/workloads/stateful-sets/object"),
		},
		{
			name:     "ingress",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateIngress("object"),
			expected: buildObjectPath("/discovery-and-load-balancing/ingresses/object"),
		},
		{
			name:     "service",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateService("object"),
			expected: buildObjectPath("/discovery-and-load-balancing/services/object"),
		},
		{
			name:     "config map",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateConfigMap("object"),
			expected: buildObjectPath("/config-and-storage/config-maps/object"),
		},
		{
			name:     "secret",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateSecret("object"),
			expected: buildObjectPath("/config-and-storage/secrets/object"),
		},
		{
			name:     "persistent volume claim",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreatePersistentVolumeClaim("object"),
			expected: buildObjectPath("/config-and-storage/persistent-volume-claims/object"),
		},
		{
			name:     "service account",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateServiceAccount("object"),
			expected: buildObjectPath("/config-and-storage/service-accounts/object"),
		},
		{
			name:     "role",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateRole("object"),
			expected: buildObjectPath("/rbac/roles/object"),
		},
		{
			name:     "role binding",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateRoleBinding("object", "roleName", nil),
			expected: buildObjectPath("/rbac/role-bindings/object"),
		},
		{
			name:     "event",
			namespace: testutil.DefaultNamespace,
			object:   testutil.CreateEvent("object"),
			expected: buildObjectPath("/events/object"),
		},
		{
			name: "unknown",
			namespace: testutil.DefaultNamespace,
			object: testutil.CreateClusterRole("object"),
			isErr: true,
		},
		{
			name: "no namespace",
			object:   testutil.CreateEvent("object"),
			isErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			op := objectPath{}

			apiVersion, kind, name := objectDetails(t, test.object)

			got, err := op.GroupVersionKindPath(test.namespace, apiVersion, kind, name)
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

func buildObjectPath(rest string) string {
	return path.Join("/content/overview/namespace", testutil.DefaultNamespace, rest)
}

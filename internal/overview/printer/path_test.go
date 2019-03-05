package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestObjectReferencePath(t *testing.T) {
	cases := []struct {
		name            string
		objectReference corev1.ObjectReference
		isErr           bool
		expected        string
	}{
		{
			name: "cron job (namespace)",
			objectReference: corev1.ObjectReference{
				APIVersion: "batch/v1beta1",
				Kind:       "CronJob",
				Name:       "cj1",
				Namespace:  "default",
			},
			expected: "/content/overview/namespace/default/workloads/cron-jobs/cj1",
		},
		{
			name: "cron job",
			objectReference: corev1.ObjectReference{
				APIVersion: "batch/v1beta1",
				Kind:       "CronJob",
				Name:       "cj1",
			},
			expected: "/content/overview/workloads/cron-jobs/cj1",
		},
		{
			name: "daemon set",
			objectReference: corev1.ObjectReference{
				APIVersion: "apps/v1",
				Kind:       "DaemonSet",
				Name:       "ds1",
			},
			expected: "/content/overview/workloads/daemon-sets/ds1",
		},
		{
			name: "deployment",
			objectReference: corev1.ObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       "d1",
			},
			expected: "/content/overview/workloads/deployments/d1",
		},
		{
			name: "job",
			objectReference: corev1.ObjectReference{
				APIVersion: "batch/v1",
				Kind:       "Job",
				Name:       "j1",
			},
			expected: "/content/overview/workloads/jobs/j1",
		},
		{
			name: "pod",
			objectReference: corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "Pod",
				Name:       "p1",
			},
			expected: "/content/overview/workloads/pods/p1",
		},
		{
			name: "replica set",
			objectReference: corev1.ObjectReference{
				APIVersion: "apps/v1",
				Kind:       "ReplicaSet",
				Name:       "rs1",
			},
			expected: "/content/overview/workloads/replica-sets/rs1",
		},
		{
			name: "replication controller",
			objectReference: corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "ReplicationController",
				Name:       "rc1",
			},
			expected: "/content/overview/workloads/replication-controllers/rc1",
		},
		{
			name: "stateful set",
			objectReference: corev1.ObjectReference{
				APIVersion: "apps/v1",
				Kind:       "StatefulSet",
				Name:       "ss1",
			},
			expected: "/content/overview/workloads/stateful-sets/ss1",
		},
		{
			name: "ingress",
			objectReference: corev1.ObjectReference{
				APIVersion: "extensions/v1beta1",
				Kind:       "Ingress",
				Name:       "i1",
			},
			expected: "/content/overview/discovery-and-load-balancing/ingresses/i1",
		},
		{
			name: "service",
			objectReference: corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "Service",
				Name:       "s1",
			},
			expected: "/content/overview/discovery-and-load-balancing/services/s1",
		},
		{
			name: "config map",
			objectReference: corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Name:       "cm1",
			},
			expected: "/content/overview/config-and-storage/config-maps/cm1",
		},
		{
			name: "persistent volume claim",
			objectReference: corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "PersistentVolumeClaim",
				Name:       "pvc1",
			},
			expected: "/content/overview/config-and-storage/persistent-volume-claims/pvc1",
		},
		{
			name: "secret",
			objectReference: corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "Secret",
				Name:       "s1",
			},
			expected: "/content/overview/config-and-storage/secrets/s1",
		},
		{
			name: "service account",
			objectReference: corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "ServiceAccount",
				Name:       "sa1",
			},
			expected: "/content/overview/config-and-storage/service-accounts/sa1",
		},
		{
			name: "role",
			objectReference: corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "Role",
				Name:       "r1",
			},
			expected: "/content/overview/rbac/roles/r1",
		},
		{
			name: "role binding",
			objectReference: corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "RoleBinding",
				Name:       "rb1",
			},
			expected: "/content/overview/rbac/role-bindings/rb1",
		},
		{
			name: "event",
			objectReference: corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "Event",
				Name:       "e1",
			},
			expected: "/content/overview/events/e1",
		},
		{
			name: "invalid",
			objectReference: corev1.ObjectReference{
				APIVersion: "v2",
				Kind:       "Event",
				Name:       "e1",
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ObjectReferencePath(tc.objectReference)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

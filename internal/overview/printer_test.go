package overview

import (
	"fmt"
	"testing"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/kubernetes/pkg/apis/batch"
)

func Test_printCronJobSummary(t *testing.T) {
	object := loadFromFile(t, "cronjob-1.yaml")
	var cronJob *batch.CronJob
	cronJob, ok := convertToInternal(t, object).(*batch.CronJob)
	require.True(t, ok)
	jobs := []*batch.Job{}

	got, err := printCronJobSummary(cronJob, jobs)
	require.NoError(t, err)

	expected := content.NewSection()
	expected.AddText("Name", "hello")
	expected.AddText("Namespace", "default")
	expected.AddLabels("Labels", map[string]string{"overview": "default"})
	expected.AddText("Annotations", "<none>")
	expected.AddTimestamp("Create Time", "2018-09-18T12:30:09Z")
	expected.AddText("Active", "0")
	expected.AddText("Schedule", "*/1 * * * *")
	expected.AddText("Suspend", "false")
	expected.AddTimestamp("Last Schedule", "2018-11-02T09:45:00Z")
	expected.AddText("Concurrency Policy", "Allow")
	expected.AddText("Starting Deadline Seconds", "<unset>")

	assert.Equal(t, expected, got)
}

func Test_gvkPath(t *testing.T) {
	cases := []struct {
		apiVersion string
		kind       string
		name       string
		expected   string
		isPanic    bool
	}{
		{
			apiVersion: "apps/v1",
			kind:       "DaemonSet",
			name:       "name",
			expected:   "/content/overview/workloads/daemon-sets/name",
		},
		{
			apiVersion: "extensions/v1beta1",
			kind:       "ReplicaSet",
			name:       "name",
			expected:   "/content/overview/workloads/replica-sets/name",
		},
		{
			apiVersion: "extensions/v1beta1",
			kind:       "Deployment",
			name:       "name",
			expected:   "/content/overview/workloads/deployments/name",
		},
		{
			apiVersion: "apps/v1",
			kind:       "StatefulSet",
			name:       "name",
			expected:   "/content/overview/workloads/stateful-sets/name",
		},
		{
			apiVersion: "batch/v1beta1",
			kind:       "CronJob",
			name:       "name",
			expected:   "/content/overview/workloads/cron-jobs/name",
		},
		{
			apiVersion: "batch/v1beta1",
			kind:       "Job",
			name:       "name",
			expected:   "/content/overview/workloads/jobs/name",
		},
		{
			apiVersion: "v1",
			kind:       "ReplicationController",
			name:       "name",
			expected:   "/content/overview/workloads/replication-controllers/name",
		},
		{
			apiVersion: "v1",
			kind:       "Service",
			name:       "name",
			expected:   "/content/overview/discovery-and-load-balancing/services/name",
		},
		{
			apiVersion: "v1",
			kind:       "Secret",
			name:       "name",
			expected:   "/content/overview/config-and-storage/secrets/name",
		},
		{
			apiVersion: "v1",
			kind:       "ServiceAccount",
			name:       "name",
			expected:   "/content/overview/config-and-storage/service-accounts/name",
		},
		{
			apiVersion: "rbac.authorization.k8s.io/v1",
			kind:       "Role",
			name:       "name",
			expected:   "/content/overview/rbac/roles/name",
		},
		{
			apiVersion: "unknown",
			kind:       "unknown",
			name:       "name",
			isPanic:    true,
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("apiVersion:%q kind:%q", tc.apiVersion, tc.kind)
		t.Run(name, func(t *testing.T) {
			if tc.isPanic {
				assert.Panics(t, func() {
					gvkPath(tc.apiVersion, tc.kind, tc.name)
				})
				return
			}

			assert.NotPanics(t, func() {
				got := gvkPath(tc.apiVersion, tc.kind, tc.name)
				assert.Equal(t, tc.expected, got)
			})
		})
	}
}

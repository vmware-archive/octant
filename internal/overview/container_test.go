package overview

import (
	"context"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func TestContainerSummary_invalid_object(t *testing.T) {
	assertViewInvalidObject(t, NewContainerSummary("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestContainerSummary(t *testing.T) {
	ctx := context.Background()
	object := loadFromFile(t, "cronjob-1.yaml")
	var cronJob *batch.CronJob
	cronJob, ok := convertToInternal(t, object).(*batch.CronJob)
	require.True(t, ok)
	cache := NewMemoryCache()

	v := NewContainerSummary("prefix", "ns", clock.NewFakeClock(time.Now()))

	got, err := v.Content(ctx, cronJob, cache)
	require.NoError(t, err)

	podTemplate := content.NewSummary("Pod Template", []content.Section{
		{
			Items: []content.Item{
				content.TextItem("Labels", "<none>"),
			},
		},
	})

	containerTemplate := content.NewSummary("Container Template", []content.Section{
		{
			Title: "hello",
			Items: []content.Item{
				content.TextItem("Image", "busybox"),
				content.TextItem("Port", "<none>"),
				content.TextItem("Host Port", "<none>"),
				content.TextItem("Args", "[/bin/sh, -c, date; echo Hello from the Kubernetes cluster]"),
				content.TextItem("Environment", "<none>"),
				content.ListItem("Mounts", map[string]string{}),
			},
		},
	})

	expected := []content.Content{
		&podTemplate,
		&containerTemplate,
	}

	assert.Equal(t, expected, got)
}

func Test_podTemplateSpec(t *testing.T) {
	pts := core.PodTemplateSpec{}

	cases := []struct {
		name     string
		object   runtime.Object
		expected *core.PodTemplateSpec
		isErr    bool
	}{
		{
			name: "cronjob",
			object: &batch.CronJob{
				Spec: batch.CronJobSpec{
					JobTemplate: batch.JobTemplateSpec{
						Spec: batch.JobSpec{
							Template: pts}}}},
			expected: &pts,
		},
		{
			name: "daemonset",
			object: &extensions.DaemonSet{
				Spec: extensions.DaemonSetSpec{
					Template: pts},
			},
			expected: &pts,
		},
		{
			name: "deployment",
			object: &extensions.Deployment{
				Spec: extensions.DeploymentSpec{
					Template: pts},
			},
			expected: &pts,
		},
		{
			name: "job",
			object: &batch.Job{
				Spec: batch.JobSpec{
					Template: pts},
			},
			expected: &pts,
		},
		{
			name: "replicaset",
			object: &extensions.ReplicaSet{
				Spec: extensions.ReplicaSetSpec{
					Template: pts},
			},
			expected: &pts,
		},
		{
			name: "replication controller",
			object: &core.ReplicationController{
				Spec: core.ReplicationControllerSpec{
					Template: &pts},
			},
			expected: &pts,
		},
		{
			name: "stateful set",
			object: &apps.StatefulSet{
				Spec: apps.StatefulSetSpec{
					Template: pts},
			},
			expected: &pts,
		},
		{
			name:  "nil object",
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := podTemplateSpec(tc.object)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

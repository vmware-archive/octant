package overview

import (
	"context"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestContainerSummary_invalid_object(t *testing.T) {
	assertViewInvalidObject(t, NewContainerSummary("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestContainerSummary(t *testing.T) {
	ctx := context.Background()
	object := loadFromFile(t, "cronjob-1.yaml")
	cronJob, ok := object.(*batchv1beta1.CronJob)
	require.True(t, ok)
	c := cache.NewMemoryCache()

	v := NewContainerSummary("prefix", "ns", clock.NewFakeClock(time.Now()))

	got, err := v.Content(ctx, cronJob, c)
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
	pts := corev1.PodTemplateSpec{}

	cases := []struct {
		name     string
		object   runtime.Object
		expected *corev1.PodTemplateSpec
		isErr    bool
	}{
		{
			name: "cronjob",
			object: &batchv1beta1.CronJob{
				Spec: batchv1beta1.CronJobSpec{
					JobTemplate: batchv1beta1.JobTemplateSpec{
						Spec: batchv1.JobSpec{
							Template: pts,
						},
					},
				},
			},
			expected: &pts,
		},
		{
			name: "daemonset",
			object: &appsv1.DaemonSet{
				Spec: appsv1.DaemonSetSpec{
					Template: pts},
			},
			expected: &pts,
		},
		{
			name: "deployment",
			object: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: pts},
			},
			expected: &pts,
		},
		{
			name: "job",
			object: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Template: pts},
			},
			expected: &pts,
		},
		{
			name: "replicaset",
			object: &appsv1.ReplicaSet{
				Spec: appsv1.ReplicaSetSpec{
					Template: pts},
			},
			expected: &pts,
		},
		{
			name: "replication controller",
			object: &corev1.ReplicationController{
				Spec: corev1.ReplicationControllerSpec{
					Template: &pts},
			},
			expected: &pts,
		},
		{
			name: "stateful set",
			object: &appsv1.StatefulSet{
				Spec: appsv1.StatefulSetSpec{
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

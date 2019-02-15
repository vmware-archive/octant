package printer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_ReplicationControllerListHandler(t *testing.T) {
	printOptions := Options{
		Cache: cache.NewMemoryCache(),
	}

	got, err := ReplicationControllerListHandler(validReplicationControllerList, printOptions)
	require.NoError(t, err)

	containers := component.NewContainers()
	containers.Add("nginx", "nginx:1.15")

	cols := component.NewTableCols("Name", "Labels", "Status", "Age", "Containers", "Selector")
	expected := component.NewTable("ReplicationControllers", cols)
	expected.Add(component.TableRow{
		"Name":       component.NewLink("", "rc-test", "/content/overview/namespace/default/workloads/replication-controllers/rc-test"),
		"Labels":     component.NewLabels(validReplicationControllerLabels),
		"Status":     component.NewText("0/3"),
		"Age":        component.NewTimestamp(validReplicationControllerCreationTime),
		"Containers": containers,
		"Selector":   component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "myapp")}),
	})

	assert.Equal(t, expected, got)
}

var (
	validReplicationControllerLabels = map[string]string{
		"foo": "bar",
	}

	validReplicationControllerCreationTime = time.Unix(1547211430, 0)

	validReplicationController = &corev1.ReplicationController{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ReplicationController",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rc-test",
			Namespace: "default",
			CreationTimestamp: metav1.Time{
				Time: validReplicationControllerCreationTime,
			},
			Labels: validReplicationControllerLabels,
		},
		Status: corev1.ReplicationControllerStatus{
			Replicas:          3,
			AvailableReplicas: 0,
		},
		Spec: corev1.ReplicationControllerSpec{
			Selector: map[string]string{
				"app": "myapp",
			},
			Template: &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "nginx",
							Image: "nginx:1.15",
						},
					},
				},
			},
		},
	}

	validReplicationControllerList = &corev1.ReplicationControllerList{
		Items: []corev1.ReplicationController{
			*validReplicationController,
		},
	}
)

func TestReplicationControllerStatus(t *testing.T) {
	c := cache.NewMemoryCache()

	pods := &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx-g7f72",
					Namespace: "default",
					Labels:    validReplicationControllerLabels,
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx-p64jr",
					Namespace: "default",
					Labels:    validReplicationControllerLabels,
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx-x8nrk",
					Namespace: "testing",
					Labels:    validReplicationControllerLabels,
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
				},
			},
		},
	}

	for _, p := range pods.Items {
		u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(p)
		if err != nil {
			return
		}

		c.Store(&unstructured.Unstructured{Object: u})
	}

	rcs := NewReplicationControllerStatus(validReplicationController)
	got, err := rcs.Create(c)
	require.NoError(t, err)

	expected := component.NewQuadrant()
	require.NoError(t, expected.Set(component.QuadNW, "Running", "3"))
	require.NoError(t, expected.Set(component.QuadNE, "Waiting", "0"))
	require.NoError(t, expected.Set(component.QuadSW, "Succeeded", "0"))
	require.NoError(t, expected.Set(component.QuadSE, "Failed", "0"))

	assert.Equal(t, expected, got)
}

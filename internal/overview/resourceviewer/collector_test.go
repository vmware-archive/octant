package resourceviewer

import (
	"testing"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func Test_Collector(t *testing.T) {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment",
			UID:  types.UID("deployment"),
		},
		Status: appsv1.DeploymentStatus{
			Replicas:          1,
			AvailableReplicas: 1,
		},
	}

	replicaSet1 := &extv1beta1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{APIVersion: "extensions/v1beta1", Kind: "ReplicaSet"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "replicaSet1",
			UID:  types.UID("replicaSet1"),
		},
		Spec: extv1beta1.ReplicaSetSpec{
			Replicas: ptrInt32(1),
		},
		Status: extv1beta1.ReplicaSetStatus{
			Replicas:          1,
			AvailableReplicas: 1,
		},
	}

	replicaSet2 := &extv1beta1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{APIVersion: "extensions/v1beta1", Kind: "ReplicaSet"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "replicaSet2",
			UID:  types.UID("replicaSet2"),
		},
	}

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod",
			UID:  types.UID("pod"),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(replicaSet1,
					schema.FromAPIVersionAndKind(replicaSet1.APIVersion,
						replicaSet1.Kind)),
			},
		},
	}

	c := NewCollector()

	err := c.Process(deployment)
	require.NoError(t, err)

	err = c.Process(replicaSet1)
	require.NoError(t, err)

	err = c.Process(replicaSet2)
	require.NoError(t, err)

	err = c.Process(pod)
	require.NoError(t, err)

	err = c.AddChild(deployment, replicaSet1, replicaSet2)
	require.NoError(t, err)

	err = c.AddChild(replicaSet1, pod)
	require.NoError(t, err)

	got, err := c.ViewComponent()
	require.NoError(t, err)

	expected := &component.ResourceViewer{
		Metadata: component.Metadata{
			Type:  "resourceViewer",
			Title: component.Title(component.NewText("Resource Viewer")),
		},
		Config: component.ResourceViewerConfig{
			Edges: component.AdjList{
				"deployment": []component.Edge{
					{
						Node: "replicaSet1",
						Type: "explicit",
					},
				},
				"replicaSet1": []component.Edge{
					{
						Node: "pods-replicaSet1",
						Type: "explicit",
					},
				},
			},
			Nodes: component.Nodes{
				"deployment": component.Node{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "deployment",
					Status:     "ok",
					Details:    component.Title(component.NewText("Deployment is OK")),
				},
				"replicaSet1": component.Node{
					APIVersion: "extensions/v1beta1",
					Kind:       "ReplicaSet",
					Name:       "replicaSet1",
					Status:     "ok",
					Details:    component.Title(component.NewText("Replica Set is OK")),
				},
				"pods-replicaSet1": component.Node{
					APIVersion: "v1",
					Kind:       "Pod",
					Name:       "replicaSet1 pods",
					Status:     "ok",
					Details:    component.Title(component.NewText("Pod count: 1")),
				},
			},
		},
	}

	assert.Equal(t, expected, got)
}

func ptrInt32(i int32) *int32 {
	return &i
}

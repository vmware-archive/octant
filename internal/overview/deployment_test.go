package overview

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func TestDeploymentSummary_InvalidObject(t *testing.T) {
	ds := NewDeploymentSummary()
	ctx := context.Background()

	object := &unstructured.Unstructured{}

	_, err := ds.Content(ctx, object, nil)
	require.Error(t, err)
}

func TestDeploymentSummary(t *testing.T) {
	ds := NewDeploymentSummary()

	ctx := context.Background()

	rhl := int32(0)

	object := &extensions.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment",
			Namespace: "default",
			Labels: map[string]string{
				"app": "myapp",
			},
			CreationTimestamp: metav1.Time{
				Time: time.Unix(1539603521, 0),
			},
		},
		Spec: extensions.DeploymentSpec{
			MinReadySeconds:      5,
			RevisionHistoryLimit: &rhl,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "myapp",
				},
			},
			Strategy: extensions.DeploymentStrategy{
				RollingUpdate: &extensions.RollingUpdateDeployment{
					MaxSurge:       intstr.FromString("25%"),
					MaxUnavailable: intstr.FromString("25%"),
				},
				Type: extensions.RollingUpdateDeploymentStrategyType,
			},
		},
		Status: extensions.DeploymentStatus{
			UpdatedReplicas:     1,
			Replicas:            2,
			AvailableReplicas:   3,
			UnavailableReplicas: 4,
		},
	}

	cache := newSpyCache()

	contents, err := ds.Content(ctx, object, cache)
	require.NoError(t, err)

	summary := content.NewSummary("Details", []content.Section{
		{
			Items: []content.Item{
				content.TextItem("Name", "deployment"),
				content.TextItem("Namespace", "default"),
				content.LabelsItem("Labels", map[string]string{"app": "myapp"}),
				content.LabelsItem("Annotations", map[string]string{}),
				content.TextItem("Creation Time", "Mon, 15 Oct 2018 11:38:41 +0000"),
				content.TextItem("Selector", "app=myapp"),
				content.TextItem("Strategy", "RollingUpdate"),
				content.TextItem("Min Ready Seconds", "5"),
				content.TextItem("Revision History Limit", "0"),
				content.TextItem("Rolling Update Strategy", "Max Surge: 25%, Max unavailable: 25%"),
				content.TextItem("Status", "1 updated, 2 total, 3 available, 4 unavailable"),
			},
		},
	})

	expected := []content.Content{
		&summary,
	}

	assert.Equal(t, expected, contents)
}

func TestDeploymentReplicaSets(t *testing.T) {
	drs := NewDeploymentReplicaSets()

	ctx := context.Background()

	rhl := int32(0)

	object := &extensions.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment",
			Namespace: "default",
			Labels: map[string]string{
				"app": "myapp",
			},
			CreationTimestamp: metav1.Time{
				Time: time.Unix(1539603521, 0),
			},
			UID: types.UID("ac833d23-c17e-11e8-9212-025000000001"),
		},
		Spec: extensions.DeploymentSpec{
			MinReadySeconds:      5,
			RevisionHistoryLimit: &rhl,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "myapp",
				},
			},
			Strategy: extensions.DeploymentStrategy{
				RollingUpdate: &extensions.RollingUpdateDeployment{
					MaxSurge:       intstr.FromString("25%"),
					MaxUnavailable: intstr.FromString("25%"),
				},
				Type: extensions.RollingUpdateDeploymentStrategyType,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "myapp",
					},
				},

				Spec: core.PodSpec{

					SecurityContext: &core.PodSecurityContext{},
					Containers: []core.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.13.6",
							Ports: []core.ContainerPort{
								{
									Protocol:      "TCP",
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
		Status: extensions.DeploymentStatus{
			UpdatedReplicas:     1,
			Replicas:            2,
			AvailableReplicas:   3,
			UnavailableReplicas: 4,
		},
	}

	cache := NewMemoryCache()

	storeFromFile(t, "replicaset-1.yaml", cache)

	contents, err := drs.Content(ctx, object, cache)
	require.NoError(t, err)

	replicaSetColumns := []content.TableColumn{
		tableCol("Name"),
		tableCol("Desired"),
		tableCol("Current"),
		tableCol("Ready"),
		tableCol("Age"),
		tableCol("Containers"),
		tableCol("Images"),
		tableCol("Selector"),
		tableCol("Labels"),
	}

	newReplicaSetTable := content.NewTable("New Replica Set")
	newReplicaSetTable.Columns = replicaSetColumns

	newReplicaSetTable.AddRow(
		content.TableRow{
			"Name":       content.NewLinkText("rs1", "/content/overview/workloads/replica-sets/rs1"),
			"Desired":    content.NewStringText("3"),
			"Current":    content.NewStringText("3"),
			"Ready":      content.NewStringText("3"),
			"Age":        content.NewStringText("2d"),
			"Containers": content.NewStringText("nginx"),
			"Images":     content.NewStringText("nginx:1.13.6"),
			"Selector":   content.NewStringText("app=myapp,pod-template-hash=2350241137"),
			"Labels":     content.NewStringText("app=myapp,pod-template-hash=2350241137"),
		},
	)

	oldReplicaSetsTable := content.NewTable("Old Replica Sets")
	oldReplicaSetsTable.Columns = replicaSetColumns

	expected := []content.Content{
		&newReplicaSetTable,
		&oldReplicaSetsTable,
	}

	assert.Equal(t, expected, contents)
}

func storeFromFile(t *testing.T, name string, cache Cache) {
	decoded := loadFromFile(t, name)

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(decoded)
	require.NoError(t, err)

	err = cache.Store(&unstructured.Unstructured{Object: m})
	require.NoError(t, err)
}

func loadFromFile(t *testing.T, name string) runtime.Object {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	rs1, err := ioutil.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)

	decoded, _, err := decode(rs1, nil, nil)
	require.NoError(t, err)

	return decoded
}

func convertToInternal(t *testing.T, in runtime.Object) runtime.Object {
	var out runtime.Object

	switch in.(type) {
	case *extensionsv1beta1.ReplicaSet:
		out = &extensions.ReplicaSet{}
	}

	err := scheme.Scheme.Convert(in, out, runtime.InternalGroupVersioner)
	require.NoError(t, err)

	return out
}

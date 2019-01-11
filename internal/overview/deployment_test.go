package overview

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestDeploymentSummary_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewDeploymentSummary("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestDeploymentSummary(t *testing.T) {
	ds := NewDeploymentSummary("prefix", "ns", clock.NewFakeClock(time.Now()))

	ctx := context.Background()

	rhl := int32(0)

	twentyFivePercent := intstr.FromString("25%")
	object := &appsv1.Deployment{
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
		Spec: appsv1.DeploymentSpec{
			MinReadySeconds:      5,
			RevisionHistoryLimit: &rhl,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "myapp",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       &twentyFivePercent,
					MaxUnavailable: &twentyFivePercent,
				},
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "containerName",
							Image: "image",
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{
			UpdatedReplicas:     1,
			Replicas:            2,
			AvailableReplicas:   3,
			UnavailableReplicas: 4,
		},
	}

	cache := newSpyCache()

	contents, err := ds.Content(ctx, object, cache)
	require.NoError(t, err)

	details := content.NewSummary("Details", []content.Section{
		{
			Items: []content.Item{
				content.TextItem("Name", "deployment"),
				content.TextItem("Namespace", "default"),
				content.LabelsItem("Labels", map[string]string{"app": "myapp"}),
				content.LabelsItem("Annotations", map[string]string{}),
				content.TimeItem("Creation Time", "2018-10-15T11:38:41Z"),
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
		&details,
	}

	assert.Equal(t, expected, contents)
}

func TestDeploymentReplicaSets(t *testing.T) {
	drs := NewDeploymentReplicaSets("prefix", "ns", clock.NewFakeClock(time.Now()))

	ctx := context.Background()

	rhl := int32(0)

	twentyFivePercent := intstr.FromString("25%")
	object := &appsv1.Deployment{
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
		Spec: appsv1.DeploymentSpec{
			MinReadySeconds:      5,
			RevisionHistoryLimit: &rhl,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "myapp",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       &twentyFivePercent,
					MaxUnavailable: &twentyFivePercent,
				},
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "myapp",
					},
				},

				Spec: corev1.PodSpec{

					SecurityContext: &corev1.PodSecurityContext{},
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.13.6",
							Ports: []corev1.ContainerPort{
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
		Status: appsv1.DeploymentStatus{
			UpdatedReplicas:     1,
			Replicas:            2,
			AvailableReplicas:   3,
			UnavailableReplicas: 4,
		},
	}

	cache := cache.NewMemoryCache()

	storeFromFile(t, "replicaset-1.yaml", cache)

	contents, err := drs.Content(ctx, object, cache)
	require.NoError(t, err)

	replicaSetColumns := []content.TableColumn{
		view.TableCol("Name"),
		view.TableCol("Desired"),
		view.TableCol("Current"),
		view.TableCol("Ready"),
		view.TableCol("Age"),
		view.TableCol("Containers"),
		view.TableCol("Images"),
		view.TableCol("Selector"),
		view.TableCol("Labels"),
	}

	newReplicaSetTable := content.NewTable("New Replica Set", "This Deployment does not have a current Replica")
	newReplicaSetTable.Columns = replicaSetColumns

	newReplicaSetTable.AddRow(
		content.TableRow{
			"Name":       content.NewLinkText("rs1", "/content/overview/workloads/replica-sets/rs1"),
			"Desired":    content.NewStringText("3"),
			"Current":    content.NewStringText("3"),
			"Ready":      content.NewStringText("3"),
			"Age":        content.NewStringText("24h"),
			"Containers": content.NewStringText("nginx"),
			"Images":     content.NewStringText("nginx:1.13.6"),
			"Selector":   content.NewStringText("app=myapp,pod-template-hash=2350241137"),
			"Labels":     content.NewStringText("app=myapp,pod-template-hash=2350241137"),
		},
	)

	oldReplicaSetsTable := content.NewTable("Old Replica Sets", "This Deployment does not have any old Replicas")
	oldReplicaSetsTable.Columns = replicaSetColumns

	expected := []content.Content{
		&newReplicaSetTable,
		&oldReplicaSetsTable,
	}

	assert.Equal(t, expected, contents)
}

func storeFromFile(t *testing.T, name string, cache cache.Cache) {
	decoded := loadFromFile(t, name)
	obj, err := meta.Accessor(decoded)
	if err != nil {
		t.Fatalf("could not create meta.Accessor")
	}
	// Override objects to be "1d" old, makes comparing in tests easier
	obj.SetCreationTimestamp(metav1.NewTime(time.Now().AddDate(0, 0, -1)))

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
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

/*
func convertToInternal(t *testing.T, in runtime.Object) runtime.Object {
	var out runtime.Object

	switch in.(type) {
	case *corev1.ConfigMap:
		out = &core.ConfigMap{}
	case *batchv1beta1.CronJob:
		out = &batch.CronJob{}
	case *extensionsv1beta1.Ingress:
		out = &extensionsv1beta1.Ingress{}
	case *extensionsv1beta1.ReplicaSet:
		out = &apps.ReplicaSet{}
	case *corev1.Secret:
		out = &core.Secret{}
	case *corev1.Pod:
		out = &core.Pod{}
	case *corev1.PodList:
		out = &core.PodList{}
	case *corev1.Service:
		out = &core.Service{}
	case *appsv1.Deployment:
		out = &apps.Deployment{}
	default:
		t.Fatalf("don't know how to convert %T to internal", in)
	}

	err := scheme.Scheme.Convert(in, out, runtime.InternalGroupVersioner)
	require.NoError(t, err)

	return out
}
*/

package testutil

import (
	"time"

	"github.com/heptio/developer-dash/internal/overview/objectvisitor"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func CreateDeployment(name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta:   genTypeMeta(objectvisitor.DeploymentGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

func genTypeMeta(gvk schema.GroupVersionKind) metav1.TypeMeta {
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	return metav1.TypeMeta{
		APIVersion: apiVersion,
		Kind:       kind,
	}
}

func genObjectMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:              name,
		Namespace:         "namespace",
		UID:               types.UID(name),
		CreationTimestamp: metav1.Time{Time: time.Now()},
	}
}

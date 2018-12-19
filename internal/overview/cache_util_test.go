package overview

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func newScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Ingress"}, &v1beta1.Ingress{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}, &corev1.Pod{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"}, &corev1.Service{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}, &corev1.Secret{})
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}, &appsv1.Deployment{})
	return scheme
}

func newUnstructured(apiVersion, kind, namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"namespace": namespace,
				"name":      name,
			},
		},
	}
}

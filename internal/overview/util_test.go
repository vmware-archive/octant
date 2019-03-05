package overview

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var resources = []*metav1.APIResourceList{
	{
		GroupVersion: "apps/v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "deployments",
				SingularName: "deployment",
				Group:        "apps",
				Version:      "v1",
				Kind:         "Deployment",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
	{
		GroupVersion: "extensions/v1beta1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "ingresses",
				SingularName: "ingress",
				Group:        "extensions",
				Version:      "v1beta1",
				Kind:         "Ingress",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
	{
		GroupVersion: "v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "services",
				SingularName: "service",
				Group:        "",
				Version:      "v1",
				Kind:         "Service",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
			metav1.APIResource{
				Name:         "secrets",
				SingularName: "secret",
				Group:        "",
				Version:      "v1",
				Kind:         "Secret",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
	{
		GroupVersion: "bar/v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "bars",
				SingularName: "bar",
				Group:        "bar",
				Version:      "v1",
				Kind:         "Bar",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
	{
		GroupVersion: "foo/v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "kinds",
				SingularName: "kind",
				Group:        "foo",
				Version:      "v1",
				Kind:         "Kind",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
			metav1.APIResource{
				Name:         "foos",
				SingularName: "foo",
				Group:        "foo",
				Version:      "v1",
				Kind:         "Foo",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
			metav1.APIResource{
				Name:         "others",
				SingularName: "other",
				Group:        "foo",
				Version:      "v1",
				Kind:         "Other",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
}

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

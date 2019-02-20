package testutil

import (
	"github.com/heptio/developer-dash/internal/conversion"
	"github.com/heptio/developer-dash/internal/gvk"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateDaemonSet creates a daemon set
func CreateDaemonSet(name string) *appsv1.DaemonSet {
	maxUnavailable := intstr.FromInt(1)

	return &appsv1.DaemonSet{
		TypeMeta:   genTypeMeta(gvk.DaemonSetGVK),
		ObjectMeta: genObjectMeta(name),
		Spec: appsv1.DaemonSetSpec{
			RevisionHistoryLimit: conversion.PtrInt32(10),
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				RollingUpdate: &appsv1.RollingUpdateDaemonSet{
					MaxUnavailable: &maxUnavailable,
				},
			},
		},
		Status: appsv1.DaemonSetStatus{
			CurrentNumberScheduled: 1,
			DesiredNumberScheduled: 1,
			NumberAvailable:        1,
			NumberReady:            1,
			UpdatedNumberScheduled: 1,
		},
	}
}

// CreateDeployment creates a deployment
func CreateDeployment(name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta:   genTypeMeta(gvk.DeploymentGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

// CreateIngress creates an ingress
func CreateIngress(name string) *extv1beta1.Ingress {
	return &extv1beta1.Ingress{
		TypeMeta:   genTypeMeta(gvk.IngressGVK),
		ObjectMeta: genObjectMeta(name),
		Spec: extv1beta1.IngressSpec{
			Backend: &extv1beta1.IngressBackend{
				ServiceName: "app",
				ServicePort: intstr.FromInt(80),
			},
		},
	}
}

// CreatePod creates a pod
func CreatePod(name string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta:   genTypeMeta(gvk.PodGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

// CreateReplicaSet creates a replica set
func CreateReplicaSet(name string) *appsv1.ReplicaSet {
	return &appsv1.ReplicaSet{
		TypeMeta:   genTypeMeta(gvk.ReplicaSetGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

// CreateSer4vice creates a service
func CreateService(name string) *corev1.Service {
	return &corev1.Service{
		TypeMeta:   genTypeMeta(gvk.ServiceGVK),
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
		Name:      name,
		Namespace: "namespace",
		UID:       types.UID(name),
	}
}

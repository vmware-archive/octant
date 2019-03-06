package testutil

import (
	"github.com/heptio/developer-dash/internal/conversion"
	"github.com/heptio/developer-dash/internal/gvk"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateClusterRoleBinding creates a cluster role binding
func CreateClusterRoleBinding(name string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta:   genTypeMeta(gvk.ClusterRoleBindingGVK),
		ObjectMeta: genObjectMeta(name),
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     "role-name",
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: "test@example.com",
			},
		},
	}
}

// CreateCRD creates a CRD
func CreateCRD(name string) *apiextv1beta1.CustomResourceDefinition {
	return &apiextv1beta1.CustomResourceDefinition{
		TypeMeta:   genTypeMeta(gvk.CustomResourceDefinitionGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

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

// CreateReplicationController creates a replication controller
func CreateReplicationController(name string) *corev1.ReplicationController {
	return &corev1.ReplicationController{
		TypeMeta:   genTypeMeta(gvk.ReplicationControllerGVK),
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

// CreateSecret creates a secret
func CreateSecret(name string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta:   genTypeMeta(gvk.SecretGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

// CreateService creates a service
func CreateService(name string) *corev1.Service {
	return &corev1.Service{
		TypeMeta:   genTypeMeta(gvk.ServiceGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

// CreateServiceAccount creates a service account
func CreateServiceAccount(name string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta:   genTypeMeta(gvk.ServiceAccountGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

// CreatePersistentVolumeClaim creates a persistent volume claim
func CreatePersistentVolumeClaim(name string) *corev1.PersistentVolumeClaim {
	storageClass := "manual"
	file := corev1.PersistentVolumeFilesystem

	return &corev1.PersistentVolumeClaim{
		TypeMeta:   genTypeMeta(gvk.PersistentVolumeClaimGVK),
		ObjectMeta: genObjectMeta(name),
		Spec: corev1.PersistentVolumeClaimSpec{
			VolumeName:       "task-pv-volume",
			StorageClassName: &storageClass,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("3Gi"),
				},
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			VolumeMode: &file,
		},
		Status: corev1.PersistentVolumeClaimStatus{
			Phase: corev1.ClaimBound,
			Capacity: corev1.ResourceList{
				corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("10Gi"),
			},
		},
	}
}

func CreateRole(name string) *rbacv1.Role {
	return &rbacv1.Role{
		TypeMeta:   genTypeMeta(gvk.RoleGVK),
		ObjectMeta: genObjectMeta(name),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods",
				},
				Verbs: []string{
					"get",
					"watch",
					"list",
				},
			},
		},
	}
}

func CreateRoleBindingSubject(kind, name string) *rbacv1.Subject {
	return &rbacv1.Subject{
		Kind: kind,
		Name: name,
	}
}

func CreateRoleBinding(roleBindingName, roleName string, subjects []rbacv1.Subject) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		TypeMeta:   genTypeMeta(gvk.RoleBindingGVK),
		ObjectMeta: genObjectMeta(roleBindingName),
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     roleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: subjects,
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

/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package testutil

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"

	"github.com/vmware/octant/internal/conversion"
	"github.com/vmware/octant/internal/gvk"
)

// DefaultNamespace is the namespace that objects will belong to.
const DefaultNamespace = "namespace"

// CreateClusterRoleBinding creates a cluster role binding
func CreateClusterRoleBinding(name, roleName string, subjects []rbacv1.Subject) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta:   genTypeMeta(gvk.ClusterRoleBindingGVK),
		ObjectMeta: genObjectMeta(name, false),
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     roleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: subjects,
	}
}

// CreateConfigMap creates a config map.
func CreateConfigMap(name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta:   genTypeMeta(gvk.ConfigMapGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateCRD creates a CRD
func CreateCRD(name string) *apiextv1beta1.CustomResourceDefinition {
	return &apiextv1beta1.CustomResourceDefinition{
		TypeMeta:   genTypeMeta(gvk.CustomResourceDefinitionGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateCustomResource creates a custom resource.
func CreateCustomResource(name string) *unstructured.Unstructured {
	m := map[string]interface{}{
		"apiVersion": "stable.example.com/v1",
		"kind":       "CronTab",
		"metadata": map[string]interface{}{
			"name": name,
		},
		"spec": map[string]interface{}{
			"cronSpec": "* * * * */5",
			"image":    "my-awesome-image",
		},
	}

	u := &unstructured.Unstructured{Object: m}
	u.SetNamespace(DefaultNamespace)

	return u
}

func CreateCronJob(name string) *batchv1beta1.CronJob {
	return &batchv1beta1.CronJob{
		TypeMeta:   genTypeMeta(gvk.CronJobGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateDaemonSet creates a daemon set
func CreateDaemonSet(name string) *appsv1.DaemonSet {
	maxUnavailable := intstr.FromInt(1)

	return &appsv1.DaemonSet{
		TypeMeta:   genTypeMeta(gvk.DaemonSetGVK),
		ObjectMeta: genObjectMeta(name, true),
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
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateEvent creates a event
func CreateEvent(name string) *corev1.Event {
	return &corev1.Event{
		TypeMeta:   genTypeMeta(gvk.Event),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateIngress creates an ingress
func CreateIngress(name string) *extv1beta1.Ingress {
	return &extv1beta1.Ingress{
		TypeMeta:   genTypeMeta(gvk.IngressGVK),
		ObjectMeta: genObjectMeta(name, true),
		Spec: extv1beta1.IngressSpec{
			Backend: &extv1beta1.IngressBackend{
				ServiceName: "app",
				ServicePort: intstr.FromInt(80),
			},
		},
	}
}

func CreateJob(name string) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta:   genTypeMeta(gvk.JobGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreatePod creates a pod
func CreatePod(name string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta:   genTypeMeta(gvk.PodGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateReplicationController creates a replication controller
func CreateReplicationController(name string) *corev1.ReplicationController {
	return &corev1.ReplicationController{
		TypeMeta:   genTypeMeta(gvk.ReplicationControllerGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateReplicaSet creates a replica set
func CreateReplicaSet(name string) *appsv1.ReplicaSet {
	return &appsv1.ReplicaSet{
		TypeMeta:   genTypeMeta(gvk.ReplicaSetGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateSecret creates a secret
func CreateSecret(name string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta:   genTypeMeta(gvk.SecretGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateService creates a service
func CreateService(name string) *corev1.Service {
	return &corev1.Service{
		TypeMeta:   genTypeMeta(gvk.ServiceGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateServiceAccount creates a service account
func CreateServiceAccount(name string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta:   genTypeMeta(gvk.ServiceAccountGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateStatefulSet creates a stateful set
func CreateStatefulSet(name string) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		TypeMeta:   genTypeMeta(gvk.StatefulSetGVK),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreatePersistentVolumeClaim creates a persistent volume claim
func CreatePersistentVolumeClaim(name string) *corev1.PersistentVolumeClaim {
	storageClass := "manual"
	file := corev1.PersistentVolumeFilesystem

	return &corev1.PersistentVolumeClaim{
		TypeMeta:   genTypeMeta(gvk.PersistentVolumeClaimGVK),
		ObjectMeta: genObjectMeta(name, true),
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

// CreateRole creates a role.
func CreateRole(name string) *rbacv1.Role {
	return &rbacv1.Role{
		TypeMeta:   genTypeMeta(gvk.RoleGVK),
		ObjectMeta: genObjectMeta(name, true),
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

// CreateClusterRole creates a cluster role.
func CreateClusterRole(name string) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta:   genTypeMeta(gvk.ClusterRoleGVK),
		ObjectMeta: genObjectMeta(name, false),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"stable.example.com",
				},
				Resources: []string{
					"crontabs",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"update",
					"patch",
					"delete",
				},
			},
		},
	}
}

func CreateRoleBindingSubject(kind, name, namespace string) *rbacv1.Subject {
	return &rbacv1.Subject{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
	}
}

func CreateRoleBinding(roleBindingName, roleName string, subjects []rbacv1.Subject) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		TypeMeta:   genTypeMeta(gvk.RoleBindingGVK),
		ObjectMeta: genObjectMeta(roleBindingName, true),
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

func genObjectMeta(name string, withNamespace bool) metav1.ObjectMeta {
	var namespace string
	if withNamespace {
		namespace = "namespace"
	}

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		UID:       types.UID(name),
	}
}

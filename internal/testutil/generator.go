/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package testutil

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"

	"github.com/vmware-tanzu/octant/internal/conversion"
	"github.com/vmware-tanzu/octant/internal/gvk"
)

// DefaultNamespace is the namespace that objects will belong to.
const DefaultNamespace = "namespace"

// CreateClusterRoleBinding creates a cluster role binding
func CreateClusterRoleBinding(name, roleName string, subjects []rbacv1.Subject) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta:   genTypeMeta(gvk.ClusterRoleBinding),
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
		TypeMeta:   genTypeMeta(gvk.ConfigMap),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CRDOption is an option for configuring CreateCRD.
type CRDOption func(definition *apiextv1beta1.CustomResourceDefinition)

// WithGenericCRD creates a crd with group/kind and one version
func WithGenericCRD() CRDOption {
	return func(crd *apiextv1beta1.CustomResourceDefinition) {
		crd.Spec.Group = "group"
		crd.Spec.Versions = []apiextv1beta1.CustomResourceDefinitionVersion{
			{
				Name:   "v1",
				Served: true,
			},
		}
		crd.Spec.Names.Kind = "kind"
	}
}

// CreateCRD creates a CRD
func CreateCRD(name string, options ...CRDOption) *apiextv1beta1.CustomResourceDefinition {
	crd := &apiextv1beta1.CustomResourceDefinition{
		TypeMeta:   genTypeMeta(gvk.CustomResourceDefinition),
		ObjectMeta: genObjectMeta(name, true),
	}

	for _, option := range options {
		option(crd)
	}

	return crd
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
		TypeMeta:   genTypeMeta(gvk.CronJob),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateDaemonSet creates a daemon set
func CreateDaemonSet(name string) *appsv1.DaemonSet {
	maxUnavailable := intstr.FromInt(1)

	return &appsv1.DaemonSet{
		TypeMeta:   genTypeMeta(gvk.DaemonSet),
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

// DeploymentOption is an option for configuration CreateDeployment.
type DeploymentOption func(d *appsv1.Deployment)

func WithGenericDeployment() DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.Spec.Template.Spec.Containers = []corev1.Container{
			{
				Name:  "container-name",
				Image: "image",
			},
		}
	}
}

// CreateDeployment creates a deployment
func CreateDeployment(name string, options ...DeploymentOption) *appsv1.Deployment {
	d := &appsv1.Deployment{
		TypeMeta:   genTypeMeta(gvk.Deployment),
		ObjectMeta: genObjectMeta(name, true),
	}

	for _, option := range options {
		option(d)
	}

	return d
}

// CreateEvent creates a event
func CreateEvent(name string) *corev1.Event {
	return &corev1.Event{
		TypeMeta:   genTypeMeta(gvk.Event),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateHorizontalPodAutoscaler creates a horizontal pod autoscaler
func CreateHorizontalPodAutoscaler(name string) *autoscalingv1.HorizontalPodAutoscaler {
	return &autoscalingv1.HorizontalPodAutoscaler{
		TypeMeta:   genTypeMeta(gvk.HorizontalPodAutoscaler),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateIngress creates an ingress
func CreateIngress(name string) *extv1beta1.Ingress {
	return &extv1beta1.Ingress{
		TypeMeta:   genTypeMeta(gvk.Ingress),
		ObjectMeta: genObjectMeta(name, true),
		Spec: extv1beta1.IngressSpec{
			Backend: &extv1beta1.IngressBackend{
				ServiceName: "app",
				ServicePort: intstr.FromInt(80),
			},
		},
	}
}

// CreateJob creates a job
func CreateJob(name string) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta:   genTypeMeta(gvk.Job),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateNamespace creates a namespace
func CreateNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta:   genTypeMeta(gvk.Namespace),
		ObjectMeta: genObjectMeta(name, false),
	}
}

// CreateNetworkPolicy creates a network policy
func CreateNetworkPolicy(name string) *networkingv1.NetworkPolicy {
	return &networkingv1.NetworkPolicy{
		TypeMeta:   genTypeMeta(gvk.Namespace),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateNode creates a node
func CreateNode(name string) *corev1.Node {
	return &corev1.Node{
		TypeMeta:   genTypeMeta(gvk.Node),
		ObjectMeta: genObjectMeta(name, false),
	}
}

// PodOption is an option for configuring CreatePod.
type PodOption func(*corev1.Pod)

// CreatePod creates a pod
func CreatePod(name string, options ...PodOption) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta:   genTypeMeta(gvk.Pod),
		ObjectMeta: genObjectMeta(name, true),
	}

	for _, option := range options {
		option(pod)
	}

	return pod
}

type PodMetricOption func(metrics *metricsv1beta1.PodMetrics)

func CreatePodMetrics(name string, options ...PodMetricOption) *metricsv1beta1.PodMetrics {
	m := &metricsv1beta1.PodMetrics{
		TypeMeta:   genTypeMeta(gvk.PodMetrics),
		ObjectMeta: genObjectMeta(name, true),
	}

	for _, option := range options {
		option(m)
	}

	return m
}

// CreateReplicationController creates a replication controller
func CreateReplicationController(name string) *corev1.ReplicationController {
	return &corev1.ReplicationController{
		TypeMeta:   genTypeMeta(gvk.ReplicationController),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateAppReplicaSet creates a replica set
func CreateAppReplicaSet(name string) *appsv1.ReplicaSet {
	return &appsv1.ReplicaSet{
		TypeMeta:   genTypeMeta(gvk.AppReplicaSet),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateExtReplicaSet creates a replica set
func CreateExtReplicaSet(name string) *extv1beta1.ReplicaSet {
	return &extv1beta1.ReplicaSet{
		TypeMeta:   genTypeMeta(gvk.ExtReplicaSet),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateSecret creates a secret
func CreateSecret(name string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta:   genTypeMeta(gvk.Secret),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateService creates a service
func CreateService(name string) *corev1.Service {
	return &corev1.Service{
		TypeMeta:   genTypeMeta(gvk.Service),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateServiceAccount creates a service account
func CreateServiceAccount(name string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta:   genTypeMeta(gvk.ServiceAccount),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreateStatefulSet creates a stateful set
func CreateStatefulSet(name string) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		TypeMeta:   genTypeMeta(gvk.StatefulSet),
		ObjectMeta: genObjectMeta(name, true),
	}
}

// CreatePersistentVolumeClaim creates a persistent volume claim
func CreatePersistentVolumeClaim(name string) *corev1.PersistentVolumeClaim {
	storageClass := "manual"
	file := corev1.PersistentVolumeFilesystem

	return &corev1.PersistentVolumeClaim{
		TypeMeta:   genTypeMeta(gvk.PersistentVolumeClaim),
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

// CreatePersistentVolume creates a persistent volume
func CreatePersistentVolume(name string) *corev1.PersistentVolume {
	return &corev1.PersistentVolume{
		TypeMeta:   genTypeMeta(gvk.PersistentVolume),
		ObjectMeta: genObjectMeta(name, true),
		Spec:       corev1.PersistentVolumeSpec{},
		Status: corev1.PersistentVolumeStatus{
			Phase: corev1.VolumeBound,
		},
	}
}

// CreateRole creates a role.
func CreateRole(name string) *rbacv1.Role {
	return &rbacv1.Role{
		TypeMeta:   genTypeMeta(gvk.Role),
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
		TypeMeta:   genTypeMeta(gvk.ClusterRole),
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
		TypeMeta:   genTypeMeta(gvk.RoleBinding),
		ObjectMeta: genObjectMeta(roleBindingName, true),
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     roleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: subjects,
	}
}

func CreateTimestamp() *metav1.Time {
	return &metav1.Time{
		Time: Time(),
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

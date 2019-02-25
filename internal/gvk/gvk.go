package gvk

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	DaemonSetGVK             = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"}
	DeploymentGVK            = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	IngressGVK               = schema.GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Ingress"}
	ServiceAccountGVK        = schema.GroupVersionKind{Version: "v1", Kind: "ServiceAccount"}
	SecretGVK                = schema.GroupVersionKind{Version: "v1", Kind: "Secret"}
	ServiceGVK               = schema.GroupVersionKind{Version: "v1", Kind: "Service"}
	PodGVK                   = schema.GroupVersionKind{Version: "v1", Kind: "Pod"}
	ReplicaSetGVK            = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"}
	PersistentVolumeClaimGVK = schema.GroupVersionKind{Version: "v1", Kind: "PersistentVolumeClaim"}
	ReplicationControllerGVK = schema.GroupVersionKind{Version: "v1", Kind: "ReplicationController"}
	StatefulSetGVK           = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}
	RoleBindingGVK           = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"}
)

package testutil

import (
	"time"

	"github.com/heptio/developer-dash/internal/conversion"
	"github.com/heptio/developer-dash/internal/overview/objectvisitor"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateDeployment creates a deployment
func CreateDeployment(name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta:   genTypeMeta(objectvisitor.DeploymentGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

// CreateDaemonSet creates a daemon set
func CreateDaemonSet(name string) *appsv1.DaemonSet {
	maxUnavailable := intstr.FromInt(1)

	return &appsv1.DaemonSet{
		TypeMeta:   genTypeMeta(objectvisitor.DaemonSetGVK),
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

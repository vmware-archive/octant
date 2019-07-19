/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectvisitor_test

import (
	"context"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	clusterFake "github.com/vmware/octant/internal/cluster/fake"
	configFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/modules/overview/objectvisitor"
	ovFake "github.com/vmware/octant/internal/modules/overview/objectvisitor/fake"
	queryerFake "github.com/vmware/octant/internal/queryer/fake"
	"github.com/vmware/octant/internal/testutil"
)

func TestDefaultVisitor_Visit_use_default_visitor(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	resourceList := &metav1.APIResourceList{}
	resourceList.APIResources = []metav1.APIResource{
		{
			Kind:       "Deployment",
			Namespaced: true,
		},
		{
			Kind:       "ReplicaSet",
			Namespaced: true,
		},
	}

	discoveryClient := clusterFake.NewMockDiscoveryInterface(controller)
	discoveryClient.EXPECT().
		ServerResourcesForGroupVersion("apps/v1").Return(resourceList, nil).AnyTimes()

	clusterClient := clusterFake.NewMockClientInterface(controller)
	clusterClient.EXPECT().DiscoveryClient().Return(discoveryClient, nil).AnyTimes()

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().ClusterClient().Return(clusterClient).AnyTimes()

	deployment := testutil.CreateDeployment("deployment")

	replicaSet := testutil.CreateAppReplicaSet("replica-set")
	replicaSet.SetOwnerReferences(testutil.ToOwnerReferences(t, deployment))

	pod := testutil.CreatePod("pod")
	pod.SetOwnerReferences(testutil.ToOwnerReferences(t, replicaSet))

	ctx := context.Background()

	q := queryerFake.NewMockQueryer(controller)
	q.EXPECT().
		Children(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, object metav1.Object) ([]runtime.Object, error) {
			switch object.GetName() {
			case deployment.Name:
				return []runtime.Object{replicaSet}, nil
			case replicaSet.Name:
				return []runtime.Object{pod}, nil
			case pod.Name:
				return []runtime.Object{}, nil
			default:
				return nil, errors.New("can't retrieve children for object")
			}
		}).AnyTimes()

	q.EXPECT().
		OwnerReference(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, namespace string, ownerReference metav1.OwnerReference) (runtime.Object, error) {
			switch ownerReference.Name {
			case deployment.Name:
				return deployment, nil
			case replicaSet.Name:
				return replicaSet, nil
			default:
				return nil, errors.Errorf("owner reference for %s", ownerReference.Name)
			}
		}).AnyTimes()

	handler := ovFake.NewMockObjectHandler(controller)

	adjList := make(map[string][]string)
	accessor := meta.NewAccessor()

	handler.EXPECT().
		AddEdge(gomock.Any(), gomock.Any()).
		DoAndReturn(func(parent runtime.Object, children ...runtime.Object) error {
			parentName, err := accessor.Name(parent)
			require.NoError(t, err)

			cur := adjList[parentName]

			for _, child := range children {
				name, err := accessor.Name(child)
				require.NoError(t, err)

				cur = append(cur, name)
			}

			adjList[parentName] = cur

			return nil
		}).AnyTimes()

	var processed []string
	handler.EXPECT().
		Process(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, object runtime.Object) error {
			name, err := accessor.Name(object)
			require.NoError(t, err)
			processed = append(processed, name)
			return nil
		}).AnyTimes()

	var tvList []objectvisitor.TypedVisitor

	dv, err := objectvisitor.NewDefaultVisitor(dashConfig, q,
		objectvisitor.SetTypedVisitors(tvList))
	require.NoError(t, err)

	err = dv.Visit(ctx, replicaSet, handler)
	require.NoError(t, err)

	for k := range adjList {
		sort.Strings(adjList[k])
	}

	expectedAdjList := map[string][]string{
		"deployment":  {"replica-set"},
		"replica-set": {"deployment", "pod"},
		"pod":         {"replica-set"},
	}
	assert.Equal(t, expectedAdjList, adjList)

	sort.Strings(processed)
	expectedProcessed := []string{"deployment", "pod", "replica-set"}
	assert.Equal(t, expectedProcessed, processed)
}

func TestDefaultVisitor_Visit_use_typed_visitor(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	pod := testutil.CreatePod("pod")
	unstructuredPod := testutil.ToUnstructured(t, pod)

	q := queryerFake.NewMockQueryer(controller)

	handler := ovFake.NewMockObjectHandler(controller)

	defaultHandler := ovFake.NewMockDefaultTypedVisitor(controller)
	defaultHandler.EXPECT().
		Visit(gomock.Any(), unstructuredPod, handler, gomock.Any()).Return(nil)

	tv := ovFake.NewMockTypedVisitor(controller)
	tv.EXPECT().Supports().Return(gvk.PodGVK).AnyTimes()
	tv.EXPECT().
		Visit(gomock.Any(), unstructuredPod, handler, gomock.Any())
	tvList := []objectvisitor.TypedVisitor{tv}

	dv, err := objectvisitor.NewDefaultVisitor(dashConfig, q,
		objectvisitor.SetDefaultHandler(defaultHandler),
		objectvisitor.SetTypedVisitors(tvList))
	require.NoError(t, err)

	ctx := context.Background()
	err = dv.Visit(ctx, pod, handler)
	require.NoError(t, err)
}

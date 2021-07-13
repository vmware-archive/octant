/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testClient "k8s.io/client-go/kubernetes/fake"

	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/action"
	actionFake "github.com/vmware-tanzu/octant/pkg/action/fake"
	"github.com/vmware-tanzu/octant/pkg/cluster"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/store/fake"
)

func Test_Cordon(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	objectStore := fake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)
	kubernetesClient := clusterFake.NewMockKubernetesInterface(controller)
	clusterClient := clusterFake.NewMockClientInterface(controller)
	fakeClientset := testClient.NewSimpleClientset()

	clusterClient.EXPECT().KubernetesClient().AnyTimes().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().CoreV1().AnyTimes().Return(fakeClientset.CoreV1())

	cases := []struct {
		name          string
		clusterClient cluster.ClientInterface
		key           store.Key
		message       string
		alertType     action.AlertType
		cordoned      bool
		isErr         bool
	}{
		{
			name: "cordon unmarked node",
			key: store.Key{
				APIVersion: "v1",
				Kind:       "Node",
				Name:       "unmarked-node",
			},
			message:   `Node "unmarked-node" marked as unschedulable`,
			alertType: action.AlertTypeInfo,
			cordoned:  false,
			isErr:     false,
		},
		{
			name: "cordon marked node",
			key: store.Key{
				APIVersion: "v1",
				Kind:       "Node",
				Name:       "marked-node",
			},
			message:   `Unable to cordon node "marked-node": node "marked-node" already marked`,
			alertType: action.AlertTypeWarning,
			cordoned:  true,
			isErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := testutil.CreateNode(tc.key.Name)
			node.Spec.Unschedulable = tc.cordoned

			objectStore.EXPECT().
				Get(ctx, gomock.Eq(tc.key)).
				Return(testutil.ToUnstructured(t, node), nil)

			alerter.EXPECT().
				SendAlert(gomock.Any()).
				DoAndReturn(func(alert action.Alert) {
					assert.Equal(t, tc.alertType, alert.Type)
					assert.Equal(t, tc.message, alert.Message)
					assert.NotNil(t, alert.Expiration)
				})

			_, err := fakeClientset.CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
			require.NoError(t, err)

			cordon := octant.NewCordon(objectStore, clusterClient)
			assert.Equal(t, octant.ActionOverviewCordon, cordon.ActionName())

			payload := action.CreatePayload(octant.ActionOverviewCordon, map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Node",
				"name":       tc.key.Name,
			})

			require.NoError(t, cordon.Handle(ctx, alerter, payload))
		})
	}
}

func Test_Uncordon(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	objectStore := fake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)
	kubernetesClient := clusterFake.NewMockKubernetesInterface(controller)
	clusterClient := clusterFake.NewMockClientInterface(controller)
	fakeClientset := testClient.NewSimpleClientset()

	clusterClient.EXPECT().KubernetesClient().AnyTimes().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().CoreV1().AnyTimes().Return(fakeClientset.CoreV1())

	cases := []struct {
		name          string
		clusterClient cluster.ClientInterface
		key           store.Key
		message       string
		alertType     action.AlertType
		cordoned      bool
		isErr         bool
	}{
		{
			name: "uncordon marked node",
			key: store.Key{
				APIVersion: "v1",
				Kind:       "Node",
				Name:       "marked-node",
			},
			message:   `Node "marked-node" marked as schedulable`,
			alertType: action.AlertTypeInfo,
			cordoned:  true,
			isErr:     false,
		},
		{
			name: "uncordon unmarked node",
			key: store.Key{
				APIVersion: "v1",
				Kind:       "Node",
				Name:       "unmarked-node",
			},
			message:   `Unable to uncordon node "unmarked-node": node "unmarked-node" already unmarked`,
			alertType: action.AlertTypeWarning,
			cordoned:  false,
			isErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := testutil.CreateNode(tc.key.Name)
			node.Spec.Unschedulable = tc.cordoned

			objectStore.EXPECT().
				Get(ctx, gomock.Eq(tc.key)).
				Return(testutil.ToUnstructured(t, node), nil)

			alerter.EXPECT().
				SendAlert(gomock.Any()).
				DoAndReturn(func(alert action.Alert) {
					assert.Equal(t, tc.alertType, alert.Type)
					assert.Equal(t, tc.message, alert.Message)
					assert.NotNil(t, alert.Expiration)
				})

			_, err := fakeClientset.CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
			require.NoError(t, err)

			uncordon := octant.NewUncordon(objectStore, clusterClient)
			assert.Equal(t, octant.ActionOverviewUncordon, uncordon.ActionName())

			payload := action.CreatePayload(octant.ActionOverviewUncordon, map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Node",
				"name":       tc.key.Name,
			})

			require.NoError(t, uncordon.Handle(ctx, alerter, payload))
		})
	}
}

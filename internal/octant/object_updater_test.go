/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

package octant

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/action"
	actionFake "github.com/vmware-tanzu/octant/pkg/action/fake"
	"github.com/vmware-tanzu/octant/pkg/store"
	storeFake "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestObjectUpdateFromPayload(t *testing.T) {
	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))
	podS, err := kubernetes.SerializeToString(pod)
	require.NoError(t, err)

	tests := []struct {
		name    string
		payload action.Payload
		wanted  *unstructured.Unstructured
		wantErr bool
	}{
		{
			name: "in general",
			payload: action.Payload{
				"update": podS,
			},
			wanted: pod,
		},
		{
			name:    "no update in payload",
			payload: action.Payload{},
			wantErr: true,
		},
		{
			name: "upload is not yaml",
			payload: action.Payload{
				"update": "<<<",
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := ObjectUpdateFromPayload(test.payload)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.wanted, actual)
		})
	}
}

func TestObjectUpdaterDispatcher_Handle(t *testing.T) {
	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))
	podKey, err := store.KeyFromObject(pod)
	require.NoError(t, err)
	podPayload := action.Payload{
		"namespace":  pod.GetNamespace(),
		"apiVersion": pod.GetAPIVersion(),
		"kind":       pod.GetKind(),
		"name":       pod.GetName(),
	}

	tests := []struct {
		name              string
		payload           action.Payload
		objectFromPayload func(action.Payload) (*unstructured.Unstructured, error)
		initStore         func(ctrl *gomock.Controller) *storeFake.MockStore
		initAlerter       func(ctrl *gomock.Controller) *actionFake.MockAlerter
		wantErr           bool
	}{
		{
			name:    "in general",
			payload: podPayload,
			objectFromPayload: func(payload action.Payload) (*unstructured.Unstructured, error) {
				return pod, nil
			},
			initStore: func(ctrl *gomock.Controller) *storeFake.MockStore {
				objectStore := storeFake.NewMockStore(ctrl)
				objectStore.EXPECT().
					Update(gomock.Any(), podKey, gomock.Any()).Return(nil)

				return objectStore
			},
			initAlerter: func(ctrl *gomock.Controller) *actionFake.MockAlerter {
				alerter := actionFake.NewMockAlerter(ctrl)
				alerter.EXPECT().
					SendAlert(gomock.Any()).
					DoAndReturn(func(alert action.Alert) {
						require.Equal(t, action.AlertTypeInfo, alert.Type)
					})
				return alerter
			},
		},
		{
			name:    "unable to load object",
			payload: podPayload,
			objectFromPayload: func(payload action.Payload) (*unstructured.Unstructured, error) {
				return nil, fmt.Errorf("error")
			},
			initStore: func(ctrl *gomock.Controller) *storeFake.MockStore {
				objectStore := storeFake.NewMockStore(ctrl)
				return objectStore
			},
			initAlerter: func(ctrl *gomock.Controller) *actionFake.MockAlerter {
				alerter := actionFake.NewMockAlerter(ctrl)
				alerter.EXPECT().
					SendAlert(gomock.Any()).
					DoAndReturn(func(alert action.Alert) {
						require.Equal(t, action.AlertTypeError, alert.Type)
					})
				return alerter
			},
		},
		{
			name:    "update failed",
			payload: podPayload,
			objectFromPayload: func(payload action.Payload) (*unstructured.Unstructured, error) {
				return pod, nil
			},
			initStore: func(ctrl *gomock.Controller) *storeFake.MockStore {
				objectStore := storeFake.NewMockStore(ctrl)
				objectStore.EXPECT().
					Update(gomock.Any(), podKey, gomock.Any()).Return(fmt.Errorf("error"))

				return objectStore
			},
			initAlerter: func(ctrl *gomock.Controller) *actionFake.MockAlerter {
				alerter := actionFake.NewMockAlerter(ctrl)
				alerter.EXPECT().
					SendAlert(gomock.Any()).
					DoAndReturn(func(alert action.Alert) {
						require.Equal(t, action.AlertTypeError, alert.Type)
					})
				return alerter
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			objectStore := test.initStore(ctrl)
			alerter := test.initAlerter(ctrl)

			o := NewObjectUpdaterDispatcher(objectStore,
				func(dispatcher *ObjectUpdaterDispatcher) {
					dispatcher.objectFromPayload = test.objectFromPayload
				})
			err := o.Handle(ctx, alerter, test.payload)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

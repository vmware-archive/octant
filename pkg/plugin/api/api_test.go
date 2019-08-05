/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/portforward"
	portForwardFake "github.com/vmware/octant/internal/portforward/fake"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/plugin/api"
	"github.com/vmware/octant/pkg/store"
	storeFake "github.com/vmware/octant/pkg/store/fake"
)

type apiMocks struct {
	objectStore *storeFake.MockStore
	pf          *portForwardFake.MockPortForwarder
}

func TestAPI(t *testing.T) {
	listKey := store.Key{
		Namespace:  "default",
		APIVersion: "apps/v1",
		Kind:       "Deployment",
	}

	objects := testutil.ToUnstructuredList(t,
		testutil.CreateDeployment("deployment"),
	)

	getKey := store.Key{
		Namespace:  "default",
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Name:       "deployment",
	}
	object := testutil.ToUnstructured(t, testutil.CreateDeployment("deployment"))

	pfRequest := api.PortForwardRequest{
		Namespace: "default",
		PodName:   "pod",
		Port:      uint16(8080),
	}

	pfResponse := api.PortForwardResponse{
		ID:   "12345",
		Port: uint16(54321),
	}

	cases := []struct {
		name     string
		initFunc func(t *testing.T, mocks *apiMocks)
		doFunc   func(t *testing.T, client *api.Client)
	}{
		{
			name: "list",
			initFunc: func(t *testing.T, mocks *apiMocks) {
				mocks.objectStore.EXPECT().
					List(gomock.Any(), gomock.Eq(listKey)).Return(objects, nil)
			},
			doFunc: func(t *testing.T, client *api.Client) {
				clientCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				got, err := client.List(clientCtx, listKey)
				require.NoError(t, err)

				expected := objects

				assert.Equal(t, expected, got)
			},
		},
		{
			name: "get",
			initFunc: func(t *testing.T, mocks *apiMocks) {
				mocks.objectStore.EXPECT().
					Get(gomock.Any(), gomock.Eq(getKey)).Return(object, nil)
			},
			doFunc: func(t *testing.T, client *api.Client) {
				clientCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				got, err := client.Get(clientCtx, getKey)
				require.NoError(t, err)

				expected := object

				assert.Equal(t, expected, got)
			},
		},
		{
			name: "port forward",
			initFunc: func(t *testing.T, mocks *apiMocks) {
				resp := portforward.CreateResponse{
					ID: "12345",
					Ports: []portforward.PortForwardPortSpec{
						{Local: uint16(54321)},
					},
				}

				mocks.pf.EXPECT().
					Create(
						gomock.Any(), gvk.Pod, "pod", "default", uint16(8080)).
					Return(resp, nil)
			},
			doFunc: func(t *testing.T, client *api.Client) {
				clientCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				got, err := client.PortForward(clientCtx, pfRequest)
				require.NoError(t, err)

				expected := pfResponse

				assert.Equal(t, expected, got)
			},
		},
		{
			name: "port forward cancel",
			initFunc: func(t *testing.T, mocks *apiMocks) {
				mocks.pf.EXPECT().
					StopForwarder("12345")
			},
			doFunc: func(t *testing.T, client *api.Client) {
				clientCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				client.CancelPortForward(clientCtx, "12345")
			},
		},
		{
			name: "port forward",
			initFunc: func(t *testing.T, mocks *apiMocks) {
				resp := portforward.CreateResponse{
					ID: "12345",
					Ports: []portforward.PortForwardPortSpec{
						{Local: uint16(54321)},
					},
				}

				mocks.pf.EXPECT().
					Create(
						gomock.Any(), gvk.Pod, "pod", "default", uint16(8080)).
					Return(resp, nil)
			},
			doFunc: func(t *testing.T, client *api.Client) {
				clientCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				got, err := client.PortForward(clientCtx, pfRequest)
				require.NoError(t, err)

				expected := pfResponse

				assert.Equal(t, expected, got)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			appObjectStore := storeFake.NewMockStore(controller)
			pf := portForwardFake.NewMockPortForwarder(controller)
			tc.initFunc(t, &apiMocks{
				objectStore: appObjectStore,
				pf:          pf})

			service := &api.GRPCService{
				ObjectStore:   appObjectStore,
				PortForwarder: pf,
			}

			a, err := api.New(service)
			require.NoError(t, err)

			ctx := context.Background()

			err = a.Start(ctx)
			require.NoError(t, err)

			checkPort(t, true, a.Addr())

			client, err := api.NewClient(a.Addr())
			require.NoError(t, err)

			tc.doFunc(t, client)

		})
	}
}

func checkPort(t *testing.T, isListen bool, addr string) {
	_, err := net.Listen("tcp", addr)
	if isListen {
		require.Error(t, err)
		return
	}

	require.NoError(t, err)
}

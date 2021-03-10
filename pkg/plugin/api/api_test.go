/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api_test

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/portforward"
	portForwardFake "github.com/vmware-tanzu/octant/internal/portforward/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/plugin/api"
	"github.com/vmware-tanzu/octant/pkg/store"
	storeFake "github.com/vmware-tanzu/octant/pkg/store/fake"
)

// matches arguments of type context.Context
var contextType gomock.Matcher = gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem())

// matches arguments of type store.Key
var storeKeyType gomock.Matcher = gomock.AssignableToTypeOf(reflect.TypeOf((*store.Key)(nil)).Elem())

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
					List(contextType, gomock.Eq(listKey)).
					Return(objects, false, nil).
					Do(func(ctx context.Context, _ store.Key) {
						require.Equal(t, "bar", ctx.Value(api.DashboardMetadataKey("foo")))
					})
			},
			doFunc: func(t *testing.T, client *api.Client) {
				clientCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				clientCtx = metadata.AppendToOutgoingContext(clientCtx, "x-octant-foo", "bar")
				got, err := client.List(clientCtx, listKey)
				require.NoError(t, err)

				expected := objects

				assert.Equal(t, expected, got)
			},
		},
		{
			name: "update",
			initFunc: func(t *testing.T, mocks *apiMocks) {
				mocks.objectStore.EXPECT().
					Update(contextType, storeKeyType, gomock.Any()).
					Return(nil).
					Do(func(ctx context.Context, _ store.Key, _ func(*unstructured.Unstructured) error) {
						require.Equal(t, "bar", ctx.Value(api.DashboardMetadataKey("foo")))
					})
			},
			doFunc: func(t *testing.T, client *api.Client) {
				clientCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				clientCtx = metadata.AppendToOutgoingContext(clientCtx, "x-octant-foo", "bar")
				err := client.Update(clientCtx, object)
				require.NoError(t, err)
			},
		},
		{
			name: "create",
			initFunc: func(t *testing.T, mocks *apiMocks) {
				mocks.objectStore.EXPECT().
					Create(contextType, gomock.Eq(object)).
					Return(nil).
					Do(func(ctx context.Context, _ *unstructured.Unstructured) {
						require.Equal(t, "bar", ctx.Value(api.DashboardMetadataKey("foo")))
					})
			},
			doFunc: func(t *testing.T, client *api.Client) {
				clientCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				clientCtx = metadata.AppendToOutgoingContext(clientCtx, "x-octant-foo", "bar")
				err := client.Create(clientCtx, object)
				require.NoError(t, err)
			},
		},
		{
			name: "get",
			initFunc: func(t *testing.T, mocks *apiMocks) {
				mocks.objectStore.EXPECT().
					Get(contextType, gomock.Eq(getKey)).
					Return(object, nil).
					Do(func(ctx context.Context, _ store.Key) {
						require.Equal(t, "bar", ctx.Value(api.DashboardMetadataKey("foo")))
					})
			},
			doFunc: func(t *testing.T, client *api.Client) {
				clientCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				clientCtx = metadata.AppendToOutgoingContext(clientCtx, "x-octant-foo", "bar")
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
						gomock.Any(), gomock.Any(), gvk.Pod, "pod", "default", uint16(8080)).
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
						gomock.Any(), gomock.Any(), gvk.Pod, "pod", "default", uint16(8080)).
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

			viper.SetDefault("client-max-recv-msg-size", 1024*1024*16)

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

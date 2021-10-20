/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"testing"

	linkFake "github.com/vmware-tanzu/octant/internal/link/fake"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	storefake "github.com/vmware-tanzu/octant/pkg/store/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_runIngressStatus(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T, *storefake.MockStore) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				mockServiceInCache(t, o, "default", "single-service", "service_single_service.yaml")
				objectFile := "ingress_single_service.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				Details:    []component.Component{component.NewText("Ingress is OK")},
				Properties: []component.Property{{Label: "Default Backend", Value: component.NewText("single-service")}},
			},
		},
		{
			name: "no matching backends",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				key := store.Key{Namespace: "default", APIVersion: "v1", Kind: "Service", Name: "no-such-service"}
				o.EXPECT().Get(gomock.Any(), gomock.Eq(key)).Return(nil, nil)

				objectFile := "ingress_no_matching_backend.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusError,
				Details:    []component.Component{component.NewText("Backend refers to service \"no-such-service\" which doesn't exist")},
				Properties: []component.Property{{Label: "Default Backend", Value: component.NewText("Not configured")}},
			},
		},
		{
			name: "no matching port",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				mockServiceInCache(t, o, "default", "service-wrong-port", "service_wrong_port.yaml")
				objectFile := "ingress_no_matching_port.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)
			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusError,
				Details:    []component.Component{component.NewText("Backend for service \"service-wrong-port\" specifies an invalid port")},
				Properties: []component.Property{{Label: "Default Backend", Value: component.NewText("Not configured")}},
			},
		},
		{
			name: "mismatched TLS host",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				mockServiceInCache(t, o, "default", "my-service", "service_my-service.yaml")
				mockSecretInCache(t, o, "default", "testsecret-tls", "secret_testsecret-tls.yaml")

				objectFile := "ingress_mismatched_tls_host.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusError,
				Details:    []component.Component{component.NewText("No matching TLS host for rule \"not-the-tls-host.com\"")},
				Properties: []component.Property{{Label: "Default Backend", Value: component.NewText("Not configured")}},
			},
		},
		{
			name: "wildcard TLS host",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				mockServiceInCache(t, o, "default", "my-service", "service_my-service.yaml")
				mockSecretInCache(t, o, "default", "testsecret-tls", "secret_testsecret-tls.yaml")

				objectFile := "ingress_wildcard_tls_host.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				Details:    []component.Component{component.NewText("Ingress is OK")},
				Properties: []component.Property{{Label: "Default Backend", Value: component.NewText("Not configured")}},
			},
		},
		{
			name: "missing TLS secret",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				mockServiceInCache(t, o, "default", "my-service", "service_my-service.yaml")

				key := store.Key{Namespace: "default", APIVersion: "v1", Kind: "Secret", Name: "no-such-secret"}
				o.EXPECT().Get(gomock.Any(), gomock.Eq(key)).Return(nil, nil)

				objectFile := "ingress_ingress-bad-tls-host.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				NodeStatus: component.NodeStatusError,
				Details:    []component.Component{component.NewText("Secret \"no-such-secret\" does not exist")},
				Properties: []component.Property{{Label: "Default Backend", Value: component.NewText("Not configured")}},
			},
		},
		{
			name: "object is nil",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				return nil
			},
			isErr: true,
		},
		{
			name: "object is not an ingress",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				return &unstructured.Unstructured{}
			},
			isErr: true,
		},
		{
			name: "service multiple port service",
			init: func(t *testing.T, o *storefake.MockStore) runtime.Object {
				mockServiceInCache(t, o, "default", "multiple-port-service", "service_multiple_port_service.yaml")
				objectFile := "ingress_multiple_port_service.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				Details:    []component.Component{component.NewText("Ingress is OK")},
				Properties: []component.Property{{Label: "Default Backend", Value: component.NewText("multiple-port-service")}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			linkInterface := linkFake.NewMockInterface(controller)
			defer controller.Finish()

			o := storefake.NewMockStore(controller)

			object := tc.init(t, o)

			ctx := context.Background()
			status, err := runIngressStatus(ctx, object, o, linkInterface)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}

func Test_hostMatcher(t *testing.T) {
	cases := []struct {
		name    string
		hosts   []string
		lookups map[string]bool
	}{
		{
			name:  "string",
			hosts: []string{"example2.com"},
			lookups: map[string]bool{
				"example2.com": true,
				"example1.com": false,
			},
		},
		{
			name:  "global",
			hosts: []string{"*.example2.com"},
			lookups: map[string]bool{
				"foo.example2.com": true,
				"bar.example2.com": true,
				"foo.example1.com": false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hm := hostMatcher{}

			for _, host := range tc.hosts {
				require.NoError(t, hm.AddHost(host))
			}

			for k, v := range tc.lookups {
				require.Equal(t, v, hm.Match(k))

			}
		})
	}

}

func mockSecretInCache(t *testing.T, o *storefake.MockStore, namespace, name, file string) runtime.Object {
	secret := testutil.LoadObjectFromFile(t, file)
	key := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Secret",
		Name:       name,
	}

	o.EXPECT().Get(gomock.Any(), gomock.Eq(key)).Return(testutil.ToUnstructured(t, secret), nil)

	return secret
}

func mockServiceInCache(t *testing.T, o *storefake.MockStore, namespace, name, file string) runtime.Object {
	secret := testutil.LoadObjectFromFile(t, file)
	key := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Service",
		Name:       name,
	}

	o.EXPECT().Get(gomock.Any(), gomock.Eq(key)).Return(testutil.ToUnstructured(t, secret), nil)

	return secret
}

/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_namespaces_list(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/namespaces", nil)

	tests := []struct {
		name     string
		init     func(*clusterfake.MockNamespaceInterface)
		expected []string
	}{
		{
			name: "general",
			init: func(ns *clusterfake.MockNamespaceInterface) {
				ns.EXPECT().Names().Return([]string{"default", "other"}, nil)
			},
			expected: []string{"default", "other"},
		},
		{
			name: "cannot list due to rbac error",
			init: func(ns *clusterfake.MockNamespaceInterface) {
				ns.EXPECT().Names().Return(nil, errors.Errorf("error"))
				ns.EXPECT().InitialNamespace().Return("initial-namespace")
			},
			expected: []string{"initial-namespace"},
		},
	}

	for _, tc := range tests {
		controller := gomock.NewController(t)
		defer controller.Finish()

		nsClient := clusterfake.NewMockNamespaceInterface(controller)
		tc.init(nsClient)

		handler := newNamespaces(nsClient, log.NopLogger())
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		var nr namespacesResponse
		err := json.NewDecoder(resp.Body).Decode(&nr)
		require.NoError(t, err)

		assert.Equal(t, tc.expected, nr.Namespaces)
	}
}

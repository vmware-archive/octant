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

	clusterfake "github.com/vmware/octant/internal/cluster/fake"
	"github.com/vmware/octant/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_clusterInfo(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/info", nil)

	tests := []struct {
		name     string
		init     func(*testing.T, *clusterfake.MockInfoInterface)
		expected clusterInfoResponse
	}{
		{
			name: "general",
			init: func(t *testing.T, infoClient *clusterfake.MockInfoInterface) {
				infoClient.EXPECT().Context().Return("main-context")
				infoClient.EXPECT().Cluster().Return("my-cluster")
				infoClient.EXPECT().Server().Return("https://localhost:6443")
				infoClient.EXPECT().User().Return("me-of-course")
			},
			expected: clusterInfoResponse{
				Context: "main-context",
				Cluster: "my-cluster",
				Server:  "https://localhost:6443",
				User:    "me-of-course",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			infoClient := clusterfake.NewMockInfoInterface(controller)

			tc.init(t, infoClient)

			handler := newClusterInfo(infoClient, log.NopLogger())
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)

			var ciResp clusterInfoResponse
			err := json.NewDecoder(resp.Body).Decode(&ciResp)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, ciResp)
		})
	}
}

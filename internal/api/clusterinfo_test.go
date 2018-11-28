package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_clusterInfo(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/info", nil)

	tests := []struct {
		name       string
		infoClient cluster.InfoInterface
		expected   clusterInfoResponse
	}{
		{
			name: "general",
			infoClient: fake.ClusterInfo{
				ContextVal: "main-context",
				ClusterVal: "my-cluster",
				ServerVal:  "https://localhost:6443",
				UserVal:    "me-of-course",
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

		})
		handler := newClusterInfo(tc.infoClient, log.NopLogger())
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		var ciResp clusterInfoResponse
		err := json.NewDecoder(resp.Body).Decode(&ciResp)
		require.NoError(t, err)

		assert.Equal(t, tc.expected, ciResp)
	}
}

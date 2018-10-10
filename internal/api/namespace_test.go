package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heptio/developer-dash/internal/module"
	modulefake "github.com/heptio/developer-dash/internal/module/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_namespace_update(t *testing.T) {
	cases := []struct {
		name              string
		ns                string
		statusCode        int
		expectedNamespace string
	}{
		{
			name:              "update ns with valid ns",
			ns:                "new-ns",
			statusCode:        http.StatusNoContent,
			expectedNamespace: "new-ns",
		},
		{
			name:              "update ns empty string",
			ns:                "",
			statusCode:        http.StatusBadRequest,
			expectedNamespace: "default",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := modulefake.NewModule("module")
			manager := modulefake.NewStubManager("default", []module.Module{m})

			handler := newNamespace(manager)

			ts := httptest.NewServer(http.HandlerFunc(handler.update))
			defer ts.Close()

			nr := namespaceRequest{Namespace: tc.ns}
			data, err := json.Marshal(&nr)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewReader(data))
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tc.statusCode, resp.StatusCode)
			assert.Equal(t, tc.expectedNamespace, manager.GetNamespace())
		})
	}
}

func Test_namespace_read(t *testing.T) {
	m := modulefake.NewModule("module")
	manager := modulefake.NewStubManager("default", []module.Module{m})

	handler := newNamespace(manager)

	ts := httptest.NewServer(http.HandlerFunc(handler.read))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	require.NoError(t, err)

	defer resp.Body.Close()

	var nr namespaceResponse
	err = json.NewDecoder(resp.Body).Decode(&nr)
	require.NoError(t, err)

	expected := namespaceResponse{Namespace: "default"}

	assert.Equal(t, expected, nr)
}

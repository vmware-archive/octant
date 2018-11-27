package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	modulefake "github.com/heptio/developer-dash/internal/module/fake"
	"github.com/heptio/go-telemetry/pkg/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var telemetryClient = &telemetry.NilClient{}

func TestAPI_routes(t *testing.T) {
	cases := []struct {
		path            string
		method          string
		body            io.Reader
		expectedCode    int
		expectedContent string
	}{
		{
			path:         "/namespaces",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
		},
		{
			path:         "/navigation",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
		},
		{
			path:         "/content/",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
		},
		{
			path:            "/content/module/",
			method:          http.MethodGet,
			expectedCode:    http.StatusOK,
			expectedContent: "root",
		},
		{
			path:            "/content/module/nested",
			method:          http.MethodGet,
			expectedCode:    http.StatusOK,
			expectedContent: "module",
		},
		{
			path:         "/missing",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%s: %s", tc.method, tc.path)
		t.Run(name, func(t *testing.T) {
			m := modulefake.NewModule("module", log.NopLogger())

			manager := modulefake.NewStubManager("default", []module.Module{m})

			nsClient := fake.NewNamespaceClient([]string{"default"}, nil, "default")
			srv := New("/", nsClient, manager, log.NopLogger(), telemetryClient)

			err := srv.RegisterModule(m)
			require.NoError(t, err)

			ts := httptest.NewServer(srv.Handler())
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			u.Path = tc.path

			req, err := http.NewRequest(tc.method, u.String(), tc.body)
			require.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			if tc.expectedContent != "" {
				assert.Equal(t, tc.expectedContent, string(data))
			}
			assert.Equal(t, tc.expectedCode, res.StatusCode)
		})
	}
}

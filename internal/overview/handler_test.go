package overview

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handler_routes(t *testing.T) {
	cases := []struct {
		name         string
		path         string
		values       url.Values
		expectedCode int
		expectedBody string
	}{
		{
			name:         "GET /",
			path:         "/",
			expectedCode: http.StatusOK,
			expectedBody: "{\"contents\":[{\"type\":\"table\",\"title\":\"/\",\"columns\":[{\"name\":\"foo\",\"accessor\":\"foo\"},{\"name\":\"bar\",\"accessor\":\"bar\"},{\"name\":\"baz\",\"accessor\":\"baz\"}],\"rows\":[{\"bar\":\"r1c2\",\"baz\":\"r1c3\",\"foo\":\"r1c1\"},{\"bar\":\"r2c2\",\"baz\":\"r2c3\",\"foo\":\"r2c1\"},{\"bar\":\"r3c2\",\"baz\":\"r3c3\",\"foo\":\"r3c1\"}]}]}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newHandler("/")

			ts := httptest.NewServer(h)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			u.Path = tc.path
			u.RawQuery = tc.values.Encode()

			resp, err := http.Get(u.String())
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
			assert.Equal(t, tc.expectedBody, strings.TrimSpace(string(body)))
		})
	}
}

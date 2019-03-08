package dash

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/api"
	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	modulefake "github.com/heptio/developer-dash/internal/module/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_dash_Run(t *testing.T) {
	cases := []struct {
		name         string
		hasCustomURL bool
		expected     string
	}{
		{
			name:     "embedded dashboard ui",
			expected: "embedded",
		},
		{
			name:         "custom dashboard ui",
			hasCustomURL: true,
			expected:     "custom",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			ctx, cancel := context.WithCancel(context.Background())
			namespace := "default"
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)

			var uiURL string
			if tc.hasCustomURL {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "custom")
				}))
				defer ts.Close()

				uiURL = ts.URL
			}

			defaultHandler := func() (http.Handler, error) {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "embedded")
				}), nil
			}

			nsClient := clusterfake.NewMockNamespaceInterface(controller)
			nsClient.EXPECT().InitialNamespace().Return("default").AnyTimes()

			o := newDashModule()
			manager := modulefake.NewStubManager("default", []module.Module{o})

			infoClient := clusterfake.NewMockInfoInterface(controller)
			service := api.New(ctx, apiPathPrefix, nsClient, infoClient, manager, log.NopLogger())
			d, err := newDash(listener, namespace, uiURL, service, log.NopLogger())
			require.NoError(t, err)

			d.willOpenBrowser = false
			d.defaultHandler = defaultHandler

			var runErr error
			ch := make(chan bool, 1)

			go func() {
				runErr = d.Run(ctx)
				ch <- true
			}()

			dashURL := fmt.Sprintf("http://%s", listener.Addr())

			resp, err := http.Get(dashURL)
			require.NoError(t, err)

			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, string(data))

			cancel()
			<-ch
			assert.NoError(t, runErr)

		})
	}
}

func Test_dash_routes(t *testing.T) {
	cases := []struct {
		path         string
		expectedCode int
		expectedBody string
	}{
		{
			path:         "/",
			expectedCode: http.StatusOK,
			expectedBody: "body",
		},
		{
			path:         "/nested",
			expectedCode: http.StatusOK,
			expectedBody: "body",
		},
		{
			path:         "/api/v1/namespaces",
			expectedCode: http.StatusOK,
			expectedBody: "{\"namespaces\":[\"default\"]}\n",
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("GET: %s", tc.path)
		t.Run(name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			namespace := "default"
			uiURL := ""
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)

			nsClient := clusterfake.NewMockNamespaceInterface(controller)
			nsClient.EXPECT().InitialNamespace().Return("default").AnyTimes()
			nsClient.EXPECT().Names().Return([]string{"default"}, nil).AnyTimes()

			o := newDashModule()
			manager := modulefake.NewStubManager("default", []module.Module{o})

			infoClient := clusterfake.NewMockInfoInterface(controller)

			ctx := context.Background()
			service := api.New(ctx, apiPathPrefix, nsClient, infoClient, manager, log.NopLogger())

			d, err := newDash(listener, namespace, uiURL, service, log.NopLogger())
			require.NoError(t, err)
			d.defaultHandler = func() (http.Handler, error) {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "body")
				}), nil
			}

			handler, err := d.handler(ctx)
			require.NoError(t, err)

			ts := httptest.NewServer(handler)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			u.Path = tc.path

			res, err := http.Get(u.String())
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)
			data, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedBody, string(data))
		})
	}
}

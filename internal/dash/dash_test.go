package dash

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

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

			d := newDash(listener, namespace, uiURL)
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

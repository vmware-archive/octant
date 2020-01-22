/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/api"
	apiFake "github.com/vmware-tanzu/octant/internal/api/fake"
	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	moduleFake "github.com/vmware-tanzu/octant/internal/module/fake"
	"github.com/vmware-tanzu/octant/internal/terminal"
	terminalFake "github.com/vmware-tanzu/octant/internal/terminal/fake"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestAPI_routes(t *testing.T) {
	cases := []struct {
		path                string
		method              string
		body                io.Reader
		expectedCode        int
		expectedContent     string
		expectedContentPath string
		expectedNamespace   string
	}{
		{
			path:         "/missing",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%s: %s", tc.method, tc.path)
		t.Run(name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashConfig := configFake.NewMockDash(controller)
			logger := log.NopLogger()
			dashConfig.EXPECT().Logger().Return(logger).AnyTimes()
			clusterClient := clusterFake.NewMockClientInterface(controller)
			dashConfig.EXPECT().ClusterClient().Return(clusterClient).AnyTimes()
			terminalManager := terminalFake.NewMockManager(controller)
			dashConfig.EXPECT().TerminalManager().Return(terminalManager).AnyTimes()

			m := moduleFake.NewMockModule(controller)
			m.EXPECT().
				Name().Return("module").AnyTimes()
			m.EXPECT().
				ContentPath().Return("/module").AnyTimes()
			m.EXPECT().
				Content(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, contentPath, prefix, namespace string, opts module.ContentOptions) (component.ContentResponse, error) {
					switch contentPath {
					case "/":
						return component.ContentResponse{
							Title: component.Title(component.NewText("/")),
						}, nil
					case "/nested":
						return component.ContentResponse{
							Title: component.Title(component.NewText("/nested")),
						}, nil
					default:
						return component.ContentResponse{}, errors.New("not found")
					}
				}).
				AnyTimes()
			m.EXPECT().
				Navigation(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, namespace, prefix string) ([]navigation.Navigation, error) {
					nav := navigation.Navigation{
						Path:  prefix,
						Title: "module",
					}

					return []navigation.Navigation{nav}, nil
				}).
				AnyTimes()

			actionDispatcher := apiFake.NewMockActionDispatcher(controller)

			ctx := context.Background()
			srv := api.New(ctx, "/", actionDispatcher, dashConfig)

			instances := make(chan terminal.Instance)
			terminalManager.EXPECT().Select(ctx).Return(instances).AnyTimes()

			handler, err := srv.Handler(ctx)
			require.NoError(t, err)

			ts := httptest.NewServer(handler)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			// Add relative section to server url
			u, err = u.Parse(tc.path)
			require.NoError(t, err)

			req, err := http.NewRequest(tc.method, u.String(), tc.body)
			require.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			defer func() {
				require.NoError(t, res.Body.Close())
			}()

			data, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			if tc.expectedContent != "" {
				assert.Equal(t, tc.expectedContent, string(data))
			}
			assert.Equal(t, tc.expectedCode, res.StatusCode)

		})
	}
}

/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

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

	apiFake "github.com/vmware/octant/internal/api/fake"
	clusterFake "github.com/vmware/octant/internal/cluster/fake"
	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/module"
	moduleFake "github.com/vmware/octant/internal/module/fake"
	"github.com/vmware/octant/pkg/view/component"
)

type testMocks struct {
	namespace *clusterFake.MockNamespaceInterface
	info      *clusterFake.MockInfoInterface
}

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
			path:         "/cluster-info",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
		},
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
			path:            "/content/module/",
			method:          http.MethodGet,
			expectedCode:    http.StatusOK,
			expectedContent: "{\"title\":[{\"metadata\":{\"type\":\"text\"},\"config\":{\"value\":\"/\"}}],\"viewComponents\":null}\n",
		},
		{
			path:                "/content/module/namespace/another-namespace/",
			method:              http.MethodGet,
			expectedCode:        http.StatusOK,
			expectedContent:     "{\"title\":[{\"metadata\":{\"type\":\"text\"},\"config\":{\"value\":\"/\"}}],\"viewComponents\":null}\n",
			expectedNamespace:   "another-namespace",
			expectedContentPath: "/",
		},
		{
			path:                "/content/module/?namespace=fromquery",
			method:              http.MethodGet,
			expectedCode:        http.StatusOK,
			expectedContent:     "{\"title\":[{\"metadata\":{\"type\":\"text\"},\"config\":{\"value\":\"/\"}}],\"viewComponents\":null}\n",
			expectedNamespace:   "fromquery",
			expectedContentPath: "/",
		},
		{
			path:                "/content/module/namespace/path-takes-precedence/?namespace=fromquery",
			method:              http.MethodGet,
			expectedCode:        http.StatusOK,
			expectedContent:     "{\"title\":[{\"metadata\":{\"type\":\"text\"},\"config\":{\"value\":\"/\"}}],\"viewComponents\":null}\n",
			expectedNamespace:   "path-takes-precedence",
			expectedContentPath: "/",
		},
		{
			path:            "/content/module/nested",
			method:          http.MethodGet,
			expectedCode:    http.StatusOK,
			expectedContent: "{\"title\":[{\"metadata\":{\"type\":\"text\"},\"config\":{\"value\":\"/nested\"}}],\"viewComponents\":null}\n",
		},
		{
			path:                "/content/module/namespace/default/nested",
			method:              http.MethodGet,
			expectedCode:        http.StatusOK,
			expectedContent:     "{\"title\":[{\"metadata\":{\"type\":\"text\"},\"config\":{\"value\":\"/nested\"}}],\"viewComponents\":null}\n",
			expectedNamespace:   "default",
			expectedContentPath: "/nested",
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
			controller := gomock.NewController(t)
			defer controller.Finish()

			mocks := &testMocks{
				namespace: clusterFake.NewMockNamespaceInterface(controller),
				info:      clusterFake.NewMockInfoInterface(controller),
			}

			mocks.info.EXPECT().Context().Return("main-context").AnyTimes()
			mocks.info.EXPECT().Cluster().Return("my-cluster").AnyTimes()
			mocks.info.EXPECT().Server().Return("https://localhost:6443").AnyTimes()
			mocks.info.EXPECT().User().Return("me-of-course").AnyTimes()

			mocks.namespace.EXPECT().Names().Return([]string{"default"}, nil).AnyTimes()

			m := moduleFake.NewMockModule(controller)
			m.EXPECT().
				Name().Return("module").AnyTimes()
			m.EXPECT().
				ContentPath().Return("/module").AnyTimes()
			m.EXPECT().
				Content(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
				Handlers(gomock.Any()).Return(make(map[string]http.Handler))
			m.EXPECT().
				Navigation(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, namespace, prefix string) ([]octant.Navigation, error) {
					nav := octant.Navigation{
						Path:  prefix,
						Title: "module",
					}

					return []octant.Navigation{nav}, nil
				}).
				AnyTimes()

			manager := moduleFake.NewMockManagerInterface(controller)

			clusterClient := apiFake.NewMockClusterClient(controller)
			clusterClient.EXPECT().NamespaceClient().Return(mocks.namespace, nil).AnyTimes()
			clusterClient.EXPECT().InfoClient().Return(mocks.info, nil).AnyTimes()

			ctx := context.Background()
			srv := New(ctx, "/", clusterClient, manager, log.NopLogger())

			err := srv.RegisterModule(m)
			require.NoError(t, err)

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

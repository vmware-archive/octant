/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dashConfigFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/kubeconfig"
	"github.com/vmware/octant/internal/kubeconfig/fake"
	"github.com/vmware/octant/internal/log"
)

func Test_updateCurrentContextHandler(t *testing.T) {
	newContextName := "my-context"

	contextUpdated := false

	h := &updateCurrentContextHandler{
		logger: log.NopLogger(),
		contextUpdateFunc: func(name string) error {
			assert.Equal(t, newContextName, name)
			contextUpdated = true
			return nil
		},
	}

	ts := httptest.NewServer(h)
	defer ts.Close()

	req := updateCurrentContextRequest{
		RequestedContext: "my-context",
	}

	reqData, err := json.Marshal(req)
	require.NoError(t, err)

	r := bytes.NewReader(reqData)

	res, err := http.Post(ts.URL, "application/json", r)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.True(t, contextUpdated)
}

func Test_kubeContextGenerator(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	kc := &kubeconfig.KubeConfig{
		CurrentContext: "current-context",
	}

	loader := fake.NewMockLoader(controller)
	loader.EXPECT().
		Load("/path").
		Return(kc, nil)

	configLoaderFuncOpt := func(x *kubeContextGenerator) {
		x.ConfigLoader = loader
	}

	dashConfig := dashConfigFake.NewMockDash(controller)
	dashConfig.EXPECT().KubeConfigPath().Return("/path")
	dashConfig.EXPECT().ContextName().Return("")

	kgc := newKubeContextGenerator(dashConfig, configLoaderFuncOpt)

	assert.Equal(t, "kubeConfig", kgc.Name())

	ctx := context.Background()
	e, err := kgc.Event(ctx)
	require.NoError(t, err)

	assert.Equal(t, eventTypeKubeConfig, e.Type)

	resp := kubeContextsResponse{
		CurrentContext: kc.CurrentContext,
	}
	expectedData, err := json.Marshal(&resp)
	require.NoError(t, err)

	assert.JSONEq(t, string(expectedData), string(e.Data))
}

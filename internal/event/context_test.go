/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dashConfigFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/kubeconfig"
	"github.com/vmware/octant/internal/kubeconfig/fake"
	"github.com/vmware/octant/internal/octant"
)

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

	configLoaderFuncOpt := func(x *ContextsGenerator) {
		x.ConfigLoader = loader
	}

	dashConfig := dashConfigFake.NewMockDash(controller)
	dashConfig.EXPECT().KubeConfigPath().Return("/path")
	dashConfig.EXPECT().ContextName().Return("")

	kgc := NewContextsGenerator(dashConfig, configLoaderFuncOpt)

	assert.Equal(t, "kubeConfig", kgc.Name())

	ctx := context.Background()
	e, err := kgc.Event(ctx)
	require.NoError(t, err)

	assert.Equal(t, octant.EventTypeKubeConfig, e.Type)

	resp := kubeContextsResponse{
		CurrentContext: kc.CurrentContext,
	}

	assert.Equal(t, resp, e.Data)
}

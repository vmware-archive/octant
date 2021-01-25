/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dashConfigFake "github.com/vmware-tanzu/octant/internal/config/fake"
)

func Test_kubeContextGenerator(t *testing.T) {
	currentContext := "current-context"
	controller := gomock.NewController(t)
	defer controller.Finish()
	dashConfig := dashConfigFake.NewMockDash(controller)
	dashConfig.EXPECT().CurrentContext().Return(currentContext)
	dashConfig.EXPECT().Contexts().Return(nil)

	kgc := NewContextsGenerator(dashConfig)

	assert.Equal(t, "kubeConfig", kgc.Name())

	ctx := context.Background()
	e, err := kgc.Event(ctx)
	require.NoError(t, err)

	assert.Equal(t, event.EventTypeKubeConfig, e.Type)

	resp := kubeContextsResponse{
		CurrentContext: currentContext,
	}

	assert.Equal(t, resp, e.Data)
}

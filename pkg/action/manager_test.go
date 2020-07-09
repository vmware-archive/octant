/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package action_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/action/fake"
)

func TestManager(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	alerter := fake.NewMockAlerter(controller)

	logger := log.NopLogger()

	m := action.NewManager(logger)

	payloadRan := false
	fn := func(context.Context, action.Alerter, action.Payload) error {
		payloadRan = true
		return nil
	}

	actionPath := "path"

	err := m.Register(actionPath, "internal", fn)
	require.NoError(t, err)

	payload := action.Payload{}

	ctx := context.Background()

	err = m.Dispatch(ctx, alerter, actionPath, payload)
	require.NoError(t, err)

	assert.True(t, payloadRan)
}

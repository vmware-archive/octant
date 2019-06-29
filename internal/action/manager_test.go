/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package action

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/log"
)

func TestManager(t *testing.T) {
	logger := log.NopLogger()

	m := NewManager(logger)

	payloadRan := false
	fn := func(context.Context, Payload) error {
		payloadRan = true
		return nil
	}

	actionPath := "path"

	err := m.Register(actionPath, fn)
	require.NoError(t, err)

	payload := Payload{}

	ctx := context.Background()

	err = m.Dispatch(ctx, actionPath, payload)
	require.NoError(t, err)

	assert.True(t, payloadRan)
}

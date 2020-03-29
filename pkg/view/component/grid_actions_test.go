/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

package component

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/pkg/action"
)

func TestGridActions_AddAction(t *testing.T) {
	ga := NewGridActions()

	payload := action.Payload{"foo": "bar"}
	ga.AddAction("name", "/path", payload)

	expected := []GridAction{
		{
			Name:       "name",
			ActionPath: "/path",
			Payload:    payload,
		},
	}
	require.Equal(t, expected, ga.Config.Actions)
}

/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"testing"

	"github.com/vmware-tanzu/octant/pkg/action"

	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_DeleteObjectConfirmation(t *testing.T) {
	pod := testutil.CreatePod("pod")
	option, err := DeleteObjectConfirmationButton(pod)
	require.NoError(t, err)

	button := component.NewButton("", action.Payload{})
	option(button)

	expected := component.NewButton("",
		action.Payload{},
		component.WithButtonConfirmation(
			"Delete Pod",
			"Are you sure you want to delete *Pod* **pod**? This action is permanent and cannot be recovered.",
		))

	require.Equal(t, expected, button)
}

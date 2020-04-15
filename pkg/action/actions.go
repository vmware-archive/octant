/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package action

const (
	// RequestSetNamespace is the action for when the current namespace in Octant changes.
	// The ActionRequest.Payload for this action contains a single string entry `namespace` with a value
	// of the new current namespace.
	RequestSetNamespace   = "action.octant.dev/setNamespace"
)

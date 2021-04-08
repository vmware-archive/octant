/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package action

const (
	// RequestSetNamespace is the action for when the current namespace in Octant changes.
	// The ActionRequest.Payload for this action contains a single string entry `namespace` with a value
	// of the new current namespace.
	RequestSetNamespace = "action.octant.dev/setNamespace"

	// RequestSetFilter is the action for when the filters are updated in Octant.
	// The ActionRequest.Payload for this action contains an array of key/values representing
	// the filters.
	RequestSetFilter = "action.octant.dev/setFilter"

	// RequestSetContext is the action for when the context in Octant changes.
	// The ActionRequest.Payload for this action contains a single string entry `contextName` with a value
	// of the new context.
	RequestSetContext = "action.octant.dev/setContext"
)

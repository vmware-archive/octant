/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

type EventType string

// Event is an event for the dash frontend.
type Event struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
	Err  error
}

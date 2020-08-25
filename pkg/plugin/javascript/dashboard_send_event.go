/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package javascript

import (
	"context"
	"fmt"

	"github.com/dop251/goja"

	ocontext "github.com/vmware-tanzu/octant/internal/context"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/event"
)

// DashboardList is a function that lists objects by key.
type DashboardSendEvent struct {
	WebsocketClientManager event.WSClientGetter
}

var _ octant.DashboardClientFunction = &DashboardSendEvent{}

// NewDashboardSendEvent creates an instance of DashboardSendEvent.
func NewDashboardSendEvent(websocketClientManager event.WSClientGetter) *DashboardSendEvent {
	d := &DashboardSendEvent{
		WebsocketClientManager: websocketClientManager,
	}
	return d
}

// Name returns the name of this function. It will always return "SendEvent".
func (d *DashboardSendEvent) Name() string {
	return "SendEvent"
}

// Call creates a function call that sends an event to a websocket client. If no ClientID
// is provided or the clientID cannot be found a javascript exception is raised.
func (d *DashboardSendEvent) Call(ctx context.Context, vm *goja.Runtime) func(c goja.FunctionCall) goja.Value {
	// clientID string, event EventType, payload action.Payload
	return func(c goja.FunctionCall) goja.Value {
		clientID := c.Argument(0).String()
		if clientID == "" {
			panic(panicMessage(vm, fmt.Errorf("clientID is empty"), ""))
		}

		eventType := event.EventType(c.Argument(1).String())
		if eventType == "" {
			panic(panicMessage(vm, fmt.Errorf("eventType is empty"), ""))
		}

		var payload action.Payload
		obj := c.Argument(2).ToObject(vm)

		// This will never error since &key is a pointer to a type.
		_ = vm.ExportTo(obj, &payload)

		event := event.CreateEvent(eventType, payload)

		if d.WebsocketClientManager == nil {
			panic(panicMessage(vm, fmt.Errorf("websocket client manager is nil"), ""))
		}

		sender := d.WebsocketClientManager.Get(clientID)
		if sender == nil {
			panic(panicMessage(vm, fmt.Errorf("unable to find ws client %s", ocontext.WebsocketClientIDFrom(ctx)), ""))
		}

		sender.Send(event)

		return goja.Undefined()
	}
}

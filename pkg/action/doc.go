package action

/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package action defines how plugins can interact with Octant actions. This includes
registration of custom actions as well as holding the public, internal Octant actions
plugins can provide handlers for.

Octant Actions

Octant actions are defined in this package (see actions.go). The actions can be used by plugin
authors to easily add custom behavior for these internal actions.

For example, if you wanted your plugin to react to when the current namespace has been changed, you would
do so using the RequestSetNamespace action.

In your main function before registering your plugin:

	// Set up the action names this plugin handles.
	capabilities := &plugin.Capabilities{
		ActionNames:           []string{action.RequestSetNamespace},
	}

	// Set up the action handler.
	options := []service.PluginOption{
		service.WithActionHandler(handleAction),
	}

Define your handleAction function:

	func handleAction(request *service.ActionRequest) error {
		switch request.ActionName {
			case action.RequestSetNamespace:
				namespace, err := request.Payload.String("namespace")
				// err check, do work
		}
		return nil
	}

*/

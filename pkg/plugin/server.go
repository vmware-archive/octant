/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"github.com/hashicorp/go-plugin"
)

var (
	// Handshake is the handshake configuration for plugins. Will
	// be used the dashboard and the plugin.
	Handshake = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "DASHBOARD_PLUGIN",
		MagicCookieValue: "dashboard",
	}
)

// Serve serves a plugin.
func Serve(service Service) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: plugin.PluginSet{
			Name: &ServicePlugin{Impl: service},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

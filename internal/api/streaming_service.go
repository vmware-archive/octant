/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"fmt"
	"net/http"

	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/config"
)

func streamService(manager api.ClientManager, dashConfig config.Dash) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveStreamingApi(manager, dashConfig, w, r)
	}
}

func serveStreamingApi(manager api.ClientManager, dashConfig config.Dash, w http.ResponseWriter, r *http.Request) {
	_, err := manager.ClientFromRequest(dashConfig, w, r)
	if err != nil {
		if dashConfig != nil {
			logger := dashConfig.Logger()
			logger.WithErr(err).Errorf("create websocket client")
		}
		return
	}
}

// Create dummy websocketService and serveWebsocket
func loadingStreamService(manager api.ClientManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveLoadingStreamingApi(manager, w, r)
	}
}

func serveLoadingStreamingApi(manager api.ClientManager, w http.ResponseWriter, r *http.Request) {
	_, err := manager.TemporaryClientFromLoadingRequest(w, r)
	if err != nil {
		fmt.Println("create loading websocket client")
	}
}

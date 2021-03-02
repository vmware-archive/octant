/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"fmt"
	"net/http"

	"github.com/vmware-tanzu/octant/internal/config"
)

func streamService(manager ClientManager, dashConfig config.Dash) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveStreamingApi(manager, dashConfig, w, r)
	}
}

func serveStreamingApi(manager ClientManager, dashConfig config.Dash, w http.ResponseWriter, r *http.Request) {
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
func loadingStreamService(manager ClientManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveLoadingStreamingApi(manager, w, r)
	}
}

func serveLoadingStreamingApi(manager ClientManager, w http.ResponseWriter, r *http.Request) {
	_, err := manager.TemporaryClientFromLoadingRequest(w, r)
	if err != nil {
		fmt.Println("create loading websocket client")
	}
}

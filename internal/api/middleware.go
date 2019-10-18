/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/vmware-tanzu/octant/internal/log"
	dashstrings "github.com/vmware-tanzu/octant/internal/util/strings"
)

// shouldAllowHost returns true if the incoming request.Host shuold be allowed
// to access the API otherwise false.
func shouldAllowHost(host string, acceptedHosts []string) bool {
	if dashstrings.Contains("0.0.0.0", acceptedHosts) {
		return true
	}
	return dashstrings.Contains(host, acceptedHosts)
}

// rebindHandler is a middleware that will only accept the supplied hosts
func rebindHandler(ctx context.Context, acceptedHosts []string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var host string
			var err error
			if strings.Contains(r.Host, ":") {
				host, _, err = net.SplitHostPort(r.Host)
			} else {
				host = r.Host
			}

			if err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			if !shouldAllowHost(host, acceptedHosts) {
				logger := log.From(ctx)
				logger.Debugf("Requester %s not in accepted hosts: %s\nTo allow this host add it to the OCTANT_ACCEPTED_HOSTS environment variable.", host, acceptedHosts)
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

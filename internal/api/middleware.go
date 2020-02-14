/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/vmware-tanzu/octant/internal/log"
	dashstrings "github.com/vmware-tanzu/octant/internal/util/strings"
)

var defaultPorts = map[string]string{"http": "80", "https": "443"}

// equalASCIIFold returns true if s is equal to t with ASCII case folding as
// defined in RFC 4790. (Source: https://github.com/gorilla/websocket/blob/master/util.go#L176)
func equalASCIIFold(s, t string) bool {
	for s != "" && t != "" {
		sr, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		tr, size := utf8.DecodeRuneInString(t)
		t = t[size:]
		if sr == tr {
			continue
		}
		if 'A' <= sr && sr <= 'Z' {
			sr = sr + 'a' - 'A'
		}
		if 'A' <= tr && tr <= 'Z' {
			tr = tr + 'a' - 'A'
		}
		if sr != tr {
			return false
		}
	}
	return s == t
}

// checkSameOrigin verifies the Host and the Origin
// (Source: https://github.com/gorilla/websocket/issues/398#issuecomment-409983240)
func checkSameOrigin(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}
	if equalASCIIFold(u.Host, r.Host) {
		return true
	}

	defaultPort, ok := defaultPorts[u.Scheme]
	if !ok {
		return false
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err == nil {
		return port == defaultPort && equalASCIIFold(host, r.Host)
	}

	host, port, err = net.SplitHostPort(r.Host)
	if err == nil {
		return port == defaultPort && equalASCIIFold(u.Host, host)
	}

	return false
}

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

			if !checkSameOrigin(r) {
				http.Error(w, "forbidden bad origin", http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

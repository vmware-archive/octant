/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/modules/overview/container"
)

type logEntry struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Message   string    `json:"message,omitempty"`
}

type logResponse struct {
	Entries []logEntry `json:"entries,omitempty"`
}

func containerLogsHandler(ctx context.Context, clusterClient cluster.ClientInterface) http.HandlerFunc {
	logger := log.From(ctx)

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		containerName := vars["container"]
		podName := vars["pod"]
		namespace := vars["namespace"]

		kubeClient, err := clusterClient.KubernetesClient()
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error(), logger)
			return
		}

		lines := make(chan string)
		done := make(chan bool)

		var entries []logEntry

		go func() {
			for line := range lines {
				parts := strings.SplitN(line, " ", 2)
				logTime, err := time.Parse(time.RFC3339, parts[0])
				if err == nil {
					entries = append(entries, logEntry{
						Timestamp: logTime,
						Message:   parts[1],
					})
				}
			}

			done <- true
		}()

		err = container.Logs(r.Context(), kubeClient, namespace, podName, containerName, lines)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error(), logger)
			return
		}

		<-done

		var lr logResponse

		if len(entries) <= 100 {
			lr.Entries = entries
		} else {
			// take last 100 entries from the slice
			lr.Entries = entries[len(entries)-100:]
		}

		if err := json.NewEncoder(w).Encode(&lr); err != nil {
			logger := log.From(ctx)
			logger.With("err", err.Error()).Errorf("unable to encode log entries")
		}
	}
}

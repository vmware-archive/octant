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

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/modules/overview/container"
)

type logEntry struct {
	Timestamp *time.Time `json:"timestamp,omitempty"`
	Message   string     `json:"message,omitempty"`
}

type logResponse struct {
	Entries []logEntry `json:"entries,omitempty"`
}

func containerLogsHandler(ctx context.Context, dashConfig config.Dash) http.HandlerFunc {
	logger := log.From(ctx)

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		containerName := vars["container"]
		podName := vars["pod"]
		namespace := vars["namespace"]

		kubeClient, err := dashConfig.ClusterClient().KubernetesClient()
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error(), logger)
			return
		}

		done := make(chan bool)

		var entries []logEntry

		logStreamer, err := container.NewLogStreamer(ctx, kubeClient, namespace, podName, containerName)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error(), logger)
			return
		}

		logCh := make(chan container.LogEntry)
		go func() {
			err = logStreamer.Stream(ctx, logCh)
		}()

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error(), logger)
			return
		}

		go func() {
			for line := range logCh {
				entry := logEntry{Message: line.Line()}
				parts := strings.SplitN(line.Line(), " ", 2)
				logTime, err := time.Parse(time.RFC3339, parts[0])
				if err == nil {
					entry = logEntry{
						Timestamp: &logTime,
						Message:   parts[1],
					}
				}
				entries = append(entries, entry)
			}
			done <- true
		}()

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

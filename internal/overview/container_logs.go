package overview

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/overview/container"
)

func containerLogsHandler(ctx context.Context, clusterClient cluster.ClientInterface) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		containerName := vars["container"]
		podName := vars["pod"]
		namespace := vars["namespace"]

		kubeClient, err := clusterClient.KubernetesClient()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

package commands

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/heptio/developer-dash/internal/dash"
	"github.com/heptio/go-telemetry/pkg/telemetry"
	"github.com/spf13/cobra"
)

func newDashCmd() *cobra.Command {
	var namespace string
	var uiURL string
	var kubeconfig string

	dashCmd := &cobra.Command{
		Use:   "dash",
		Short: "Show dashboard",
		Long:  `Heptio Kubernetes dashboard`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt)

			runCh := make(chan bool, 1)

			telemetryClient := newTelemetry()
			startTime := time.Now()
			go func() {
				if err := dash.Run(ctx, namespace, uiURL, kubeconfig, telemetryClient); err != nil {
					log.Print(err)
					os.Exit(1)
				}

				runCh <- true
			}()

			select {
			case <-sigCh:
				msDuration := int64(time.Since(startTime) / time.Millisecond)
				telemetryClient.With(telemetry.Labels{"type": "signal"}).SendEvent("dash.shutdown", telemetry.Measurements{
					"duration": msDuration,
					"count":    1,
				})
				log.Print("Shutting dashboard down due to interrupt")
				telemetryClient.Close()
			case <-runCh:
				msDuration := int64(time.Since(startTime) / time.Millisecond)
				telemetryClient.With(telemetry.Labels{"type": "normal"}).SendEvent("dash.shutdown", telemetry.Measurements{
					"duration": msDuration,
					"count":    1,
				})
				log.Print("Dashboard has exited")
				telemetryClient.Close()
			}
		},
	}

	dashCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
	dashCmd.Flags().StringVar(&uiURL, "ui-url", "", "UI URL")

	if home := homeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	dashCmd.Flags().StringVar(&kubeconfig, "kubeconfig", kubeconfig, "absolute path to kubeconfig file")

	return dashCmd
}

func newTelemetry() telemetry.Interface {
	if _, ok := os.LookupEnv("DASH_DISABLE_TELEMETRY"); ok {
		return &telemetry.NilClient{}
	}

	telemetryAddress := os.Getenv("DASH_TELEMETRY_ADDRESS")
	if telemetryAddress == "" {
		telemetryAddress = telemetry.DefaultAddress
	}

	telemetryClient, err := telemetry.NewClient(telemetryAddress, 10*time.Second, log.New(os.Stderr, "TELEMETRY: ", log.LstdFlags))
	if err != nil {
		log.Print("failed creating telemetry client", err)
		return &telemetry.NilClient{}
	}

	return telemetryClient
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	return os.Getenv("USERPROFILE")
}

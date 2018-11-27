package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/heptio/developer-dash/internal/dash"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/go-telemetry/pkg/telemetry"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/clientcmd"
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

			// Set to 9 for verbose.
			// TODO control via cobra verbose flag along with zap.
			flag.CommandLine.Parse([]string{"-logtostderr", "-v", "0"}) // Set glog to verbose

			z, err := zap.NewDevelopment()
			if err != nil {
				fmt.Printf("failed to initialize logger: %v\n", err)
				os.Exit(1)
			}
			defer z.Sync()
			logger := log.Wrap(z.Sugar())

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt)

			runCh := make(chan bool, 1)

			telemetryClient := newTelemetry(logger)
			startTime := time.Now()

			go func() {
				if err := dash.Run(ctx, namespace, uiURL, kubeconfig, logger, telemetryClient); err != nil {
					logger.Errorf("running dashboard: %v", err)
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
				logger.Debugf("Shutting dashboard down due to interrupt")
				telemetryClient.Close()
			case <-runCh:
				msDuration := int64(time.Since(startTime) / time.Millisecond)
				telemetryClient.With(telemetry.Labels{"type": "normal"}).SendEvent("dash.shutdown", telemetry.Measurements{
					"duration": msDuration,
					"count":    1,
				})
				logger.Debugf("Dashboard has exited")
				telemetryClient.Close()
			}
		},
	}

	dashCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Initial namespace")
	dashCmd.Flags().StringVar(&uiURL, "ui-url", "", "UI URL")

	kubeconfig = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()

	dashCmd.Flags().StringVar(&kubeconfig, "kubeconfig", kubeconfig, "absolute path to kubeconfig file")

	return dashCmd
}

func newTelemetry(logger log.Logger) telemetry.Interface {
	if _, ok := os.LookupEnv("DASH_DISABLE_TELEMETRY"); ok {
		return &telemetry.NilClient{}
	}

	telemetryAddress := os.Getenv("DASH_TELEMETRY_ADDRESS")
	if telemetryAddress == "" {
		telemetryAddress = telemetry.DefaultAddress
	}

	telemetryClient, err := telemetry.NewClient(telemetryAddress, 10*time.Second, logger.Named("telemetry"))
	if err != nil {
		logger.Errorf("failed creating telemetry client", err)
		return &telemetry.NilClient{}
	}

	return telemetryClient
}

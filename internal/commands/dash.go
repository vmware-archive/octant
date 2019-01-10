package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/heptio/developer-dash/internal/dash"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/go-telemetry/pkg/telemetry"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/tools/clientcmd"
)

func newDashCmd() *cobra.Command {
	var namespace string
	var uiURL string
	var kubeconfig string
	var verboseLevel int

	dashCmd := &cobra.Command{
		Use:   "dash",
		Short: "Show dashboard",
		Long:  `Heptio Kubernetes dashboard`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Configure glog verbosity (in client-go)
			// flag.CommandLine.Parse([]string{"-logtostderr", "-v", strconv.Itoa(verboseLevel)}) // Set glog to verbose
			// TODO how does this work in klog??

			z, err := newZapLogger(verboseLevel)
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

	dashCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "initial namespace")
	dashCmd.Flags().StringVar(&uiURL, "ui-url", "", "dashboard url")
	dashCmd.Flags().CountVarP(&verboseLevel, "verbose", "v", "verbosity level")

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

// Returns a new zap logger, setting level according to the provided
// verbosity level as an offset of the base level, Info.
// i.e. verboseLevel==0, level==Info
//      verboseLevel==1, level==Debug
func newZapLogger(verboseLevel int) (*zap.Logger, error) {
	level := zapcore.InfoLevel - zapcore.Level(verboseLevel)
	if level < zapcore.DebugLevel || level > zapcore.FatalLevel {
		level = zapcore.DebugLevel
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}

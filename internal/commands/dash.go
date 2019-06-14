package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/heptio/developer-dash/internal/dash"
	"github.com/heptio/developer-dash/internal/log"
)

func newClusterEyeCmd() *cobra.Command {
	var namespace string
	var uiURL string
	var kubeConfig string
	var verboseLevel int
	var enableOpenCensus bool
	var initialContext string

	clusterEyeCmd := &cobra.Command{
		Use:   "clustereye",
		Short: "clustereye kubernetes dashboard",
		Long:  "clustereye is a dashboard for high bandwidth cluster analysis operations",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// TODO enable support for klog

			z, err := newZapLogger(verboseLevel)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
				os.Exit(1)
			}
			defer z.Sync()
			logger := log.Wrap(z.Sugar())

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt)

			runCh := make(chan bool, 1)

			shutdownCh := make(chan bool, 1)

			go func() {
				options := dash.Options{
					EnableOpenCensus: enableOpenCensus,
					KubeConfig:       kubeConfig,
					Namespace:        namespace,
					FrontendURL:      uiURL,
					Context:          initialContext,
				}

				if err := dash.Run(ctx, logger, shutdownCh, options); err != nil {
					logger.WithErr(err).Errorf("dashboard failed")
					os.Exit(1)
				}

				runCh <- true
			}()

			select {
			case <-sigCh:
				logger.Debugf("Shutting dashboard down due to interrupt")
				cancel()
				// TODO implement graceful shutdown semantics

				<-shutdownCh
			case <-runCh:
				logger.Debugf("Dashboard has exited")
			}
		},
	}

	clusterEyeCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "initial namespace")
	clusterEyeCmd.Flags().StringVar(&uiURL, "ui-url", "", "dashboard url")
	clusterEyeCmd.Flags().CountVarP(&verboseLevel, "verbose", "v", "verbosity level")
	clusterEyeCmd.Flags().BoolVarP(&enableOpenCensus, "enable-opencensus", "c", false, "enable open census")
	clusterEyeCmd.Flags().StringVarP(&initialContext, "context", "", "", "initial context")

	kubeConfig = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()

	clusterEyeCmd.Flags().StringVar(&kubeConfig, "kubeConfig", kubeConfig, "absolute path to kubeConfig file")

	return clusterEyeCmd
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

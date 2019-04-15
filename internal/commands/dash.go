package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/heptio/developer-dash/internal/dash"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/tools/clientcmd"
)

func newSugarloafCmd() *cobra.Command {
	var namespace string
	var uiURL string
	var kubeconfig string
	var verboseLevel int
	var enableOpenCensus bool

	sugarloafCmd := &cobra.Command{
		Use:   "sugarloaf",
		Short: "sugarloaf kubernetes dashboard",
		Long:  "sugarloaf is a dashboard for high bandwidth cluster analysis operations",
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
					KubeConfig:       kubeconfig,
					Namespace:        namespace,
					FrontendURL:      uiURL,
				}

				if err := dash.Run(ctx, logger, shutdownCh, options); err != nil {
					logger.Errorf("running dashboard: %v", err)
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

	sugarloafCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "initial namespace")
	sugarloafCmd.Flags().StringVar(&uiURL, "ui-url", "", "dashboard url")
	sugarloafCmd.Flags().CountVarP(&verboseLevel, "verbose", "v", "verbosity level")
	sugarloafCmd.Flags().BoolVarP(&enableOpenCensus, "enable-opencensus", "c", false, "enable open census")

	kubeconfig = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()

	sugarloafCmd.Flags().StringVar(&kubeconfig, "kubeconfig", kubeconfig, "absolute path to kubeconfig file")

	return sugarloafCmd
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

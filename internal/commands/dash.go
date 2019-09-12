/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"context"
	"flag"
	"fmt"
	golog "log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	"github.com/vmware/octant/internal/dash"
	"github.com/vmware/octant/internal/log"
)

func newOctantCmd() *cobra.Command {
	var namespace string
	var uiURL string
	var kubeConfig string
	var verboseLevel int
	var enableOpenCensus bool
	var initialContext string
	var klogVerbosity int
	var clientQPS float32
	var clientBurst int

	octantCmd := &cobra.Command{
		Use:   "octant",
		Short: "octant kubernetes dashboard",
		Long:  "octant is a dashboard for high bandwidth cluster analysis operations",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// TODO enable support for klog
			z, err := newZapLogger(verboseLevel)
			if err != nil {
				golog.Printf("failed to initialize logger: %v", err)
				os.Exit(1)
			}
			defer func() {
				// this fails, but it should be safe to ignore according
				// to https://github.com/uber-go/zap/issues/328
				_ = z.Sync()
			}()

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
					ClientQPS:        clientQPS,
					ClientBurst:      clientBurst,
				}

				if klogVerbosity > 0 {
					klog.InitFlags(nil)
					verbosityOpt := fmt.Sprintf("-v=%d", klogVerbosity)
					if err := flag.CommandLine.Parse([]string{verbosityOpt, "-logtostderr=true"}); err != nil {
						logger.WithErr(err).Errorf("unable to parse klog flags")
					}

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

	octantCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "initial namespace")
	octantCmd.Flags().StringVar(&uiURL, "ui-url", "", "dashboard url")
	octantCmd.Flags().CountVarP(&verboseLevel, "verbosity", "v", "verbosity level")
	octantCmd.Flags().BoolVarP(&enableOpenCensus, "enable-opencensus", "c", false, "enable open census")
	octantCmd.Flags().StringVarP(&initialContext, "context", "", "", "initial context")
	octantCmd.Flags().IntVarP(&klogVerbosity, "klog-verbosity", "", 0, "klog verbosity level")
	octantCmd.Flags().Float32VarP(&clientQPS, "client-qps", "", 200, "maximum QPS for client")
	octantCmd.Flags().IntVarP(&clientBurst, "client-burst", "", 400, "maximum burst for client throttle")

	kubeConfig = os.Getenv("KUBECONFIG")
	if kubeConfig == "" {
		kubeConfig = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	}

	octantCmd.Flags().StringVar(&kubeConfig, "kubeconfig", kubeConfig, "absolute path to kubeConfig file")

	return octantCmd
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

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/dash"
)

func newOctantCmd(version string) *cobra.Command {
	octantCmd := &cobra.Command{
		Use:   "octant",
		Short: "octant kubernetes dashboard",
		Long:  "octant is a dashboard for high bandwidth cluster analysis operations",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := bindViper(cmd); err != nil {
				golog.Printf("unable to bind flags: %v", err)
				os.Exit(1)
			}

			logLevel := 0
			if viper.GetBool("verbose") {
				logLevel = 1
			}

			z, err := newZapLogger(logLevel)
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

			logger.Debugf("disable-open-browser: %s", viper.Get("disable-open-browser"))

			if viper.GetString("kubeconfig") == "" {
				viper.Set("kubeconfig", clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename())
			}

			go func() {
				options := dash.Options{
					DisableClusterOverview: viper.GetBool("disable-cluster-overview"),
					EnableOpenCensus:       viper.GetBool("enable-opencensus"),
					KubeConfig:             viper.GetString("kubeconfig"),
					Namespace:              viper.GetString("namespace"),
					FrontendURL:            viper.GetString("ui-url"),
					Context:                viper.GetString("context"),
					ClientQPS:              float32(viper.GetFloat64("client-qps")),
					ClientBurst:            viper.GetInt("client-burst"),
					UserAgent:              fmt.Sprintf("octant/%s", version),
				}

				klogVerbosity := viper.GetString("klog-verbosity")
				var klogOpts []string

				klogFlagSet := flag.NewFlagSet("klog", flag.ContinueOnError)
				if klogVerbosity == "" {
					// klog's output is not helpful to Octant, so send it to the ether.
					klogOpts = append(klogOpts,
						fmt.Sprintf("-logtostderr=false"),
						fmt.Sprintf("-alsologtostderr=false"),
					)
				} else {
					klogOpts = append(klogOpts,
						fmt.Sprintf("-v=%s", klogVerbosity),
						fmt.Sprintf("-logtostderr=true"),
						fmt.Sprintf("-alsologtostderr=true"),
					)
				}

				klog.InitFlags(klogFlagSet)

				_ = klogFlagSet.Parse(klogOpts)

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
				// TODO implement graceful shutdown semantics (GH#494)

				<-shutdownCh
			case <-runCh:
				logger.Debugf("Dashboard has exited")
			}
		},
	}

	// All flags can also be environment variables by adding the OCTANT_ prefix
	// and replacing - with _. Example: OCTANT_DISABLE_CLUSTER_OVERVIEW
	octantCmd.Flags().SortFlags = false

	octantCmd.Flags().StringP("context", "", "", "initial context")
	octantCmd.Flags().BoolP("disable-cluster-overview", "", false, "disable cluster overview")
	octantCmd.Flags().BoolP("enable-feature-applications", "", false, "enable applications feature")
	octantCmd.Flags().String("kubeconfig", "", "absolute path to kubeConfig file")
	octantCmd.Flags().StringP("namespace", "n", "", "initial namespace")
	octantCmd.Flags().StringP("plugin-path", "", "", "plugin path")
	octantCmd.Flags().BoolP("verbose", "v", false, "turn on debug logging")

	octantCmd.Flags().StringP("accepted-hosts", "", "", "accepted hosts list [DEV]")
	octantCmd.Flags().Float32P("client-qps", "", 200, "maximum QPS for client [DEV]")
	octantCmd.Flags().IntP("client-burst", "", 400, "maximum burst for client throttle [DEV]")
	octantCmd.Flags().BoolP("disable-open-browser", "", false, "disable automatic launching of the browser [DEV]")
	octantCmd.Flags().BoolP("enable-opencensus", "c", false, "enable open census [DEV]")
	octantCmd.Flags().IntP("klog-verbosity", "", 0, "klog verbosity level [DEV]")
	octantCmd.Flags().StringP("listener-addr", "", "", "listener address for the octant frontend [DEV]")
	octantCmd.Flags().StringP("local-content", "", "", "local content path [DEV]")
	octantCmd.Flags().StringP("proxy-frontend", "", "", "url to send frontend request to [DEV]")
	octantCmd.Flags().String("ui-url", "", "dashboard url [DEV]")

	return octantCmd
}

func bindViper(cmd *cobra.Command) error {
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("OCTANT")
	viper.AutomaticEnv()

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	if err := viper.BindEnv("xdg-config-home", "XDG_CONFIG_HOME"); err != nil {
		return err
	}
	if err := viper.BindEnv("home", "HOME"); err != nil {
		return err
	}
	if err := viper.BindEnv("local-app-data", "LOCALAPPDATA"); err != nil {
		return err
	}
	if err := viper.BindEnv("kubeconfig", "KUBECONFIG"); err != nil {
		return err
	}

	return nil
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

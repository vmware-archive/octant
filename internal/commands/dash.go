package commands

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/heptio/developer-dash/internal/dash"
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

			go func() {
				if err := dash.Run(ctx, namespace, uiURL, kubeconfig); err != nil {
					log.Print(err)
					os.Exit(1)
				}

				runCh <- true
			}()

			select {
			case <-sigCh:
				log.Print("Shutting dashboard down due to interrupt")
			case <-runCh:
				log.Print("Dashboard has exited")
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

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	return os.Getenv("USERPROFILE")
}

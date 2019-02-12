package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/overview/container"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	namespace := "default"
	podName := "wonderful-wallaby-mysql-d67d9d78b-sqzlf"
	containerName := "wonderful-wallaby-mysql"

	if err := run(namespace, podName, containerName); err != nil {
		log.Panic(fmt.Sprintf("run error: %v", err))
	}

}

func run(namespace, podName, containerName string) error {
	kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()

	clusterClient, err := cluster.FromKubeconfig(kubeconfig)
	if err != nil {
		return errors.Wrap(err, "create cluster client")
	}

	kubeClient, err := clusterClient.KubernetesClient()
	if err != nil {
		return errors.Wrap(err, "build kubernetes client")
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	lines := make(chan string)
	done := make(chan bool)

	go func() {
		for line := range lines {
			fmt.Println(line)
		}

		done <- true
	}()

	err = container.Logs(ctx, kubeClient, namespace, podName, containerName, lines)
	if err != nil {
		return errors.Wrap(err, "do logs")
	}

	<-done

	return nil
}

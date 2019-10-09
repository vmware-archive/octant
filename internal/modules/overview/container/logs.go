/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package container

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// durContainerUpWait is the time to wait before checking if a container has started
	durContainerUpWait = 1 * time.Second
)

func Logs(ctx context.Context, client kubernetes.Interface, namespace, podName, container string, logCh chan<- string) error {
	lp := logPrinter{
		client:    client,
		namespace: namespace,
		podName:   podName,
		container: container,
	}

	return lp.logs(ctx, logCh)
}

type logPrinter struct {
	client kubernetes.Interface

	namespace string
	podName   string
	container string
}

func (lp *logPrinter) logs(ctx context.Context, ch chan<- string) error {
	if ch == nil {
		return errors.New("channel is nil")
	}

	defer close(ch)

	for ctx.Err() == nil {
		hasStarted, err := lp.containerHasStarted()
		if err != nil {
			return errors.Wrap(err, "check if container has started")
		}

		if hasStarted {
			break
		}

		time.Sleep(durContainerUpWait)
	}

	stream, err := lp.stream()
	if err != nil {
		return errors.Wrap(err, "stream container logs")
	}
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	for ctx.Err() == nil && scanner.Scan() {
		ch <- scanner.Text()
	}

	if scanner.Err() != nil {
		return errors.Wrap(err, "scanner error")
	}

	return nil
}

func (lp *logPrinter) containerHasStarted() (bool, error) {
	pod, err := lp.client.CoreV1().Pods(lp.namespace).Get(lp.podName, metav1.GetOptions{})
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("get pod %s in %s", lp.podName, lp.namespace))
	}

	for _, status := range pod.Status.ContainerStatuses {
		if status.Name == lp.container && status.State.Waiting == nil {
			return true, nil
		}
	}

	return false, nil
}

func (lp *logPrinter) stream() (io.ReadCloser, error) {
	return lp.client.CoreV1().Pods(lp.namespace).GetLogs(lp.podName, &corev1.PodLogOptions{
		Container:  lp.container,
		Follow:     false,
		Timestamps: true,
	}).Stream()
}

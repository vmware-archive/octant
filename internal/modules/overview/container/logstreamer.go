/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package container

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// durContainerUpWait is the time to wait before checking if a container has started
	durContainerUpWait = 1 * time.Second
)

type logStreamer struct {
	namespace  string
	pod        string
	containers []string
	stream     chan LogEntry

	ctx      context.Context
	cancelFn context.CancelFunc
	client   kubernetes.Interface
}

var _ LogStreamer = (*logStreamer)(nil)

func NewLogStreamer(ctx context.Context, client kubernetes.Interface, namespace string, podName string, containerNames ...string) (*logStreamer, error) {
	ctx, cancelFn := context.WithCancel(ctx)

	// Remove and replace with a container lookup method if containerNames length is 0.
	if len(containerNames) == 0 {
		return nil, errors.New("no container names provided")
	}

	return &logStreamer{
		namespace:  namespace,
		pod:        podName,
		containers: containerNames,
		client:     client,
		ctx:        ctx,
		cancelFn:   cancelFn,
	}, nil
}

func (s *logStreamer) Names() []string {
	return s.containers
}

func (s *logStreamer) Stream(ctx context.Context, logCh chan<- LogEntry) error {
	defer close(logCh)

	for ctx.Err() == nil {
		hasStarted, err := s.containerHasStarted(s.Names()[0])
		if err != nil {
			return fmt.Errorf("check if container has started: %w", err)
		}

		if hasStarted {
			break
		}

		time.Sleep(durContainerUpWait)
	}

	container := s.containers[0]
	stream, err := s.containerStream(container)

	if err != nil {
		return fmt.Errorf("stream container logs: %w", err)
	}
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	for ctx.Err() == nil && scanner.Scan() {
		entry := NewLogEntry(container, scanner.Text())
		logCh <- entry
	}

	if scanner.Err() != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

func (s *logStreamer) Close() {
	s.cancelFn()
}

func (s *logStreamer) containerStream(container string) (io.ReadCloser, error) {
	return s.client.CoreV1().Pods(s.namespace).GetLogs(s.pod, &corev1.PodLogOptions{
		Container: container,
		// Change this to true when we start to actually stream.
		Follow:     false,
		Timestamps: true,
	}).Stream()
}

func (s *logStreamer) containerHasStarted(container string) (bool, error) {
	pod, err := s.client.CoreV1().Pods(s.namespace).Get(s.pod, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("get pod %s in %s: %w", s.pod, s.namespace, err)
	}

	for _, status := range pod.Status.ContainerStatuses {
		if status.Name == container && status.State.Waiting == nil {
			return true, nil
		}
	}

	return false, nil
}

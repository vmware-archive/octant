/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package container

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/store"
)

type logStreamer struct {
	namespace    string
	pod          string
	containers   []string
	sinceSeconds *int64
	creationTime *v1.Time
	stream       chan LogEntry

	ctx      context.Context
	cancelFn context.CancelFunc
	config   config.Dash
	wg       sync.WaitGroup
}

var _ LogStreamer = (*logStreamer)(nil)

// NewLogStreamer returns an instance of a logStream configured to stream logs for the given namespace/pod/container(s).
func NewLogStreamer(ctx context.Context, dashConfig config.Dash, key store.Key, sinceSeconds int64, containerNames ...string) (*logStreamer, error) {
	ctx, cancelFn := context.WithCancel(ctx)

	if shouldFetchContainerNames(containerNames) {
		// reset containerNames since it contains only empty entries
		containerNames = []string{}

		object, err := dashConfig.ObjectStore().Get(ctx, key)
		if err != nil {
			cancelFn()
			return nil, fmt.Errorf("getting pod from objectstore: %w", err)
		}

		var pod corev1.Pod
		err = kubernetes.FromUnstructured(object, &pod)

		if err != nil {
			cancelFn()
			return nil, fmt.Errorf("converting unstructured: %w", err)
		}

		for _, container := range pod.Spec.Containers {
			containerNames = append(containerNames, container.Name)
		}
	}

	var creationTime *v1.Time
	if sinceSeconds < 0 {
		object, err := dashConfig.ObjectStore().Get(ctx, key)
		if err != nil {
			cancelFn()
			return nil, fmt.Errorf("getting pod from objectstore: %w", err)
		}

		var pod corev1.Pod
		err = kubernetes.FromUnstructured(object, &pod)

		if err != nil {
			cancelFn()
			return nil, fmt.Errorf("converting unstructured: %w", err)
		}

		creationTime = &v1.Time{Time: pod.CreationTimestamp.Time}
	}

	return &logStreamer{
		namespace:    key.Namespace,
		pod:          key.Name,
		containers:   containerNames,
		sinceSeconds: &sinceSeconds,
		creationTime: creationTime,
		config:       dashConfig,
		ctx:          ctx,
		cancelFn:     cancelFn,
	}, nil
}

// shouldFetchContainerNames checks if the list of containerNames is empty
// or if the containerNames contains a singular entry that is equal to empty string
// and is used to determine if the LogStream should fetch the list of containers to stream logs for.
func shouldFetchContainerNames(containerNames []string) bool {
	if len(containerNames) == 0 {
		return true
	}

	if len(containerNames) == 1 && containerNames[0] == "" {
		return true
	}

	return false
}

// Names returns a list of container names that the log streamer is streaming logs for.
func (s *logStreamer) Names() []string {
	if s.containers == nil {
		return []string{}
	}
	return s.containers
}

// Stream takes a context and a log channel for writing log entries to. Stream
// will handle closing any open streams and closing the log channel when an error
// or EOF is encountered.
func (s *logStreamer) Stream(ctx context.Context, logCh chan<- LogEntry) {
	for _, container := range s.containers {
		container := container
		stream, err := s.containerStream(container)

		if err != nil {
			s.config.Logger().Errorf("unable to stream logs for %s: %w", container, err)
			continue
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			scanner := bufio.NewScanner(stream)
			for ctx.Err() == nil && scanner.Scan() {
				entry := NewLogEntry(container, scanner.Text())
				logCh <- entry
			}
			return
		}()
	}

	go func() {
		s.wg.Wait()
		s.Close(logCh)
		return
	}()

	return
}

// Close calls the cancel function and closes the stream.
func (s *logStreamer) Close(logCh chan<- LogEntry) {
	close(logCh)
	s.cancelFn()
}

func (s *logStreamer) containerStream(container string) (io.ReadCloser, error) {
	client, err := s.config.ClusterClient().KubernetesClient()
	if err != nil {
		return nil, err
	}

	options := &corev1.PodLogOptions{
		Container:  container,
		Follow:     true,
		Timestamps: true,
	}
	if s.creationTime != nil {
		options.SinceTime = s.creationTime
	} else {
		options.SinceSeconds = s.sinceSeconds
	}

	request := client.CoreV1().Pods(s.namespace).GetLogs(s.pod, options)
	return request.Stream(s.ctx)
}

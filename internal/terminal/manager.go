/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"
	"strings"
	"unicode"

	corev1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/store"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

//go:generate mockgen -source=manager.go -destination=./fake/mock_interface.go -package=fake github.com/vmware-tanzu/octant/internal/terminal TerminalManager

type manager struct {
	restClient  rest.Interface
	config      *rest.Config
	objectStore store.Store
	instances   map[string]Instance
}

var _ Manager = (*manager)(nil)

// NewTerminalManager creates a concrete TerminalMananger
func NewTerminalManager(ctx context.Context, client cluster.ClientInterface, objectStore store.Store) (*manager, error) {
	restClient, err := client.RESTClient()
	if err != nil {
		return nil, errors.Wrap(err, "fetching RESTClient")
	}

	tm := &manager{
		restClient:  restClient,
		config:      client.RESTConfig(),
		objectStore: objectStore,
		instances:   map[string]Instance{},
	}
	return tm, nil
}

func (tm *manager) Create(ctx context.Context, logger log.Logger, key store.Key, container string, command string) (Instance, error) {
	logger.Debugf("create")

	t := NewTerminalInstance(ctx, key, container, command)
	tm.instances[t.ID()] = t

	pod, ok, err := tm.objectStore.Get(ctx, key)
	if err != nil {
		logger.Errorf("objectStore: %s", err)
		return nil, err
	}
	if !ok {
		return nil, errors.New("pod not found")
	}

	logger.Debugf("prePOST")
	req := tm.restClient.Post().
		Resource("pods").
		Name(pod.GetName()).
		Namespace(pod.GetNamespace()).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: container,
		Command:   parseCommand(command),
		Stdin:     t.Stdin() != nil,
		Stdout:    t.Stdout() != nil,
		Stderr:    t.Stderr() != nil,
		TTY:       false,
	}, scheme.ParameterCodec)

	rc, err := remotecommand.NewSPDYExecutor(tm.config, "POST", req.URL())
	if err != nil {
		logger.Errorf("executor: %s", err)
		return nil, err
	}

	logger.Debugf("postPOST")

	opts := remotecommand.StreamOptions{
		Stdin:  t.Stdin(),
		Stdout: t.Stdout(),
		Stderr: t.Stderr(),
		Tty:    false,
		//TerminalSizeQueue: remotecommand.TerminalSizeQueue,
	}

	err = rc.Stream(opts)
	if err != nil {
		logger.Errorf("streaming: %s", err)
		return nil, err
	}

	t.Stream(ctx, logger)
	return t, nil
}

func (tm *manager) Get(ctx context.Context, id string) (Instance, bool) {
	v, ok := tm.instances[id]
	return v, ok
}

func (tm *manager) List(ctx context.Context) []Instance {
	instances := make([]Instance, 0, len(tm.instances))
	for _, instance := range tm.instances {
		instances = append(instances, instance)
	}
	return instances
}

func (tm *manager) StopAll(ctx context.Context) error {
	for _, instance := range tm.instances {
		instance.Stop(ctx)
	}
	return nil
}

func parseCommand(command string) []string {
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}
	return strings.FieldsFunc(command, f)
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
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

// Manager defines the interface for querying terminal instance.
type Manager interface {
	List(namespace string) []Instance
	Get(id string) (Instance, bool)
	Delete(id string)
	Create(ctx context.Context, logger log.Logger, key store.Key, container string, command string, tty bool) (Instance, error)
	StopAll() error
}

type manager struct {
	restClient  rest.Interface
	config      *rest.Config
	objectStore store.Store
	instances   sync.Map
}

var _ Manager = (*manager)(nil)

// NewTerminalManager creates a concrete TerminalMananger
func NewTerminalManager(ctx context.Context, client cluster.ClientInterface, objectStore store.Store) (Manager, error) {
	restClient, err := client.RESTClient()
	if err != nil {
		return nil, errors.Wrap(err, "fetching RESTClient")
	}

	tm := &manager{
		restClient:  restClient,
		config:      client.RESTConfig(),
		objectStore: objectStore,
	}
	return tm, nil
}

func (tm *manager) Create(ctx context.Context, logger log.Logger, key store.Key, container, command string, tty bool) (Instance, error) {
	t := NewTerminalInstance(ctx, logger, key, container, command, tty)
	tm.instances.Store(t.ID(), t)

	pod, ok, err := tm.objectStore.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("pod not found")
	}

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
		TTY:       tty,
	}, scheme.ParameterCodec)

	rc, err := remotecommand.NewSPDYExecutor(tm.config, "POST", req.URL())
	if err != nil {
		t.SetExitMessage(fmt.Sprintf("%v", err))
		return nil, err
	}

	opts := remotecommand.StreamOptions{
		Stdin:             t.Stdin(),
		Stdout:            t.Stdout(),
		Stderr:            t.Stderr(),
		Tty:               tty,
		TerminalSizeQueue: t.SizeQueue(),
	}

	go func() {
		logger.Debugf("running stream command")
		err = rc.Stream(opts)
		if err != nil {
			t.SetExitMessage(fmt.Sprintf("%s", err))
			logger.Errorf("streaming: %+v", err)
		}
		t.Stop()
		logger.Debugf("stopping stream command")
	}()

	return t, nil
}

func (tm *manager) Get(id string) (Instance, bool) {
	v, ok := tm.instances.Load(id)
	if !ok {
		return nil, ok
	}
	return v.(Instance), ok
}

func (tm *manager) List(namespace string) []Instance {
	instances := []Instance{}
	tm.instances.Range(func(k interface{}, v interface{}) bool {
		instance := v.(Instance)
		if namespace == "" {
			instances = append(instances, instance)

		} else if instance.Key().Namespace == namespace {
			instances = append(instances, instance)
		}
		return true
	})
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].CreatedAt().Before(instances[j].CreatedAt())
	})
	return instances
}

func (tm *manager) Delete(id string) {
	t, ok := tm.Get(id)
	if ok {
		t.Stop()
		tm.instances.Delete(id)
	}
}

func (tm *manager) StopAll() error {
	tm.instances.Range(func(k interface{}, v interface{}) bool {
		v.(Instance).Stop()
		return true
	})
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

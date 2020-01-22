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
	Create(ctx context.Context, logger log.Logger, key store.Key, container string, command string) (Instance, error)
	Select(ctx context.Context) chan Instance
	ForceUpdate(id string)
	SendScrollback(id string) bool
	SetScrollback(id string, send bool)
	StopAll() error
}

type manager struct {
	restClient   rest.Interface
	config       *rest.Config
	objectStore  store.Store
	instances    sync.Map
	scrollback   sync.Map
	chanInstance chan Instance
}

var _ Manager = (*manager)(nil)

// NewTerminalManager creates a concrete TerminalMananger
func NewTerminalManager(ctx context.Context, client cluster.ClientInterface, objectStore store.Store) (Manager, error) {
	restClient, err := client.RESTClient()
	if err != nil {
		return nil, errors.Wrap(err, "fetching RESTClient")
	}

	tm := &manager{
		restClient:   restClient,
		config:       client.RESTConfig(),
		objectStore:  objectStore,
		chanInstance: make(chan Instance, 100),
	}
	return tm, nil
}

// SetScrollback sets the current scrollback state for a terminal instance.
func (tm *manager) SetScrollback(id string, sendScrollback bool) {
	tm.scrollback.Store(id, sendScrollback)
}

// SendScrollback returns the current scrollback state for a terminal instance.
func (tm *manager) SendScrollback(id string) bool {
	v, ok := tm.scrollback.Load(id)
	if !ok {
		return false
	}
	return v.(bool)
}

// ForceUpdate sends the instance for the given ID on to the instance work channel.
func (tm *manager) ForceUpdate(id string) {
	i, ok := tm.Get(id)
	if ok {
		tm.chanInstance <- i
	}
}

func (tm *manager) Select(ctx context.Context) chan Instance {
	return tm.chanInstance
}

func (tm *manager) Create(ctx context.Context, logger log.Logger, key store.Key, container, command string) (Instance, error) {
	t := NewTerminalInstance(ctx, logger, key, container, command, tm.chanInstance)
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
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	rc, err := remotecommand.NewSPDYExecutor(tm.config, "POST", req.URL())
	if err != nil {
		t.SetExitMessage(fmt.Sprintf("%v", err))
		return nil, err
	}

	pty := t.PTY()
	opts := remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
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

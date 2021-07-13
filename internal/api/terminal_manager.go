/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu/octant/internal/util/kubernetes"

	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/terminal"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/store"
)

const (
	readBufferSize         = 4096
	RequestTerminalCommand = "action.octant.dev/sendTerminalCommand"
	RequestTerminalResize  = "action.octant.dev/sendTerminalResize"
	RequestActiveTerminal  = "action.octant.dev/setActiveTerminal"
)

type terminalStateManager struct {
	client   api.OctantClient
	config   config.Dash
	ctx      context.Context
	instance terminal.Instance

	chanInstance          chan terminal.Instance
	terminalSubscriptions sync.Map
	existingInstance      bool
}

type terminalOutput struct {
	Scrollback  []byte `json:"scrollback,omitempty"`
	Line        []byte `json:"line,omitempty"`
	ExitMessage []byte `json:"exitMessage,omitempty"`
}

var _ StateManager = (*terminalStateManager)(nil)

// NewTerminalStateManager returns a terminal state manager.
func NewTerminalStateManager(dashConfig config.Dash) StateManager {
	return &terminalStateManager{
		config:                dashConfig,
		terminalSubscriptions: sync.Map{},
		chanInstance:          make(chan terminal.Instance, 10),
	}
}

// Handlers returns a slice of handlers.
func (s *terminalStateManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: RequestTerminalCommand,
			Handler:     s.SendTerminalCommand,
		},
		{
			RequestType: RequestTerminalResize,
			Handler:     s.SendTerminalResize,
		},
		{
			RequestType: RequestActiveTerminal,
			Handler:     s.SetActiveTerminal,
		},
	}
}

func (s *terminalStateManager) SetActiveTerminal(state octant.State, payload action.Payload) error {
	namespace, err := payload.String("namespace")
	if err != nil {
		return fmt.Errorf("getting namespace from payload: %w", err)
	}
	podName, err := payload.String("podName")
	if err != nil {
		return fmt.Errorf("getting podName from payload: %w", err)
	}

	containerName, err := payload.String("containerName")
	if err != nil {
		return fmt.Errorf("getting containerName from payload: %w", err)
	}

	eventType := event.NewTerminalEventType(namespace, podName, containerName)
	key := store.KeyFromGroupVersionKind(gvk.Pod)
	key.Name = podName
	key.Namespace = namespace

	if s.instance != nil {
		if s.instance.Key() == key && s.instance.Active() && s.instance.Container() == containerName {
			s.existingInstance = true
			s.chanInstance <- s.instance
			return nil
		}
		// Remove old terminal instance
		prevEventType := event.NewTerminalEventType(s.instance.Key().Namespace, s.instance.Key().Name, s.instance.Container())
		val, ok := s.terminalSubscriptions.Load(eventType)
		if ok {
			cancelFn, ok := val.(context.CancelFunc)
			if !ok {
				return fmt.Errorf("bad cancelFn conversion for %s", eventType)
			}
			s.terminalSubscriptions.Delete(prevEventType)
			cancelFn()
		}
	}

	val, ok := s.terminalSubscriptions.Load(eventType)
	if ok {
		cancelFn, ok := val.(context.CancelFunc)
		if !ok {
			return fmt.Errorf("bad cancelFn conversion for %s", eventType)
		}
		cancelFn()
	}

	objectStore := s.config.ObjectStore()
	po := &corev1.Pod{}
	pod, err := objectStore.Get(s.ctx, key)
	if err != nil {
		return err
	}
	if pod != nil {
		err := kubernetes.FromUnstructured(pod, po)
		if err != nil {
			return err
		}
	}

	cancelFn := s.startStream(po, key, containerName)
	s.terminalSubscriptions.Store(eventType, cancelFn)
	return nil
}

func (s *terminalStateManager) startStream(pod *corev1.Pod, key store.Key, container string) context.CancelFunc {
	ctx, cancelFn := context.WithCancel(s.ctx)
	logger := log.From(s.ctx).With("startStream", container)

	eventType := event.NewTerminalEventType(pod.Namespace, pod.Name, container)

	commands := []string{"bash", "sh"}
	if IsWindowsContainer(pod) {
		commands = []string{"powershell", "cmd"}
	}

	for _, command := range commands {
		validInstance, err := terminal.NewTerminalInstance(ctx, s.config.ClusterClient(), logger, key, container, command, s.chanInstance)
		if err != nil {
			logger.Debugf("streaming: %+v", err)
			continue
		}

		if validInstance != nil {
			s.instance = validInstance
			go s.sendTerminalEvents(ctx, eventType, s.instance, s.chanInstance)
			break
		}
	}

	if s.instance == nil {
		cancelFn()
	}
	return cancelFn
}

func (s *terminalStateManager) SendTerminalResize(state octant.State, payload action.Payload) error {
	if s.instance == nil {
		return errors.New("terminal instance not found")
	}

	rows, err := payload.Uint16("rows")
	if err != nil {
		return errors.Wrap(err, "extract rows from payload")
	}

	cols, err := payload.Uint16("cols")
	if err != nil {
		return errors.Wrap(err, "extract cols from payload")
	}

	if s.instance.Active() {
		s.instance.Resize(cols, rows)
	}
	return nil
}

func (s *terminalStateManager) SendTerminalCommand(state octant.State, payload action.Payload) error {
	if s.instance == nil {
		return errors.New("terminal instance not found")
	}

	key, err := payload.String("key")
	if err != nil {
		return errors.Wrap(err, "extract key from payload")
	}

	return s.instance.Write([]byte(key))
}

func (s *terminalStateManager) Start(ctx context.Context, state octant.State, client api.OctantClient) {
	s.client = client
	s.ctx = ctx
}

func (s *terminalStateManager) sendTerminalEvents(ctx context.Context, terminalEventType event.EventType, instance terminal.Instance, terminalCh <-chan terminal.Instance) {
	ctx, cancelFn := context.WithCancel(s.ctx)
	for {
		select {
		case <-ctx.Done():
			cancelFn()
			return
		case t := <-terminalCh:
			event, err := newEvent(ctx, t, !s.instance.Active() || s.existingInstance)
			if err != nil {
				break
			}
			s.client.Send(event)
			s.existingInstance = false
		case <-time.After(25 * time.Millisecond):
			break
		}
	}
}

func newEvent(ctx context.Context, t terminal.Instance, sendScrollback bool) (event.Event, error) {
	line, err := t.Read(readBufferSize)
	if err != nil {
		t.SetExitMessage(fmt.Sprintf("%v\n", err))
		t.Stop()
		return event.Event{}, errors.Wrap(err, "read error")
	}

	if line == nil && !sendScrollback {
		return event.Event{}, errors.New("no scrollback or line")
	}

	key := t.Key()
	eventType := event.NewTerminalEventType(key.Namespace, key.Name, t.Container())
	data := terminalOutput{Line: line}

	if sendScrollback {
		data.Scrollback = t.Scrollback()
		if !t.Active() {
			msg := t.ExitMessage()
			if msg != "" {
				data.Scrollback = append(data.Scrollback, []byte("\n"+msg)...)
				data.ExitMessage = []byte(msg)
			} else {
				data.Scrollback = []byte("\n" + "(process exited)")
				data.ExitMessage = data.Scrollback
			}
			t.Stop()
		}
	}

	return event.Event{
		Type: eventType,
		Data: data,
		Err:  nil,
	}, nil
}

func IsWindowsContainer(pod *corev1.Pod) bool {
	var hasNodeSelector, hasToleration bool
	windowsToleration := corev1.Toleration{
		Key:      "os",
		Operator: corev1.TolerationOpEqual,
		Value:    "windows",
		Effect:   corev1.TaintEffectNoSchedule,
	}
	for k, v := range pod.Spec.NodeSelector {
		if k == "kubernetes.io/os" && v == "windows" {
			hasNodeSelector = true
		}
	}
	for _, toleration := range pod.Spec.Tolerations {
		if toleration == windowsToleration {
			hasToleration = true
		}
	}
	return hasNodeSelector && hasToleration
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminalviewer

import (
	"context"
	"github.com/vmware-tanzu/octant/internal/terminal"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// ToComponent converts an object into a terminal component.
func ToComponent(ctx context.Context, object runtime.Object, terminalManager terminal.Manager, logger log.Logger) (*component.Terminal, error) {
	yv, err := new(ctx, object, terminalManager, logger )
	if err != nil {
		return nil, errors.Wrap(err, "create Terminal viewer")
	}

	return yv.ToComponent()
}

// Terminal Viewer is a terminal viewer for objects.
type terminalViewer struct {
	object runtime.Object
	terminalManager terminal.Manager
	context context.Context
	logger log.Logger
}

// New creates an instance of TerminalViewer.
func new(context context.Context, object runtime.Object, terminalManager terminal.Manager, logger log.Logger) (*terminalViewer, error) {
	if object == nil {
		return nil, errors.New("can't create Terminal view for nil object")
	}

	return &terminalViewer{
		context: context,
		object: object,
		terminalManager: terminalManager,
		logger: logger.With("Terminal Viewer", context),
	}, nil
}

// ToComponent converts the Terminal Viewer to a component.
func (tv *terminalViewer) ToComponent() (*component.Terminal, error) {
	pod, err := getPod(tv)
	if err != nil {
		return nil, err
	}

	key, err:=  store.KeyFromObject(tv.object)
	if err != nil {
		return nil, err
	}

	container:= ""
	if len(pod.Spec.Containers) > 0 {
		container= pod.Spec.Containers[0].Name
	}
	t, err := tv.terminalManager.Create(context.Background(), tv.logger, key, container, "/bin/sh", "")

	if err != nil {
		return nil, err
	}

	id := ""
	if t != nil {
		id= t.ID()
	}

	details := component.TerminalDetails{
		Container: container,
		Command:   "/bin/sh",
		UUID:      id,
		Active:    true,
	}
	term := component.NewTerminal(pod.Namespace, "Terminal", pod.Name, details)
	return term, nil
}

func getPod(tv *terminalViewer) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	switch t := tv.object.(type) {
	case *unstructured.Unstructured:
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(t.Object, pod); err != nil {
			return nil, err
		}
	case *corev1.Pod:
		pod = t
	default:
		pod = nil
	}

	if pod == nil {
		return nil, errors.Errorf("can't fetch Terminal from a %T", tv.object)
	}

	return pod, nil
}

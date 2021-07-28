/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminalviewer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// ToComponent converts an object into a terminal component.
func ToComponent(ctx context.Context, object runtime.Object, logger log.Logger, dashConfig config.Dash) (*component.Terminal, error) {
	ecg, err := NewEphemeralContainerGenerator(ctx, dashConfig, logger, object)
	if err != nil {
		return nil, fmt.Errorf("ephemeral container: %w", err)
	}

	if ecg.FeatureEnabled() {
		if err := ecg.UpdateObject(ctx, object); err != nil {
			return nil, err
		}
	}

	tv, err := new(ctx, object, logger)
	if err != nil {
		return nil, errors.Wrap(err, "create Terminal viewer")
	}

	return tv.ToComponent()
}

// Terminal Viewer is a terminal viewer for objects.
type terminalViewer struct {
	object  runtime.Object
	context context.Context
	logger  log.Logger
}

// New creates an instance of TerminalViewer.
func new(context context.Context, object runtime.Object, logger log.Logger) (*terminalViewer, error) {
	if object == nil {
		return nil, errors.New("can't create Terminal view for nil object")
	}

	return &terminalViewer{
		context: context,
		object:  object,
		logger:  logger.With("Terminal Viewer", context),
	}, nil
}

// ToComponent converts the Terminal Viewer to a component.
func (tv *terminalViewer) ToComponent() (*component.Terminal, error) {
	pod, err := getPod(tv)
	if err != nil {
		return nil, err
	}

	container := ""

	var containers []string
	if len(pod.Spec.EphemeralContainers) > 0 {
		for _, ec := range pod.Spec.EphemeralContainers {
			for _, s := range pod.Status.EphemeralContainerStatuses {
				if s.Name == ec.Name {
					if s.State.Terminated == nil {
						containers = append(containers, ec.Name)
					}
				}
			}
		}
		container = pod.Spec.EphemeralContainers[0].Name
	}

	if len(pod.Spec.Containers) > 0 {
		if container == "" {
			container = getFirstContainer(pod).Name
		}
		for _, c := range pod.Spec.Containers {
			for _, s := range pod.Status.ContainerStatuses {
				if s.Name == c.Name {
					if s.State.Terminated == nil {
						containers = append(containers, c.Name)
					}
				}
			}
		}
	}

	details := component.TerminalDetails{
		Container: container,
		Command:   "/bin/sh",
		Active:    true,
	}
	term := component.NewTerminal(pod.Namespace, "Terminal", pod.Name, containers, details)
	return term, nil
}

func getFirstContainer(pod *corev1.Pod) corev1.Container {
	var selected = pod.Spec.Containers[0]

	for _, c := range pod.Spec.Containers {
		var isInit = false
		for _, d := range pod.Spec.InitContainers {
			if d.Name == c.Name {
				isInit = true
				break
			}
		}
		if !isInit {
			selected = c
			break
		}
	}

	return selected
}

func getPod(tv *terminalViewer) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	switch t := tv.object.(type) {
	case *unstructured.Unstructured:
		if err := kubernetes.FromUnstructured(t, pod); err != nil {
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

package logviewer

import (
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ToComponent converts an object into a logviewer component.
func ToComponent(object runtime.Object) (component.Component, error) {
	if object == nil {
		return nil, errors.Errorf("object is nil")
	}

	pod, ok := object.(*corev1.Pod)
	if !ok {
		return nil, errors.New("object is not a pod")
	}

	var containerNames []string

	for _, c := range pod.Spec.InitContainers {
		containerNames = append(containerNames, c.Name)
	}

	for _, c := range pod.Spec.Containers {
		containerNames = append(containerNames, c.Name)
	}

	logsComponent := component.NewLogs(pod.Namespace, pod.Name, containerNames)

	return logsComponent, nil
}

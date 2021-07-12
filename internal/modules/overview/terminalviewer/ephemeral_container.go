package terminalviewer

import (
	"context"
	"errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/log"
)

type EphemeralContainerGenerator struct {
	ctx        context.Context
	dashConfig config.Dash
	logger     log.Logger
	object     runtime.Object
}

func NewEphemeralContainerGenerator(ctx context.Context, dashConfig config.Dash, logger log.Logger, object runtime.Object) (*EphemeralContainerGenerator, error) {
	if object == nil || dashConfig == nil {
		return nil, errors.New("cannot create ephemeral container generator")
	}

	return &EphemeralContainerGenerator{
		ctx:        ctx,
		logger:     logger,
		dashConfig: dashConfig,
		object:     object,
	}, nil
}

func (e *EphemeralContainerGenerator) UpdateObject(ctx context.Context, object runtime.Object) error {
	if object == nil {
		return errors.New("object is nil")
	}

	client, err := e.dashConfig.ClusterClient().KubernetesClient()
	if err != nil {
		return err
	}

	pod := object.(*corev1.Pod)
	pods := client.CoreV1().Pods(pod.Namespace)

	if len(pod.Spec.EphemeralContainers) == 0 {
		ec, err := pods.GetEphemeralContainers(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		container := pod.Spec.Containers[0].Name

		debugContainer := corev1.EphemeralContainer{
			TargetContainerName: container,
			EphemeralContainerCommon: corev1.EphemeralContainerCommon{
				Name:                     "debug-" + container,
				Image:                    "debian",
				Stdin:                    true,
				TTY:                      true,
				TerminationMessagePolicy: corev1.TerminationMessageReadFile,
				ImagePullPolicy:          corev1.PullIfNotPresent,
			},
		}

		ec.EphemeralContainers = append(ec.EphemeralContainers, debugContainer)

		e.logger.Debugf("Creating ephemeral container for: %s", container)
		_, err = pods.UpdateEphemeralContainers(ctx, pod.Name, ec, metav1.UpdateOptions{})
		if err != nil {
			e.logger.Debugf("pod update for ephemeral container: %+v", err)
		}
	}
	return nil
}

func (e *EphemeralContainerGenerator) FeatureEnabled() bool {
	discoveryClient, err := e.dashConfig.ClusterClient().DiscoveryClient()
	if err != nil {
		return false
	}

	_, resourceList, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		return true
	}
	for _, resource := range resourceList {
		for _, apiResource := range resource.APIResources {
			if apiResource.Kind == "EphemeralContainers" {
				return true
			}
		}
	}

	return false
}

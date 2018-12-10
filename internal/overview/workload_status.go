package overview

import (
	"strconv"

	"github.com/heptio/developer-dash/internal/content"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

// ResourceStatus represents the status of a resource.
type ResourceStatus struct {
	Status     content.NodeStatus
	Reason     string
	RelatedUID types.UID
}

// ResourceStatusList is a list of resource statuses (validation results)
type ResourceStatusList []ResourceStatus

// Collapse returns a single overall status for the list
func (l ResourceStatusList) Collapse() content.NodeStatus {
	var hasWarning bool
	for _, s := range l {
		switch s.Status {
		case content.NodeStatusError:
			return content.NodeStatusError
		case content.NodeStatusWarning:
			hasWarning = true
		}
	}
	if hasWarning {
		return content.NodeStatusWarning
	}
	return content.NodeStatusOK
}

var (
	ingressStatusNoBackendsDefined = ResourceStatus{
		Status: content.NodeStatusWarning,
		Reason: "No backends defined. All traffic will be sent to the default backend configured for the ingress controller.",
	}
	ingressStatusNoMatchingBackend = ResourceStatus{
		Status: content.NodeStatusError,
		Reason: "No matching backends - check service name.",
	}
	ingressStatusNoMatchingPort = ResourceStatus{
		Status: content.NodeStatusError,
		Reason: "No matching backend ports - check service port definitions.",
	}
	ingressStatusMismatchedTLSHost = ResourceStatus{
		Status: content.NodeStatusWarning,
		Reason: "No matching TLS host for rule.",
	}
	ingressStatusNoTLSSecretDefined = ResourceStatus{
		Status: content.NodeStatusError,
		Reason: "TLS secret was not defined.",
	}
	ingressStatusNoMatchingTLSSecret = ResourceStatus{
		Status: content.NodeStatusError,
		Reason: "No matching TLS secret could be found.",
	}
)

// TODO
func statusForPod(pod *core.Pod) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForPodGroup(grp *podGroup) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForService(svc *core.Service) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForIngress(ingress *v1beta1.Ingress, c Cache) (ResourceStatusList, error) {
	if ingress == nil {
		return nil, nil
	}

	result := make([]ResourceStatus, 0)

	backends, err := listIngressBackends(ingress, c)
	if err != nil {
		return nil, err
	}
	// Validation - no backends defined
	if len(backends) == 0 {
		result = append(result, ingressStatusNoBackendsDefined)
	}

	for _, b := range backends {
		if b.ServiceName == "" {
			continue
		}
		svc, err := loadService(b.ServiceName, ingress.Namespace, c)
		if err != nil {
			return nil, err
		}

		// Validation - no matching backend
		if svc == nil {
			result = append(result, ingressStatusNoMatchingBackend)
			continue
		}
		// Validation - no matching port
		if !matchPort(b, svc.Spec.Ports) {
			status := ingressStatusNoMatchingPort
			status.RelatedUID = svc.UID
			result = append(result, status)
			continue
		}
	}

	// Validation - mismatched TLS host
	tlsHosts := tlsHostMap(ingress)
	if len(tlsHosts) > 0 {
		for _, rule := range ingress.Spec.Rules {
			if rule.Host == "" {
				continue
			}
			if ok := tlsHosts[rule.Host]; !ok {
				result = append(result, ingressStatusMismatchedTLSHost)
			}
		}
	}

	// Validation - TLS configuration with no matching secret
	for _, tls := range ingress.Spec.TLS {
		if tls.SecretName == "" {
			result = append(result, ingressStatusNoTLSSecretDefined)
			continue
		}

		key := CacheKey{
			Namespace:  ingress.Namespace,
			APIVersion: "v1",
			Kind:       "Secret",
			Name:       tls.SecretName,
		}
		secrets, err := loadSecrets(key, c)
		if err != nil {
			// Special case - we assume if there was an error it was an access error
			// (the user may not be allowed to see secrets) - and will skip validating TLS.
			break
		}
		if len(secrets) > 0 {
			continue
		}

		result = append(result, ingressStatusNoMatchingTLSSecret)
	}

	return result, nil
}

// matchPort returns true if a matching port is founded for the provided backend
// in the slice of service ports.
func matchPort(b v1beta1.IngressBackend, ports []core.ServicePort) bool {
	for _, p := range ports {
		switch b.ServicePort.Type {
		case intstr.String:
			if i, err := strconv.Atoi(b.ServicePort.StrVal); err == nil {
				if int32(i) == p.Port {
					return true
				}
			}
			if b.ServicePort.StrVal == p.Name {
				return true
			}
		case intstr.Int:
			if int32(b.ServicePort.IntVal) == p.Port {
				return true
			}
		default:
			continue
		}
	}

	return false
}

func statusForReplicaSet(replicaSet *extensions.ReplicaSet) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForDeployment(deployment *extensions.Deployment) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForStatefulSet(s *apps.StatefulSet) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForReplicationController(rc *core.ReplicationController) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForDaemonSet(ds *extensions.DaemonSet) content.NodeStatus {
	return content.NodeStatusOK
}

// tlsHostMap returns a map whose keys are the defined TLS hosts for an ingress.
func tlsHostMap(ingress *v1beta1.Ingress) map[string]bool {
	if ingress == nil {
		return nil
	}

	result := make(map[string]bool)
	for _, tls := range ingress.Spec.TLS {
		for _, host := range tls.Hosts {
			result[host] = true
		}
	}
	return result
}

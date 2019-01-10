package overview

import (
	"github.com/heptio/developer-dash/internal/cache"
	"strconv"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	deploymentReplicasUnavailable = ResourceStatus{
		Status: content.NodeStatusWarning,
		Reason: "One or replicas is unavailable.",
	}
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
	replicaSetAvailableReplicas = ResourceStatus{
		Status: content.NodeStatusWarning,
		Reason: "Replicas count does not match expected.",
	}
)

// TODO
func statusForPod(pod *corev1.Pod) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForPodGroup(grp *podGroup) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForService(svc *corev1.Service) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForIngress(ingress *v1beta1.Ingress, c cache.Cache) (ResourceStatusList, error) {
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

		key := cache.CacheKey{
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
func matchPort(b v1beta1.IngressBackend, ports []corev1.ServicePort) bool {
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

func statusForStatefulSet(s *appsv1.StatefulSet) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForReplicationController(rc *corev1.ReplicationController) content.NodeStatus {
	return content.NodeStatusOK
}

func statusForDaemonSet(ds *appsv1.DaemonSet) content.NodeStatus {
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

// nodeStatusCheck is a function that checks a runtime object and returns
// a list of statuses.
type nodeStatusCheck func(runtime.Object) (ResourceStatusList, error)

// nodeStatus runs zero or more status checks for a runtime object.
type nodeStatus struct {
	checks []nodeStatusCheck
}

// newNodeStatus creates an instance of node status given a list of checks.
func newNodeStatus(checks ...nodeStatusCheck) *nodeStatus {
	return &nodeStatus{
		checks: checks,
	}
}

// check determines the status for an object using the predefined checks.
func (ns *nodeStatus) check(obj runtime.Object) (ResourceStatusList, error) {
	if obj == nil {
		return nil, errors.New("node is nil")
	}

	var list ResourceStatusList

	for _, checkFn := range ns.checks {
		statuses, err := checkFn(obj)
		if err != nil {
			return nil, err
		}

		list = append(list, statuses...)

	}

	return list, nil
}

// deploymentCheckUnavailable returns a warning if the deployment unavailable replicas
// count is greater than 1.
func deploymentCheckUnavailable(obj runtime.Object) (ResourceStatusList, error) {
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, errors.Errorf("expected Deployment; received %T", obj)
	}

	var list ResourceStatusList

	if deployment.Status.UnavailableReplicas != 0 {
		list = append(list, deploymentReplicasUnavailable)
	}

	return list, nil
}

// replicasSetCheckAvailableReplicas returns a warning if the available replicas
// does not match the number of total replicas.
func replicasSetCheckAvailableReplicas(obj runtime.Object) (ResourceStatusList, error) {
	replicaSet, ok := obj.(*appsv1.ReplicaSet)
	if !ok {
		return nil, errors.Errorf("expected ReplicaSet; received %T", obj)
	}

	status := replicaSet.Status

	var list ResourceStatusList

	if status.AvailableReplicas != status.Replicas {
		list = append(list, replicaSetAvailableReplicas)
	}

	return list, nil
}

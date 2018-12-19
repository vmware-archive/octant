package overview

import (
	"testing"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

func Test_matchPort(t *testing.T) {
	tests := []struct {
		name     string
		backend  v1beta1.IngressBackend
		ports    []corev1.ServicePort
		expected bool
	}{
		{
			name: "match name",
			backend: v1beta1.IngressBackend{
				ServicePort: intstr.IntOrString{
					Type:   intstr.String,
					StrVal: "grpc",
				},
			},
			ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name: "nope",
				},
				corev1.ServicePort{
					Name: "grpc",
				},
			},
			expected: true,
		},
		{
			name: "match port (int)",
			backend: v1beta1.IngressBackend{
				ServicePort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 80,
				},
			},
			ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name: "nope",
				},
				corev1.ServicePort{
					Name: "http",
					Port: 80,
				},
			},
			expected: true,
		},
		{
			name: "match port (string)",
			backend: v1beta1.IngressBackend{
				ServicePort: intstr.IntOrString{
					Type:   intstr.String,
					StrVal: "80",
				},
			},
			ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name: "nope",
				},
				corev1.ServicePort{
					Name: "http",
					Port: 80,
				},
			},
			expected: true,
		},
		{
			name: "no match",
			backend: v1beta1.IngressBackend{
				ServicePort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 80,
				},
			},
			ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name: "nope",
				},
				corev1.ServicePort{
					Name: "https",
					Port: 443,
				},
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := matchPort(tc.backend, tc.ports)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_statusForIngress(t *testing.T) {
	tests := []struct {
		name     string
		objects  []string
		expected ResourceStatusList
	}{
		{
			name: "Single service ingress",
			objects: []string{
				`---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
    name: single-service-ingress
    annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  backend:
    serviceName: single-service
    servicePort: 80
`,
				`---
apiVersion: v1
kind: Service
metadata:
  name: single-service
spec:
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9376
`,
			},
			expected: ResourceStatusList{},
		},
		{
			name: "No matching backends",
			objects: []string{
				`---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
    name: ingress-no-service-found
    annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
    rules:
    - http:
        paths:
        - path: /testpath
          backend:
            serviceName: no-such-service
            servicePort: 80
`,
			},
			expected: ResourceStatusList{
				ingressStatusNoMatchingBackend,
			},
		},
		{
			name: "No defined backends",
			objects: []string{
				`---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
    name: ingress-no-service-found
    annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
    rules:
    - http:
        paths:
        - path: /testpath
`,
			},
			expected: ResourceStatusList{
				ingressStatusNoBackendsDefined,
			},
		},
		{
			name: "No matching port",
			objects: []string{
				`---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
    name: ingress-bad-port
    annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
    rules:
    - http:
        paths:
        - path: /testpath
          backend:
            serviceName: service-wrong-port
            servicePort: 80
`,
				`---
apiVersion: v1
kind: Service
metadata:
  name: service-wrong-port
spec:
  ports:
    - protocol: TCP
      port: 8888
      targetPort: 9376
`,
			},
			expected: ResourceStatusList{
				ingressStatusNoMatchingPort,
			},
		},
		{
			name: "Mismatched TLS host",
			objects: []string{
				`---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
    name: ingress-bad-tls-host
    annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
    tls:
    - hosts:
      - sslexample.foo.com
      secretName: testsecret-tls
    rules:
    - host: not-the-tls-host.com
      http:
        paths:
        - path: /testpath
          backend:
            serviceName: my-service
            servicePort: grpc
`,
				`---
apiVersion: v1
kind: Secret
metadata:
  name: testsecret-tls
type: Opaque
data:
  tls.crt: Zm9vCg==
  tls.key: YmFyCg==
`,
				`---
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  ports:
    - protocol: TCP
      name: grpc
      port: 8888
      targetPort: 9376
`,
			},
			expected: ResourceStatusList{
				ingressStatusMismatchedTLSHost,
			},
		},
		{
			name: "Missing TLS secret",
			objects: []string{
				`---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
    name: ingress-bad-tls-host
    annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
    tls:
    - hosts:
      - sslexample.foo.com
      secretName: no-such-secret
    rules:
    - host: sslexample.foo.com
      http:
        paths:
        - path: /testpath
          backend:
            serviceName: my-service
            servicePort: grpc
`,
				`---
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  ports:
    - protocol: TCP
      name: grpc
      port: 8888
      targetPort: 9376
`,
			},
			expected: ResourceStatusList{
				ingressStatusNoMatchingTLSSecret,
			},
		},
		{
			name: "No TLS secret defined",
			objects: []string{
				`---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
    name: ingress-bad-tls-host
    annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
    tls:
    - hosts:
      - sslexample.foo.com
    rules:
    - host: sslexample.foo.com
      http:
        paths:
        - path: /testpath
          backend:
            serviceName: my-service
            servicePort: grpc
`,
				`---
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  ports:
    - protocol: TCP
      name: grpc
      port: 8888
      targetPort: 9376
`,
			},
			expected: ResourceStatusList{
				ingressStatusNoTLSSecretDefined,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			objects := make([]runtime.Object, 0)
			uObjs := make([]runtime.Object, 0)

			// Parse objects
			for _, raw := range tc.objects {
				obj, _, err := scheme.Codecs.UniversalDeserializer().Decode([]byte(raw), nil, nil)
				require.NoError(t, err, "decoding object: %v", raw)
				if err != nil {
					return
				}
				acc := meta.NewAccessor()
				ns, err := acc.Namespace(obj)
				require.NoError(t, err)
				if ns == "" {
					err = acc.SetNamespace(obj, "default")
					require.NoError(t, err)
				}

				objects = append(objects, obj)

				m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
				require.NoError(t, err)
				uObjs = append(uObjs, &unstructured.Unstructured{
					Object: m,
				})
			}

			c, cancel, err := newCache(t, uObjs)
			require.NoError(t, err)
			defer cancel()

			switch v := objects[0].(type) {
			case *v1beta1.Ingress:
				actual, err := statusForIngress(v, c)
				require.NoError(t, err)
				if err != nil {
					return
				}
				assert.Equal(t, tc.expected, actual)

			default:
				t.Errorf("Invalid type: %T", objects[0])
				return
			}
		})
	}
}

func Test_nodeStatus_check(t *testing.T) {
	errorList := ResourceStatusList{
		{
			Status: content.NodeStatusError,
			Reason: "failed",
		},
	}
	okList := ResourceStatusList{
		{
			Status: content.NodeStatusOK,
			Reason: "ok",
		},
	}

	cases := []struct {
		name     string
		checks   []nodeStatusCheck
		object   runtime.Object
		isErr    bool
		expected ResourceStatusList
	}{
		{
			name:   "nil object",
			object: nil,
			isErr:  true,
		},
		{
			name: "check failed",
			checks: []nodeStatusCheck{
				func(runtime.Object) (ResourceStatusList, error) {
					return nil, errors.New("failed")
				},
			},
			object: &appsv1.Deployment{},
			isErr:  true,
		},
		{
			name: "check has non ok status",
			checks: []nodeStatusCheck{
				func(runtime.Object) (ResourceStatusList, error) {
					return errorList, nil
				},
			},
			object:   &appsv1.Deployment{},
			expected: errorList,
		},
		{
			name: "all checks are ok",
			checks: []nodeStatusCheck{
				func(runtime.Object) (ResourceStatusList, error) {
					return okList, nil
				},
			},
			object:   &appsv1.Deployment{},
			expected: okList,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ns := newNodeStatus(tc.checks...)
			status, err := ns.check(tc.object)
			if tc.isErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, status)
		})
	}
}

func Test_deploymentCheckUnavailable(t *testing.T) {

	cases := []struct {
		name     string
		obj      runtime.Object
		isErr    bool
		expected ResourceStatusList
	}{
		{
			name:     "in general",
			obj:      &appsv1.Deployment{},
			expected: nil,
		},
		{
			name:  "not a deployment",
			obj:   nil,
			isErr: true,
		},
		{
			name:     "has unavailable replicas",
			obj:      &appsv1.Deployment{Status: appsv1.DeploymentStatus{UnavailableReplicas: 1}},
			expected: ResourceStatusList{deploymentReplicasUnavailable},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, err := deploymentCheckUnavailable(tc.obj)
			if tc.isErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, status)
		})
	}
}

func Test_replicaSetCheckAvailable(t *testing.T) {
	cases := []struct {
		name     string
		obj      runtime.Object
		isErr    bool
		expected ResourceStatusList
	}{
		{
			name:     "in general",
			obj:      &appsv1.ReplicaSet{},
			expected: nil,
		},
		{
			name:  "not a replica set",
			obj:   nil,
			isErr: true,
		},
		{
			name: "has available replicas",
			obj: &appsv1.ReplicaSet{Status: appsv1.ReplicaSetStatus{
				Replicas: 1, AvailableReplicas: 0}},
			expected: ResourceStatusList{replicaSetAvailableReplicas},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, err := replicasSetCheckAvailableReplicas(tc.obj)
			if tc.isErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, status)
		})
	}
}

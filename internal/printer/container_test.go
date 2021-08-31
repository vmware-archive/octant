/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/portforward"
	pffake "github.com/vmware-tanzu/octant/internal/portforward/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_ContainerConfiguration(t *testing.T) {
	var (
		propagation    = corev1.MountPropagationHostToContainer
		validContainer = &corev1.Container{
			Name:  "nginx",
			Image: "nginx:1.15",
			Ports: []corev1.ContainerPort{
				{
					Name:     "http",
					HostPort: 80,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     "metrics",
					HostPort: 8080,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:          "application",
					ContainerPort: 8443,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          "dtls",
					ContainerPort: 443,
					Protocol:      corev1.ProtocolUDP,
				},
			},
			Command: []string{"/usr/bin/nginx"},
			Args:    []string{"-v", "-p", "80"},

			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "config",
					ReadOnly:  true,
					MountPath: "/etc/nginx",
				},
				{
					Name:             "data",
					MountPath:        "/var/www",
					SubPath:          "/content",
					MountPropagation: &propagation,
				},
			},
			Env: []corev1.EnvVar{
				{
					Name:  "tier",
					Value: "prod",
				},
				{
					Name: "fieldref",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.name"},
					},
				},
				{
					Name: "resourcefieldref",
					ValueFrom: &corev1.EnvVarSource{
						ResourceFieldRef: &corev1.ResourceFieldSelector{
							Resource: "requests.cpu",
						},
					},
				},
				{
					Name: "configmapref",
					ValueFrom: &corev1.EnvVarSource{
						ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{Name: "myconfig"},
							Key:                  "somekey",
						},
					},
				},
				{
					Name: "secretref",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{Name: "mysecret"},
							Key:                  "somesecretkey",
						},
					},
				},
			},
			EnvFrom: []corev1.EnvFromSource{
				{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{Name: "fromconfig"},
					},
				},
				{
					SecretRef: &corev1.SecretEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{Name: "fromsecret"},
					},
				},
			},
		}
		validInitContainer = &corev1.Container{
			Name:    "busybox",
			Image:   "busybox:1.28",
			Command: []string{"sh"},
			Args:    []string{"-c", "until nslookup mydb; do echo waiting for mydb; sleep 2; done;"},
		}
	)

	now := time.Now()

	envTable := component.NewTable("Environment", "There are no defined environment variables!",
		component.NewTableCols("Name", "Value", "Source"))
	envTable.Add(
		component.TableRow{
			"Name":   component.NewText("configmapref"),
			"Value":  component.NewText(""),
			"Source": component.NewLink("", "myconfig:somekey", "/configMap"),
		},
		component.TableRow{
			"Name":   component.NewText("fieldref"),
			"Value":  component.NewText(""),
			"Source": component.NewText("metadata.name"),
		},
		component.TableRow{
			"Name":   component.NewText("resourcefieldref"),
			"Value":  component.NewText(""),
			"Source": component.NewText("requests.cpu"),
		},
		component.TableRow{
			"Name":   component.NewText("secretref"),
			"Value":  component.NewText(""),
			"Source": component.NewLink("", "mysecret:somesecretkey", "/secret"),
		},
		component.TableRow{
			"Name":   component.NewText("tier"),
			"Value":  component.NewText("prod"),
			"Source": component.NewText(""),
		},
		// EnvFromSource
		component.TableRow{
			"Source": component.NewLink("", "fromconfig", "/fromConfig"),
		},
		component.TableRow{
			"Source": component.NewLink("", "fromsecret", "/fromSecret"),
		},
	)

	volTable := component.NewTable("Volume Mounts", "There are no volume mounts!",
		component.NewTableCols("Name", "Mount Path", "Propagation"))
	volTable.Add(
		component.TableRow{
			"Name":        component.NewText("config"),
			"Mount Path":  component.NewText("/etc/nginx (ro)"),
			"Propagation": component.NewText(""),
		},
		component.TableRow{
			"Name":        component.NewText("data"),
			"Mount Path":  component.NewText("/var/www/content (rw)"),
			"Propagation": component.NewText("HostToContainer"),
		},
	)

	cases := []struct {
		name      string
		container *corev1.Container
		isErr     bool
		isInit    bool
		action    component.Action
		expected  *component.Summary
	}{
		{
			name:      "in general",
			container: validContainer,
			expected: component.NewSummary("Container nginx", []component.SummarySection{
				{
					Header:  "Image",
					Content: component.NewText("nginx:1.15"),
				},
				{
					Header:  "Image ID",
					Content: component.NewText("nginx-image-id"),
				},
				{
					Header:  "Image Manifest",
					Content: component.NewJSONEditor("{\"manifests\":[{\"digest\":\"sha256:e770165fef9e36b990882a4083d8ccf5e29e469a8609bb6b2e3b47d9510e2c8d\",\"mediaType\":\"application\\/vnd.docker.distribution.manifest.v2+json\",\"platform\":{\"architecture\":\"amd64\",\"os\":\"linux\"},\"size\":948},{\"digest\":\"sha256:26687467368eba1745b3af5f673156e5598b0d3609ddc041d4afb3000a7c97c4\",\"mediaType\":\"application\\/vnd.docker.distribution.manifest.v2+json\",\"platform\":{\"architecture\":\"arm\",\"os\":\"linux\",\"variant\":\"v7\"},\"size\":948},{\"digest\":\"sha256:322d209ca0e9dcd69cf1bb9354cb2c573255e96689f31b0964753389b780269c\",\"mediaType\":\"application\\/vnd.docker.distribution.manifest.v2+json\",\"platform\":{\"architecture\":\"arm64\",\"os\":\"linux\",\"variant\":\"v8\"},\"size\":948},{\"digest\":\"sha256:2393dbb3ac0f27a4b097908f78510aa20dce07c029540762447ab4731119bab7\",\"mediaType\":\"application\\/vnd.docker.distribution.manifest.v2+json\",\"platform\":{\"architecture\":\"386\",\"os\":\"linux\"},\"size\":948},{\"digest\":\"sha256:16f53d8a8fcef518bfc7ad0b87f572c036eedc5307a2539e4c73741a7fe8ea76\",\"mediaType\":\"application\\/vnd.docker.distribution.manifest.v2+json\",\"platform\":{\"architecture\":\"ppc64le\",\"os\":\"linux\"},\"size\":948},{\"digest\":\"sha256:a89d88340baf686e95076902c5f89bd54755cbb324eaae5a2a470f98db342f55\",\"mediaType\":\"application\\/vnd.docker.distribution.manifest.v2+json\",\"platform\":{\"architecture\":\"s390x\",\"os\":\"linux\"},\"size\":948}],\"mediaType\":\"application\\/vnd.docker.distribution.manifest.list.v2+json\",\"schemaVersion\":2}", true),
				},
				{
					Header:  "Image Configuration",
					Content: component.NewJSONEditor("{\n  \"created\": \"2019-05-08T03:01:41.947151778Z\",\n  \"architecture\": \"amd64\",\n  \"os\": \"linux\",\n  \"config\": {\n    \"ExposedPorts\": {\n      \"80/tcp\": {}\n    },\n    \"Env\": [\n      \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\n      \"NGINX_VERSION=1.15.12-1~stretch\",\n      \"NJS_VERSION=1.15.12.0.3.1-1~stretch\"\n    ],\n    \"Cmd\": [\n      \"nginx\",\n      \"-g\",\n      \"daemon off;\"\n    ],\n    \"Labels\": {\n      \"maintainer\": \"NGINX Docker Maintainers <docker-maint@nginx.com>\"\n    },\n    \"StopSignal\": \"SIGTERM\"\n  },\n  \"rootfs\": {\n    \"type\": \"layers\",\n    \"diff_ids\": [\n      \"sha256:6270adb5794c6987109e54af00ab456977c5d5cc6f1bc52c1ce58d32ec0f15f4\",\n      \"sha256:6ba094226eea86e21761829b88bdfdc9feb14bd83d60fb7e666f0943253657e8\",\n      \"sha256:332fa54c58864e2dcd3df0ad88c69b2707d45f2d8121dad6278a15148900e490\"\n    ]\n  },\n  \"history\": [\n    {\n      \"created\": \"2019-05-08T00:33:32.152758355Z\",\n      \"created_by\": \"/bin/sh -c #(nop) ADD file:fcb9328ea4c1156709f3d04c3d9a5f3667e77fb36a4a83390ae2495555fc0238 in / \"\n    },\n    {\n      \"created\": \"2019-05-08T00:33:32.718284983Z\",\n      \"created_by\": \"/bin/sh -c #(nop)  CMD [\\\"bash\\\"]\",\n      \"empty_layer\": true\n    },\n    {\n      \"created\": \"2019-05-08T03:01:16.010671568Z\",\n      \"created_by\": \"/bin/sh -c #(nop)  LABEL maintainer=NGINX Docker Maintainers <docker-maint@nginx.com>\",\n      \"empty_layer\": true\n    },\n    {\n      \"created\": \"2019-05-08T03:01:16.175452264Z\",\n      \"created_by\": \"/bin/sh -c #(nop)  ENV NGINX_VERSION=1.15.12-1~stretch\",\n      \"empty_layer\": true\n    },\n    {\n      \"created\": \"2019-05-08T03:01:16.36342084Z\",\n      \"created_by\": \"/bin/sh -c #(nop)  ENV NJS_VERSION=1.15.12.0.3.1-1~stretch\",\n      \"empty_layer\": true\n    },\n    {\n      \"created\": \"2019-05-08T03:01:40.497446007Z\",\n      \"created_by\": \"/bin/sh -c set -x \\t&& apt-get update \\t&& apt-get install --no-install-recommends --no-install-suggests -y gnupg1 apt-transport-https ca-certificates \\t&& \\tNGINX_GPGKEY=573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62; \\tfound=''; \\tfor server in \\t\\tha.pool.sks-keyservers.net \\t\\thkp://keyserver.ubuntu.com:80 \\t\\thkp://p80.pool.sks-keyservers.net:80 \\t\\tpgp.mit.edu \\t; do \\t\\techo \\\"Fetching GPG key $NGINX_GPGKEY from $server\\\"; \\t\\tapt-key adv --keyserver \\\"$server\\\" --keyserver-options timeout=10 --recv-keys \\\"$NGINX_GPGKEY\\\" && found=yes && break; \\tdone; \\ttest -z \\\"$found\\\" && echo >&2 \\\"error: failed to fetch GPG key $NGINX_GPGKEY\\\" && exit 1; \\tapt-get remove --purge --auto-remove -y gnupg1 && rm -rf /var/lib/apt/lists/* \\t&& dpkgArch=\\\"$(dpkg --print-architecture)\\\" \\t&& nginxPackages=\\\" \\t\\tnginx=${NGINX_VERSION} \\t\\tnginx-module-xslt=${NGINX_VERSION} \\t\\tnginx-module-geoip=${NGINX_VERSION} \\t\\tnginx-module-image-filter=${NGINX_VERSION} \\t\\tnginx-module-njs=${NJS_VERSION} \\t\\\" \\t&& case \\\"$dpkgArch\\\" in \\t\\tamd64|i386) \\t\\t\\techo \\\"deb https://nginx.org/packages/mainline/debian/ stretch nginx\\\" >> /etc/apt/sources.list.d/nginx.list \\t\\t\\t&& apt-get update \\t\\t\\t;; \\t\\t*) \\t\\t\\techo \\\"deb-src https://nginx.org/packages/mainline/debian/ stretch nginx\\\" >> /etc/apt/sources.list.d/nginx.list \\t\\t\\t\\t\\t\\t&& tempDir=\\\"$(mktemp -d)\\\" \\t\\t\\t&& chmod 777 \\\"$tempDir\\\" \\t\\t\\t\\t\\t\\t&& savedAptMark=\\\"$(apt-mark showmanual)\\\" \\t\\t\\t\\t\\t\\t&& apt-get update \\t\\t\\t&& apt-get build-dep -y $nginxPackages \\t\\t\\t&& ( \\t\\t\\t\\tcd \\\"$tempDir\\\" \\t\\t\\t\\t&& DEB_BUILD_OPTIONS=\\\"nocheck parallel=$(nproc)\\\" \\t\\t\\t\\t\\tapt-get source --compile $nginxPackages \\t\\t\\t) \\t\\t\\t\\t\\t\\t&& apt-mark showmanual | xargs apt-mark auto > /dev/null \\t\\t\\t&& { [ -z \\\"$savedAptMark\\\" ] || apt-mark manual $savedAptMark; } \\t\\t\\t\\t\\t\\t&& ls -lAFh \\\"$tempDir\\\" \\t\\t\\t&& ( cd \\\"$tempDir\\\" && dpkg-scanpackages . > Packages ) \\t\\t\\t&& grep '^Package: ' \\\"$tempDir/Packages\\\" \\t\\t\\t&& echo \\\"deb [ trusted=yes ] file://$tempDir ./\\\" > /etc/apt/sources.list.d/temp.list \\t\\t\\t&& apt-get -o Acquire::GzipIndexes=false update \\t\\t\\t;; \\tesac \\t\\t&& apt-get install --no-install-recommends --no-install-suggests -y \\t\\t\\t\\t\\t\\t$nginxPackages \\t\\t\\t\\t\\t\\tgettext-base \\t&& apt-get remove --purge --auto-remove -y apt-transport-https ca-certificates && rm -rf /var/lib/apt/lists/* /etc/apt/sources.list.d/nginx.list \\t\\t&& if [ -n \\\"$tempDir\\\" ]; then \\t\\tapt-get purge -y --auto-remove \\t\\t&& rm -rf \\\"$tempDir\\\" /etc/apt/sources.list.d/temp.list; \\tfi\"\n    },\n    {\n      \"created\": \"2019-05-08T03:01:41.355881721Z\",\n      \"created_by\": \"/bin/sh -c ln -sf /dev/stdout /var/log/nginx/access.log \\t&& ln -sf /dev/stderr /var/log/nginx/error.log\"\n    },\n    {\n      \"created\": \"2019-05-08T03:01:41.538214273Z\",\n      \"created_by\": \"/bin/sh -c #(nop)  EXPOSE 80\",\n      \"empty_layer\": true\n    },\n    {\n      \"created\": \"2019-05-08T03:01:41.740886057Z\",\n      \"created_by\": \"/bin/sh -c #(nop)  STOPSIGNAL SIGTERM\",\n      \"empty_layer\": true\n    },\n    {\n      \"created\": \"2019-05-08T03:01:41.947151778Z\",\n      \"created_by\": \"/bin/sh -c #(nop)  CMD [\\\"nginx\\\" \\\"-g\\\" \\\"daemon off;\\\"]\",\n      \"empty_layer\": true\n    }\n  ]\n}", true),
				},
				{
					Header:  "Host Ports",
					Content: component.NewText("80/TCP, 8080/TCP"),
				},
				{
					Header: "Container Ports",
					Content: component.NewPorts([]component.Port{
						*component.NewPort("namespace", "v1", "Pod", "pod", 8443, "TCP", component.PortForwardState{IsForwardable: true, IsForwarded: true}),
						*component.NewPort("namespace", "v1", "Pod", "pod", 443, "UDP", component.PortForwardState{IsForwardable: false, IsForwarded: false}),
					}),
				},
				{
					Header:  "Last State",
					Content: component.NewText(fmt.Sprintf("terminated with 255 at %s: reason", now)),
				},
				{
					Header:  "Current State",
					Content: component.NewText(fmt.Sprintf("started at %s", now)),
				},
				{
					Header:  "Ready",
					Content: component.NewText("true"),
				},
				{
					Header:  "Restart Count",
					Content: component.NewText("2"),
				},
				{
					Header:  "Environment",
					Content: envTable,
				},
				{
					Header:  "Command",
					Content: component.NewText("['/usr/bin/nginx']"),
				},
				{
					Header:  "Args",
					Content: component.NewText("['-v', '-p', '80']"),
				},
				{
					Header:  "Volume Mounts",
					Content: volTable,
				},
			}...),
		},
		{
			name:      "init containers",
			container: validInitContainer,
			isInit:    true,
			expected: component.NewSummary("Init Container busybox", []component.SummarySection{
				{
					Header:  "Image",
					Content: component.NewText("busybox:1.28"),
				},
				{
					Header:  "Image ID",
					Content: component.NewText("busybox-image-id"),
				},
				{
					Header:  "Image Manifest",
					Content: component.NewJSONEditor("{\n   \"schemaVersion\": 2,\n   \"mediaType\": \"application/vnd.docker.distribution.manifest.list.v2+json\",\n   \"manifests\": [\n      {\n         \"mediaType\": \"application/vnd.docker.distribution.manifest.v2+json\",\n         \"size\": 527,\n         \"digest\": \"sha256:74f634b1bc1bd74535d5209589734efbd44a25f4e2dc96d78784576a3eb5b335\",\n         \"platform\": {\n            \"architecture\": \"amd64\",\n            \"os\": \"linux\"\n         }\n      },\n      {\n         \"mediaType\": \"application/vnd.docker.distribution.manifest.v2+json\",\n         \"size\": 527,\n         \"digest\": \"sha256:35e28b647bd4976b7cacfaa32b7b253817d0881d77b6cda731ad46a29d08c2cb\",\n         \"platform\": {\n            \"architecture\": \"arm\",\n            \"os\": \"linux\",\n            \"variant\": \"v5\"\n         }\n      },\n      {\n         \"mediaType\": \"application/vnd.docker.distribution.manifest.v2+json\",\n         \"size\": 527,\n         \"digest\": \"sha256:420befcb0c197618f0252108d553d8a112e291e2a6a75d8a2b4933f511480ea3\",\n         \"platform\": {\n            \"architecture\": \"arm\",\n            \"os\": \"linux\",\n            \"variant\": \"v6\"\n         }\n      },\n      {\n         \"mediaType\": \"application/vnd.docker.distribution.manifest.v2+json\",\n         \"size\": 527,\n         \"digest\": \"sha256:4df1e7dbe58b7fe24145291700e4fdf89a80677ffeb9b972840b42e3ec065e1f\",\n         \"platform\": {\n            \"architecture\": \"arm\",\n            \"os\": \"linux\",\n            \"variant\": \"v7\"\n         }\n      },\n      {\n         \"mediaType\": \"application/vnd.docker.distribution.manifest.v2+json\",\n         \"size\": 527,\n         \"digest\": \"sha256:859d41e4316c182cb559f9ae3c5ffcac8602ee1179794a1707c06cd092a008d3\",\n         \"platform\": {\n            \"architecture\": \"arm64\",\n            \"os\": \"linux\",\n            \"variant\": \"v8\"\n         }\n      },\n      {\n         \"mediaType\": \"application/vnd.docker.distribution.manifest.v2+json\",\n         \"size\": 527,\n         \"digest\": \"sha256:19f468f7dde9dc85d1576e6eb244b190661764199e21fcb53d84378bef16e334\",\n         \"platform\": {\n            \"architecture\": \"386\",\n            \"os\": \"linux\"\n         }\n      },\n      {\n         \"mediaType\": \"application/vnd.docker.distribution.manifest.v2+json\",\n         \"size\": 528,\n         \"digest\": \"sha256:2d8967e4a68583a4bb2d7e236c60a1d72a585439b41e7a77555edad8df0f2bf4\",\n         \"platform\": {\n            \"architecture\": \"ppc64le\",\n            \"os\": \"linux\"\n         }\n      },\n      {\n         \"mediaType\": \"application/vnd.docker.distribution.manifest.v2+json\",\n         \"size\": 528,\n         \"digest\": \"sha256:67510360fd7c837d71ecfbd9f7991d72a2d2cbda3b383115a0dda0f0936b57f6\",\n         \"platform\": {\n            \"architecture\": \"s390x\",\n            \"os\": \"linux\"\n         }\n      }\n   ]\n}", true),
				},
				{
					Header:  "Image Configuration",
					Content: component.NewJSONEditor("{\n  \"created\": \"2018-05-23T21:19:31.132152818Z\",\n  \"architecture\": \"amd64\",\n  \"os\": \"linux\",\n  \"config\": {\n    \"Env\": [\n      \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"\n    ],\n    \"Cmd\": [\n      \"sh\"\n    ]\n  },\n  \"rootfs\": {\n    \"type\": \"layers\",\n    \"diff_ids\": [\n      \"sha256:432b65032b9466b4dadcc5c7b11701e71d21c18400aae946b101ad16be62333a\"\n    ]\n  },\n  \"history\": [\n    {\n      \"created\": \"2018-05-23T21:19:30.902651601Z\",\n      \"created_by\": \"/bin/sh -c #(nop) ADD file:5f0439d8328ab58c087cd067c91ce92765da98916d91b083df6590477b7b9f19 in / \"\n    },\n    {\n      \"created\": \"2018-05-23T21:19:31.132152818Z\",\n      \"created_by\": \"/bin/sh -c #(nop)  CMD [\\\"sh\\\"]\",\n      \"empty_layer\": true\n    }\n  ]\n}", true),
				},
				{
					Header:  "Command",
					Content: component.NewText("['sh']"),
				},
				{
					Header:  "Args",
					Content: component.NewText("['-c', 'until nslookup mydb; do echo waiting for mydb; sleep 2; done;']"),
				},
			}...),
		},
		{
			name:      "container is nil",
			container: nil,
			isErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			ctx := context.Background()
			configMap := testutil.CreateConfigMap("myconfig")

			pf := pffake.NewMockPortForwarder(controller)
			gvk := schema.GroupVersionKind{Version: "v1", Kind: "Pod"}

			states := []portforward.State{
				{
					CreatedAt: testutil.Time(),
					Ports: []portforward.ForwardedPort{
						{
							Local:  uint16(45275),
							Remote: uint16(8443),
						},
					},
					Pod: portforward.Target{
						GVK:       gvk,
						Namespace: "namespace",
						Name:      "pod",
					},
				},
			}

			state := createPortForwardState("stateid", "namespace", "pod", gvk)

			pf.EXPECT().FindPod("namespace", gomock.Eq(gvk), "pod").Return(states, nil).AnyTimes()
			pf.EXPECT().FindTarget("namespace", gomock.Eq(gvk), "pod").Return(states, nil).AnyTimes()
			pf.EXPECT().Get(gomock.Any()).Return(state, true).AnyTimes()

			tpo.PathForGVK("namespace", "v1", "Secret", "mysecret", "mysecret:somesecretkey", "/secret")
			tpo.PathForGVK("namespace", "v1", "ConfigMap", "myconfig", "myconfig:somekey", "/configMap")
			tpo.PathForGVK("namespace", "v1", "Secret", "fromsecret", "fromsecret", "/fromSecret")
			tpo.PathForGVK("namespace", "v1", "ConfigMap", "fromconfig", "fromconfig", "/fromConfig")

			parentPod := testutil.CreatePod("pod")
			parentPod.Namespace = "namespace"
			parentPod.Status = corev1.PodStatus{
				InitContainerStatuses: []corev1.ContainerStatus{
					{
						Name:    "busybox",
						ImageID: "busybox-image-id",
					},
				},
				ContainerStatuses: []corev1.ContainerStatus{
					{
						Name:         "nginx",
						ImageID:      "nginx-image-id",
						Ready:        true,
						RestartCount: 2,
						State: corev1.ContainerState{
							Running: &corev1.ContainerStateRunning{
								StartedAt: metav1.Time{Time: now},
							},
						},
						LastTerminationState: corev1.ContainerState{
							Terminated: &corev1.ContainerStateTerminated{
								FinishedAt: metav1.Time{Time: now},
								Reason:     "reason",
								ExitCode:   255,
							},
						},
					},
				},
			}

			if tc.container != nil {
				key := store.Key{
					Kind:       configMap.Kind,
					APIVersion: configMap.APIVersion,
					Name:       configMap.Name,
					Namespace:  configMap.Namespace,
				}
				tpo.objectStore.EXPECT().Get(ctx, gomock.Eq(key)).Return(testutil.ToUnstructured(t, configMap), nil).AnyTimes()
			}

			cc := NewContainerConfiguration(ctx, parentPod, tc.container, pf, IsInit(tc.isInit), WithPrintOptions(printOptions))
			summary, err := cc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_containerNotFoundError(t *testing.T) {
	e := containerNotFoundError{name: "name"}

	expected := fmt.Sprintf("container %q not found", "name")
	assert.Equal(t, expected, e.Error())

	assert.False(t, e.isContainerFound())
}

func Test_findContainerStatus(t *testing.T) {
	tests := []struct {
		name            string
		podFactory      func(statusName string) *corev1.Pod
		statusName      string
		isInit          bool
		isErr           bool
		expectedErrType reflect.Type
		expectedStatus  *corev1.ContainerStatus
	}{
		{
			name:       "container with status",
			statusName: "name",
			podFactory: func(statusName string) *corev1.Pod {
				pod := testutil.CreatePod("pod")

				pod.Status.ContainerStatuses = append(
					pod.Status.ContainerStatuses,
					corev1.ContainerStatus{Name: statusName})

				return pod
			},
			expectedStatus: &corev1.ContainerStatus{Name: "name"},
		},
		{
			name:       "init container with status",
			isInit:     true,
			statusName: "name",
			podFactory: func(statusName string) *corev1.Pod {
				pod := testutil.CreatePod("pod")

				pod.Status.InitContainerStatuses = append(
					pod.Status.ContainerStatuses,
					corev1.ContainerStatus{Name: statusName})

				return pod
			},
			expectedStatus: &corev1.ContainerStatus{Name: "name"},
		},
		{
			name: "no containers",
			podFactory: func(statusName string) *corev1.Pod {
				pod := testutil.CreatePod("pod")
				return pod
			},
			isErr:           true,
			expectedErrType: reflect.TypeOf(&containerNotFoundError{}),
		},
		{
			name: "pod is nil",
			podFactory: func(statusName string) *corev1.Pod {
				return nil
			},
			isErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pod := test.podFactory(test.statusName)

			status, err := findContainerStatus(pod, test.statusName, test.isInit)
			if test.isErr {
				require.Error(t, err)

				if test.expectedErrType != nil {
					errType := reflect.TypeOf(err)
					require.Equal(t, test.expectedErrType, errType)
				}

				return
			}
			require.NoError(t, err)

			require.Equal(t, test.expectedStatus, status)
		})
	}
}

func Test_editContainerAction(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment", testutil.WithGenericDeployment())
	container := deployment.Spec.Template.Spec.Containers[0]

	got, err := editContainerAction(deployment, &container)
	require.NoError(t, err)

	form, err := component.CreateFormForObject(octant.ActionOverviewContainerEditor, deployment,
		component.NewFormFieldText("Image", "containerImage", container.Image),
		component.NewFormFieldHidden("containersPath", `["spec","template","spec","containers"]`),
		component.NewFormFieldHidden("containerName", container.Name),
	)
	require.NoError(t, err)

	expected := component.Action{
		Name:  "Edit",
		Title: "Container container-name Editor",
		Form:  form,
	}
	require.Equal(t, expected, got)
}

func createPortForwardState(id, namespace, targetName string, gvk schema.GroupVersionKind) portforward.State {
	return portforward.State{
		ID:        id,
		CreatedAt: testutil.Time(),
		Pod: portforward.Target{
			GVK:       gvk,
			Namespace: namespace,
			Name:      targetName,
		},
		Target: portforward.Target{
			GVK:       gvk,
			Namespace: namespace,
			Name:      targetName,
		}}
}

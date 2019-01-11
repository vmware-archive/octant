package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view"

	"github.com/heptio/developer-dash/internal/content"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type resourceViewer struct{}

func (rv *resourceViewer) IsEmpty() bool {
	return false
}

func (rv *resourceViewer) MarshalJSON() ([]byte, error) {
	return rv1, nil
}

func (rv *resourceViewer) ViewComponent() content.ViewComponent {
	return content.ViewComponent{}
}

type ResourceViewerStub struct{}

func NewResourceViewerStub(prefix, namespace string, c clock.Clock) view.View {
	return &ResourceViewerStub{}
}

func (rss *ResourceViewerStub) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	return []content.Content{&resourceViewer{}}, nil
}

func (rss *ResourceViewerStub) ViewComponent() content.ViewComponent {
	return content.ViewComponent{}
}

var (
	rv0 = []byte(`{
	"type": "resourceviewer",
	"selected": "pod",
	"adjacencyList": {
		"ingress-uid": [
		{
			"node": "service-uid",
			"edge": "explicit"
		}
		],
		"service-uid": [
		{
			"node": "pod",
			"edge": "implicit"
		},
		{
			"node": "pod2",
			"edge": "implicit"
		}
		],
		"deployment-uid": [
		{
			"node": "replica-set-uid",
			"edge": "explicit"
		},
		{
			"node": "rs2",
			"edge": "explicit"
		}
		],
		"replica-set-uid": [
		{
			"node": "pod",
			"edge": "explicit"
		}
		],
		"rs2": [
		{
			"node": "pod",
			"edge": "explicit"
		}
		]
	},
	"objects": {
		"ingress-uid": {
		"name": "ingress",
		"apiVersion": "extensions/v1beta1",
		"kind": "Ingress",
		"status": "ok",
		"isNetwork": true,
		"views": []
		},
		"service-uid": {
		"name": "service",
		"apiVersion": "v1",
		"kind": "Service",
		"status": "ok",
		"isNetwork": true,
		"views": []
		},
		"deployment-uid": {
		"name": "deployment",
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"status": "ok",
		"views": []
		},
		"replica-set-uid": {
		"name": "deployment-abc-12345",
		"apiVersion": "apps/v1",
		"kind": "ReplicaSet",
		"status": "ok",
		"views": []
		},
		"rs2": {
		"name": "deployment-abc-23456",
		"apiVersion": "apps/v1",
		"kind": "ReplicaSet",
		"status": "ok",
		"views": []
		},
		"pod": {
		"name": "deployment-abc pods",
		"apiVersion": "v1",
		"kind": "Pod",
		"status": "ok",
		"views": []
		},
		"pod2": {
		"name": "other",
		"apiVersion": "v1",
		"kind": "Pod",
		"status": "warning",
		"views": []
		}
	}
	}`)

	rv1 = []byte(` {
	"type": "resourceviewer",
	"selected": "d1",
	"adjacencyList": {
		"rs0": [
			{
				"node": "pods-rs0",
				"edge": "explicit"
			}
		],
		"rs1": [
			{
				"node": "pods-rs0",
				"edge": "explicit"
			}
		],
		"rs2": [
			{
				"node": "pods-rs0",
				"edge": "explicit"
			}
		],
		"rs3": [
			{
				"node": "pods-rs0",
				"edge": "explicit"
			}
		],
		"s1": [
			{
				"node": "pods-rs0",
				"edge": "implicit"
			}
		],
		"d1": [
			{
				"node": "rs0",
				"edge": "explicit"
			},
			{
				"node": "rs1",
				"edge": "explicit"
			},
			{
				"node": "rs2",
				"edge": "explicit"
			}
		]
	},
	"objects": {
		"rs0": {
			"name": "grafana-6d4fd8c49",
			"apiVersion": "apps/v1",
			"kind": "ReplicaSet",
			"status": "ok"
		},
		"rs1": {
			"name": "grafana-99c8784f6",
			"apiVersion": "apps/v1",
			"kind": "ReplicaSet",
			"status": "ok"
		},
		"rs2": {
			"name": "grafana-6b5b79d6cf",
			"apiVersion": "apps/v1",
			"kind": "ReplicaSet",
			"status": "ok"
		},
		"s1": {
			"name": "grafana",
			"apiVersion": "v1",
			"kind": "Service",
			"status": "ok",
			"isNetwork": true
		},
		"d1": {
			"name": "grafana",
			"apiVersion": "apps/v1",
			"kind": "Deployment",
			"status": "ok"
		},
		"rs3": {
			"name": "grafana-d69f77cc4",
			"apiVersion": "apps/v1",
			"kind": "ReplicaSet",
			"status": "ok"
		},
		"pods-rs0": {
			"name": "pods-rs0",
			"apiVersion": "v1",
			"kind": "pods",
			"status": "ok"
		}
	}
}`)
)

package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/content"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type resourceViewer struct{}

func (rv *resourceViewer) IsEmpty() bool {
	return false
}

func (rv *resourceViewer) MarshalJSON() ([]byte, error) {
	b := []byte(`{
		"type": "resourceviewer",
		"selected": "pod",
		"dag": {
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

	return b, nil
}

type ResourceViewerStub struct{}

func NewResourceViewerStub(prefix, namespace string, c clock.Clock) View {
	return &ResourceViewerStub{}
}

func (rss *ResourceViewerStub) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	return []content.Content{&resourceViewer{}}, nil
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ResourceViewer_Marshal(t *testing.T) {
	cases := []struct {
		name         string
		input        *ResourceViewer
		expectedPath string
		isErr        bool
	}{
		{
			name: "in general",
			input: &ResourceViewer{
				Config: ResourceViewerConfig{
					Edges: AdjList{
						"69e4ea11-2985-11e9-b356-42010a8000e5": []Edge{
							{
								Node: "bf4800b5b6602c4c78ba3b654af02b3b",
								Type: "explicit",
							},
						},
						"71c2b4eb-2949-11e9-b356-42010a8000e5": []Edge{
							{
								Node: "8682460a-29b5-11e9-b356-42010a8000e5",
								Type: "explicit",
							},
						},
						"8682460a-29b5-11e9-b356-42010a8000e5": []Edge{
							{
								Node: "bf4800b5b6602c4c78ba3b654af02b3b",
								Type: "explicit",
							},
						},
					},
					Nodes: Nodes{
						"69e4ea11-2985-11e9-b356-42010a8000e5": Node{
							Name:       "my-nginx",
							APIVersion: "v1",
							Kind:       "Service",
							Status:     "ok",
							Path:       NewLink("", "my-nginx", "/overview/namespace/default/discovery-and-load-balancing/services/my-nginx"),
						},
						"71c2b4eb-2949-11e9-b356-42010a8000e5": Node{
							Name:       "nginx-deployment",
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Status:     "ok",
							Path:       NewLink("", "nginx-deployment", "/overview/namespace/default/workloads/deployments/nginx-deployment"),
						},
						"8682460a-29b5-11e9-b356-42010a8000e5": Node{
							Name:       "nginx-deployment-56c74bb7cd",
							APIVersion: "extensions/v1beta1",
							Kind:       "ReplicaSet",
							Status:     "ok",
							Path:       NewLink("", "nginx-deployment-56c74bb7cd", "/overview/namespace/default/workloads/replica-sets/nginx-deployment-56c74bb7cd"),
						},
						"bf4800b5b6602c4c78ba3b654af02b3b": Node{
							Name:       "nginx-deployment-56c74bb7cd pods",
							APIVersion: "v1",
							Kind:       "Pod",
							Status:     "ok",
							Path:       NewLink("", "nginx-deployment-56c74bb7cd pods", "/overview/namespace/default/workloads/pods/nginx-deployment-56c74bb7cd pods"),
						},
					},
				},
				Base: newBase(TypeResourceViewer, TitleFromString("Resource Viewer")),
			},
			expectedPath: "resource_viewer.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := err != nil
			if isErr != tc.isErr {
				t.Fatalf("Unexpected error: %v", err)
			}

			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err, "reading test fixtures")
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

func Test_ResourceViewer_AddEdge(t *testing.T) {
	rv := NewResourceViewer("Resource Viewer")

	node := Node{}
	childNode := Node{}
	rv.AddNode("nodeID", node)
	rv.AddNode("childID", childNode)

	require.NoError(t, rv.AddEdge("nodeID", "childID", EdgeTypeExplicit))

	expected := ResourceViewerConfig{
		Edges: AdjList{
			"nodeID": []Edge{
				{Node: "childID", Type: EdgeTypeExplicit},
			},
		},
		Nodes: Nodes{
			"nodeID":  node,
			"childID": childNode,
		},
	}

	assert.Equal(t, expected, rv.Config)
}

func Test_ResourceViewer_AddEdge_missing_node(t *testing.T) {
	rv := NewResourceViewer("Resource Viewer")

	require.Error(t, rv.AddEdge("nodeID", "childID", EdgeTypeExplicit))

	node := Node{}
	rv.AddNode("nodeID", node)

	require.Error(t, rv.AddEdge("nodeID", "childID", EdgeTypeExplicit))
}

func Test_ResourceViewer_AddNode(t *testing.T) {
	rv := NewResourceViewer("Resource Viewer")

	node := Node{
		Name:       "my-nginx",
		APIVersion: "v1",
		Kind:       "Service",
		Status:     "ok",
	}

	rv.AddNode("nodeID", node)

	expected := ResourceViewerConfig{
		Edges: AdjList{},
		Nodes: Nodes{"nodeID": node},
	}

	assert.Equal(t, expected, rv.Config)
}

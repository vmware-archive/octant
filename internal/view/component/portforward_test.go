package component

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PortForward_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    ViewComponent
		expected string
		isErr    bool
	}{
		{
			name: "initial state",
			input: &PortForward{
				Config: PortForwardConfig{
					Text:   "9090/TCP",
					Action: PortForwardActionCreate,
					Status: PortForwardStatusInitial,
					Ports:  []PortForwardPortSpec{},
					Target: PortForwardTarget{
						APIVersion: "v1",
						Kind:       "Pod",
						Namespace:  "default",
						Name:       "mypod",
					},
				},
			},
			expected: `
           {
             "metadata": {
               "type": "portforward"
             },
             "config": {
               "text": "9090/TCP",
               "id": "",
               "action": "create",
               "status": "initial",
               "ports": [],
               "target": {
                 "apiVersion": "v1",
                 "kind": "Pod",
                 "namespace": "default",
                 "name": "mypod"
               }
             }
           }
`,
		},
		{
			name: "running state",
			input: &PortForward{
				Config: PortForwardConfig{
					Text:   "9090/TCP",
					Action: PortForwardActionDelete,
					Status: PortForwardStatusRunning,
					ID:     "f9d6db9d-c2d8-4346-a115-ad9dcad08dca",
					Ports: []PortForwardPortSpec{
						PortForwardPortSpec{
							Local:  54356,
							Remote: 9090,
						},
					},
					Target: PortForwardTarget{
						APIVersion: "v1",
						Kind:       "Pod",
						Namespace:  "default",
						Name:       "mypod",
					},
				},
			},
			expected: `
            {
              "metadata": {
                "type": "portforward"
              },
              "config": {
                "text": "9090/TCP",
                "id": "f9d6db9d-c2d8-4346-a115-ad9dcad08dca",
                "action": "delete",
                "status": "running",
                "ports": [
                  {
                    "local": 54356,
                    "remote": 9090
                  }
                ],
                "target": {
                  "apiVersion": "v1",
                  "kind": "Pod",
                  "namespace": "default",
                  "name": "mypod"
                }
              }
            }
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := (err != nil)
			if isErr != tc.isErr {
				t.Fatalf("Unexpected error: %v", err)
			}

			t.Logf("actual\n%v", string(actual))

			assert.JSONEq(t, tc.expected, string(actual))
		})
	}
}

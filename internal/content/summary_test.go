package content

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummary(t *testing.T) {
	sections := []Section{
		{
			Title: "Primary",
			Items: []Item{
				TextItem("Name", "nginx"),
				TextItem("Namespace", "default"),
			},
		},
		{
			Title: "Network",
			Items: []Item{
				LinkItem("docker-for-desktop", "click-here", "/api/node/blah"),
				TextItem("IP", "10.1.68.108"),
				JSONItem("health", map[string]interface{}{
					"status":      "OK",
					"lastChecked": "Yesterday",
					"details": map[string]interface{}{
						"cluster": "Not broken",
						"demo":    "Welp",
					},
				}),
			},
		},
	}

	summary := NewSummary("details", sections)

	mockB, err := ioutil.ReadFile(filepath.Join("testdata", "summary_mock.json"))
	require.NoError(t, err)

	var expected Summary
	err = json.Unmarshal(mockB, &expected)
	require.NoError(t, err)

	assert.Equal(t, expected, summary)
}

func TestSummary_IsEmpty(t *testing.T) {
	summary := NewSummary("summary", nil)
	assert.True(t, summary.IsEmpty())
}

func TestLabelsItem(t *testing.T) {
	cases := []struct {
		name     string
		labels   map[string]string
		expected Item
	}{
		{
			name: "empty",
			expected: Item{
				Type:  "text",
				Label: "name",
				Data: map[string]interface{}{
					"value": "<none>",
				},
			},
		},
		{
			name: "with content",
			labels: map[string]string{
				"z": "content",
				"b": "content",
			},
			expected: Item{
				Type:  "labels",
				Label: "name",
				Data: map[string]interface{}{
					"items": []Item{
						Item{
							Type:  "text",
							Label: "name",
							Data: map[string]interface{}{
								"value": "b=content",
							},
						},
						Item{
							Type:  "text",
							Label: "name",
							Data: map[string]interface{}{
								"value": "z=content",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			item := LabelsItem("name", tc.labels)
			assert.Equal(t, tc.expected, item)
		})
	}
}

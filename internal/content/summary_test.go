package content

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
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

	mockB, err := ioutil.ReadFile("./summary_mock.json")
	if err != nil {
		panic(err)
	}

	var expected Summary
	if err := json.Unmarshal(mockB, &expected); err != nil {
		panic(err)
	}

	if diff := cmp.Diff(summary, expected); diff != "" {
		require.Fail(t, diff)
	}
}

package content

import (
	"encoding/json"
	"io/ioutil"
	"testing"

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
				LinkItem("docker-for-desktop", "/api/node/blah"),
				TextItem("IP", "10.1.68.108"),
				JSONItem("health", map[string]interface{}{
					"status":      "OK",
					"lastChecked": "Yesterday",
					"details": map[string]string{
						"cluster": "Not broken",
						"demo":    "Welp",
					},
				}),
			},
		},
	}

	summary := NewSummary("details", sections)

	expectedB, err := ioutil.ReadFile("./summary_mock.json")
	if err != nil {
		panic(err)
	}

	outputB, err := json.Marshal(summary)
	if err != nil {
		t.Error(err)
	}

	expected := string(expectedB)
	output := string(outputB)

	require.Equal(t, expected, output)
}

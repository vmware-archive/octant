package api

import (
	"encoding/json"
	"net/http"
)

type pathDescription struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
	Path string `json:"path,omitempty"`
}

type description interface{}

type tableResponse struct {
	Type    string          `json:"type,omitempty"`
	Title   string          `json:"title,omitempty"`
	Columns []string        `json:"columns,omitempty"`
	Rows    [][]description `json:"rows,omitempty"`
}

type term struct {
	Term        string      `json:"term,omitempty"`
	Description description `json:"description,omitempty"`
}

type cardResponse struct {
	Name  string `json:"name,omitempty"`
	Terms []term `json:"terms,omitempty"`
}

type cardsResponse struct {
	Type  string         `json:"type,omitempty"`
	Title string         `json:"title,omitempty"`
	Cards []cardResponse `json:"cards,omitempty"`
}

type contentResponse interface {
}

type contentsResponse struct {
	Contents []contentResponse `json:"contents,omitempty"`
}

type content struct{}

var _ http.Handler = (*content)(nil)

func (c *content) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cr := &contentsResponse{
		Contents: []contentResponse{
			&cardsResponse{
				Type:  "cards",
				Title: "Details",
				Cards: []cardResponse{
					{
						Name: "Pod",
						Terms: []term{
							{
								Term:        "Name",
								Description: "nginx",
							},
							{
								Term:        "Namespace",
								Description: "default",
							},
							{
								Term:        "Creation Time",
								Description: "2018-09-26T00:42UTC",
							},
							{
								Term:        "Status",
								Description: "Running",
							},
							{
								Term:        "QoS Class",
								Description: "BestEffort",
							},
						},
					},
					{
						Name: "Network",
						Terms: []term{
							{
								Term: "Node",
								Description: &pathDescription{
									Type: "path",
									Path: "/node/node-1",
									Text: "node-1",
								},
							},
							{
								Term:        "IP",
								Description: "10.1.68.108",
							},
						},
					},
				},
			},
			&tableResponse{
				Type:  "table",
				Title: "Conditions",
				Columns: []string{
					"Type",
					"Status",
					"Last heartbeat time",
					"Last transition time",
					"Reason",
					"Message",
				},
				Rows: [][]description{
					{"Initialized", "True", "", "2 minutes", "", ""},
					{
						"Ready",
						"False",
						"",
						"2 minutes",
						"ContainersNotReady",
						"containers with unready status: [debian-container]",
					},
					{"PodScheduled", "True", "", "2 minutes", "", ""},
				},
			},
		},
	}

	json.NewEncoder(w).Encode(cr)
}

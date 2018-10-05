package overview

import (
	"encoding/json"
)

type contentResponse struct {
	Contents []content `json:"contents,omitempty"`
}

type content interface {
}

type table struct {
	Type    string        `json:"type,omitempty"`
	Title   string        `json:"title,omitempty"`
	Columns []tableColumn `json:"columns,omitempty"`
	Rows    []tableRow    `json:"rows,omitempty"`
}

func newTable(title string) table {
	return table{
		Type:  "table",
		Title: title,
	}
}

func (t *table) AddRow(row tableRow) {
	t.Rows = append(t.Rows, row)
}

type tableColumn struct {
	Name     string `json:"name,omitempty"`
	Accessor string `json:"accessor,omitempty"`
}

type tableRow map[string]text

// text is something than can be presented as text.
type text interface {
	json.Marshaler
}

// stringText is text that is presented as a string.
type stringText string

// newStringText creates an instance of stringText.
func newStringText(s string) *stringText {
	st := stringText(s)
	return &st
}

func (s stringText) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"type": "string",
		"text": string(s),
	}

	return json.Marshal(&m)
}

// linkText is text that contains a link.
type linkText struct {
	Text string
	Ref  string
}

// newLinkText create an instance of linkText.
func newLinkText(s string, ref string) *linkText {
	return &linkText{
		Text: s,
		Ref:  ref,
	}
}

func (t *linkText) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"type": "link",
		"text": t.Text,
		"ref":  t.Ref,
	}

	return json.Marshal(&m)
}

// labelsText is text that contains labels.
type labelsText struct {
	Labels map[string]string
}

// newLabelsText create an instance of labelsText.
func newLabelsText(labels map[string]string) *labelsText {
	return &labelsText{
		Labels: labels,
	}
}

func (t *labelsText) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"type":   "labels",
		"labels": t.Labels,
	}

	return json.Marshal(&m)
}

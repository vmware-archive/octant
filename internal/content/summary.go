package content

import (
	"fmt"
	"sort"
)

var _ Content = (*Summary)(nil)

type Summary struct {
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Sections []Section `json:"sections"`
}

func (s *Summary) IsEmpty() bool {
	return len(s.Sections) == 0
}

type Section struct {
	Title string `json:"title"`
	Items []Item `json:"items"`
}

func NewSection() Section {
	return Section{}
}

func (s *Section) AddText(label, text string) {
	s.Items = append(s.Items, TextItem(label, text))
}

func (s *Section) AddLabels(label string, labels map[string]string) {
	s.Items = append(s.Items, LabelsItem(label, labels))
}

func (s *Section) AddList(label string, kv map[string]string) {
	s.Items = append(s.Items, ListItem(label, kv))
}

func (s *Section) AddLink(label, value, link string) {
	s.Items = append(s.Items, LinkItem(label, value, link))
}

func (s *Section) AddJSON(label string, blob interface{}) {
	s.Items = append(s.Items, JSONItem(label, blob))
}

type Item struct {
	Type  string      `json:"type"`
	Label string      `json:"label"`
	Data  interface{} `json:"data"`
}

func TextItem(label, text string) Item {
	return Item{
		Type:  "text",
		Label: label,
		Data: map[string]interface{}{
			"value": text,
		},
	}
}

func LabelsItem(label string, labels map[string]string) Item {
	if len(labels) == 0 {
		return TextItem(label, "<none>")
	}

	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var out []Item
	for _, key := range keys {
		labelSelector := fmt.Sprintf("%s=%s", key, labels[key])
		out = append(out, TextItem(label, labelSelector))
	}

	return Item{
		Type:  "labels",
		Label: label,
		Data: map[string]interface{}{
			"items": out,
		},
	}
}

func ListItem(label string, kv map[string]string) Item {
	if len(kv) == 0 {
		return TextItem(label, "<none>")
	}

	keys := make([]string, 0, len(kv))
	for key := range kv {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var out []Item
	for _, key := range keys {
		labelSelector := fmt.Sprintf("%s=%s", key, kv[key])
		out = append(out, TextItem("", labelSelector))
	}

	return Item{
		Type:  "list",
		Label: label,
		Data: map[string]interface{}{
			"items": out,
		},
	}
}

func LinkItem(label, value, link string) Item {
	return Item{
		Type:  "link",
		Label: label,
		Data: map[string]interface{}{
			"value": value,
			"ref":   link,
		},
	}
}

func JSONItem(label string, blob interface{}) Item {
	return Item{
		Type:  "json",
		Label: label,
		Data:  blob,
	}
}

func NewSummary(title string, sections []Section) Summary {
	return Summary{
		Type:     "summary",
		Title:    title,
		Sections: sections,
	}
}

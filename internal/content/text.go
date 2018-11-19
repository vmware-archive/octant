package content

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// Text is something than can be presented as text.
type Text interface {
	json.Marshaler
}

// StringText is text that is presented as a string.
type StringText string

// NewStringText creates an instance of StringText.
func NewStringText(s string) *StringText {
	st := StringText(s)
	return &st
}

func (s StringText) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"type": "string",
		"text": string(s),
	}

	return json.Marshal(&m)
}

type TimeText string

// NewTimeText creates an instance of TimeText
func NewTimeText(s string) *TimeText {
	tt := TimeText(s)
	return &tt
}

// ParseTimeText returns empty string if not in RFC3339 format
func ParseTimeText(t TimeText) (*string, error) {
	ts := string(t)
	parsedTime, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, errors.Errorf("Timestamp is not RFC3339: %v", ts)
	}

	if parsedTime.IsZero() {
		ts = ""
	}
	return &ts, nil
}

func (t TimeText) MarshalJSON() ([]byte, error) {
	ts, err := ParseTimeText(t)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{
		"type": "time",
		"time": ts,
	}

	return json.Marshal(&m)
}

// LinkText is text that contains a link.
type LinkText struct {
	Text string
	Ref  string
}

// NewLinkText create an instance of linkText.
func NewLinkText(s string, ref string) *LinkText {
	return &LinkText{
		Text: s,
		Ref:  ref,
	}
}

func (t *LinkText) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"type": "link",
		"text": t.Text,
		"ref":  t.Ref,
	}

	return json.Marshal(&m)
}

// LabelsText is text that contains labels.
type LabelsText struct {
	Labels map[string]string
}

// NewLabelsText create an instance of LabelsText.
func NewLabelsText(labels map[string]string) *LabelsText {
	return &LabelsText{
		Labels: labels,
	}
}

func (t *LabelsText) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"type":   "labels",
		"labels": t.Labels,
	}

	return json.Marshal(&m)
}

// ListText is text that contains a list.
type ListText struct {
	List []string
}

// NewListText create an instance of listText.
func NewListText(list []string) *ListText {
	return &ListText{
		List: list,
	}
}

func (t *ListText) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"type": "list",
		"list": t.List,
	}

	return json.Marshal(&m)
}

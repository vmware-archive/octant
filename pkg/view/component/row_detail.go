package component

import (
	"github.com/vmware-tanzu/octant/internal/util/json"
)

const (
	// ExpandableRowKey is the key for an expandable row
	ExpandableRowKey = "_expand"
)

// AddExpandableDetail allows a row to be expandable
func (t TableRow) AddExpandableDetail(body Component, replaceContent bool) {
	er, ok := t[ExpandableRowKey].(*ExpandableRowDetail)
	if !ok {
		er = NewExpandableRowDetail(body, replaceContent)
	}
	t[ExpandableRowKey] = er
}

// ExpandableRowDetail is a component hidden by a toggle under each table row.
//
// +octant:component
type ExpandableRowDetail struct {
	Base

	Config ExpandableDetailConfig `json:"config"`
}

// ExpandableDetailConfig is a configuration for an expandable row detail.
type ExpandableDetailConfig struct {
	Replace bool      `json:"replace"`
	Body    Component `json:"body"`
}

// Unmarshal unmarshals an expandable row detail config from JSON.
func (e *ExpandableDetailConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Body    *TypedObject `json:"body"`
		Replace bool         `json:"replace"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x.Body != nil {
		var err error
		e.Body, err = x.Body.ToComponent()
		if err != nil {
			return err
		}
	}
	e.Replace = x.Replace
	return nil
}

var _ Component = &ExpandableRowDetail{}

// NewExpandableRowDetail creates an expandable detail for a table row.
func NewExpandableRowDetail(body Component, replaceContent bool) *ExpandableRowDetail {
	e := ExpandableRowDetail{
		Base: newBase(TypeExpandableRowDetail, nil),
		Config: ExpandableDetailConfig{
			Body:    body,
			Replace: replaceContent,
		},
	}
	return &e
}

type expandableDetailMarshal ExpandableRowDetail

func (e ExpandableRowDetail) MarshalJSON() ([]byte, error) {
	k := expandableDetailMarshal{
		Base:   e.Base,
		Config: e.Config,
	}
	k.Metadata.Type = TypeExpandableRowDetail
	return json.Marshal(&k)
}

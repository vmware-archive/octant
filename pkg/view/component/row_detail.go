package component

import (
	"github.com/vmware-tanzu/octant/internal/util/json"
)

const (
	// ExpandableRowKey is the key for an expandable row
	ExpandableRowKey = "_expand"
)

// AddExpandableDetail allows a row to be expandable
func (t TableRow) AddExpandableDetail(details *ExpandableRowDetail) {
	t[ExpandableRowKey] = details
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
	Replace bool        `json:"replace"`
	Body    []Component `json:"body"`
}

// UnmarshalJSON unmarshals an expandable row detail config from JSON.
func (e *ExpandableDetailConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Body    []TypedObject `json:"body"`
		Replace bool          `json:"replace"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	e.Replace = x.Replace

	if x.Body == nil {
		return nil
	}

	for _, typedObject := range x.Body {
		component, err := typedObject.ToComponent()
		if err != nil {
			return err
		}

		e.Body = append(e.Body, component)
	}
	return nil
}

var _ Component = &ExpandableRowDetail{}

// NewExpandableRowDetail creates an expandable detail for a table row.
func NewExpandableRowDetail(body ...Component) *ExpandableRowDetail {
	e := ExpandableRowDetail{
		Base: newBase(TypeExpandableRowDetail, nil),
		Config: ExpandableDetailConfig{
			Body:    body,
			Replace: false,
		},
	}
	return &e
}

func (e *ExpandableRowDetail) SetReplace(replace bool) {
	e.Config.Replace = replace
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

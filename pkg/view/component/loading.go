package component

import "encoding/json"

// Loading is a component for a spinner
type Loading struct {
	base
	Config LoadingConfig `json:"config"`
}

// LoadingConfig is the contents of Loading
type LoadingConfig struct {
	Text string `json:"value"`
}

// NewLoading creates a loading component
func NewLoading(title []TitleComponent, message string) *Loading {
	return &Loading{
		base: newBase(typeLoading, title),
		Config: LoadingConfig{
			Text: message,
		},
	}
}

// SupportsTitle denotes this is a LoadingComponent.
func (t *Loading) SupportsTitle() {}

type loadingMarshal Loading

// MarshalJSON implements json.Marshaler
func (t *Loading) MarshalJSON() ([]byte, error) {
	m := loadingMarshal(*t)
	m.Metadata.Type = typeLoading
	return json.Marshal(&m)
}

// String returns the text content of the component.
func (t *Loading) String() string {
	return t.Config.Text
}

// LessThan returns true if this component's value is less than the argument supplied.
func (t *Loading) LessThan(i interface{}) bool {
	v, ok := i.(*Loading)
	if !ok {
		return false
	}

	return t.Config.Text < v.Config.Text

}

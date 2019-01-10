package component

import "encoding/json"

// ExpressionSelector is a component for a single expression within a selector
type ExpressionSelector struct {
	Metadata Metadata                 `json:"metadata"`
	Config   ExpressionSelectorConfig `json:"config"`
}

// ExpressionSelectorConfig is the contents of ExpressionSelector
type ExpressionSelectorConfig struct {
	Key      string   `json:"key"`
	Operator Operator `json:"operator"`
	Values   []string `json:"values"`
}

// Operator represents a key's relationship to a set of values.
// Valid operators are In, NotIn, Exists and DoesNotExist.
type Operator string

const (
	// OperatorIn means a key value is in a set of possible values
	OperatorIn Operator = "In"
	// OperatorNotIn means a key value is not in a set of exclusionary values
	OperatorNotIn Operator = "NotIn"
	// OperatorExists means a key exists on the selected resource
	OperatorExists Operator = "Exists"
	// OperatorDoesNotExist means a key does not exists on the selected resource
	OperatorDoesNotExist Operator = "DoesNotExist"
)

// NewExpressionSelector creates a expressionSelector component
func NewExpressionSelector(k string, o Operator, values []string) *ExpressionSelector {
	return &ExpressionSelector{
		Metadata: Metadata{
			Type: "expressionSelector",
		},
		Config: ExpressionSelectorConfig{
			Key:      k,
			Operator: o,
			Values:   values,
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *ExpressionSelector) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifes whether the component is considered empty. Implements ViewComponent.
func (t *ExpressionSelector) IsEmpty() bool {
	return t.Config.Key == ""
}

// IsSelector marks the component as selector flavor. Implements Selector.
func (t *ExpressionSelector) IsSelector() {
}

type expressionSelectorMarshal ExpressionSelector

// MarshalJSON implements json.Marshaler
func (t *ExpressionSelector) MarshalJSON() ([]byte, error) {
	m := expressionSelectorMarshal(*t)
	m.Metadata.Type = "expressionSelector"
	return json.Marshal(&m)
}

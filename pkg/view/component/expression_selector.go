/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"

	"github.com/pkg/errors"
)

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

// MatchOperator matches an operator.
func MatchOperator(s string) (Operator, error) {
	operators := []Operator{OperatorIn, OperatorNotIn, OperatorExists, OperatorExists}
	for _, o := range operators {
		if string(o) == s {
			return o, nil
		}
	}

	return Operator("invalid"), errors.Errorf("operator %q is not valid", s)
}

// ExpressionSelector is a component for a single expression within a selector
type ExpressionSelector struct {
	base
	Config ExpressionSelectorConfig `json:"config"`
}

// NewExpressionSelector creates a expressionSelector component
func NewExpressionSelector(k string, o Operator, values []string) *ExpressionSelector {
	return &ExpressionSelector{
		base: newBase(typeExpressionSelector, nil),
		Config: ExpressionSelectorConfig{
			Key:      k,
			Operator: o,
			Values:   values,
		},
	}
}

// Name is the name of the ExpressionSelector.
func (t *ExpressionSelector) Name() string {
	return t.Config.Key
}

// ExpressionSelectorConfig is the contents of ExpressionSelector
type ExpressionSelectorConfig struct {
	Key      string   `json:"key"`
	Operator Operator `json:"operator"`
	Values   []string `json:"values"`
}

// GetMetadata accesses the components metadata. Implements Component.
func (t *ExpressionSelector) GetMetadata() Metadata {
	return t.Metadata
}

// IsSelector marks the component as selector flavor. Implements Selector.
func (t *ExpressionSelector) IsSelector() {
}

type expressionSelectorMarshal ExpressionSelector

// MarshalJSON implements json.Marshaler
func (t *ExpressionSelector) MarshalJSON() ([]byte, error) {
	m := expressionSelectorMarshal(*t)
	m.Metadata.Type = typeExpressionSelector
	return json.Marshal(&m)
}

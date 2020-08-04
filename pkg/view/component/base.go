/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

const (
	typeAnnotations        = "annotations"
	typeButtonGroup        = "buttonGroup"
	typeCard               = "card"
	typeCardList           = "cardList"
	typeCodeBlock          = "codeBlock"
	typeContainers         = "containers"
	typeDonutChart         = "donutChart"
	typeEditor             = "editor"
	typeError              = "error"
	typeExtension          = "extension"
	typeExpressionSelector = "expressionSelector"
	typeFlexLayout         = "flexlayout"
	typeGraphviz           = "graphviz"
	typeGridActions        = "gridActions"
	typeIFrame             = "iframe"
	typeLabels             = "labels"
	typeLabelSelector      = "labelSelector"
	typeLink               = "link"
	typeList               = "list"
	typeLoading            = "loading"
	typeLogs               = "logs"
	typePodStatus          = "podStatus"
	typePort               = "port"
	typePorts              = "ports"
	typeQuadrant           = "quadrant"
	typeResourceViewer     = "resourceViewer"
	typeSelectors          = "selectors"
	typeSingleStat         = "singleStat"
	typeStepper            = "stepper"
	typeSummary            = "summary"
	typeTable              = "table"
	typeTerminal           = "terminal"
	typeText               = "text"
	typeTimestamp          = "timestamp"
	typeYAML               = "yaml"
)

// base is a base component.
type base struct {
	Metadata `json:"metadata"`
}

func newBase(componentType string, title []TitleComponent) base {
	return base{
		Metadata: Metadata{
			Type:  componentType,
			Title: title,
		},
	}
}

// GetMetadata returns the component's metadata.
func (b *base) GetMetadata() Metadata {
	return b.Metadata
}

func (b *base) SetMetadata(metadata Metadata) {
	b.Metadata = metadata
}

// SetAccessor sets the accessor for the object.
func (b *base) SetAccessor(accessor string) {
	b.Metadata.Accessor = accessor
}

// IsEmpty returns false by default. Let the components that wrap base
// determine if they are empty or not if they wish.
func (b *base) IsEmpty() bool {
	return false
}

// String returns an empty string. If a component wants to provide a value
// it can override this function.
func (b *base) String() string {
	return ""
}

// LessThan returns false.
func (b *base) LessThan(_ interface{}) bool {
	return false
}

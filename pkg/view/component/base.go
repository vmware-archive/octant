/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

const (
	// TypeAccordion is an accordion component.
	TypeAccordion = "accordion"
	// TypeAnnotations is an annotations component.
	TypeAnnotations = "annotations"
	// TypeButtonGroup is a button group component.
	TypeButtonGroup = "buttonGroup"
	// TypeCard is a card component.
	TypeCard = "card"
	// TypeCardList is a card list component.
	TypeCardList = "cardList"
	// TypeCode is a code block component.
	TypeCode = "codeBlock"
	// TypeContainers is a container component.
	TypeContainers = "containers"
	// TypeDonutChart is a donut chart component.
	TypeDonutChart = "donutChart"
	// TypeDropdown is a dropdown component.
	TypeDropdown = "dropdown"
	// TypeEditor is an editor component.
	TypeEditor = "editor"
	// TypeError is an error component.
	TypeError = "error"
	// TypeExpandableRowDetail is an expandable detail component for table rows.
	TypeExpandableRowDetail = "expandableRowDetail"
	// TypeExtension is an extension component.
	TypeExtension = "extension"
	// TypeExpressionSelector is an expression selector component.
	TypeExpressionSelector = "expressionSelector"
	// TypeFlexLayout is a flex layout component.
	TypeFlexLayout = "flexlayout"
	// TypeFormField is a form field component.
	TypeFormField = "formField"
	// TypeGraphviz is a graphviz component.
	TypeGraphviz = "graphviz"
	// TypeGridActions is a grid actions component.
	TypeGridActions = "gridActions"
	// TypeIFrame is an iframe component.
	TypeIFrame = "iframe"
	// TypeJSONEditor is a JSON Editor component.
	TypeJSONEditor = "jsonEditor"
	// TypeLabels is a labels component.
	TypeLabels = "labels"
	// TypeLabelSelector is a label selector component.
	TypeLabelSelector = "labelSelector"
	// TypeLink is a link component.
	TypeLink = "link"
	// TypeList is a list component.
	TypeList = "list"
	// TypeLoading is a loading component.
	TypeLoading = "loading"
	// TypeLogs is a logs component.
	TypeLogs = "logs"
	// TypeModal is a modal component.
	TypeModal = "modal"
	// TypePodStatus is a pod status component.
	TypePodStatus = "podStatus"
	// TypePort is a port component.
	TypePort = "port"
	// TypePorts is a ports component.
	TypePorts = "ports"
	// TypeQuadrant is a quadrant component.
	TypeQuadrant = "quadrant"
	// TypeResourceViewer is a resource viewer component.
	TypeResourceViewer = "resourceViewer"
	// TypeSelectFile is a SelectFile component.
	TypeSelectFile = "selectFile"
	// TypeSelectors is a selectors component.
	TypeSelectors = "selectors"
	// TypeSingleStat is a single stat component.
	TypeSingleStat = "singleStat"
	// TypeStepper is a stepper component.
	TypeStepper = "stepper"
	// TypeSummary is a summary component.
	TypeSummary = "summary"
	// TypeTable is a table component.
	TypeTable = "table"
	// TypeTerminal is a terminal component.
	TypeTerminal = "terminal"
	// TypeText is a text component.
	TypeText = "text"
	// TypeTimeline is a timeline component.
	TypeTimeline = "timeline"
	// TypeTimestamp is a timestamp component.
	TypeTimestamp = "timestamp"
	// TypeYAML is a YAML component.
	TypeYAML = "yaml"
	// TypeIcon is a Icon component.
	TypeIcon = "icon"
	// TypeSignpost is a SignPost component.
	TypeSignpost = "signpost"
	// TypeButton is a Button component.
	TypeButton = "button"
)

// Base is an abstract base for components..
type Base struct {
	Metadata `json:"metadata"`
}

func newBase(componentType string, title []TitleComponent) Base {
	return Base{
		Metadata: Metadata{
			Type:  componentType,
			Title: title,
		},
	}
}

// GetMetadata returns the component's metadata.
func (b *Base) GetMetadata() Metadata {
	return b.Metadata
}

func (b *Base) SetMetadata(metadata Metadata) {
	b.Metadata = metadata
}

// SetAccessor sets the accessor for the object.
func (b *Base) SetAccessor(accessor string) {
	b.Metadata.Accessor = accessor
}

// IsEmpty returns false by default. Let the components that wrap Base
// determine if they are empty or not if they wish.
func (b *Base) IsEmpty() bool {
	return false
}

// String returns an empty string. If a component wants to provide a value
// it can override this function.
func (b *Base) String() string {
	return ""
}

// LessThan returns false.
func (b *Base) LessThan(_ interface{}) bool {
	return false
}

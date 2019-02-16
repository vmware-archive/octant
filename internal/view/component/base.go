package component

const (
	typeAnnotations        = "annotations"
	typeContainers         = "containers"
	typeExpressionSelector = "expressionSelector"
	typeFlexLayout         = "flexlayout"
	typeGrid               = "grid"
	typeLabels             = "labels"
	typeLabelSelector      = "labelSelector"
	typeLink               = "link"
	typeList               = "list"
	typeLogs               = "logs"
	typePanel              = "panel"
	typeQuadrant           = "quadrant"
	typeResourceViewer     = "resourceViewer"
	typeSelectors          = "selectors"
	typeSummary            = "summary"
	typeTable              = "table"
	typeText               = "text"
	typeTimestamp          = "timestamp"
	typeYAML               = "yaml"
)

// base is a base component.
type base struct {
	Metadata `json:"metadata"`
}

func newBase(componentType string, title []TitleViewComponent) base {
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

// SetAccessor sets the accessor for the object.
func (b *base) SetAccessor(accessor string) {
	b.Metadata.Accessor = accessor
}

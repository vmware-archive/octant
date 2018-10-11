package overview

// SectionDescriber is a wrapper to combine content from multiple describers.
type SectionDescriber struct {
	describers []Describer
}

// NewSectionDescriber creates a SectionDescriber.
func NewSectionDescriber(describers ...Describer) *SectionDescriber {
	return &SectionDescriber{
		describers: describers,
	}
}

// Describe generates content.
func (d *SectionDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error) {
	var contents []Content

	for _, child := range d.describers {
		childContent, err := child.Describe(prefix, namespace, cache, fields)
		if err != nil {
			return nil, err
		}

		contents = append(contents, childContent...)
	}

	return contents, nil
}

package describer

import (
	"context"

	"github.com/heptio/developer-dash/pkg/view/component"
)

// Section is a wrapper to combine content from multiple describers.
type Section struct {
	path       string
	title      string
	describers []Describer
}

// NewSection creates a Section.
func NewSection(p, title string, describers ...Describer) *Section {
	return &Section{
		path:       p,
		title:      title,
		describers: describers,
	}
}

// Describe generates content.
func (d *Section) Describe(ctx context.Context, prefix, namespace string, options Options) (component.ContentResponse, error) {
	list := component.NewList(d.title, nil)

	for _, child := range d.describers {
		cResponse, err := child.Describe(ctx, prefix, namespace, options)
		if err != nil {
			return EmptyContentResponse, err
		}

		for _, vc := range cResponse.Components {
			if nestedList, ok := vc.(*component.List); ok {
				for i := range nestedList.Config.Items {
					item := nestedList.Config.Items[i]
					if !item.IsEmpty() {
						list.Add(item)
					}
				}
			}
		}
	}

	cr := component.ContentResponse{
		Components: []component.Component{list},
		Title:      component.Title(component.NewText(d.title)),
	}

	return cr, nil
}

// PathFilters returns path filters for the section.
func (d *Section) PathFilters() []PathFilter {
	PathFilters := []PathFilter{
		*NewPathFilter(d.path, d),
	}

	for _, child := range d.describers {
		PathFilters = append(PathFilters, child.PathFilters()...)
	}

	return PathFilters
}

package overview

import "github.com/heptio/developer-dash/internal/content"

type ContentResponse struct {
	Title          string                  `json:"title,omitempty"`
	Views          []Content               `json:"views,omitempty"`
	ViewComponents []content.ViewComponent `json:"viewComponents"`
}

var emptyContentResponse = ContentResponse{}

type Content struct {
	Contents []content.Content `json:"contents,omitempty"`
	Title    string            `json:"title,omitempty"`
}

type ViewComponentable interface {
	ViewComponent() content.ViewComponent
}

type List struct {
	title  string
	config ListConfig
}

type ListConfig struct {
	Items []content.ViewComponent `json:"items"`
}

func NewList(title string) *List {
	return &List{
		title: title,
		config: ListConfig{
			Items: make([]content.ViewComponent, 0),
		},
	}
}

func (l *List) Add(item ViewComponentable) {
	l.config.Items = append(l.config.Items, item.ViewComponent())
}

func (l *List) ViewComponent() content.ViewComponent {
	return content.ViewComponent{
		Metadata: content.Metadata{
			Type:  "list",
			Title: l.title,
		},
		Config: l.config,
	}
}

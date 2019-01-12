package overview

import (
	"github.com/heptio/developer-dash/internal/content"
	"github.com/heptio/developer-dash/internal/view/component"
)

type ContentResponse struct {
	Title          string                    `json:"title,omitempty"`
	ViewComponents []component.ViewComponent `json:"viewComponents"`
	Views          []Content                 `json:"-"`
}

var emptyContentResponse = ContentResponse{}

type Content struct {
	Contents []content.Content `json:"contents,omitempty"`
	Title    string            `json:"title,omitempty"`
}

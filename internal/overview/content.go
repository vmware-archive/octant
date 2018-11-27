package overview

import "github.com/heptio/developer-dash/internal/content"

type ContentResponse struct {
	Views       map[string]Content `json:"views,omitempty"`
	DefaultView string             `json:"default_view,omitempty"`
}

var emptyContentResponse = ContentResponse{}

type Content struct {
	Contents []content.Content `json:"contents,omitempty"`
	Title    string            `json:"title,omitempty"`
}

package overview

import "github.com/heptio/developer-dash/internal/content"

type ContentResponse struct {
	Title string    `json:"title,omitempty"`
	Views []Content `json:"views,omitempty"`
}

var emptyContentResponse = ContentResponse{}

type Content struct {
	Contents []content.Content `json:"contents,omitempty"`
	Title    string            `json:"title,omitempty"`
}

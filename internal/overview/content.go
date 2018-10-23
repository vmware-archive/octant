package overview

import "github.com/heptio/developer-dash/internal/content"

type ContentResponse struct {
	Contents []content.Content `json:"contents,omitempty"`
	Title    string            `json:"title,omitempty"`
}

var emptyContentResponse = ContentResponse{}

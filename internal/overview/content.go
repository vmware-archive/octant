package overview

import "github.com/heptio/developer-dash/internal/content"

type contentResponse struct {
	Contents []content.Content `json:"contents,omitempty"`
}

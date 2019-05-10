package clustereye_test

import (
	"testing"

	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/stretchr/testify/assert"
)

func Test_NewNavigation(t *testing.T) {
	path := "/path"
	title := "title"

	nav := clustereye.NewNavigation(title, path)

	assert.Equal(t, path, nav.Path)
	assert.Equal(t, title, nav.Title)
}

package fake

import (
	"testing"

	"github.com/heptio/developer-dash/internal/module"
	"github.com/stretchr/testify/assert"
)

func TestStubManager(t *testing.T) {
	m := NewModule("module")

	sm := NewStubManager("default", []module.Module{m})

	assert.Equal(t, []module.Module{m}, sm.Modules())
	assert.Equal(t, "default", sm.GetNamespace())

	sm.SetNamespace("other")
	assert.Equal(t, "other", sm.GetNamespace())
}

package fake

import (
	"context"
	"testing"

	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestModule_Name(t *testing.T) {
	m := NewModule("module", log.NopLogger())
	assert.Equal(t, "module", m.Name())
}

func TestModule_ContentPath(t *testing.T) {
	m := NewModule("module", log.NopLogger())
	assert.Equal(t, "/module", m.ContentPath())
}

func TestModule_Navigation(t *testing.T) {
	m := NewModule("module", log.NopLogger())

	expected := &clustereye.Navigation{
		Path:  "/module",
		Title: "module",
	}

	ctx := context.Background()
	got, err := m.Navigation(ctx, "", "/module")
	require.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestModule_Content(t *testing.T) {
	m := NewModule("module", log.NopLogger())

	cases := []struct {
		path     string
		expected component.ContentResponse
		isErr    bool
	}{
		{
			path: "/",
			expected: component.ContentResponse{
				Title: component.Title(component.NewText("/")),
			},
		},
		{
			path: "/nested",
			expected: component.ContentResponse{
				Title: component.Title(component.NewText("/nested")),
			},
		},
		{
			path:  "/missing",
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			ctx := context.Background()

			got, err := m.Content(ctx, tc.path, "/prefix", "namespace", module.ContentOptions{})
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

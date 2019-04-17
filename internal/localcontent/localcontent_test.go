package localcontent_test

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/heptio/developer-dash/internal/localcontent"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/sugarloaf"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LocalContent_Name(t *testing.T) {
	withLocalContent(t, func(lc *localcontent.LocalContent) {
		assert.Equal(t, "local", lc.Name())
	})
}

func Test_LocalContent_Content_root(t *testing.T) {
	withLocalContent(t, func(lc *localcontent.LocalContent) {
		ctx := context.Background()
		content, err := lc.Content(ctx, "/", "prefix", "namespace", module.ContentOptions{})
		require.NoError(t, err)

		assert.Equal(t, component.Title(component.NewText("Local Contents")), content.Title)
		assert.Len(t, content.Components, 1)

		table, ok := content.Components[0].(*component.Table)
		if assert.True(t, ok, "component is not a table") {
			expectedCols := component.NewTableCols("Title", "File")
			assert.Equal(t, expectedCols, table.Config.Columns)

			expectedRows := []component.TableRow{
				{
					"Title": component.NewLink("", "Sample content", "/content/local/table"),
					"File":  component.NewText("table.json"),
				},
			}
			assert.Equal(t, expectedRows, table.Rows())
		}
	})
}

func Test_LocalContent_Content_file(t *testing.T) {
	withLocalContent(t, func(lc *localcontent.LocalContent) {
		ctx := context.Background()
		content, err := lc.Content(ctx, "/table", "prefix", "namespace", module.ContentOptions{})
		require.NoError(t, err)

		assert.Equal(t, component.Title(component.NewText("Sample content")),
			content.Title)
		assert.Len(t, content.Components, 1)

		list, ok := content.Components[0].(*component.List)
		if assert.Truef(t, ok, "component is not a list (%T)", list) {
			require.Len(t, list.Config.Items, 1)
			table, ok := list.Config.Items[0].(*component.Table)
			assert.Truef(t, ok, "component is not a table (%T)", table)
		}
	})
}

func Test_LocalContent_Content_invalid_file(t *testing.T) {
	withLocalContent(t, func(lc *localcontent.LocalContent) {
		ctx := context.Background()
		_, err := lc.Content(ctx, "/invalid", "prefix", "namespace", module.ContentOptions{})
		require.Error(t, err)
	})
}

func Test_LocalContent_Navigation(t *testing.T) {
	withLocalContent(t, func(lc *localcontent.LocalContent) {
		ctx := context.Background()
		nav, err := lc.Navigation(ctx, "", "/root")
		require.NoError(t, err)

		expectedNav := &sugarloaf.Navigation{
			Title: "Local Content",
			Path:  "/root/",
			Children: []*sugarloaf.Navigation{
				{
					Title: "Sample content",
					Path:  "/root/table",
				},
			},
		}

		assert.Equal(t, expectedNav, nav)
	})
}

func withLocalContent(t *testing.T, fn func(lc *localcontent.LocalContent)) {
	lc := initLocalContent(t)
	defer os.RemoveAll(lc.Root())

	fn(lc)
}

func initLocalContent(t *testing.T) *localcontent.LocalContent {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	_, err = copyFile(filepath.Join("localdata", "table.json"),
		filepath.Join(dir, "table.json"))
	require.NoError(t, err)

	return localcontent.New(dir)
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, errors.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

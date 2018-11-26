package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	table := NewTable("title", "this table is empty")

	require.True(t, table.IsEmpty())
	require.Equal(t, "this table is empty", table.EmptyContent)

	table.Columns = []TableColumn{
		{Name: "col1", Accessor: "col1"},
		{Name: "col2", Accessor: "col2"},
	}

	table.AddRow(TableRow{
		"col1": NewStringText("c1r1"),
		"col2": NewStringText("c2r1"),
	})

	assert.False(t, table.IsEmpty())
	expectedColumns := []string{"col1", "col2"}
	assert.Equal(t, expectedColumns, table.ColumnNames())
}

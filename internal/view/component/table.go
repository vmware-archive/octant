package component

import "encoding/json"

// Table contains other ViewComponents
type Table struct {
	Metadata Metadata    `json:"metadata"`
	Config   TableConfig `json:"config"`
}

// TableConfig is the contents of a Table
type TableConfig struct {
	Columns []TableCol `json:"columns"`
	Rows    []TableRow `json:"rows"`
}

// TableCol describes a column from a table. Accessor is the key this
// column will appear as in table rows, and must be unique within a table.
type TableCol struct {
	Name     string `json:"name"`
	Accessor string `json:"accessor"`
}

// TableRow is a row in table. Each key->value represents a particular column in the row.
type TableRow map[string]ViewComponent

// NewTable creates a table component
func NewTable(title string, cols []TableCol) *Table {
	return &Table{
		Metadata: Metadata{
			Type: "table",
		},
		Config: TableConfig{
			Columns: cols,
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Table) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifes whether the component is considered empty. Implements ViewComponent.
func (t *Table) IsEmpty() bool {
	return len(t.Config.Rows) == 0 || len(t.Config.Columns) == 0
}

// Add adds additional items to the tail of the table.
func (t *Table) Add(rows ...TableRow) {
	t.Config.Rows = append(t.Config.Rows, rows...)
}

type tableMarshal Table

// MarshalJSON implements json.Marshaler
func (t *Table) MarshalJSON() ([]byte, error) {
	m := tableMarshal(*t)
	m.Metadata.Type = "table"
	return json.Marshal(&m)
}

package component

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

// TableConfig is the contents of a Table
type TableConfig struct {
	Columns      []TableCol `json:"columns"`
	Rows         []TableRow `json:"rows"`
	EmptyContent string     `json:"emptyContent"`
}

// TableCol describes a column from a table. Accessor is the key this
// column will appear as in table rows, and must be unique within a table.
type TableCol struct {
	Name     string `json:"name"`
	Accessor string `json:"accessor"`
}

// TableRow is a row in table. Each key->value represents a particular column in the row.
type TableRow map[string]Component

func (t *TableRow) UnmarshalJSON(data []byte) error {
	*t = make(TableRow)

	x := map[string]TypedObject{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	for k, v := range x {
		vc, err := v.ToComponent()
		if err != nil {
			return err
		}

		(*t)[k] = vc
	}

	return nil
}

// Table contains other Components
type Table struct {
	base
	Config TableConfig `json:"config"`

	mu sync.Mutex
}

// NewTable creates a table component
func NewTable(title string, cols []TableCol) *Table {
	return &Table{
		base: newBase(typeTable, TitleFromString(title)),
		Config: TableConfig{
			Columns: cols,
		},
	}
}

// NewTableWithRows creates a table with rows.
func NewTableWithRows(title string, cols []TableCol, rows []TableRow) *Table {
	table := NewTable(title, cols)
	table.Add(rows...)
	return table
}

// NewTableCols returns a slice of table columns, each with name/accessor
// set according to the provided keys arguments.
func NewTableCols(keys ...string) []TableCol {
	if len(keys) == 0 {
		return make([]TableCol, 0)
	}

	cols := make([]TableCol, len(keys))

	for i, k := range keys {
		cols[i].Name = k
		cols[i].Accessor = k
	}
	return cols
}

// IsEmpty returns true if there is one or more rows.
func (t *Table) IsEmpty() bool {
	return len(t.Config.Rows) < 1
}

func (t *Table) Sort(name string, reverse bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	sort.Slice(t.Rows(), func(i, j int) bool {
		a, ok := t.Config.Rows[i][name]
		if !ok {
			spew.Dump(fmt.Sprintf("%s:%d/%d", name, i, j), t.Config.Rows)
			return false
		}

		b, ok := t.Config.Rows[j][name]
		if !ok {
			spew.Dump(fmt.Sprintf("%s:%d/%d", name, i, j), t.Config.Rows)
			return false
		}

		if reverse {
			return !a.LessThan(b)
		}

		return a.LessThan(b)
	})
}

// Add adds additional items to the tail of the table. Use this function to
// add rows in a concurrency safe fashion.
func (t *Table) Add(rows ...TableRow) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Config.Rows = append(t.Config.Rows, rows...)
}

// AddColumn adds a column to the table.
func (t *Table) AddColumn(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Config.Columns = append(t.Config.Columns, TableCol{
		Name:     name,
		Accessor: name,
	})
}

// Columns returns the table columns.
func (t *Table) Columns() []TableCol {
	return t.Config.Columns
}

// Rows returns the table rows.
func (t *Table) Rows() []TableRow {
	return t.Config.Rows
}

type tableMarshal Table

// MarshalJSON implements json.Marshaler
func (t *Table) MarshalJSON() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	m := tableMarshal{
		base:   t.base,
		Config: t.Config,
	}

	m.Metadata.Type = typeTable
	return json.Marshal(&m)
}

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

// TableFilter describer a text filter for a table.
type TableFilter struct {
	Values   []string `json:"values"`
	Selected []string `json:"selected"`
}

// TableConfig is the contents of a Table
type TableConfig struct {
	Columns      []TableCol             `json:"columns"`
	Rows         []TableRow             `json:"rows"`
	EmptyContent string                 `json:"emptyContent"`
	Loading      bool                   `json:"loading"`
	Filters      map[string]TableFilter `json:"filters"`
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
func NewTable(title, placeholder string, cols []TableCol) *Table {
	return &Table{
		base: newBase(typeTable, TitleFromString(title)),
		Config: TableConfig{
			Columns:      cols,
			EmptyContent: placeholder,
			Filters:      make(map[string]TableFilter),
		},
	}
}

// NewTableWithRows creates a table with rows.
func NewTableWithRows(title, placeholder string, cols []TableCol, rows []TableRow) *Table {
	table := NewTable(title, placeholder, cols)
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

func (t *Table) SetPlaceholder(placeholder string) {
	t.Config.EmptyContent = placeholder
}

func (t *Table) Sort(name string, reverse bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	sort.SliceStable(t.Rows(), func(i, j int) bool {
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

// AddFilter adds a filter to the table. Each column can only have a
// single filter.
func (t *Table) AddFilter(columnName string, filter TableFilter) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Config.Filters[columnName] = filter
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

func (t *Table) SetIsLoading(isLoading bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Config.Loading = isLoading

}

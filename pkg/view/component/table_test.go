/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TableCols(t *testing.T) {
	cases := []struct {
		name     string
		in       []string
		expected []TableCol
	}{
		{
			name: "in general",
			in:   []string{"a"},
			expected: []TableCol{
				{Name: "a", Accessor: "a"},
			},
		},
		{
			name:     "empty list",
			in:       []string{},
			expected: []TableCol{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewTableCols(tc.in...)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_Table_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        Component
		expectedPath string
		isErr        bool
	}{
		{
			name: "general",
			input: &Table{
				base: newBase(typeTable, TitleFromString("my table")),
				Config: TableConfig{
					Filters: map[string]TableFilter{},
					Columns: []TableCol{
						{Name: "Name", Accessor: "Name"},
						{Name: "Description", Accessor: "Description"},
					},
					Rows: []TableRow{
						{
							"Name": &Text{
								Config: TextConfig{
									Text: "First",
								},
							},
							"Description": &Text{
								Config: TextConfig{
									Text: "The first row",
								},
							},
						},
						{
							"Name": &Text{
								Config: TextConfig{
									Text: "Last",
								},
							},
							"Description": &Text{
								Config: TextConfig{
									Text: "The last row",
								},
							},
						},
					},
					EmptyContent: "",
					Loading:      false,
				},
			},
			expectedPath: "table.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := err != nil
			if isErr != tc.isErr {
				t.Fatalf("Unexpected error: %v", err)
			}

			expected, err := ioutil.ReadFile(filepath.Join("testdata", tc.expectedPath))
			require.NoError(t, err, "reading test fixtures")
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

func Test_Table_isEmpty(t *testing.T) {
	cases := []struct {
		name    string
		table   *Table
		isEmpty bool
	}{
		{
			name:    "empty",
			table:   NewTable("my table", "placeholder", NewTableCols("col1")),
			isEmpty: true,
		},
		{
			name: "not empty",
			table: NewTableWithRows("my table", "placeholder", NewTableCols("col1"), []TableRow{
				{"col1": NewText("cell1")},
			}),
			isEmpty: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isEmpty, tc.table.IsEmpty())
		})
	}
}

func Test_Table_AddColumn(t *testing.T) {
	table := NewTable("table", "placeholder", NewTableCols("a"))
	table.AddColumn("b")
	expected := NewTableCols("a", "b")

	assert.Equal(t, expected, table.Columns())
}

func Test_Table_Sort(t *testing.T) {
	cases := []struct {
		name     string
		rows     []TableRow
		reverse  bool
		expected []TableRow
	}{
		{
			name: "asc sort",
			rows: []TableRow{
				{"a": NewText("2")},
				{"a": NewText("1")},
				{"a": NewText("3")},
			},
			reverse: false,
			expected: []TableRow{
				{"a": NewText("1")},
				{"a": NewText("2")},
				{"a": NewText("3")},
			},
		},
		{
			name: "desc sort",
			rows: []TableRow{
				{"a": NewText("2")},
				{"a": NewText("1")},
				{"a": NewText("3")},
			},
			reverse: true,
			expected: []TableRow{
				{"a": NewText("3")},
				{"a": NewText("2")},
				{"a": NewText("1")},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			table := NewTableWithRows("table", "placeholder", NewTableCols("a"), tc.rows)
			table.Sort("a", tc.reverse)
			expected := NewTableWithRows("table", "placeholder", NewTableCols("a"), tc.expected)

			assert.Equal(t, expected, table)
		})
	}
}

func TestTable_AddFilter(t *testing.T) {
	table := NewTable("table", "placeholder", NewTableCols("a"))
	filter := TableFilter{
		Values:   []string{"foo", "bar"},
		Selected: []string{"foo"},
	}
	table.AddFilter("a", filter)

	expected := map[string]TableFilter{"a": filter}

	assert.Equal(t, expected, table.Config.Filters)
}

package component

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Dropdown_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        Component
		expectedFile string
		isErr        bool
	}{
		{
			name: "standard dropdown",
			input: &Dropdown{
				Base: newBase(TypeDropdown, TitleFromString("dropdown")),
				Config: DropdownConfig{
					DropdownPosition: TopLeft,
					DropdownType:     DropdownButton,
					Action:           "action.octant.dev/dropdownTest",
					UseSelection:     false,
					Items: []DropdownItemConfig{
						{
							Name:        "first",
							Type:        PlainText,
							Label:       "First Item",
							Description: "This is the first item",
						},
						{
							Name:        "second",
							Type:        Url,
							Label:       "Second Item",
							Url:         "/items/second",
							Description: "This is the second item",
						},
					},
				},
			},
			expectedFile: "dropdown.json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := json.Marshal(test.input)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expected, err := ioutil.ReadFile(filepath.Join("testdata", test.expectedFile))
			require.NoError(t, err)

			assert.JSONEq(t, string(expected), string(got))

		})
	}
}

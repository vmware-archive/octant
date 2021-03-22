package component

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/action"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ButtonGroup_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        func() *ButtonGroup
		expectedFile string
		isErr        bool
	}{
		{
			name: "empty button group",
			input: func() *ButtonGroup {
				return NewButtonGroup()
			},
			expectedFile: "button_group_empty.json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := json.Marshal(test.input())
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

func TestButtonGroup_AddButton(t *testing.T) {
	bg := NewButtonGroup()
	button := NewButton("button", action.Payload{})
	bg.AddButton(button)
	expected := ButtonGroup{
		Base: newBase(TypeButtonGroup, nil),
		Config: ButtonGroupConfig{
			Buttons: []Button{
				button,
			},
		},
	}
	require.Equal(t, expected, *bg)
}

func Test_Button_Options(t *testing.T) {
	cases := []struct {
		name     string
		optsFunc []ButtonOption
		expected Button
	}{
		{
			name: "modal",
			optsFunc: []ButtonOption{
				WithModal(&Modal{
					Base:   newBase(TypeModal, TitleFromString("modal title")),
					Config: ModalConfig{},
				}),
			},
			expected: Button{
				Modal: NewModal(TitleFromString("modal title")),
			},
		},
		{
			name: "link",
			optsFunc: []ButtonOption{
				WithButtonLink("example.com"),
			},
			expected: Button{
				Ref: "example.com",
			},
		},
		{
			name: "status",
			optsFunc: []ButtonOption{
				WithButtonStatus(ButtonStatusDanger),
			},
			expected: Button{
				Status: ButtonStatusDanger,
			},
		},
		{
			name: "size",
			optsFunc: []ButtonOption{
				WithButtonSize(ButtonSizeLarge),
			},
			expected: Button{
				Size: ButtonSizeLarge,
			},
		},
		{
			name: "style",
			optsFunc: []ButtonOption{
				WithButtonStyle(ButtonStyleOutline),
			},
			expected: Button{
				Style: ButtonStyleOutline,
			},
		},
		{
			name: "multiple opts",
			optsFunc: []ButtonOption{
				WithButtonStatus(ButtonStatusSuccess),
				WithButtonSize(ButtonSizeBlock),
				WithButtonStyle(ButtonStyleFlat),
			},
			expected: Button{
				Status: ButtonStatusSuccess,
				Size:   ButtonSizeBlock,
				Style:  ButtonStyleFlat,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewButton("", nil, tc.optsFunc...)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestNewButtonLink(t *testing.T) {
	buttonLink := NewButtonLink("link button", "example.com")
	expected := Button{
		Name: "link button",
		Ref:  "example.com",
	}
	require.Equal(t, expected, buttonLink)
}

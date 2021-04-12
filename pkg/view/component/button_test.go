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

func Test_Button_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		input        func() *Button
		expectedFile string
		isErr        bool
	}{
		{
			name: "empty button",
			input: func() *Button {
				return NewButton("test", action.Payload{"foo": "bar"})
			},
			expectedFile: "button.json",
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

func Test_Button_Options(t *testing.T) {
	cases := []struct {
		name     string
		optsFunc []ButtonOption
		expected *Button
	}{
		{
			name: "modal",
			optsFunc: []ButtonOption{
				WithModal(&Modal{
					Base:   newBase(TypeModal, TitleFromString("modal title")),
					Config: ModalConfig{},
				}),
			},
			expected: &Button{
				Base: newBase(TypeButton, nil),
				Config: ButtonConfig{
					Modal: NewModal(TitleFromString("modal title")),
				},
			},
		},
		{
			name: "status",
			optsFunc: []ButtonOption{
				WithButtonStatus(ButtonStatusDanger),
			},
			expected: &Button{
				Base: newBase(TypeButton, nil),
				Config: ButtonConfig{
					Status: ButtonStatusDanger,
				},
			},
		},
		{
			name: "size",
			optsFunc: []ButtonOption{
				WithButtonSize(ButtonSizeMedium),
			},
			expected: &Button{
				Base: newBase(TypeButton, nil),
				Config: ButtonConfig{
					Size: ButtonSizeMedium,
				},
			},
		},
		{
			name: "style",
			optsFunc: []ButtonOption{
				WithButtonStyle(ButtonStyleOutline),
			},
			expected: &Button{
				Base: newBase(TypeButton, nil),
				Config: ButtonConfig{
					Style: ButtonStyleOutline,
				},
			},
		},
		{
			name: "multiple opts",
			optsFunc: []ButtonOption{
				WithButtonStatus(ButtonStatusSuccess),
				WithButtonSize(ButtonSizeBlock),
				WithButtonStyle(ButtonStyleFlat),
			},
			expected: &Button{
				Base: newBase(TypeButton, nil),
				Config: ButtonConfig{
					Status: ButtonStatusSuccess,
					Size:   ButtonSizeBlock,
					Style:  ButtonStyleFlat,
				},
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

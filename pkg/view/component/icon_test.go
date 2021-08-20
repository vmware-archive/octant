package component

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

func Test_Icon_Marshal(t *testing.T) {
	test := []struct {
		name         string
		input        *Icon
		expectedPath string
		isErr        bool
	}{
		{
			name: "in general",
			input: &Icon{
				Base: newBase(TypeIcon, nil),
				Config: IconConfig{
					Shape:      "user",
					Size:       "16",
					Direction:  DirectionDown,
					Flip:       FlipHorizontal,
					Solid:      true,
					Status:     StatusDanger,
					Inverse:    false,
					Badge:      BadgeDanger,
					Color:      "#add8e6",
					BadgeColor: "purple",
					Label:      "example icon",
				},
			},
			expectedPath: "icon.json",
			isErr:        false,
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := err != nil
			if isErr != tc.isErr {
				t.Fatalf("UnExpected error: %v", err)
			}

			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err)
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

func Test_Icon_With_SVG(t *testing.T) {
	test := []struct {
		name         string
		input        *Icon
		expectedPath string
		isErr        bool
	}{
		{
			name: "in general with svg",
			input: &Icon{
				Base: newBase(TypeIcon, nil),
				Config: IconConfig{
					Shape:      "user",
					Size:       "16",
					Direction:  DirectionDown,
					Flip:       FlipHorizontal,
					Solid:      true,
					Status:     StatusDanger,
					Inverse:    false,
					Badge:      BadgeDanger,
					Color:      "#add8e6",
					BadgeColor: "purple",
					Label:      "example icon",
					CustomSvg:  "<svg width=\"9px\" height=\"9px\" viewBox=\"0 0 9 9\"><g stroke-width=\"1\"><g transform=\"translate(-496.000000, -443.000000)\" stroke=\"#111111\"></g></g></svg>",
				},
			},
			expectedPath: "icon_with_svg.json",
			isErr:        false,
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := err != nil
			if isErr != tc.isErr {
				t.Fatalf("UnExpected error: %v", err)
			}

			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err)
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

func Test_Icon_Options(t *testing.T) {
	test := []struct {
		name     string
		optsFunc IconOption
		expected *Icon
	}{
		{
			name:     "Icon with tooltip",
			optsFunc: WithTooltip("hello", TooltipLarge, TooltipTopRight),
			expected: &Icon{
				Base: newBase(TypeIcon, nil),
				Config: IconConfig{
					Tooltip: &TooltipConfig{
						Message:  "hello",
						Size:     TooltipLarge,
						Position: TooltipTopRight,
					},
				},
			},
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			result := NewIcon("", tc.optsFunc)
			assert.Equal(t, tc.expected, result)
		})
	}
}

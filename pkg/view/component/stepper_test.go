package component

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_stepper_Marshal(t *testing.T) {
	form := Form{}
	form.Fields = append(form.Fields, NewFormFieldText("test", "test", "test"))

	tests := []struct {
		name         string
		input        Component
		expectedPath string
		isErr        bool
	}{
		{
			name: "general",
			input: &Stepper{
				base: newBase(typeStepper, TitleFromString("my stepper")),
				Config: StepperConfig{
					Action: "action.octant.dev/stepperTest",
					Steps: []StepConfig{
						{
							Name:        "Step 1",
							Title:       "First Step",
							Description: "Setup step",
							Form:        form,
						},
						{
							Name:        "Step 2",
							Title:       "Second Step",
							Description: "Confirmation step",
							Form:        form,
						},
					},
				},
			},
			expectedPath: "stepper.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			isErr := err != nil
			if isErr != tc.isErr {
				t.Fatalf("Unexpected error: %v", err)
			}

			expected, err := ioutil.ReadFile(path.Join("testdata", tc.expectedPath))
			require.NoError(t, err, "reading test fixtures")
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

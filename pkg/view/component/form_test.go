package component

import (
	"testing"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestForm_UnmarshalJSON(t *testing.T) {
	fieldNumber := NewFormFieldNumber("label", "name", "7")
	fieldNumber.AddValidator("Max value should be 7", map[FormValidator]interface{}{"required": "", "minLength": "7"})

	fieldText := NewFormFieldNumber("label", "name", "text")

	fl := NewFormFieldLayout(fieldNumber, fieldText)

	tests := []struct {
		name      string
		formField FormField
	}{
		{
			name:      "text field",
			formField: NewFormFieldText("label", "name", "value"),
		},
		{
			name: "check box field",
			formField: NewFormFieldCheckBox("label", "name", []InputChoice{
				{
					Label:   "foo",
					Value:   "foo",
					Checked: true,
				},
				{
					Label:   "bar",
					Value:   "bar",
					Checked: false,
				},
			}),
		},
		{
			name: "radio field",
			formField: NewFormFieldRadio("label", "name", []InputChoice{
				{
					Label:   "foo",
					Value:   "foo",
					Checked: true,
				},
			}),
		},
		{
			name: "select field",
			formField: NewFormFieldSelect("label", "name", []InputChoice{
				{
					Label: "baz",
					Value: "baz",
				},
			}, true),
		},
		{
			name:      "password field",
			formField: NewFormFieldPassword("label", "name", "value"),
		},
		{
			name:      "number field",
			formField: fieldNumber,
		},
		{
			name:      "text area field",
			formField: NewFormFieldTextarea("label", "name", "7"),
		},
		{
			name:      "hidden field",
			formField: NewFormFieldHidden("name", "7"),
		},
		{
			name:      "form layout",
			formField: fl,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := Form{
				Fields: []FormField{test.formField},
			}

			data, err := json.Marshal(&form)
			require.NoError(t, err)

			var got Form
			require.NoError(t, json.Unmarshal(data, &got))

			dataGot, err := json.Marshal(&got)
			assert.JSONEq(t, string(data), string(dataGot))

		})
	}
}

func TestCreateFormForObject(t *testing.T) {
	object := testutil.CreatePod("pod")
	got, err := CreateFormForObject("action", object,
		NewFormFieldNumber("number", "name", "0"))
	require.NoError(t, err)

	expected := Form{
		Fields: []FormField{
			NewFormFieldNumber("number", "name", "0"),
			NewFormFieldHidden("apiVersion", object.APIVersion),
			NewFormFieldHidden("kind", object.Kind),
			NewFormFieldHidden("name", object.Name),
			NewFormFieldHidden("namespace", object.Namespace),
			NewFormFieldHidden("action", "action"),
		},
	}
	require.Equal(t, expected, got)
}

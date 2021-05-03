package component

import (
	"fmt"
	"testing"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestFormFieldLayout_UnmarshalJSON(t *testing.T) {
	textFieldOne := NewFormFieldText("Test Horizontal", "Test Horizontal", "")
	textFieldTwo := NewFormFieldText("Cluster Name", "clusterName", "")
	expected := NewFormFieldLayout("", []FormField{textFieldOne, textFieldTwo})

	data, err := json.Marshal(&expected)
	require.NoError(t, err)

	fmt.Println(string(data))

	var got FormFieldLayout

	require.NoError(t, json.Unmarshal(data, &got))

	assertFormFieldEqual(t, expected, &got)
}

func TestFormFieldCheckBox_UnmarshalJSON(t *testing.T) {
	choices := []InputChoice{
		{Label: "foo", Value: "foo", Checked: false},
		{Label: "bar", Value: "bar", Checked: true},
		{Label: "baz", Value: "baz", Checked: false},
	}

	expected := NewFormFieldCheckBox("label", "name", choices)
	expected.AddValidator("error", map[FormValidator]interface{}{})

	data, err := json.Marshal(&expected)
	require.NoError(t, err)

	var got FormFieldCheckBox

	require.NoError(t, json.Unmarshal(data, &got))

	assertFormFieldEqual(t, expected, &got)
}

func TestFormFieldRadio_UnmarshalJSON(t *testing.T) {
	choices := []InputChoice{
		{Label: "foo", Value: "foo", Checked: false},
		{Label: "bar", Value: "bar", Checked: true},
		{Label: "baz", Value: "baz", Checked: false},
	}

	expected := NewFormFieldRadio("label", "name", choices)
	expected.AddValidator("error", map[FormValidator]interface{}{})

	data, err := json.Marshal(&expected)
	require.NoError(t, err)

	var got FormFieldRadio

	require.NoError(t, json.Unmarshal(data, &got))

	assertFormFieldEqual(t, expected, &got)
}

func TestFormFieldText_UnmarshalJSON(t *testing.T) {
	expected := NewFormFieldText("label", "name", "text")
	expected.AddValidator("placeholder", "error", map[FormValidator]interface{}{})

	data, err := json.Marshal(&expected)
	require.NoError(t, err)

	var got FormFieldText

	require.NoError(t, json.Unmarshal(data, &got))

	assertFormFieldEqual(t, expected, &got)
}

func TestFormFieldPassword_UnmarshalJSON(t *testing.T) {
	expected := NewFormFieldPassword("label", "name", "text")
	expected.AddValidator("placeholder", "error", map[FormValidator]interface{}{})

	data, err := json.Marshal(&expected)
	require.NoError(t, err)

	var got FormFieldPassword

	require.NoError(t, json.Unmarshal(data, &got))

	assertFormFieldEqual(t, expected, &got)
}

func TestFormFieldNumber_UnmarshalJSON(t *testing.T) {
	expected := NewFormFieldNumber("label", "name", "999")
	expected.AddValidator("error", map[FormValidator]interface{}{})

	data, err := json.Marshal(&expected)
	require.NoError(t, err)

	var got FormFieldNumber

	require.NoError(t, json.Unmarshal(data, &got))

	assertFormFieldEqual(t, expected, &got)
}

func TestFormFieldSelect_UnmarshalJSON(t *testing.T) {
	choices := []InputChoice{
		{Label: "foo", Value: "foo", Checked: false},
		{Label: "bar", Value: "bar", Checked: true},
		{Label: "baz", Value: "baz", Checked: false},
	}

	expected := NewFormFieldSelect("label", "name", choices, true)
	expected.AddValidator("error", map[FormValidator]interface{}{})

	data, err := json.Marshal(&expected)
	require.NoError(t, err)

	var got FormFieldSelect

	require.NoError(t, json.Unmarshal(data, &got))

	assertFormFieldEqual(t, expected, &got)
}

func TestFormFieldTextarea_UnmarshalJSON(t *testing.T) {
	expected := NewFormFieldTextarea("label", "name", "999")
	expected.AddValidator("placeholder", "error", map[FormValidator]interface{}{})

	data, err := json.Marshal(&expected)
	require.NoError(t, err)

	var got FormFieldTextarea

	require.NoError(t, json.Unmarshal(data, &got))

	assertFormFieldEqual(t, expected, &got)
}

func TestFormFieldHidden_UnmarshalJSON(t *testing.T) {
	expected := NewFormFieldHidden("label", "name")
	expected.AddValidator("placeholder", "error", map[FormValidator]interface{}{})

	data, err := json.Marshal(&expected)
	require.NoError(t, err)

	var got FormFieldHidden

	require.NoError(t, json.Unmarshal(data, &got))

	assertFormFieldEqual(t, expected, &got)
}

func TestForm_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		formField FormField
	}{
		{
			name:      "text field",
			formField: NewFormFieldText("label", "name", "value"),
		},
		{
			name:      "check box field",
			formField: NewFormFieldCheckBox("label", "name", []InputChoice{}),
		},
		{
			name:      "radio field",
			formField: NewFormFieldRadio("label", "name", []InputChoice{}),
		},
		{
			name:      "select field",
			formField: NewFormFieldSelect("label", "name", []InputChoice{}, true),
		},
		{
			name:      "password field",
			formField: NewFormFieldPassword("label", "name", "value"),
		},
		{
			name:      "number field",
			formField: NewFormFieldNumber("label", "name", "7"),
		},
		{
			name:      "text area field",
			formField: NewFormFieldTextarea("label", "name", "7"),
		},
		{
			name:      "hidden field",
			formField: NewFormFieldHidden("name", "7"),
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

			assert.Equal(t, form, got)

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

func assertFormFieldEqual(t *testing.T, expected, got FormField) {
	assert.Equal(t, expected.Value(), got.Value())
	assert.Equal(t, expected.Name(), got.Name())
	assert.Equal(t, expected.Type(), got.Type())
	assert.Equal(t, expected.Configuration(), got.Configuration())
}

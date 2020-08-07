/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	FieldTypeCheckBox = "checkbox"
	FieldTypeRadio    = "radio"
	FieldTypeText     = "text"
	FieldTypePassword = "password"
	FieldTypeNumber   = "number"
	FieldTypeSelect   = "select"
	FieldTypeTextarea = "textarea"
	FieldTypeHidden   = "hidden"
)

type InputChoice struct {
	Label   string `json:"label"`
	Value   string `json:"value"`
	Checked bool   `json:"checked"`
}

type BaseFormField struct {
	label     string
	name      string
	fieldType string
}

func newBaseFormField(label, name, fieldType string) *BaseFormField {
	return &BaseFormField{
		label:     label,
		name:      name,
		fieldType: fieldType,
	}
}

func (bff *BaseFormField) Label() string {
	return bff.label
}

func (bff *BaseFormField) Name() string {
	return bff.name
}

func (bff *BaseFormField) Type() string {
	return bff.fieldType
}

// FormField is a form field interface.
// TODO: make this more json friendly by converting it to a struct.
type FormField interface {
	Label() string
	Name() string
	Type() string
	Configuration() map[string]interface{}
	Value() interface{}

	json.Unmarshaler
	json.Marshaler
}

// marshalFormField marshals a form field to JSON.
func marshalFormField(ff FormField) ([]byte, error) {
	m := map[string]interface{}{
		"label":         ff.Label(),
		"name":          ff.Name(),
		"type":          ff.Type(),
		"configuration": ff.Configuration(),
		"value":         ff.Value(),
	}

	return json.Marshal(&m)
}

type FormFieldCheckBox struct {
	*BaseFormField

	choices []InputChoice
}

func (ff *FormFieldCheckBox) MarshalJSON() ([]byte, error) {
	return marshalFormField(ff)
}

var _ FormField = (*FormFieldCheckBox)(nil)

func NewFormFieldCheckBox(label, name string, choices []InputChoice) *FormFieldCheckBox {
	return &FormFieldCheckBox{
		BaseFormField: newBaseFormField(label, name, FieldTypeCheckBox),
		choices:       choices,
	}
}

func (ff *FormFieldCheckBox) Configuration() map[string]interface{} {
	return map[string]interface{}{
		"choices": ff.choices,
	}
}

func (ff *FormFieldCheckBox) Value() interface{} {
	var selected []string
	for _, choice := range ff.choices {
		if choice.Checked {
			selected = append(selected, choice.Value)
		}
	}

	return selected
}

func (ff *FormFieldCheckBox) UnmarshalJSON(data []byte) error {
	x := struct {
		Label         string `json:"label"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		Configuration struct {
			Choices []InputChoice
		} `json:"configuration"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	ff.BaseFormField = newBaseFormField(x.Label, x.Name, x.Type)
	ff.choices = x.Configuration.Choices

	return nil
}

type FormFieldRadio struct {
	*BaseFormField

	choices []InputChoice
}

func NewFormFieldRadio(label, name string, choices []InputChoice) *FormFieldRadio {
	return &FormFieldRadio{
		BaseFormField: newBaseFormField(label, name, FieldTypeRadio),
		choices:       choices,
	}
}

var _ FormField = (*FormFieldRadio)(nil)

func (ff *FormFieldRadio) Configuration() map[string]interface{} {
	return map[string]interface{}{
		"choices": ff.choices,
	}
}

func (ff *FormFieldRadio) Value() interface{} {
	var selected []string
	for _, choice := range ff.choices {
		if choice.Checked {
			selected = append(selected, choice.Value)
		}
	}

	value := ""
	if len(selected) > 0 {
		value = selected[0]
	}

	return value
}

func (ff *FormFieldRadio) UnmarshalJSON(data []byte) error {
	x := struct {
		Label         string `json:"label"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		Configuration struct {
			Choices []InputChoice
		} `json:"configuration"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	ff.BaseFormField = newBaseFormField(x.Label, x.Name, x.Type)
	ff.choices = x.Configuration.Choices

	return nil
}

func (ff *FormFieldRadio) MarshalJSON() ([]byte, error) {
	return marshalFormField(ff)
}

type FormFieldText struct {
	*BaseFormField

	value string
}

func NewFormFieldText(label, name, value string) *FormFieldText {
	return &FormFieldText{
		BaseFormField: newBaseFormField(label, name, FieldTypeText),
		value:         value,
	}
}

var _ FormField = (*FormFieldText)(nil)

func (ff *FormFieldText) Configuration() map[string]interface{} {
	return map[string]interface{}{}
}

func (ff *FormFieldText) Value() interface{} {
	return ff.value
}

func (ff *FormFieldText) MarshalJSON() ([]byte, error) {
	return marshalFormField(ff)
}

func (ff *FormFieldText) UnmarshalJSON(data []byte) error {
	x := struct {
		Label         string                 `json:"label"`
		Name          string                 `json:"name"`
		Type          string                 `json:"type"`
		Configuration map[string]interface{} `json:"configuration"`
		Value         string                 `json:"value"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	ff.BaseFormField = newBaseFormField(x.Label, x.Name, x.Type)
	ff.value = x.Value

	return nil

}

type FormFieldPassword struct {
	*BaseFormField

	value string
}

func NewFormFieldPassword(label, name, value string) *FormFieldPassword {
	return &FormFieldPassword{
		BaseFormField: newBaseFormField(label, name, FieldTypePassword),
		value:         value,
	}
}

var _ FormField = (*FormFieldPassword)(nil)

func (ff *FormFieldPassword) Configuration() map[string]interface{} {
	return map[string]interface{}{}
}

func (ff *FormFieldPassword) Value() interface{} {
	return ff.value
}

func (ff *FormFieldPassword) MarshalJSON() ([]byte, error) {
	return marshalFormField(ff)
}

func (ff *FormFieldPassword) UnmarshalJSON(data []byte) error {
	x := struct {
		Label         string                 `json:"label"`
		Name          string                 `json:"name"`
		Type          string                 `json:"type"`
		Configuration map[string]interface{} `json:"configuration"`
		Value         string                 `json:"value"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	ff.BaseFormField = newBaseFormField(x.Label, x.Name, x.Type)
	ff.value = x.Value

	return nil
}

type FormFieldNumber struct {
	*BaseFormField

	value string
}

func NewFormFieldNumber(label, name, value string) *FormFieldNumber {
	return &FormFieldNumber{
		BaseFormField: newBaseFormField(label, name, FieldTypeNumber),
		value:         value,
	}
}

var _ FormField = (*FormFieldNumber)(nil)

func (ff *FormFieldNumber) Configuration() map[string]interface{} {
	return map[string]interface{}{}
}

func (ff *FormFieldNumber) Value() interface{} {
	return ff.value
}

func (ff *FormFieldNumber) MarshalJSON() ([]byte, error) {
	return marshalFormField(ff)
}

func (ff *FormFieldNumber) UnmarshalJSON(data []byte) error {
	x := struct {
		Label         string                 `json:"label"`
		Name          string                 `json:"name"`
		Type          string                 `json:"type"`
		Configuration map[string]interface{} `json:"configuration"`
		Value         string                 `json:"value"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	ff.BaseFormField = newBaseFormField(x.Label, x.Name, x.Type)
	ff.value = x.Value

	return nil
}

type FormFieldSelect struct {
	*BaseFormField

	choices  []InputChoice
	multiple bool
}

func NewFormFieldSelect(label, name string, choices []InputChoice, multiple bool) *FormFieldSelect {
	return &FormFieldSelect{
		BaseFormField: newBaseFormField(label, name, FieldTypeSelect),
		choices:       choices,
		multiple:      multiple,
	}
}

var _ FormField = (*FormFieldSelect)(nil)

func (ff *FormFieldSelect) Configuration() map[string]interface{} {
	var value []string
	for _, choice := range ff.choices {
		if choice.Checked {
			value = append(value, choice.Value)
		}
	}

	return map[string]interface{}{
		"choices":  ff.choices,
		"multiple": ff.multiple,
		"value":    value,
	}
}

func (ff *FormFieldSelect) Value() interface{} {
	var value []string
	for _, choice := range ff.choices {
		if choice.Checked {
			value = append(value, choice.Value)
		}
	}

	return value
}

func (ff *FormFieldSelect) MarshalJSON() ([]byte, error) {
	return marshalFormField(ff)
}

func (ff *FormFieldSelect) UnmarshalJSON(data []byte) error {
	x := struct {
		Label         string `json:"label"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		Configuration struct {
			Choices  []InputChoice `json:"choices"`
			Multiple bool          `json:"multiple"`
		} `json:"configuration"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	ff.BaseFormField = newBaseFormField(x.Label, x.Name, x.Type)
	ff.choices = x.Configuration.Choices
	ff.multiple = x.Configuration.Multiple

	return nil
}

type FormFieldTextarea struct {
	*BaseFormField

	value string
}

func NewFormFieldTextarea(label, name, value string) *FormFieldTextarea {
	return &FormFieldTextarea{
		BaseFormField: newBaseFormField(label, name, FieldTypeTextarea),
		value:         value,
	}
}

var _ FormField = (*FormFieldTextarea)(nil)

func (ff *FormFieldTextarea) Configuration() map[string]interface{} {
	return map[string]interface{}{}
}

func (ff *FormFieldTextarea) Value() interface{} {
	return ff.value
}

func (ff *FormFieldTextarea) MarshalJSON() ([]byte, error) {
	return marshalFormField(ff)
}

func (ff *FormFieldTextarea) UnmarshalJSON(data []byte) error {
	x := struct {
		Label         string                 `json:"label"`
		Name          string                 `json:"name"`
		Type          string                 `json:"type"`
		Configuration map[string]interface{} `json:"configuration"`
		Value         string                 `json:"value"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	ff.BaseFormField = newBaseFormField(x.Label, x.Name, x.Type)
	ff.value = x.Value

	return nil
}

func NewFormFieldHidden(name, value string) *FormFieldHidden {
	return &FormFieldHidden{
		BaseFormField: newBaseFormField("", name, FieldTypeHidden),
		value:         value,
	}
}

type FormFieldHidden struct {
	*BaseFormField

	value string
}

var _ FormField = (*FormFieldHidden)(nil)

func (ff *FormFieldHidden) Configuration() map[string]interface{} {
	return map[string]interface{}{}
}

func (ff *FormFieldHidden) Value() interface{} {
	return ff.value
}

func (ff *FormFieldHidden) MarshalJSON() ([]byte, error) {
	return marshalFormField(ff)
}

func (ff *FormFieldHidden) UnmarshalJSON(data []byte) error {
	x := struct {
		Label         string                 `json:"label"`
		Name          string                 `json:"name"`
		Type          string                 `json:"type"`
		Configuration map[string]interface{} `json:"configuration"`
		Value         string                 `json:"value"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	ff.BaseFormField = newBaseFormField(x.Label, x.Name, x.Type)
	ff.value = x.Value

	return nil
}

type Form struct {
	Fields []FormField `json:"fields"`
}

func (f *Form) MarshalJSON() ([]byte, error) {
	t := struct {
		Fields []map[string]interface{} `json:"fields"`
	}{}

	for _, field := range f.Fields {
		m := map[string]interface{}{
			"label":         field.Label(),
			"name":          field.Name(),
			"type":          field.Type(),
			"configuration": field.Configuration(),
			"value":         field.Value(),
		}

		t.Fields = append(t.Fields, m)
	}

	return json.Marshal(t)
}

func (f *Form) UnmarshalJSON(data []byte) error {
	x := struct {
		Fields []struct {
			Label         string                 `json:"label"`
			Name          string                 `json:"name"`
			Type          string                 `json:"type"`
			Configuration map[string]interface{} `json:"configuration"`
			Value         interface{}            `json:"value"`
		} `json:"fields"`
	}{}

	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}

	for i := range x.Fields {
		field := x.Fields[i]
		var ff FormField

		fieldData, err := json.Marshal(field)
		if err != nil {
			return err
		}

		switch field.Type {
		case FieldTypeCheckBox:
			ff = &FormFieldCheckBox{}
		case FieldTypeRadio:
			ff = &FormFieldRadio{}
		case FieldTypeText:
			ff = &FormFieldText{}
		case FieldTypePassword:
			ff = &FormFieldPassword{}
		case FieldTypeNumber:
			ff = &FormFieldNumber{}
		case FieldTypeSelect:
			ff = &FormFieldSelect{}
		case FieldTypeTextarea:
			ff = &FormFieldTextarea{}
		case FieldTypeHidden:
			ff = &FormFieldHidden{}
		default:
			return errors.Errorf("unknown form field type %q", field)
		}

		if err := ff.UnmarshalJSON(fieldData); err != nil {
			return err
		}

		f.Fields = append(f.Fields, ff)
	}

	return nil
}

// CreateFormForObject creates a form for an object with additional fields.
func CreateFormForObject(actionName string, object runtime.Object, fields ...FormField) (Form, error) {
	if object == nil {
		return Form{}, errors.New("object is nil")
	}

	apiVersion, kind := object.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
	accessor, err := meta.Accessor(object)
	if err != nil {
		return Form{}, err
	}

	fields = append(fields,
		NewFormFieldHidden("apiVersion", apiVersion),
		NewFormFieldHidden("kind", kind),
		NewFormFieldHidden("name", accessor.GetName()),
		NewFormFieldHidden("namespace", accessor.GetNamespace()),
		NewFormFieldHidden("action", actionName),
	)

	return Form{Fields: fields}, nil
}

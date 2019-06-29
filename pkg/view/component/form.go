/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
)

type InputChoice struct {
	Label   string `json:"label"`
	Value   string `json:"value"`
	Checked bool   `json:"checked"`
}

type baseFormField struct {
	label     string
	name      string
	fieldType string
}

func newBaseFormField(label, name, fieldType string) *baseFormField {
	return &baseFormField{
		label:     label,
		name:      name,
		fieldType: fieldType,
	}
}

func (bff *baseFormField) Label() string {
	return bff.label
}

func (bff *baseFormField) Name() string {
	return bff.name
}

func (bff *baseFormField) Type() string {
	return bff.fieldType
}

type FormField interface {
	Label() string
	Name() string
	Type() string
	Configuration() map[string]interface{}
	Value() interface{}
}

type FormFieldCheckBox struct {
	*baseFormField

	choices []InputChoice
}

func NewFormFieldCheckBox(label, name string, choices []InputChoice) *FormFieldCheckBox {
	return &FormFieldCheckBox{
		baseFormField: newBaseFormField(label, name, "checkbox"),
		choices:       choices,
	}
}

var _ FormField = (*FormFieldCheckBox)(nil)

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

type FormFieldRadio struct {
	*baseFormField

	choices []InputChoice
}

func NewFormFieldRadio(label, name string, choices []InputChoice) *FormFieldRadio {
	return &FormFieldRadio{
		baseFormField: newBaseFormField(label, name, "radio"),
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

type FormFieldText struct {
	*baseFormField

	value string
}

func NewFormFieldText(label, name, value string) *FormFieldText {
	return &FormFieldText{
		baseFormField: newBaseFormField(label, name, "text"),
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

type FormFieldPassword struct {
	*baseFormField

	value string
}

func NewFormFieldFieldPassword(label, name, value string) *FormFieldPassword {
	return &FormFieldPassword{
		baseFormField: newBaseFormField(label, name, "password"),
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

type FormFieldNumber struct {
	*baseFormField

	value string
}

func NewFormFieldNumber(label, name, value string) *FormFieldNumber {
	return &FormFieldNumber{
		baseFormField: newBaseFormField(label, name, "number"),
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

type FormFieldSelect struct {
	*baseFormField

	choices  []InputChoice
	multiple bool
}

func NewFormFieldSelect(label, name string, choices []InputChoice, multiple bool) *FormFieldSelect {
	return &FormFieldSelect{
		baseFormField: newBaseFormField(label, name, "select"),
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

type FormFieldTextarea struct {
	*baseFormField

	value string
}

func NewFormFieldTextarea(label, name, value string) *FormFieldTextarea {
	return &FormFieldTextarea{
		baseFormField: newBaseFormField(label, name, "textarea"),
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

func NewFormFieldHidden(name, value string) *FormFieldTextarea {
	return &FormFieldTextarea{
		baseFormField: newBaseFormField("", name, "hidden"),
		value:         value,
	}
}

type FormFieldHidden struct {
	*baseFormField

	value string
}

var _ FormField = (*FormFieldHidden)(nil)

func (ff *FormFieldHidden) Configuration() map[string]interface{} {
	return map[string]interface{}{}
}

func (ff *FormFieldHidden) Value() interface{} {
	return ff.value
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

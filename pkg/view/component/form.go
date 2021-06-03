/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"fmt"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

type FormFieldType string

const (
	FieldTypeCheckBox FormFieldType = "checkbox"
	FieldTypeRadio    FormFieldType = "radio"
	FieldTypeText     FormFieldType = "text"
	FieldTypePassword FormFieldType = "password"
	FieldTypeNumber   FormFieldType = "number"
	FieldTypeSelect   FormFieldType = "select"
	FieldTypeTextarea FormFieldType = "textarea"
	FieldTypeHidden   FormFieldType = "hidden"
	FieldTypeLayout   FormFieldType = "layout"
)

type FormValidator string

const (
	FormValidatorMin           FormValidator = "min"
	FormValidatorMax           FormValidator = "max"
	FormValidatorRequired      FormValidator = "required"
	FormValidatorRequiredTrue  FormValidator = "requiredTrue"
	FormValidatorEmail         FormValidator = "email"
	FormValidatorMinLength     FormValidator = "minLength"
	FormValidatorMaxLength     FormValidator = "maxLength"
	FormValidatorPattern       FormValidator = "pattern"
	FormValidatorNullValidator FormValidator = "nullValidator"
)

type InputChoice struct {
	Label   string `json:"label"`
	Value   string `json:"value"`
	Checked bool   `json:"checked"`
}

// FormFieldOptions provides additional configuration for form fields with multiple inputs
type FormFieldOptions struct {
	Choices  []InputChoice `json:"choices,omitempty"`
	Multiple bool          `json:"multiple,omitempty"`
	Fields   []FormField   `json:"fields,omitempty"`
}

// FormField is a component for fields within a Form
// +octant:component
type FormField struct {
	Base
	Config FormFieldConfig `json:"config"`
}

var _ Component = (*FormField)(nil)

type formFieldMarshal FormField

type FormFieldConfig struct {
	Type          FormFieldType                 `json:"type"`
	Label         string                        `json:"label"`
	Name          string                        `json:"name"`
	Value         interface{}                   `json:"value"`
	Configuration *FormFieldOptions             `json:"configuration,omitempty"`
	Placeholder   string                        `json:"placeholder,omitempty"`
	Error         string                        `json:"error,omitempty"`
	Validators    map[FormValidator]interface{} `json:"validators,omitempty"`
	Width         int                           `json:"width,omitempty"`
}

func (ffc *FormFieldConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Type          FormFieldType `json:"type"`
		Label         string        `json:"label"`
		Name          string        `json:"name"`
		Value         interface{}   `json:"value"`
		Configuration struct {
			Fields   []TypedObject `json:"fields"`
			Choices  []InputChoice `json:"choices,omitempty"`
			Multiple bool          `json:"multiple,omitempty"`
		} `json:"configuration,omitempty"`
		Placeholder string                        `json:"placeholder,omitempty"`
		Error       string                        `json:"error,omitempty"`
		Validators  map[FormValidator]interface{} `json:"validators,omitempty"`
		Width       int                           `json:"width,omitempty"`
	}{}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	ffc.Type = x.Type
	ffc.Label = x.Label
	ffc.Name = x.Name
	ffc.Value = x.Value
	ffc.Width = x.Width
	if &x.Configuration == nil || x.Configuration.Fields != nil || x.Configuration.Choices != nil {
		ffc.Configuration = &FormFieldOptions{}
		ffc.Configuration.Choices = x.Configuration.Choices
		ffc.Configuration.Multiple = x.Configuration.Multiple
	}

	if x.Configuration.Fields != nil {
		for _, typedObject := range x.Configuration.Fields {
			component, err := typedObject.ToComponent()
			if err != nil {
				return err
			}

			field, ok := component.(*FormField)
			if !ok {
				return fmt.Errorf("item was not a form field")
			}

			ffc.Configuration.Fields = append(ffc.Configuration.Fields, *field)
		}
	}
	ffc.Placeholder = x.Placeholder
	ffc.Error = x.Error
	ffc.Validators = x.Validators
	return nil
}

func (ff *FormField) MarshalJSON() ([]byte, error) {
	m := formFieldMarshal(*ff)
	m.Metadata.Type = TypeFormField
	return json.Marshal(&m)
}

// NewFormFieldLayout creates a group of form fields
func NewFormFieldLayout(fields ...FormField) FormField {
	return FormField{
		Base: newBase(TypeFormField, nil),
		Config: FormFieldConfig{
			Type: FieldTypeLayout,
			Configuration: &FormFieldOptions{
				Fields: fields,
			},
		},
	}
}

func NewFormFieldCheckBox(label, name string, choices []InputChoice) FormField {
	return FormField{
		Base: newBase(TypeFormField, nil),
		Config: FormFieldConfig{
			Type:  FieldTypeCheckBox,
			Label: label,
			Name:  name,
			Configuration: &FormFieldOptions{
				Choices:  choices,
				Multiple: true,
			},
		},
	}
}

// AddValidator adds validator(s)
func (ff *FormField) AddValidator(errorMessage string, validators map[FormValidator]interface{}) {
	ff.Config.Error = errorMessage
	ff.Config.Validators = validators
}

func (ff *FormField) SetWidth(width int) {
	ff.Config.Width = width
}

func (ff *FormField) SetPlaceHolder(placeholder string) {
	ff.Config.Placeholder = placeholder
}

func (ff *FormField) SetLabel(label string) {
	ff.Config.Label = label
}

func NewFormFieldRadio(label, name string, choices []InputChoice) FormField {
	return FormField{
		Base: newBase(TypeFormField, nil),
		Config: FormFieldConfig{
			Type:  FieldTypeRadio,
			Label: label,
			Name:  name,
			Configuration: &FormFieldOptions{
				Choices:  choices,
				Multiple: false,
			},
		},
	}
}

func NewFormFieldText(label, name, value string) FormField {
	return FormField{
		Base: newBase(TypeFormField, nil),
		Config: FormFieldConfig{
			Type:  FieldTypeText,
			Label: label,
			Name:  name,
			Value: value,
		},
	}
}

func NewFormFieldPassword(label, name, value string) FormField {
	return FormField{
		Base: newBase(TypeFormField, nil),
		Config: FormFieldConfig{
			Type:  FieldTypePassword,
			Label: label,
			Name:  name,
			Value: value,
		},
	}
}

func NewFormFieldNumber(label, name, value string) FormField {
	return FormField{
		Base: newBase(TypeFormField, nil),
		Config: FormFieldConfig{
			Type:  FieldTypeNumber,
			Label: label,
			Name:  name,
			Value: value,
		},
	}
}

func NewFormFieldSelect(label, name string, choices []InputChoice, multiple bool) FormField {
	return FormField{
		Base: newBase(TypeFormField, nil),
		Config: FormFieldConfig{
			Type:  FieldTypeSelect,
			Label: label,
			Name:  name,
			Configuration: &FormFieldOptions{
				Choices:  choices,
				Multiple: multiple,
			},
		},
	}
}

func NewFormFieldTextarea(label, name, value string) FormField {
	return FormField{
		Base: newBase(TypeFormField, nil),
		Config: FormFieldConfig{
			Type:  FieldTypeTextarea,
			Label: label,
			Name:  name,
			Value: value,
		},
	}
}

func NewFormFieldHidden(name, value string) FormField {
	return FormField{
		Base: newBase(TypeFormField, nil),
		Config: FormFieldConfig{
			Type:  FieldTypeHidden,
			Name:  name,
			Value: value,
		},
	}
}

type Form struct {
	Fields []FormField `json:"fields"`
	Action string      `json:"action,omitempty"`
}

type formMarshal Form

func (f *Form) MarshalJSON() ([]byte, error) {
	m := formMarshal{
		Fields: f.Fields,
		Action: f.Action,
	}
	return json.Marshal(&m)
}

func (f *Form) UnmarshalJSON(data []byte) error {
	x := struct {
		Fields []TypedObject `json:"fields"`
		Action string        `json:"action,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	for _, typedObject := range x.Fields {
		component, err := typedObject.ToComponent()
		if err != nil {
			return err
		}

		field, ok := component.(*FormField)
		if !ok {
			return fmt.Errorf("item was not a form field")
		}

		f.Fields = append(f.Fields, *field)
	}
	f.Action = x.Action
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

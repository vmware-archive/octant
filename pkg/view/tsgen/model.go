/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package tsgen

import (
	"fmt"
	"reflect"
	"sort"
	gostrings "strings"

	"github.com/stoewer/go-strcase"

	"github.com/vmware-tanzu/octant/internal/util/strings"
)

// ConvertField converts struct field a to a tsgen Field.  If the
// field cannot be converted to a type script type, it will panic.
// If the field should be skipped, it will return false the second
// return value.
func ConvertField(xType reflect.Type, i int, componentNames []string) (Field, bool) {
	field := xType.Field(i)
	t := getTag(field)
	if t.skip {
		return Field{}, false
	}

	typ, refNames, err := tsType(field, componentNames)
	if err != nil {
		panic(fmt.Sprintf("unable to get type for %s: %v", field.Name, err))
	}

	f := Field{
		Name:     t.name,
		Optional: t.omitempty,
		Type:     typ,
		RefNames: refNames,
	}

	return f, true
}

// Model is a model contains information about components and their names.
type Model struct {
	Components     []Component
	ComponentNames []string
}

// Component is a component in the model.
type Component struct {
	Name   string
	TSName string
	Fields []Field
}

// ImportReference is a reference for an import that is used to build typescript imports.
type ImportReference struct {
	// Name is the name of the component.
	Name string
	// Import name is the name of the file that will be imported.
	ImportName string
}

// Referenced returns import references for this Component.
func (c Component) Referenced() []ImportReference {
	m := map[string]bool{}

	for _, f := range c.Fields {
		for _, n := range f.RefNames {
			m[n] = true
		}
	}

	var list []string
	for k := range m {
		list = append(list, k)
	}

	sort.Strings(list)

	var refs []ImportReference
	for _, r := range list {
		if c.Name != r { // Prevent nested/recursive component imports
			refs = append(refs, ImportReference{
				Name:       r,
				ImportName: strcase.KebabCase(r),
			})
		}
	}

	return refs
}

// Field is a struct field in a Component.
type Field struct {
	Name     string
	Type     string
	Optional bool
	RefNames []string
}

// TSType returns the typescript type for a field.
func (f Field) TSType() string {
	switch f.Type {
	case "component":
		return "Component<any>"
	default:
		return f.Type
	}
}

// TSFactoryType returns the typescript factory type for a field.
func (f Field) TSFactoryType() string {
	switch f.Type {
	case "component":
		return "ComponentFactory<any>"
	default:
		return f.Type
	}
}

// TSNameToComponent converts a field to a component.
func (f Field) TSNameToComponent() string {
	switch f.Type {
	case "component":
		return f.Name + ".toComponent()"
	default:
		return f.Name
	}
}

type tag struct {
	name      string
	omitempty bool
	skip      bool
}

func getTag(f reflect.StructField) tag {
	raw := f.Tag.Get("json")
	parts := gostrings.Split(raw, ",")

	if len(parts) == 0 {
		panic("unable to parse json struct tag")
	} else if len(parts) == 1 {
		if parts[0] == "-" {
			return tag{
				skip: true,
			}
		}

		return tag{
			name:      parts[0],
			omitempty: false,
		}
	}

	return tag{
		name:      parts[0],
		omitempty: strings.Contains("omitempty", parts),
	}
}

func tsType(f reflect.StructField, componentNames []string) (string, []string, error) {
	c := NewConverter(componentNames)
	return c.Convert(f.Type)
}

/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package tsgen

import (
	"fmt"
	"reflect"
	gostrings "strings"

	"github.com/vmware-tanzu/octant/internal/util/strings"
)

// Converter converts types to typescript.
type Converter struct {
	componentNames []string
}

// NewConverter creates an instance of Converter.
func NewConverter(componentNames []string) *Converter {
	c := &Converter{
		componentNames: componentNames,
	}
	return c
}

// Convert converts a type to Typescript. It returns the typescript, a list of
// components that were reference.
func (c *Converter) Convert(t reflect.Type) (string, []string, error) {
	var sb gostrings.Builder
	names, err := c.visit(&sb, t, 0)
	if err != nil {
		return "", names, err
	}

	return sb.String(), names, nil
}

// visit visits a type and then searches for nested types. It builds up typescript in the
// supplied string builder and returns components that it finds.
func (c *Converter) visit(w *gostrings.Builder, t reflect.Type, depth int) ([]string, error) {
	var componentNames []string

	switch t.Kind() {
	case reflect.Ptr:
		// if a pointer is found, visit the element it points to.
		names, err := c.visit(w, t.Elem(), depth+1)
		if err != nil {
			return nil, err
		}
		componentNames = append(componentNames, names...)
	case reflect.Map:
		var keyString string = t.Key().String()
		switch keyString {
		// typing for validation steps causes errors with reflection, so we just lock string keys.
		case "component.FormValidator":
			keyString = "string"
		}

		// if a map is found, key the key and element to build a typescript object definition.
		w.WriteString("{[key:" + keyString + "]:")

		names, err := c.visit(w, t.Elem(), depth+1)
		if err != nil {
			return nil, err
		}
		w.WriteString("}")
		componentNames = append(componentNames, names...)
	case reflect.Slice:
		// if a slice is found, visit the element type.
		names, err := c.visit(w, t.Elem(), depth+1)
		if err != nil {
			return nil, err
		}
		w.WriteString("[]")
		componentNames = append(componentNames, names...)
	case reflect.Struct:
		// if a struct is found, visit each of the fields.
		if depth > 0 && strings.Contains(t.String(), c.componentNames) {
			// if this is a component, specify the type and don't try to visit it.
			w.WriteString(fmt.Sprintf("Component<%sConfig>", t.Name()))
			return []string{t.Name()}, nil
		}

		// hard code these to ensure we don't try to visit them.
		switch t.String() {
		case "time.Time":
			// time will distill to a number, so don't visit it.
			w.WriteString("number")
			return nil, nil
		}

		w.WriteString("{\n")
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			// look in the json tag to get the name of the field.
			jsonTag := f.Tag.Get("json")
			jsonTagParts := gostrings.Split(jsonTag, ",")

			separator := ":"
			if strings.Contains("omitempty", jsonTagParts[1:]) {
				// if the json tag contains omitempty, then make this field optional.
				separator = "?:"
			}

			w.WriteString(gostrings.Repeat("  ", depth+1) + jsonTagParts[0] + separator + " ")

			names, err := c.visit(w, f.Type, depth)
			if err != nil {
				return nil, err
			}
			componentNames = append(componentNames, names...)
			w.WriteString(";\n")
		}

		w.WriteString(fmt.Sprintf("%s}", gostrings.Repeat("  ", depth)))
	case reflect.String:
		w.WriteString("string")
	case reflect.Bool:
		w.WriteString("boolean")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		w.WriteString("number")
	case reflect.Interface:
		switch t.String() {
		case "component.Component":
			// If we detect an Octant component, hard code a Component<any>.
			w.WriteString("Component<any>")
		default:
			w.WriteString("any")
		}
	default:
		return nil, fmt.Errorf("unable to handle %s", t.String())
	}

	return componentNames, nil
}

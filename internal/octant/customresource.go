/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	crdAPIVersionV1beta1 = "apiextensions.k8s.io/v1beta1"
	crdAPIVersionV1      = "apiextensions.k8s.io/v1"
)

type CustomResourceDefinitionPrinterColumn struct {
	Name        string
	Type        string
	Description string
	JSONPath    string
}

type CustomResourceDefinitionVersion struct {
	Version        string
	PrinterColumns []CustomResourceDefinitionPrinterColumn
}

type CustomResourceDefinition struct {
	object *unstructured.Unstructured
}

func NewCustomResourceDefinition(object *unstructured.Unstructured) (*CustomResourceDefinition, error) {
	if object == nil {
		return nil, fmt.Errorf("object is nil")
	}

	crd := &CustomResourceDefinition{
		object: object,
	}

	return crd, nil
}

func (crd *CustomResourceDefinition) Versions() ([]string, error) {
	switch apiVersion := crd.object.GetAPIVersion(); apiVersion {
	case crdAPIVersionV1:
		return crd.versionNames()
	case crdAPIVersionV1beta1:
		// if .spec.version exists, return that version (old format)
		version, found, err := unstructured.NestedString(crd.object.Object, "spec", "version")
		if err != nil {
			return nil, fmt.Errorf("unable to read crd .spec.version: %w", err)
		}

		if found {
			return []string{version}, nil
		}

		return crd.versionNames()
	default:
		return nil, fmt.Errorf("crd with API version '%s' is not supported", apiVersion)
	}
}

func (crd *CustomResourceDefinition) Version(version string) (CustomResourceDefinitionVersion, error) {
	switch crd.object.GetAPIVersion() {
	case crdAPIVersionV1:
		return crd.v1Version(version)
	case crdAPIVersionV1beta1:
		return crd.v1beta1Version(version)
	default:
		return CustomResourceDefinitionVersion{}, fmt.Errorf("crd with API version '%s' is not supported", version)
	}

}

func (crd *CustomResourceDefinition) v1Version(version string) (CustomResourceDefinitionVersion, error) {
	versions, err := crd.versions()
	if err != nil {
		return CustomResourceDefinitionVersion{}, err
	}
	for i := range versions {
		name, ok := versions[i]["name"].(string)
		if !ok {
			return CustomResourceDefinitionVersion{}, fmt.Errorf("unable to find CRD with version '%s'", version)
		}

		if name != version {
			continue
		}

		columns, err := crdV1PrinterColumns(versions[i]["additionalPrinterColumns"])
		if err != nil {
			return CustomResourceDefinitionVersion{}, fmt.Errorf("collect CRD printer columns: %w", err)
		}

		customResourceDefinitionVersion := CustomResourceDefinitionVersion{
			Version:        name,
			PrinterColumns: columns,
		}
		return customResourceDefinitionVersion, nil
	}

	return CustomResourceDefinitionVersion{}, fmt.Errorf("unable to find version '%s'", version)
}

func (crd *CustomResourceDefinition) v1beta1Version(version string) (CustomResourceDefinitionVersion, error) {
	raw, _, err := unstructured.NestedSlice(crd.object.Object, "spec", "additionalPrinterColumns")
	if err != nil {
		return CustomResourceDefinitionVersion{}, fmt.Errorf("unable to read crd .spec.additionalPrinterColumns: %w", err)
	}

	columns, err := crdV1beta1PrinterColumns(raw)
	if err != nil {
		return CustomResourceDefinitionVersion{}, fmt.Errorf("collect CRD printer columns: %w", err)
	}

	customResourceDefinitionVersion := CustomResourceDefinitionVersion{
		Version:        version,
		PrinterColumns: columns,
	}
	return customResourceDefinitionVersion, nil

}

func (crd *CustomResourceDefinition) versionNames() ([]string, error) {
	objects, err := crd.versions()
	if err != nil {
		return nil, err
	}

	var versions []string
	for i := range objects {
		versions = append(versions, objects[i]["name"].(string))
	}
	return versions, nil
}

func (crd *CustomResourceDefinition) versions() ([]map[string]interface{}, error) {
	versionsRaw, found, err := unstructured.NestedSlice(crd.object.Object, "spec", "versions")
	if err != nil {
		return nil, fmt.Errorf("unable to read crd .spec.versions: %w", err)
	}

	if !found {
		return nil, nil
	}

	var versions []map[string]interface{}

	for i := range versionsRaw {
		versions = append(versions, versionsRaw[i].(map[string]interface{}))
	}

	return versions, nil
}

func crdV1PrinterColumns(in interface{}) ([]CustomResourceDefinitionPrinterColumn, error) {
	if in == nil {
		return []CustomResourceDefinitionPrinterColumn{}, nil
	}

	rawList, ok := in.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unknown format for additional printer columns (%T)", in)
	}

	var columns []CustomResourceDefinitionPrinterColumn
	for i := range rawList {
		obj := rawList[i].(map[string]interface{})

		column := CustomResourceDefinitionPrinterColumn{
			Name:        mapString(obj, "name"),
			Type:        mapString(obj, "type"),
			Description: mapString(obj, "description"),
			JSONPath:    mapString(obj, "jsonPath"),
		}
		columns = append(columns, column)
	}

	return columns, nil
}

func crdV1beta1PrinterColumns(in interface{}) ([]CustomResourceDefinitionPrinterColumn, error) {
	if in == nil {
		return []CustomResourceDefinitionPrinterColumn{}, nil
	}

	rawList, ok := in.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unknown format for additional printer columns (%T)", in)
	}

	var columns []CustomResourceDefinitionPrinterColumn
	for i := range rawList {
		obj := rawList[i].(map[string]interface{})

		column := CustomResourceDefinitionPrinterColumn{
			Name:        mapString(obj, "name"),
			Type:        mapString(obj, "type"),
			Description: mapString(obj, "description"),
			JSONPath:    mapString(obj, "JSONPath"),
		}
		columns = append(columns, column)
	}

	return columns, nil
}

func mapString(m map[string]interface{}, key string) string {
	if m[key] == nil {
		return ""
	}

	if s, ok := m[key].(string); ok {
		return s
	}

	return ""

}

/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

package kubernetes

import (
	"fmt"
	"io"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	k8sJSON "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/clientcmd/api/latest"
)

// ReadObject reads an unstructured object from a reader.
func ReadObject(r io.Reader) (*unstructured.Unstructured, error) {
	d := yaml.NewYAMLOrJSONDecoder(r, 4096)
	ext := runtime.RawExtension{}
	if err := d.Decode(&ext); err != nil {
		return nil, fmt.Errorf("decode YAML: %w", err)
	}

	obj, _, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("decode YAML into object: %w", err)
	}

	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("obect is not unstructured (%T)", obj)
	}

	return u, nil
}

// SerializeToString serializes an object to a YAML string.
func SerializeToString(object runtime.Object) (string, error) {
	if object == nil {
		return "", fmt.Errorf("object is nil")
	}

	options := k8sJSON.SerializerOptions{
		Yaml:   true,
		Pretty: true,
		Strict: false,
	}
	yamlSerializer := k8sJSON.NewSerializerWithOptions(k8sJSON.DefaultMetaFactory, latest.Scheme, latest.Scheme, options)

	var sb strings.Builder
	if _, err := sb.WriteString("---\n"); err != nil {
		return "", err
	}
	if err := yamlSerializer.Encode(object, &sb); err != nil {
		return "", fmt.Errorf("encoding object as YAML: %w", err)
	}

	return sb.String(), nil
}

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

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	k8sJSON "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd/api/latest"
	k8sYAML "sigs.k8s.io/yaml"
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

// FromUnstructured converts an unstructured to a runtime object.
func FromUnstructured(u *unstructured.Unstructured, as interface{}) error {
	if _, ok := as.(*apiextv1.CustomResourceDefinition); ok {
		// CRDs do not work well with the scheme converter, so:
		//   1. convert in to JSON
		//   2. YAML unmarshal to object

		y, err := k8sYAML.Marshal(u)
		if err != nil {
			return fmt.Errorf("marshal unstructured to bytes: %w", err)
		}

		if err := k8sYAML.Unmarshal(y, as); err != nil {
			return err
		}

		return nil
	}

	if err := scheme.Scheme.Convert(u, as, nil); err != nil {
		return fmt.Errorf("scheme convert: %w", err)
	}

	if err := copyObjectMeta(as, u); err != nil {
		return fmt.Errorf("copy object metadata from unstructured: %w", err)
	}

	return nil
}

func copyObjectMeta(to interface{}, from *unstructured.Unstructured) error {
	object, ok := to.(metav1.Object)
	if !ok {
		return fmt.Errorf("%T is not a runtime object", to)
	}

	t, err := meta.TypeAccessor(object)
	if err != nil {
		return fmt.Errorf("accessing type meta: %w", err)
	}
	t.SetAPIVersion(from.GetAPIVersion())
	t.SetKind(from.GetObjectKind().GroupVersionKind().Kind)

	object.SetNamespace(from.GetNamespace())
	object.SetName(from.GetName())
	object.SetGenerateName(from.GetGenerateName())
	object.SetUID(from.GetUID())
	object.SetResourceVersion(from.GetResourceVersion())
	object.SetGeneration(from.GetGeneration())
	object.SetSelfLink(from.GetSelfLink())
	object.SetCreationTimestamp(from.GetCreationTimestamp())
	object.SetDeletionTimestamp(from.GetDeletionTimestamp())
	object.SetDeletionGracePeriodSeconds(from.GetDeletionGracePeriodSeconds())
	object.SetLabels(from.GetLabels())
	object.SetAnnotations(from.GetAnnotations())
	object.SetOwnerReferences(from.GetOwnerReferences())
	object.SetClusterName(from.GetClusterName())
	object.SetFinalizers(from.GetFinalizers())

	return nil
}

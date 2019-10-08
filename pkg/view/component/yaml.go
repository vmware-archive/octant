/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd/api/latest"

	"k8s.io/apimachinery/pkg/runtime"
	k8sJSON "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type YAMLConfig struct {
	Data string `json:"data,omitempty"`
}

type YAML struct {
	base
	Config YAMLConfig `json:"config,omitempty"`
}

func NewYAML(title []TitleComponent, data string) *YAML {
	return &YAML{
		base: newBase(typeYAML, title),
		Config: YAMLConfig{
			Data: data,
		},
	}
}

func (y *YAML) Data(object runtime.Object) error {
	yamlSerializer := k8sJSON.NewYAMLSerializer(k8sJSON.DefaultMetaFactory, latest.Scheme, latest.Scheme)

	var sb strings.Builder
	if _, err := sb.WriteString("---\n"); err != nil {
		return err
	}
	if err := yamlSerializer.Encode(object, &sb); err != nil {
		return errors.Wrap(err, "encoding object as YAML")
	}

	y.Config.Data = sb.String()

	return nil
}

// GetMetadata returns the component's metadata.
func (y *YAML) GetMetadata() Metadata {
	return y.Metadata
}

type yamlMarshal YAML

func (y *YAML) MarshalJSON() ([]byte, error) {
	m := yamlMarshal(*y)
	m.Metadata.Type = typeYAML
	return json.Marshal(&m)
}

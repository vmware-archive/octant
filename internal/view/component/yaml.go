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
	Metadata Metadata   `json:"metadata,omitempty"`
	Config   YAMLConfig `json:"config,omitempty"`
}

func NewYAML(title []TitleViewComponent) *YAML {
	return &YAML{
		Metadata: Metadata{
			Type:  "yaml",
			Title: title,
		},
		Config: YAMLConfig{},
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
	m.Metadata.Type = "yaml"
	return json.Marshal(&m)
}

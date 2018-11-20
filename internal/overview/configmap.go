package overview

import (
	"context"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/core"
)

type ConfigMapSummary struct{}

var _ View = (*ConfigMapSummary)(nil)

func NewConfigMapSummary() *ConfigMapSummary {
	return &ConfigMapSummary{}
}

func (cms *ConfigMapSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	configMap, err := retrieveConfigMap(object)
	if err != nil {
		return nil, err
	}

	detail, err := printConfigMapSummary(configMap)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

// ConfigMapDetails describe the details of a kubernetes core.ConfigMap
type ConfigMapDetails struct{}

// NewConfigMapDetails constructs a new ConfigMapDetails object
func NewConfigMapDetails() *ConfigMapDetails {
	return &ConfigMapDetails{}
}

// Content describes human readable object details
func (cm *ConfigMapDetails) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	configMap, ok := object.(*core.ConfigMap)
	if !ok {
		return nil, errors.Errorf("expected object to be a ConfigMap, it was %T", object)
	}

	table := content.NewTable("ConfigMap Data")
	table.Columns = []content.TableColumn{
		tableCol("Key"),
		tableCol("Value"),
	}

	for k, v := range configMap.Data {
		row := content.TableRow{
			"Key":   content.NewStringText(k),
			"Value": content.NewStringText(v),
		}
		table.AddRow(row)
	}

	return []content.Content{&table}, nil
}

func retrieveConfigMap(object runtime.Object) (*core.ConfigMap, error) {
	rc, ok := object.(*core.ConfigMap)
	if !ok {
		return nil, errors.Errorf("expected object to be a ConfigMap, it was %T", object)
	}

	return rc, nil
}
